# kwtsms-cli

A command-line interface for the [kwtSMS](https://www.kwtsms.com) SMS gateway. Send SMS messages, check your account balance, and manage sender IDs, all from your terminal or as part of an automated workflow.

Designed for developers, DevOps engineers, and AI agents that need to integrate SMS delivery into scripts, pipelines, and automation without writing custom code.

## About kwtSMS

[kwtSMS](https://www.kwtsms.com) is a Kuwait-based SMS gateway trusted by businesses to deliver messages across Kuwait (Zain, Ooredoo, STC, Virgin) and internationally. It offers private Sender IDs, free API testing, non-expiring credits, and competitive flat-rate pricing. Open a free account in under one minute at [kwtsms.com/signup](https://www.kwtsms.com/signup/), no paperwork or payment required.

- Website: [www.kwtsms.com](https://www.kwtsms.com)
- Dashboard: [www.kwtsms.com/login](https://www.kwtsms.com/login)
- Support: [www.kwtsms.com/support.html](https://www.kwtsms.com/support.html)
- FAQ: [www.kwtsms.com/faq/](https://www.kwtsms.com/faq/)

## Use Cases

- **SMS automation:** Send notifications, alerts, and reminders from shell scripts and cron jobs.
- **AI agent integration:** Provide AI agents and LLM-powered tools with the ability to send SMS messages as a tool call.
- **CI/CD pipelines:** Trigger SMS alerts on deployment, test failure, or system events.
- **OTP delivery:** Send one-time passwords from any script or backend service.
- **Balance monitoring:** Check remaining credits as part of a scheduled health check.

## Installation

### Direct download

Download the pre-built binary for your platform. No runtime required, just download and run.

| Platform | Download |
|---|---|
| Linux x64 | [kwtsms-cli-linux-x64](https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-linux-x64) |
| Linux ARM64 (Raspberry Pi 4/5) | [kwtsms-cli-linux-arm64](https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-linux-arm64) |
| Linux ARMv7 (Raspberry Pi 2/3) | [kwtsms-cli-linux-armv7](https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-linux-armv7) |
| macOS Intel | [kwtsms-cli-macos-x64](https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-macos-x64) |
| macOS Apple Silicon | [kwtsms-cli-macos-arm64](https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-macos-arm64) |
| Windows x64 | [kwtsms-cli-windows-x64.exe](https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-windows-x64.exe) |

**Linux/macOS — install to PATH:**
```bash
curl -Lo kwtsms-cli https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-linux-x64
chmod +x kwtsms-cli
sudo mv kwtsms-cli /usr/local/bin/
```

**Raspberry Pi (ARMv7):**
```bash
curl -Lo kwtsms-cli https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/kwtsms-cli-linux-armv7
chmod +x kwtsms-cli
sudo mv kwtsms-cli /usr/local/bin/
```

**Windows:** Download the `.exe` file and place it in a folder that is in your `PATH`.

### Compile from source

If you prefer to build from source, you only need [Go](https://go.dev/dl/) installed. No other tools or system libraries are required.

```bash
git clone https://github.com/boxlinknet/kwtsms-cli
cd kwtsms-cli
go build -o kwtsms-cli .
```

This works on any platform that Go supports, including Linux, macOS, Windows, and Raspberry Pi.

## Quick Start

```bash
# Step 1: configure your API credentials
kwtsms-cli setup

# Step 2: send your first message
kwtsms-cli send --to 96598765432 --message "Hello from kwtsms-cli"

# Check your balance
kwtsms-cli balance
```

Your API credentials are stored in `~/.config/kwtsms-cli/kwtsms-cli.toml` with restricted permissions. You can get your API credentials from your [kwtSMS account](https://www.kwtsms.com/account/api/).

## Commands

### `kwtsms-cli setup`

Interactive setup wizard. Prompts for your API username and password, verifies them against the API, lets you choose a default sender ID, and configures a log file. Writes the config file to the correct location for your operating system.

```bash
kwtsms-cli setup
```

Run this once before using any other command.

### `kwtsms-cli balance`

Display your current SMS credit balance.

```bash
kwtsms-cli balance
kwtsms-cli balance --json
```

```
Available:  1,234
Purchased:  5,000
```

### `kwtsms-cli senderid`

List all sender IDs approved on your account.

```bash
kwtsms-cli senderid
kwtsms-cli senderid --json
```

### `kwtsms-cli coverage`

List active country prefixes available for sending on your account.

```bash
kwtsms-cli coverage
kwtsms-cli coverage --json
```

### `kwtsms-cli send`

Send an SMS to one or more recipients. Accepts any number of recipients: batches larger than 200 are split automatically and sent with a short delay between batches.

```bash
# Single recipient
kwtsms-cli send --to 96598765432 --message "Your verification code is 4821"

# Multi-line message (use $'...' in bash for \n to be interpreted as a newline)
kwtsms-cli send --to 96598765432 --message $'Order confirmed\nTracking: TRK-12345\nExpected: Tomorrow'

# Multiple recipients (comma-separated)
kwtsms-cli send --to 96598765432,96512345678 --message "System maintenance tonight at 10pm"

# Bulk send: more than 200 numbers are batched automatically (200 per API call, 500ms between batches)
# Practical limits per OS: Linux/macOS ~150,000 numbers, Windows ~2,500 numbers
kwtsms-cli send --to 96550000001,96550000002,...,96550000250 --message "Announcement"

# Specify a sender ID
kwtsms-cli send --to 96598765432 --message "Your order is ready" --sender MY-BRAND

# Test mode: message is queued but not delivered, no credits consumed
kwtsms-cli send --to 96598765432 --message "Test" --test
```

```
Sent
Numbers:    1
Charged:    1
Balance:    1,234
MsgID:      f4c841adee210f31307633ceaebff2ec
```

For bulk sends spanning multiple batches, Numbers and Charged are aggregated and a MsgID is shown for each batch.

**Bulk recipient limits (command-line argument size):**

| Platform | Estimated max recipients |
|---|---|
| Linux / macOS | ~150,000 |
| Windows | ~2,500 |

For lists larger than these limits, use a script to split input and call `kwtsms-cli send` in chunks.

**Flags:**

| Flag | Required | Description |
|---|---|---|
| `-t, --to` | Yes | Recipient number(s), comma-separated. Duplicates removed automatically. |
| `-m, --message` | Yes | Message text |
| `-s, --sender` | No | Sender ID (overrides config default) |
| `--test` | No | Queue without delivery, no credits used |

## Global Flags

These flags work on every command.

| Flag | Description |
|---|---|
| `--json` | Output the raw API response as JSON, suitable for piping and scripting |
| `--config PATH` | Use a different config file instead of the default |
| `--username VALUE` | Override the API username for this call only |
| `--password VALUE` | Override the API password for this call only |
| `--version` | Print version and exit |

## Configuration

### Config file

`kwtsms-cli setup` creates the config file automatically. You can also create or edit it manually.

| Platform | Location |
|---|---|
| Linux / macOS | `~/.config/kwtsms-cli/kwtsms-cli.toml` |
| Windows | `%APPDATA%\kwtsms-cli\kwtsms-cli.toml` |

```toml
username = "myapiuser"
password = "myapipass"
sender   = "MY-SENDER"
log_file = "kwtsms-cli.log"
```

The `log_file` value is a filename relative to the directory where you run the binary. Omit it or leave it empty to disable logging.

### Environment variables

Useful for CI/CD pipelines and containerised environments where you want to inject credentials at runtime rather than storing them in a file.

```bash
export KWTSMS_USERNAME=myapiuser
export KWTSMS_PASSWORD=myapipass
export KWTSMS_SENDER=MY-SENDER
```

### Credential priority

When the same value is set in multiple places, the highest priority source wins:

```
--username / --password flags  (highest)
KWTSMS_USERNAME / KWTSMS_PASSWORD env vars
config file                    (lowest)
```

## Phone Number Formats

All of the following formats are accepted. Numbers are automatically normalised before sending.

```
+96598765432     accepted, + prefix stripped
0096598765432    accepted, 00 prefix stripped
965 9876 5432    accepted, spaces stripped
965-9876-5432    accepted, dashes stripped
96598765432      correct format
```

Arabic-Indic digits (`٩٦٥...`) are also accepted and converted automatically.

## AI Agent and Automation Usage

kwtsms-cli is built for use in automated environments. Pair `--json` with any command for structured output that is easy to parse.

**Send and capture the message ID:**
```bash
result=$(kwtsms-cli send --to 96598765432 --message "Alert: deploy complete" --json)
msg_id=$(echo "$result" | jq -r '."msg-id"')
```

**Check balance in a script:**
```bash
available=$(kwtsms-cli balance --json | jq '.available')
if [ "$available" -lt 100 ]; then
  echo "Low balance: $available credits remaining"
fi
```

**Use in a Makefile, cron job, or GitHub Actions workflow:**
```yaml
- name: Notify team via SMS
  run: kwtsms-cli send --to 96598765432 --message "Deploy complete: ${{ env.VERSION }}"
  env:
    KWTSMS_USERNAME: ${{ secrets.KWTSMS_USERNAME }}
    KWTSMS_PASSWORD: ${{ secrets.KWTSMS_PASSWORD }}
    KWTSMS_SENDER: ${{ secrets.KWTSMS_SENDER }}
```

## Support

- Dashboard: [www.kwtsms.com/login](https://www.kwtsms.com/login)
- General support and account questions: [www.kwtsms.com/support.html](https://www.kwtsms.com/support.html)
- FAQ: [www.kwtsms.com/faq/](https://www.kwtsms.com/faq/)
- Sender ID help: [www.kwtsms.com/sender-id-help.html](https://www.kwtsms.com/sender-id-help.html)

For security vulnerability reports, see [SECURITY.md](SECURITY.md).

## License

MIT
