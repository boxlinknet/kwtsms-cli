# Changelog

All notable changes to kwtsms-cli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-13

### Added
- Initial release of kwtsms-cli
- Commands: `setup`, `balance`, `senderid`, `coverage`, `send`
- Pre-built binaries for Linux x64, ARM64, ARMv7, macOS Intel, macOS Apple Silicon, Windows x64
- Compile from source with plain `go build`, no CGO, no external dependencies
- Interactive setup wizard with credential verification and sender ID selection
- Config file at platform-appropriate location (`~/.config/kwtsms-cli/kwtsms-cli.toml`)
- Environment variable overrides: `KWTSMS_USERNAME`, `KWTSMS_PASSWORD`, `KWTSMS_SENDER`
- Inline flag overrides: `--username`, `--password`
- JSON output mode via `--json` flag on all commands
- Test mode via `--test` flag on send (queued, not delivered, no credits consumed)
- Bulk send: any number of recipients, auto-batched at 200 per API call with 500ms delay
- Duplicate phone number removal
- Append-only JSON send log, configured during setup
- Input sanitization: phone numbers, message text, sender IDs, config values
- Arabic-Indic digit conversion for phone numbers
- Phone number normalisation: strips `+`, `00` prefix, spaces, dashes
- Network error messages with no internal details exposed
- `--version` flag
