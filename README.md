# kwtsms-cli

[![Release](https://img.shields.io/github/v/release/boxlinknet/kwtsms-cli)](https://github.com/boxlinknet/kwtsms-cli/releases/latest)
[![Build](https://img.shields.io/github/actions/workflow/status/boxlinknet/kwtsms-cli/release.yml?label=build)](https://github.com/boxlinknet/kwtsms-cli/actions)
[![Downloads](https://img.shields.io/github/downloads/boxlinknet/kwtsms-cli/total)](https://github.com/boxlinknet/kwtsms-cli/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/boxlinknet/kwtsms-cli)](https://go.dev)
[![License](https://img.shields.io/github/license/boxlinknet/kwtsms-cli)](https://github.com/boxlinknet/kwtsms-cli/blob/main/LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-blue)](https://github.com/boxlinknet/kwtsms-cli/releases/latest)

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

Requires only [Go](https://go.dev/dl/) installed. No other tools or system libraries needed.

```bash
git clone https://github.com/boxlinknet/kwtsms-cli
cd kwtsms-cli
go build -o kwtsms-cli .
```

Works on any platform Go supports: Linux, macOS, Windows, and Raspberry Pi.

## Quick Start

```bash
# Step 1: configure your credentials
kwtsms-cli setup

# Step 2: send your first message
kwtsms-cli send --to 96598765432 --message "Hello from kwtsms-cli"

# Check your balance
kwtsms-cli balance
```

## Commands

### `kwtsms-cli setup`

Interactive wizard. Run once before using any other command.

Prompts for API username, password, default sender ID, and log file. Verifies credentials against the API before saving anything. Writes the config file to the platform-appropriate location.

```
kwtSMS CLI Setup
----------------
API Username: myapiuser
API Password: ••••••••••••••

Verifying credentials...

Available sender IDs:
  [1] KWT-SMS
  [2] MY-BRAND

Select default sender ID [1]: 2

Log file [kwtsms-cli.log]
(Enter for default, type path to change, "none" to disable):

Config saved to: /home/user/.config/kwtsms-cli/kwtsms-cli.toml
Log file:        kwtsms-cli.log
```

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

Send an SMS to one or more recipients.

```bash
# Single recipient
kwtsms-cli send --to 96598765432 --message "Your verification code is 4821"

# Multiple recipients (comma-separated)
kwtsms-cli send --to 96598765432,96512345678 --message "System maintenance tonight at 10pm"

# Multi-line message — use $'...' in bash so \n becomes a real newline
kwtsms-cli send --to 96598765432 --message $'Order confirmed\nTracking: TRK-12345\nExpected: Tomorrow'

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

**Bulk sending:** any number of recipients is accepted. Batches larger than 200 are split automatically into groups of 200 with a 500ms delay between batches. Numbers and Charged are aggregated in the output, with one MsgID per batch.

```
Sent
Numbers:    250
Charged:    250
Balance:    1,734
MsgID:      bb6ceabbf187d0479a24eb0ea79edace
            7793a77bc56ed1ff1bc0979c332cb98d
```

Estimated maximum recipients per platform (command-line argument size limit):

| Platform | Estimated max |
|---|---|
| Linux / macOS | ~150,000 |
| Windows | ~2,500 |

**Flags:**

| Flag | Required | Description |
|---|---|---|
| `-t, --to` | Yes | Recipient phone number(s), comma-separated. Duplicates removed automatically. |
| `-m, --message` | Yes | Message text |
| `-s, --sender` | No | Sender ID (overrides config default) |
| `--test` | No | Queue without delivery, no credits used |

## Global Flags

Available on every command.

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
sender   = "MY-BRAND"
log_file = "kwtsms-cli.log"
```

### Environment variables

Useful for CI/CD pipelines and containerised environments.

```bash
export KWTSMS_USERNAME=myapiuser
export KWTSMS_PASSWORD=myapipass
export KWTSMS_SENDER=MY-BRAND
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

## Logging

When a log file is configured, every `send` call appends one JSON line to the file. The log records the timestamp, number of recipients, credits charged, balance after, and all message IDs. Credentials, phone numbers, and message text are never written to the log.

```json
{"time":"2026-03-13T09:45:00Z","numbers":250,"charged":250,"balance":1,734,"msg_ids":["bb6ceabb...","7793a77b..."]}
{"time":"2026-03-13T09:50:00Z","error":"[ERR011] Insufficient balance."}
```

The `log_file` value in the config is a filename relative to the directory where you run the binary. To disable logging, remove the `log_file` line from the config or set it to an empty string. You can also reconfigure it at any time by running `kwtsms-cli setup` again.

## AI Agent and Automation Usage

### Agent Skill

A ready-to-use skill file for Claude Code and other AI coding agents is included at [`skill/kwtsms-cli/SKILL.md`](skill/kwtsms-cli/SKILL.md). It covers:

- Auto-detecting the correct binary for the running OS and architecture
- Downloading and installing the binary
- Configuring credentials via environment variables
- All commands with expected output examples
- Rules for safe agent use (always `--test` during development, credential handling, phone format)
- Full error code reference

To use it in Claude Code, add the skill to your project and the agent will know how to install and operate kwtsms-cli without additional instructions.

### Scripting

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
