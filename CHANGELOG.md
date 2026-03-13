# Changelog

All notable changes to kwtsms-cli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-03-13

### Added
- Bulk send: recipients exceeding 200 are split into batches of 200 automatically, with a 500ms delay between batches. Output aggregates totals across all batches.
- `--version` flag: prints the current version and exits.
- Duplicate phone number deduplication: repeated numbers in `--to` are sent only once.

### Fixed
- Phone numbers with internal spaces (e.g. `965 9876 5432`) are now handled correctly as a single number rather than being split into multiple invalid tokens.
- Error messages no longer appear multiple times on API or network failures.
- Usage text is no longer printed on API or runtime errors.

## [0.1.0] - 2026-03-13

### Added
- Initial release of kwtsms-cli
- Commands: setup, balance, senderid, coverage, send
- Pre-built binaries for Linux x64, ARM64, ARMv7, macOS Intel, macOS Apple Silicon, Windows x64
- Compile from source support with plain `go build`
- Config file with platform-appropriate location
- Environment variable and inline flag credential overrides
- JSON output mode with `--json` flag
- Test mode with `--test` flag on send command
- Input sanitization for phone numbers, message text, and sender IDs
- Arabic-Indic digit conversion for phone numbers
