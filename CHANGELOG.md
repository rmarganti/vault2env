# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

-   refactor: Less redundant redundancy in piping detection.

## [1.0.0] - 2023-04-27

### Added

-   feat: JSON schema definition for config file.
-   feat: Input and output can be piped.
-   ci: linting

### Changed

-   feat!: Sources are referenced by URIs.
-   feat!: V2 of the config schema.
-   ci: Release GH action uses a PAT instead of a GITHUB_TOKEN.

### Fixed

-   fix: Correct param name in missing `secretPath` error.

## [0.1.1] - 2023-04-19

### Added

-   feat: Homebrew tap

## [0.1.0] - 2023-04-19

### Added

-   feat: `--config` option for specifying a config file

## [0.0.1] - 2023-04-18

### Added

-   feat: Fetch secrets from a Vault path and write them to an `.env` file.
-   feat: Push secrets from an `.env` file to a Vault path.
-   ci: GH actions and GoReleaser
