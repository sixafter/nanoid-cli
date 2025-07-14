# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Date format: `YYYY-MM-DD`

---
## [Unreleased]

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.24.0] - 2025-07-13

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanodid@v1.37.0](https://github.com/sixater/nanoid/releases/tag/v1.37.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.23.0] - 2025-07-12

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.36.0](https://github.com/sixafter/nanoid/releases/tag/v1.36.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.22.0] - 2025-07-10

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.35.0](https://github.com/sixafter/nanoid/releases/tag/v1.35.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.21.0] - 2025-07-09

### Added
- **debt:** Modified [README](../README.md) to include a note for macOS users regarding the `killed` error when running `nanoid` due to quarantine issues.

### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.34.0](https://github.com/sixafter/nanoid/releases/tag/v1.34.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.20.0] - 2025-06-29

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.33.0](https://github.com/sixafter/nanoid/releases/tag/v1.33.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.19.0] - 2025-06-28

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.32.1](https://github.com/sixafter/nanoid/releases/tag/v1.32.1).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.18.4] - 2025-06-15

### Added
### Changed
### Deprecated
- **debt:** Support for both Homebrew formula and casks given Go Releaser deprecation of the [`brews:`](https://goreleaser.com/customization/homebrew_formulas/) stanza in favor of [`homebrew_casks:`](https://goreleaser.com/customization/homebrew_casks/). 

### Removed
### Fixed
### Security

---
## [1.18.3] - 2025-06-14

### Added
### Changed
### Deprecated
- **debt:** Need to evaluate the Go Releaser deprecation of the `brews:` stanza in favor of [`homebrew_casks:`](https://goreleaser.com/customization/homebrew_casks/). Has to revert back to the `brews:` stanza due to compatibility issues.

### Removed
### Fixed
### Security

---
## [1.18.0] - 2025-06-14

### Added
- **debt:** Added [cosign](https://github.com/sigstore/cosign-installer) signature verification steps in [README.md](../README.md).
- **debt:** Removed deprecated `brews` stanza in go-releaser workflow.

### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.17.0] - 2025-06-14

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
- **risk:** Add signature file to each release.


---
## [1.16.0] - 2025-06-14

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
- **risk:** Sign releases with [cosign](https://github.com/sigstore/cosign-installer).

---
## [1.15.0] - 2025-06-13

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.30.0](https://github.com/sixafter/nanoid/releases/tag/v1.30.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.14.0] - 2025-06-12

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.29.0](https://github.com/sixafter/nanoid/releases/tag/v1.29.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.13.0] - 2025-05-31

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.28.0](https://github.com/sixafter/nanoid/releases/tag/v1.28.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.12.0] - 2025-05-16

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.27.0](https://github.com/sixafter/nanoid/releases/tag/v1.27.0).
- **debt:** Upgraded all Go dependencies to the latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.11.0] - 2025-05-16

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.26.0](https://github.com/sixafter/nanoid/releases/tag/v1.26.0).
- **debt:** Upgraded all Go dependencies to the latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.10.0] - 2025-05-13

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.25.0](https://github.com/sixafter/nanoid/releases/tag/v1.25.0).
- **debt:** Upgraded all Go dependencies to the latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.9.0] - 2025-05-04

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.24.1](https://github.com/sixafter/nanoid/releases/tag/v1.24.1).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.8.1] - 2025-04-14

### Added
### Changed
- **debt:** Address the deprecation notice for the Go Releaser `archives.format` as emitted in the release pipeline [here](https://github.com/sixafter/nanoid-cli/actions/runs/14455567265/job/40537917390#step:7:27).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.8.0] - 2025-04-14

### Added
- **feature:** Added various statistics to the `generate` command output when using the `--verbose` flag.

### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.24.0](https://github.com/sixafter/nanoid/releases/tag/v1.24.0).
- **debt:** Upgraded all dependencies to the latest versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.7.1] - 2025-03-20

### Added
### Changed
- **debt:** Upgraded all Go dependencies to the latest versions.
- **debt:** Upgraded all CI dependencies to the latest versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.7.0] - 2025-02-13

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.23.0](https://github.com/sixafter/nanoid/releases/tag/v1.23.0).
- **debt:** Upgraded all Go dependencies to the latest versions.
- **debt:** Upgraded to Go 1.24.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.6.0] - 2024-12-26

### Added
### Changed
- **debt:** Upgraded to [sixafter/nanoid@v1.22.0](https://github.com/sixafter/nanoid/releases/tag/v1.22.0).
- **debt:** Upgraded all Go dependencies to the latest versions.
- **debt:** Upgraded the CI pipeline to use the new GitHub Action for SonarQube Cloud analysis.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.5.0] - 2024-12-07

### Added
### Changed
- **DEBT:** Upgraded to [sixafter/nanoid@v1.21.0](https://github.com/sixafter/nanoid/releases/tag/v1.21.0).
- **DEBT:** Upgraded all Go dependencies to the latest versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.4.1] - 2024-11-24

### Added
### Changed
- **DEBT:** Upgraded to [sixafter/nanoid@v1.20.1](https://github.com/sixafter/nanoid/releases/tag/v1.20.1).
- **DEBT:** Upgraded all Go dependencies to the latest versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.4.0] - 2024-11-16

### Added
### Changed
- **DEBT:** Upgraded to [sixafter/nanoid@v1.20.0](https://github.com/sixafter/nanoid/releases/tag/v1.20.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.3.0] - 2024-11-16

### Added
### Changed
- **DEBT:** Upgraded to [sixafter/nanoid@v1.19.0](https://github.com/sixafter/nanoid/releases/tag/v1.19.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.2.0] - 2024-11-15

### Added
### Changed
- **DEBT:** Added missing license header to the CodeQL analysis configuration file.
- **DEBT:** Refactored CHANGELOG date format to `YYYY-MM-DD`.
- **DEBT:** Upgraded to [sixafter/nanoid@v1.18.1](https://github.com/sixafter/nanoid/releases/tag/v1.18.1).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.1.0] - 2024-11-14

### Added
### Changed
- **DEBT:** Upgraded to [sixafter/nanoid@v1.17.3](https://github.com/sixafter/nanoid/releases/tag/v1.17.3).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.1.0] - 2024-11-14

### Added
### Changed
- **DEBT:** Upgraded to [sixafter/nanoid@v1.17.3](https://github.com/sixafter/nanoid/releases/tag/v1.17.3).
### Deprecated
### Removed
### Fixed
### Security

---
## [1.0.1] - 2024-11-14

### Added
- **FEATURE:** Added support for Homebrew.
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.0.0] - 2024-11-13

### Added
- **FEATURE:** Initial commit.
### Changed
### Deprecated
### Removed
### Fixed
### Security

[Unreleased]: https://github.com/sixafter/nanoid-cli/compare/v1.24.0...HEAD
[1.24.0]: https://github.com/sixafter/nanoid-cli/compare/v1.23.0...v1.24.0
[1.23.0]: https://github.com/sixafter/nanoid-cli/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/sixafter/nanoid-cli/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/sixafter/nanoid-cli/compare/v1.20.0...v1.21.0
[1.20.0]: https://github.com/sixafter/nanoid-cli/compare/v1.19.0...v1.20.0
[1.19.0]: https://github.com/sixafter/nanoid-cli/compare/v1.18.4...v1.19.0
[1.18.4]: https://github.com/sixafter/nanoid-cli/compare/v1.18.3...v1.18.4
[1.18.3]: https://github.com/sixafter/nanoid-cli/compare/v1.18.0...v1.18.3
[1.18.0]: https://github.com/sixafter/nanoid-cli/compare/v1.17.0...v1.18.0
[1.17.0]: https://github.com/sixafter/nanoid-cli/compare/v1.16.0...v1.17.0
[1.16.0]: https://github.com/sixafter/nanoid-cli/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/sixafter/nanoid-cli/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/sixafter/nanoid-cli/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/sixafter/nanoid-cli/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/sixafter/nanoid-cli/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/sixafter/nanoid-cli/compare/v1.10.0...v1.11.0
[1.10.0]: https://github.com/sixafter/nanoid-cli/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/sixafter/nanoid-cli/compare/v1.8.1...v1.9.0
[1.8.1]: https://github.com/sixafter/nanoid-cli/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/sixafter/nanoid-cli/compare/v1.7.1...v1.8.0
[1.7.1]: https://github.com/sixafter/nanoid-cli/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/sixafter/nanoid-cli/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/sixafter/nanoid-cli/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/sixafter/nanoid-cli/compare/v1.4.1...v1.5.0
[1.4.1]: https://github.com/sixafter/nanoid-cli/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/sixafter/nanoid-cli/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/sixafter/nanoid-cli/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/sixafter/nanoid-cli/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/sixafter/nanoid-cli/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/sixafter/nanoid-cli/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/sixafter/nanoid-cli/compare/a6a1eb74b61e518fd0216a17dfe5c9b4c432e6e8...v1.0.0

[MUST]: https://datatracker.ietf.org/doc/html/rfc2119
[MUST NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[SHOULD]: https://datatracker.ietf.org/doc/html/rfc2119
[SHOULD NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[MAY]: https://datatracker.ietf.org/doc/html/rfc2119
[SHALL]: https://datatracker.ietf.org/doc/html/rfc2119
[SHALL NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[REQUIRED]: https://datatracker.ietf.org/doc/html/rfc2119
[RECOMMENDED]: https://datatracker.ietf.org/doc/html/rfc2119
[NOT RECOMMENDED]: https://datatracker.ietf.org/doc/html/rfc2119
