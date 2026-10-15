# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-10-15

### Added

- Initial release.
- Read a PNG image from the macOS clipboard and save it to a dated path
  (`<base>/<static>/<img>/<YYYY>/<MM>/image-<NN>.png`) with an auto-incremented
  file number.
- Copy a Markdown image reference (`![<NN>](/<img>/<YYYY>/<MM>/image-<NN>.png)`)
  back to the clipboard.
- Configurable `base`, `static`, and `img` directory components via Viper, with
  precedence: environment variables (`PNGBLOG_BASE`, `PNGBLOG_STATIC`,
  `PNGBLOG_IMG`) > `pngblog.yml` (current directory, then `$HOME`) > defaults.
- `pngblog.yml.example` template; local `pngblog.yml` is gitignored.
- `--version` flag reporting the build version.

[0.1.0]: https://github.com/megaproaktiv/pngblog/releases/tag/v0.1.0
