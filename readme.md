# pngblog

Save images from the clipboard on macOS and paste the URL for Markdown.

1) Make a screenshot that is stored in the clipboard.
2) Call `pngblog`.
3) Command-Paste in your Markdown file.

## Install

### With `go install`

Requires Go 1.27 or newer:

```bash
go install github.com/megaproaktiv/pngblog@latest
```

This installs the `pngblog` binary into `$(go env GOPATH)/bin`. Make sure that
directory is on your `PATH`.

Install a specific version:

```bash
go install github.com/megaproaktiv/pngblog@v0.1.0
```

### From a release archive

Download the archive for your architecture from the
[releases page](https://github.com/megaproaktiv/pngblog/releases), extract it,
and move the `pngblog` binary onto your `PATH`.

### From source

```bash
git clone https://github.com/megaproaktiv/pngblog.git
cd pngblog
task install   # builds and copies to /usr/local/bin
```

## Usage

```bash
pngblog
Saved: /path/to/your/blog/static/img/2026/10/image-07.png
Clipboard: ![07](/img/2026/10/image-07.png)
```

Check the version:

```bash
pngblog --version
```

## Configuration

The save path and Markdown URL are composed from three configurable directory
components:

```
save path: <base>/<static>/<img>/<YYYY>/<MM>/image-<NN>.png
markdown:  ![<NN>](/<img>/<YYYY>/<MM>/image-<NN>.png)
```

Settings are resolved in this order of precedence:

1. Environment variables: `PNGBLOG_BASE`, `PNGBLOG_STATIC`, `PNGBLOG_IMG`.
2. A `pngblog.yml` file in the current directory, then `$HOME`.
3. Built-in defaults.

Copy the example and edit it (your local `pngblog.yml` is gitignored):

```bash
cp pngblog.yml.example pngblog.yml
```

```yaml
# pngblog.yml
base: /path/to/your/blog
static: static
img: img
```
