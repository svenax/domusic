<!-- markdownlint-disable MD024 -->

# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [2.1.0] - 2025-10-19

### Added

- Add `sync` command to sync changes to `svenax.net`.

## [2.0.0] - 2025-10-06

### Changed

- Replace `cobra` with `urfave/cli`.
- Replace `viper` with custom configuration code.

### Removed

- Remove `copy` command. Forscore is better handled by using generated PDFs directly.
- Remove remaining Evernote artifacts since the `upload` command was removed previously.

## [1.7.1] - 2025-01-13

### Changed

- Use `os` instead of `io/ioutil`.
- Update requirements.

### Fixed

- Update `Makefile` and `version.go` so the correct information is injected
  at build time.

## [1.6.0] - 2023-03-12

### Added

- Finished the `copy` command that copies and renames PDF files for Forscore.

## [1.5.1] - 2023-01-28

### Changed

- Updated to Lilypond 2.24.

## [1.4.3] - 2020-01-12

### Fixed

- Updated build dependencies.

## [1.4.2] - 2018-11-05

### Changed

- `upload` now only searches in the set notebook, not everywhere.

## [1.4.1] - 2018-11-05

### Changed

- `upload` can now take multiple files and wildcards in the same way as `collection` and `make`.

## [1.4.0] - 2018-11-03

### Added

- New command `collection` that generates a .ly file with all the files given as argument.

## [1.3.0] - 2018-09-04

### Added

- New command `upload` that given a pdf file creates or updates a note on Evernote.

## [1.2.0] - 2018-05-22

### Added

- New flag `--post` as alias for `-tpng --root --crop`.

## [1.1.0] - 2018-05-22

### Added

- New flag `--keep` to keep all generated files in the music root.

## [1.0.0] - 2018-05-22

### Added

- Initial commit of all functionality.
