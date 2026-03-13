# kwtsms-cli Skill

Use this skill when you need to send SMS messages, check account balance, list sender IDs, or check coverage via the kwtSMS gateway using the `kwtsms-cli` binary.

---

## Step 1: Install the Binary

Detect the platform and download the correct binary. Run this once per environment.

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux)
    case "$ARCH" in
      x86_64)          BINARY="kwtsms-cli-linux-x64" ;;
      aarch64|arm64)   BINARY="kwtsms-cli-linux-arm64" ;;
      armv7l)          BINARY="kwtsms-cli-linux-armv7" ;;
      *) echo "Unsupported arch: $ARCH" && exit 1 ;;
    esac
    ;;
  darwin)
    case "$ARCH" in
      x86_64)  BINARY="kwtsms-cli-macos-x64" ;;
      arm64)   BINARY="kwtsms-cli-macos-arm64" ;;
      *) echo "Unsupported arch: $ARCH" && exit 1 ;;
    esac
    ;;
  *)
    echo "Windows: download kwtsms-cli-windows-x64.exe from https://github.com/boxlinknet/kwtsms-cli/releases/latest"
    exit 1
    ;;
esac

curl -Lo kwtsms-cli "https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/$BINARY"
chmod +x kwtsms-cli
```

To install system-wide (Linux/macOS):
```bash
sudo mv kwtsms-cli /usr/local/bin/
```

---

## Step 2: Configure Credentials

**Option A — Environment variables (recommended for agents):**
```bash
export KWTSMS_USERNAME=your_api_username
export KWTSMS_PASSWORD=your_api_password
export KWTSMS_SENDER=YOUR-SENDER
```

**Option B — Interactive setup (for human-assisted configuration):**
```bash
kwtsms-cli setup
```

**Option C — Inline flags (one-off calls):**
```bash
kwtsms-cli balance --username myuser --password mypass
```

Credentials are resolved in this order (highest wins): `--username/--password` flags > `KWTSMS_*` env vars > config file.

---

## Step 3: Use the Tool

### Check balance
```bash
kwtsms-cli balance
kwtsms-cli balance --json
```

```
Available:  1,234
Purchased:  5,000
```

### List sender IDs
```bash
kwtsms-cli senderid
```

### List coverage (active country prefixes)
```bash
kwtsms-cli coverage
```

### Send SMS
```bash
# Single recipient
kwtsms-cli send --to 96598765432 --message "Your code is 4821" --sender MY-SENDER

# Multiple recipients (comma-separated, duplicates removed automatically)
kwtsms-cli send --to 96598765432,96512345678 --message "System alert" --sender MY-SENDER

# Bulk send — any number of recipients, auto-batched at 200 per API call
kwtsms-cli send --to 96550000001,96550000002,...,96550000500 --message "Announcement" --sender MY-SENDER

# Multi-line message
kwtsms-cli send --to 96598765432 --message $'Line one\nLine two' --sender MY-SENDER

# Test mode — queued but not delivered, no credits consumed
kwtsms-cli send --to 96598765432 --message "Test" --sender MY-SENDER --test
```

Send output:
```
Sent
Numbers:    1
Charged:    1
Balance:    1,234
MsgID:      f4c841adee210f31307633ceaebff2ec
```

### JSON output (for scripting)
All commands support `--json`:
```bash
# Capture message ID after send
result=$(kwtsms-cli send --to 96598765432 --message "Alert" --sender MY-SENDER --json)
msg_id=$(echo "$result" | jq -r '."msg-id"')

# Check balance in a script
available=$(kwtsms-cli balance --json | jq '.available')
if [ "$available" -lt 100 ]; then
  echo "Low balance: $available credits"
fi
```

---

## Rules for Agents

- **Always use `--test` during development.** Test messages are queued but never delivered and consume no credits. Set `test: 0` only when sending for real.
- **Never log or print credentials.** Do not echo `KWTSMS_PASSWORD` or pass it as a visible argument in logs.
- **Phone number format:** digits only, international format, no `+` or `00` prefix. The tool normalises common formats automatically (`+965...`, `00965...`, spaces, dashes, Arabic-Indic digits).
- **Sender ID:** must be pre-approved on the kwtSMS account. Use `kwtsms-cli senderid` to list available IDs. Max 11 characters: alphanumeric, hyphens, dots, spaces.
- **Bulk send:** batches >200 are split automatically. No action needed.
- **Bulk JSON output:** single batch returns a JSON object, multiple batches return a JSON array.
- **Delete test messages** from the kwtSMS queue after testing to recover any held credits.

---

## Error Reference

| Code | Meaning | Action |
|---|---|---|
| ERR003 | Wrong username or password | Check credentials |
| ERR006 | No valid numbers submitted | Check phone format |
| ERR008 | Sender ID banned | Use a different sender ID |
| ERR009 | Empty message | Provide message text |
| ERR010 | Zero balance | Top up account |
| ERR011 | Insufficient balance | Top up account |
| ERR013 | Send queue full | Retry after a short wait |
| ERR025 | Invalid number (non-digit chars) | Strip `+`, spaces, dashes |
| ERR026 | No route for country | Contact kwtSMS to activate |
| ERR028 | 15s minimum between sends to same number | Wait before resending |

---

## Resources

- Repository: [github.com/boxlinknet/kwtsms-cli](https://github.com/boxlinknet/kwtsms-cli)
- Releases: [github.com/boxlinknet/kwtsms-cli/releases/latest](https://github.com/boxlinknet/kwtsms-cli/releases/latest)
- kwtSMS dashboard: [www.kwtsms.com/login](https://www.kwtsms.com/login)
- kwtSMS support: [www.kwtsms.com/support.html](https://www.kwtsms.com/support.html)
