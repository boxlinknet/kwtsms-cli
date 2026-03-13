# Changelog

All notable changes to kwtsms-cli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of kwtsms-cli
- Commands: setup, balance, senderid, coverage, send, validate
- Pre-built binaries for Linux x64, ARM64, ARMv7, macOS Intel, macOS Apple Silicon, Windows x64
- Compile from source support with plain `go build`
- Config file with platform-appropriate location
- Environment variable and inline flag credential overrides
- JSON output mode with `--json` flag
- Test mode with `--test` flag on send command
- Input sanitization for phone numbers, message text, and sender IDs
- Arabic-Indic digit conversion for phone numbers
