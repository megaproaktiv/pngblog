// Command pngblog reads a PNG image from the system clipboard, saves it into the
// blog's static image tree using a dated path and an auto-incremented file name,
// and copies a Markdown image reference back to the clipboard.
//
// It is a self-contained Go binary: clipboard access is handled by
// golang.design/x/clipboard, so there are no runtime dependencies on external
// tools such as pngpaste, osascript, or pbcopy.
//
// Save path:
//
//	<base>/<static>/<img>/<YYYY>/<MM>/image-<NN>.png
//
// Clipboard output (Markdown):
//
//	![<NN>](/<img>/<YYYY>/<MM>/image-<NN>.png)
//
// Configuration:
//
// The base, static, and img directory components are configurable via Viper.
// Settings are read, in order of precedence, from:
//
//  1. Environment variables (prefix PNGBLOG, e.g. PNGBLOG_BASE, PNGBLOG_STATIC,
//     PNGBLOG_IMG).
//  2. A pngblog.yml file in the current directory or $HOME.
//  3. Built-in defaults.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/spf13/viper"
	"golang.design/x/clipboard"
)

// version is the build version. It is overridden at release time via
// -ldflags "-X main.version=...". For `go install` builds it falls back to the
// module version embedded in the binary (see resolveVersion).
var version = "dev"

// Config holds the resolved directory components used to build the save path
// and the Markdown reference.
type Config struct {
	Base   string // filesystem root of the blog's assets
	Static string // directory under base holding static assets
	Img    string // directory under static holding images
}

func main() {
	for _, a := range os.Args[1:] {
		if a == "--version" || a == "-v" || a == "version" {
			fmt.Printf("pngblog %s\n", resolveVersion())
			return
		}
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "pngblog: "+err.Error())
		os.Exit(1)
	}
}

// resolveVersion returns the ldflags-injected version when set, otherwise the
// module version recorded by `go install module@version`.
func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

// loadConfig resolves the base/static/img directory components using Viper.
// Precedence: environment variables (PNGBLOG_*) > pngblog.yml > defaults.
func loadConfig() (Config, error) {
	v := viper.New()

	// Built-in defaults.
	v.SetDefault("base", "/Users/gernotglawe/letsbuild/aws-blog-de-base")
	v.SetDefault("static", "static")
	v.SetDefault("img", "img")

	// Environment variables, e.g. PNGBLOG_BASE, PNGBLOG_STATIC, PNGBLOG_IMG.
	v.SetEnvPrefix("pngblog")
	v.AutomaticEnv()

	// Config file: pngblog.yml in the current directory or $HOME.
	v.SetConfigName("pngblog")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(home)
	}

	// A missing config file is not an error; defaults and env vars still apply.
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, fmt.Errorf("reading config file: %w", err)
		}
	}

	return Config{
		Base:   v.GetString("base"),
		Static: v.GetString("static"),
		Img:    v.GetString("img"),
	}, nil
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Initialize clipboard access. Fails if there is no clipboard available
	// (for example a headless session with no display).
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("initializing clipboard: %w", err)
	}
	ctx := context.Background()

	// 1. Read PNG bytes from the clipboard. FmtImage yields PNG-encoded data.
	png, err := clipboard.Read(ctx, clipboard.FmtImage)
	if err != nil {
		if errors.Is(err, clipboard.ErrNoData) {
			return fmt.Errorf("clipboard does not contain an image")
		}
		return fmt.Errorf("reading image from clipboard: %w", err)
	}
	if len(png) == 0 {
		return fmt.Errorf("clipboard does not contain an image")
	}

	// 2. Compute the dated target directory and the next free file name.
	now := time.Now()
	year := now.Format("2006") // 4-digit year, e.g. 2026
	month := now.Format("01")  // zero-padded month, e.g. 01, 10

	destDir := filepath.Join(cfg.Base, cfg.Static, cfg.Img, year, month)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	nn, err := nextFreeNumber(destDir)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("image-%02d.png", nn)
	destPath := filepath.Join(destDir, fileName)

	// 3. Save the PNG.
	if err := os.WriteFile(destPath, png, 0o644); err != nil {
		return fmt.Errorf("writing image: %w", err)
	}

	// 4. Build the Markdown reference and copy it to the clipboard. The URL path
	// omits the base and static components (Hugo serves static/ at the site
	// root) and uses the configured img segment.
	label := fmt.Sprintf("%02d", nn)
	markdown := fmt.Sprintf("![%s](/%s/%s/%s/%s)", label, cfg.Img, year, month, fileName)

	if _, err := clipboard.Write(ctx, clipboard.FmtText, []byte(markdown)); err != nil {
		return fmt.Errorf("copying markdown to clipboard: %w", err)
	}

	fmt.Printf("Saved: %s\n", destPath)
	fmt.Printf("Clipboard: %s\n", markdown)
	return nil
}

// nextFreeNumber returns the lowest NN >= 1 for which image-<NN>.png does not
// yet exist in dir.
func nextFreeNumber(dir string) (int, error) {
	for nn := 1; nn <= 9999; nn++ {
		p := filepath.Join(dir, fmt.Sprintf("image-%02d.png", nn))
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return nn, nil
		}
		if err != nil {
			return 0, fmt.Errorf("checking %s: %w", p, err)
		}
	}
	return 0, fmt.Errorf("no free image number found (exceeded 9999)")
}
