# kwtsms-cli

Command-line interface for the [kwtSMS](https://kwtsms.com) SMS gateway. Send SMS, check balances, validate phone numbers, and manage sender IDs from your terminal or scripts.

## Installation

### Download binary (recommended)

Download the pre-built binary for your platform from the [releases page](https://github.com/boxlinknet/kwtsms-cli/releases/latest):

| Platform | File |
|---|---|
| Linux x64 | `kwtsms-cli-linux-x64` |
| Linux ARM64 (Raspberry Pi 4/5) | `kwtsms-cli-linux-arm64` |
| Linux ARMv7 (Raspberry Pi 2/3) | `kwtsms-cli-linux-armv7` |
| macOS Intel | `kwtsms-cli-macos-x64` |
| macOS Apple Silicon | `kwtsms-cli-macos-arm64` |
| Windows x64 | `kwtsms-cli-windows-x64.exe` |

**Linux/macOS:**
```bash
chmod +x kwtsms-cli-linux-x64
sudo mv kwtsms-cli-linux-x64 /usr/local/bin/kwtsms-cli
```

**Windows:** Download `.exe`, add to a folder in your `PATH`.

### Compile from source

Requires [Go](https://go.dev/dl/) installed (any platform):

```bash
git clone https://github.com/boxlinknet/kwtsms-cli
cd kwtsms-cli
go build -o kwtsms-cli .
```

No CGO. No external system libraries. Works on any Go-supported OS and architecture.

## Quick Start

```bash
# Configure credentials (interactive)
kwtsms-cli setup

# Send an SMS
kwtsms-cli send --to 96598765432 --message "Hello from kwtsms-cli"

# Check balance
kwtsms-cli balance
```

## Commands

### `kwtsms-cli setup`

Interactive wizard to create your config file. Verifies credentials and lets you pick a default sender ID.

```bash
kwtsms-cli setup
```

### `kwtsms-cli balance`

Show current account balance.

```bash
kwtsms-cli balance
kwtsms-cli balance --json
```

Output:
```
Available:  1,234
Purchased:  5,000
```

### `kwtsms-cli senderid`

List approved sender IDs on your account.

```bash
kwtsms-cli senderid
kwtsms-cli senderid --json
```

### `kwtsms-cli coverage`

List active country coverage and prefixes.

```bash
kwtsms-cli coverage
kwtsms-cli coverage --json
```

### `kwtsms-cli send`

Send an SMS message.

```bash
# Single recipient
kwtsms-cli send --to 96598765432 --message "Your message here"

# Multiple recipients (comma-separated)
kwtsms-cli send --to 96598765432,96512345678 --message "Broadcast message"

# Custom sender ID
kwtsms-cli send --to 96598765432 --message "Hello" --sender MY-BRAND

# Test mode (queued but not delivered, no credits consumed)
kwtsms-cli send --to 96598765432 --message "Test" --test
```

Output:
```
Sent
Numbers:    1
Charged:    1
Balance:    1,234
MsgID:      f4c841adee210f31307633ceaebff2ec
```

**Flags:**
- `-t, --to` - phone number(s), comma-separated (required)
- `-m, --message` - message text (required)
- `-s, --sender` - sender ID (overrides config default)
- `--test` - test mode, no delivery

### `kwtsms-cli validate`

Validate one or more phone numbers before sending.

```bash
# Space-separated
kwtsms-cli validate 96598765432 96512345678

# Comma-separated
kwtsms-cli validate 96598765432,96512345678

# Mixed
kwtsms-cli validate 96598765432,96512345678 96599999

# From file
cat numbers.txt | xargs kwtsms-cli validate
```

Output:
```
Valid:    96598765432
Invalid:  123
NoRoute:  none
```

## Global Flags

Available on all commands:

| Flag | Description |
|---|---|
| `--json` | Output raw API response as JSON |
| `--config PATH` | Override config file path |
| `--username VALUE` | Override API username |
| `--password VALUE` | Override API password |

## Configuration

### Config file

Created by `kwtsms-cli setup`. Location:
- **Linux/macOS:** `~/.config/kwtsms-cli/config.toml`
- **Windows:** `%APPDATA%\kwtsms-cli\config.toml`

```toml
username = "myapiuser"
password = "myapipass"
sender   = "MY-SENDER"
```

### Environment variables

Override config file values:

```bash
export KWTSMS_USERNAME=myapiuser
export KWTSMS_PASSWORD=myapipass
export KWTSMS_SENDER=MY-SENDER
```

### Priority

Inline flags override env vars, env vars override config file.

## Phone Number Formats

All of these are accepted and normalized automatically:

```
+96598765432    (+ prefix stripped)
0096598765432   (00 prefix stripped)
965 9876 5432   (spaces stripped)
965-9876-5432   (dashes stripped)
96598765432     (correct format)
```

## Security

All credentials are stored in a config file with user-only read permissions (`0600` on Unix).
Credentials are never logged or exposed in output.

To report a security vulnerability, contact: support@kwtsms.com

## License

MIT
