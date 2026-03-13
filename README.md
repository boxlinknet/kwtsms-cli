# kwtsms-cli

A command-line interface for the [kwtSMS](https://www.kwtsms.com) SMS gateway. Send SMS messages, check your account balance, validate phone numbers, and manage sender IDs, all from your terminal or as part of an automated workflow.

Designed for developers, DevOps engineers, and AI agents that need to integrate SMS delivery into scripts, pipelines, and automation without writing custom code.

## About kwtSMS

[kwtSMS](https://www.kwtsms.com) is a Kuwait-based SMS gateway providing reliable local and international SMS delivery. It supports Arabic messaging, transactional and promotional sender IDs, and a simple REST API. Used by businesses across Kuwait and the GCC region for OTP delivery, customer notifications, and bulk campaigns.

- Website: [www.kwtsms.com](https://www.kwtsms.com)
- Support: [www.kwtsms.com/support.html](https://www.kwtsms.com/support.html)
- API documentation: [www.kwtsms.com/doc/KwtSMS.com_API_Documentation_v41.pdf](https://www.kwtsms.com/doc/KwtSMS.com_API_Documentation_v41.pdf)
- FAQ: [www.kwtsms.com/faq/](https://www.kwtsms.com/faq/)

## Use Cases

- **SMS automation:** Send notifications, alerts, and reminders from shell scripts and cron jobs.
- **AI agent integration:** Provide AI agents and LLM-powered tools with the ability to send SMS messages as a tool call.
- **CI/CD pipelines:** Trigger SMS alerts on deployment, test failure, or system events.
- **OTP delivery:** Send one-time passwords from any script or backend service.
- **Bulk validation:** Validate large lists of phone numbers before running a campaign.
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

Your API credentials are stored in `~/.config/kwtsms-cli/config.toml` with restricted permissions. You can get your API credentials from your [kwtSMS account](https://www.kwtsms.com/account/api/).

## Commands

### `kwtsms-cli setup`

Interactive setup wizard. Prompts for your API username and password, verifies them against the API, then lets you choose a default sender ID. Writes the config file to the correct location for your operating system.

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

Send an SMS to one or more recipients.

```bash
# Single recipient
kwtsms-cli send --to 96598765432 --message "Your verification code is 4821"

# Multiple recipients (comma-separated)
kwtsms-cli send --to 96598765432,96512345678 --message "System maintenance tonight at 10pm"

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

**Flags:**

| Flag | Required | Description |
|---|---|---|
| `-t, --to` | Yes | Recipient number(s), comma-separated |
| `-m, --message` | Yes | Message text |
| `-s, --sender` | No | Sender ID (overrides config default) |
| `--test` | No | Queue without delivery, no credits used |

### `kwtsms-cli validate`

Check whether phone numbers are valid and routable before sending. Useful for cleaning contact lists before a bulk campaign.

```bash
# Space-separated
kwtsms-cli validate 96598765432 96512345678

# Comma-separated
kwtsms-cli validate 96598765432,96512345678

# From a file (one number per line)
cat numbers.txt | xargs kwtsms-cli validate
```

```
Valid:    96598765432
Invalid:  123
NoRoute:  none
```

- **Valid:** accepted and routable
- **Invalid:** format error (will be auto-corrected on send)
- **NoRoute:** number format is valid but the country is not activated on your account

## Global Flags

These flags work on every command.

| Flag | Description |
|---|---|
| `--json` | Output the raw API response as JSON, suitable for piping and scripting |
| `--config PATH` | Use a different config file instead of the default |
| `--username VALUE` | Override the API username for this call only |
| `--password VALUE` | Override the API password for this call only |

## Configuration

### Config file

`kwtsms-cli setup` creates the config file automatically. You can also create or edit it manually.

| Platform | Location |
|---|---|
| Linux / macOS | `~/.config/kwtsms-cli/config.toml` |
| Windows | `%APPDATA%\kwtsms-cli\config.toml` |

```toml
username = "myapiuser"
password = "myapipass"
sender   = "MY-SENDER"
```

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

**Validate a list before a campaign:**
```bash
kwtsms-cli validate $(cat numbers.txt | tr '\n' ',') --json
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

- General support and account questions: [www.kwtsms.com/support.html](https://www.kwtsms.com/support.html)
- FAQ: [www.kwtsms.com/faq/](https://www.kwtsms.com/faq/)
- Sender ID help: [www.kwtsms.com/sender-id-help.html](https://www.kwtsms.com/sender-id-help.html)
- API documentation: [www.kwtsms.com/doc/KwtSMS.com_API_Documentation_v41.pdf](https://www.kwtsms.com/doc/KwtSMS.com_API_Documentation_v41.pdf)
- WhatsApp support: [+965 9922 0322](https://wa.me/96599220322)

For security vulnerability reports, see [SECURITY.md](SECURITY.md).

## License

MIT
