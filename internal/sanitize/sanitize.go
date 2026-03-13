// Package sanitize provides input cleaning for all user-supplied values.
// It is called before any API call or config use to enforce strict input rules.
// All functions accept a raw string and return a cleaned string or an error.
// Related files: cmd/send.go, cmd/validate.go, internal/config/config.go
package sanitize

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// arabicToLatin maps Arabic-Indic and Extended Arabic-Indic digits to Latin digits.
var arabicToLatin = strings.NewReplacer(
	"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4",
	"٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
	"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
	"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
)

var (
	reNonDigit    = regexp.MustCompile(`\D`)
	reHTMLTag     = regexp.MustCompile(`<[^>]*>`)
	reHexEntity   = regexp.MustCompile(`&#x[0-9a-fA-F]+;|&#[0-9]+;|%[0-9a-fA-F]{2}`)
	reSQLMeta    = regexp.MustCompile("['`\";]|--|/\\*|\\*/")
	reSenderIDOK = regexp.MustCompile(`^[A-Za-z0-9\-\. ]+$`)
)

// SanitizePhone cleans a single phone number to kwtSMS-accepted format.
// Allowed: digits only (0-9), international format without leading zeros.
// Steps: convert Arabic-Indic digits to Latin, strip all non-digit chars,
// strip leading zeros. Returns an error if the result is empty.
func SanitizePhone(input string) (string, error) {
	// Convert Arabic-Indic digits to Latin
	s := arabicToLatin.Replace(input)
	// Strip all non-digit characters (+, spaces, dashes, dots, parentheses, etc.)
	s = reNonDigit.ReplaceAllString(s, "")
	// Strip leading zeros (handles 00-prefixed country codes)
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "", fmt.Errorf("invalid phone number: %q contains no valid digits", input)
	}
	return s, nil
}

// SanitizePhones splits a raw input on commas, spaces, and newlines,
// applies SanitizePhone to each token, and returns the cleaned list.
// Tokens that produce empty strings after cleaning are skipped.
// Returns an error if any token is invalid or the resulting list is empty.
func SanitizePhones(input string) ([]string, error) {
	// Split on commas, spaces, newlines, and tabs
	tokens := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t' || r == '\r'
	})
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no phone numbers provided")
	}
	var results []string
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		cleaned, err := SanitizePhone(token)
		if err != nil {
			return nil, err
		}
		results = append(results, cleaned)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no valid phone numbers after cleaning")
	}
	return results, nil
}

// SanitizeMessage cleans message body text before sending.
// Strips: HTML tags, hex-encoded characters, null bytes, most control
// characters (except \n and \t), zero-width spaces, BOM, soft hyphens,
// and SQL meta-characters. Returns an error if the result is empty.
func SanitizeMessage(input string) (string, error) {
	s := input

	// Strip HTML tags
	s = reHTMLTag.ReplaceAllString(s, "")

	// Strip hex-encoded and numeric HTML entities and URL-encoded chars
	s = reHexEntity.ReplaceAllString(s, "")

	// Strip null bytes and most control characters, preserve \n and \t
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\t' {
			b.WriteRune(r)
			continue
		}
		// Skip control characters (category Cc) and non-printable runes
		if unicode.IsControl(r) {
			continue
		}
		// Strip zero-width space (U+200B), BOM (U+FEFF), soft hyphen (U+00AD),
		// word joiner (U+2060), and other invisible formatting characters
		switch r {
		case '\u200B', '\uFEFF', '\u00AD', '\u2060', '\u200C', '\u200D',
			'\u200E', '\u200F', '\u202A', '\u202B', '\u202C', '\u202D', '\u202E':
			continue
		}
		b.WriteRune(r)
	}
	s = b.String()

	// Strip SQL meta-characters (defense in depth: values go into JSON, not SQL,
	// but we clean anyway to prevent injection into any downstream system)
	s = reSQLMeta.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("message is empty after sanitization")
	}
	return s, nil
}

// SanitizeSenderID validates and cleans a sender ID.
// Allowed characters: alphanumeric (A-Z, a-z, 0-9), hyphens (-), dots (.), spaces.
// Maximum length: 11 characters (Kuwait telecom limit).
// Returns an error if disallowed characters are present, length is exceeded, or result is empty.
func SanitizeSenderID(input string) (string, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "", fmt.Errorf("sender ID is empty")
	}
	if len([]rune(s)) > 11 {
		return "", fmt.Errorf("sender ID %q exceeds maximum length of 11 characters", s)
	}
	if !reSenderIDOK.MatchString(s) {
		return "", fmt.Errorf("sender ID %q contains invalid characters (allowed: A-Z, a-z, 0-9, hyphen, dot, space)", s)
	}
	return s, nil
}

// SanitizeConfigValue trims whitespace, strips surrounding quotes, and rejects
// empty values. Used for all config file and environment variable string values.
func SanitizeConfigValue(input string) (string, error) {
	s := strings.TrimSpace(input)
	// Strip surrounding double quotes
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	} else if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		// Strip surrounding single quotes
		s = s[1 : len(s)-1]
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("config value is empty")
	}
	return s, nil
}
