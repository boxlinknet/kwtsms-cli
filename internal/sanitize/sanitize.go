// Package sanitize provides input cleaning for all user-supplied values.
// It is called before any API call or config use to enforce strict input rules.
// All functions accept a raw string and return a cleaned string or an error.
// Related files: cmd/send.go, internal/config/config.go
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
	reNonDigit   = regexp.MustCompile(`\D`)
	reHTMLTag    = regexp.MustCompile(`<[^>]*>`)
	reHexEntity  = regexp.MustCompile(`&#x[0-9a-fA-F]+;|&#[0-9]+;|%[0-9a-fA-F]{2}`)
	reSQLMeta    = regexp.MustCompile("['`\";]|--|/\\*|\\*/")
	reSenderIDOK = regexp.MustCompile(`^[A-Za-z0-9\-\. ]+$`)
)

// phoneRule holds country-specific mobile number format rules.
type phoneRule struct {
	// localLengths lists valid digit counts after the country code.
	localLengths []int
	// mobileStartDigits lists valid first digits of the local number.
	// If empty, any starting digit is accepted.
	mobileStartDigits []string
}

// phoneRules is a table of country code -> format rules.
// Longest-match wins: 3-digit codes are checked before 2-digit, then 1-digit.
//
// Sources: ITU-T E.164 / National Numbering Plans, Wikipedia "Telephone numbers
// in [Country]", HowToCallAbroad.com, CountryCode.com.
//
// localLengths: digit count AFTER the country code.
// mobileStartDigits: valid first character(s) of the local number.
var phoneRules = map[string]phoneRule{
	// GCC
	"965": {[]int{8}, []string{"4", "5", "6", "9"}},          // Kuwait: Virgin/STC,Zain,Ooredoo,Zain
	"966": {[]int{9}, []string{"5"}},                          // Saudi Arabia
	"971": {[]int{9}, []string{"5"}},                          // UAE
	"973": {[]int{8}, []string{"3", "6"}},                     // Bahrain
	"974": {[]int{8}, []string{"3", "5", "6", "7"}},           // Qatar
	"968": {[]int{8}, []string{"7", "9"}},                     // Oman
	// Levant
	"962": {[]int{9}, []string{"7"}},                          // Jordan
	"961": {[]int{7, 8}, []string{"3", "7", "8"}},             // Lebanon
	"970": {[]int{9}, []string{"5"}},                          // Palestine
	"964": {[]int{10}, []string{"7"}},                         // Iraq
	"963": {[]int{9}, []string{"9"}},                          // Syria
	// Other Arab
	"967": {[]int{9}, []string{"7"}},                          // Yemen
	"20":  {[]int{10}, []string{"1"}},                         // Egypt
	"218": {[]int{9}, []string{"9"}},                          // Libya
	"216": {[]int{8}, []string{"2", "4", "5", "9"}},           // Tunisia
	"212": {[]int{9}, []string{"6", "7"}},                     // Morocco
	"213": {[]int{9}, []string{"5", "6", "7"}},                // Algeria
	"249": {[]int{9}, []string{"9"}},                          // Sudan
	// Non-Arab Middle East
	"98":  {[]int{10}, []string{"9"}},                         // Iran
	"90":  {[]int{10}, []string{"5"}},                         // Turkey
	"972": {[]int{9}, []string{"5"}},                          // Israel
	// South Asia
	"91":  {[]int{10}, []string{"6", "7", "8", "9"}},          // India
	"92":  {[]int{10}, []string{"3"}},                         // Pakistan
	"880": {[]int{10}, []string{"1"}},                         // Bangladesh
	"94":  {[]int{9}, []string{"7"}},                          // Sri Lanka
	"960": {[]int{7}, []string{"7", "9"}},                     // Maldives
	// East Asia
	"86":  {[]int{11}, []string{"1"}},                         // China
	"81":  {[]int{10}, []string{"7", "8", "9"}},               // Japan
	"82":  {[]int{10}, []string{"1"}},                         // South Korea
	"886": {[]int{9}, []string{"9"}},                          // Taiwan
	// Southeast Asia
	"65":  {[]int{8}, []string{"8", "9"}},                     // Singapore
	"60":  {[]int{9, 10}, []string{"1"}},                      // Malaysia
	"62":  {[]int{9, 10, 11, 12}, []string{"8"}},              // Indonesia
	"63":  {[]int{10}, []string{"9"}},                         // Philippines
	"66":  {[]int{9}, []string{"6", "8", "9"}},                // Thailand
	"84":  {[]int{9}, []string{"3", "5", "7", "8", "9"}},      // Vietnam
	"95":  {[]int{9}, []string{"9"}},                          // Myanmar
	"855": {[]int{8, 9}, []string{"1", "6", "7", "8", "9"}},   // Cambodia
	"976": {[]int{8}, []string{"6", "8", "9"}},                // Mongolia
	// Europe
	"44":  {[]int{10}, []string{"7"}},                         // UK
	"33":  {[]int{9}, []string{"6", "7"}},                     // France
	"49":  {[]int{10, 11}, []string{"1"}},                     // Germany
	"39":  {[]int{10}, []string{"3"}},                         // Italy
	"34":  {[]int{9}, []string{"6", "7"}},                     // Spain
	"31":  {[]int{9}, []string{"6"}},                          // Netherlands
	"32":  {[]int{9}, nil},                                     // Belgium
	"41":  {[]int{9}, []string{"7"}},                          // Switzerland
	"43":  {[]int{10}, []string{"6"}},                         // Austria
	"47":  {[]int{8}, []string{"4", "9"}},                     // Norway
	"48":  {[]int{9}, nil},                                     // Poland
	"30":  {[]int{10}, []string{"6"}},                         // Greece
	"420": {[]int{9}, []string{"6", "7"}},                     // Czech Republic
	"46":  {[]int{9}, []string{"7"}},                          // Sweden
	"45":  {[]int{8}, nil},                                     // Denmark
	"40":  {[]int{9}, []string{"7"}},                          // Romania
	"36":  {[]int{9}, nil},                                     // Hungary
	"380": {[]int{9}, nil},                                     // Ukraine
	// Americas
	"1":   {[]int{10}, nil},                                    // USA/Canada
	"52":  {[]int{10}, nil},                                    // Mexico
	"55":  {[]int{11}, nil},                                    // Brazil
	"57":  {[]int{10}, []string{"3"}},                         // Colombia
	"54":  {[]int{10}, []string{"9"}},                         // Argentina
	"56":  {[]int{9}, []string{"9"}},                          // Chile
	"58":  {[]int{10}, []string{"4"}},                         // Venezuela
	"51":  {[]int{9}, []string{"9"}},                          // Peru
	"593": {[]int{9}, []string{"9"}},                          // Ecuador
	"53":  {[]int{8}, []string{"5", "6"}},                     // Cuba
	// Africa
	"27":  {[]int{9}, []string{"6", "7", "8"}},                // South Africa
	"234": {[]int{10}, []string{"7", "8", "9"}},               // Nigeria
	"254": {[]int{9}, []string{"1", "7"}},                     // Kenya
	"233": {[]int{9}, []string{"2", "5"}},                     // Ghana
	"251": {[]int{9}, []string{"7", "9"}},                     // Ethiopia
	"255": {[]int{9}, []string{"6", "7"}},                     // Tanzania
	"256": {[]int{9}, []string{"7"}},                          // Uganda
	"237": {[]int{9}, []string{"6"}},                          // Cameroon
	"225": {[]int{10}, nil},                                    // Ivory Coast
	"221": {[]int{9}, []string{"7"}},                          // Senegal
	"252": {[]int{9}, []string{"6", "7"}},                     // Somalia
	"250": {[]int{9}, []string{"7"}},                          // Rwanda
	// Oceania
	"61":  {[]int{9}, []string{"4"}},                          // Australia
	"64":  {[]int{8, 9, 10}, []string{"2"}},                   // New Zealand
}

// countryNames maps country code to display name for error messages.
var countryNames = map[string]string{
	"965": "Kuwait", "966": "Saudi Arabia", "971": "UAE", "973": "Bahrain",
	"974": "Qatar", "968": "Oman", "962": "Jordan", "961": "Lebanon",
	"970": "Palestine", "964": "Iraq", "963": "Syria", "967": "Yemen",
	"20": "Egypt", "218": "Libya", "216": "Tunisia", "212": "Morocco",
	"213": "Algeria", "249": "Sudan", "98": "Iran", "90": "Turkey",
	"972": "Israel", "91": "India", "92": "Pakistan", "880": "Bangladesh",
	"94": "Sri Lanka", "960": "Maldives", "86": "China", "81": "Japan",
	"82": "South Korea", "886": "Taiwan", "65": "Singapore", "60": "Malaysia",
	"62": "Indonesia", "63": "Philippines", "66": "Thailand", "84": "Vietnam",
	"95": "Myanmar", "855": "Cambodia", "976": "Mongolia", "44": "UK",
	"33": "France", "49": "Germany", "39": "Italy", "34": "Spain",
	"31": "Netherlands", "32": "Belgium", "41": "Switzerland", "43": "Austria",
	"47": "Norway", "48": "Poland", "30": "Greece", "420": "Czech Republic",
	"46": "Sweden", "45": "Denmark", "40": "Romania", "36": "Hungary",
	"380": "Ukraine", "1": "USA/Canada", "52": "Mexico", "55": "Brazil",
	"57": "Colombia", "54": "Argentina", "56": "Chile", "58": "Venezuela",
	"51": "Peru", "593": "Ecuador", "53": "Cuba", "27": "South Africa",
	"234": "Nigeria", "254": "Kenya", "233": "Ghana", "251": "Ethiopia",
	"255": "Tanzania", "256": "Uganda", "237": "Cameroon", "225": "Ivory Coast",
	"221": "Senegal", "252": "Somalia", "250": "Rwanda", "61": "Australia",
	"64": "New Zealand",
}

// findCountryCode returns the matching country code prefix from a normalized
// phone number. Tries 3-digit codes first, then 2-digit, then 1-digit
// (longest match wins). Returns empty string if no rule is found.
func findCountryCode(normalized string) string {
	if len(normalized) >= 3 {
		if _, ok := phoneRules[normalized[:3]]; ok {
			return normalized[:3]
		}
	}
	if len(normalized) >= 2 {
		if _, ok := phoneRules[normalized[:2]]; ok {
			return normalized[:2]
		}
	}
	if len(normalized) >= 1 {
		if _, ok := phoneRules[normalized[:1]]; ok {
			return normalized[:1]
		}
	}
	return ""
}

// validatePhoneFormat checks a normalized phone number against country-specific
// rules for local number length and mobile starting digits.
// Numbers with no matching country rule pass through (generic E.164 only).
func validatePhoneFormat(normalized string) error {
	cc := findCountryCode(normalized)
	if cc == "" {
		return nil // no specific rule — generic E.164 length already checked
	}

	rule := phoneRules[cc]
	local := normalized[len(cc):]
	country := countryNames[cc]
	if country == "" {
		country = "+" + cc
	}

	// Validate local number length
	validLen := false
	for _, l := range rule.localLengths {
		if len(local) == l {
			validLen = true
			break
		}
	}
	if !validLen {
		expected := make([]string, len(rule.localLengths))
		for i, l := range rule.localLengths {
			expected[i] = fmt.Sprintf("%d", l)
		}
		return fmt.Errorf("invalid %s number: expected %s digits after +%s, got %d",
			country, strings.Join(expected, " or "), cc, len(local))
	}

	// Validate mobile starting digit (if rules exist)
	if len(rule.mobileStartDigits) > 0 && len(local) > 0 {
		validPrefix := false
		for _, prefix := range rule.mobileStartDigits {
			if strings.HasPrefix(local, prefix) {
				validPrefix = true
				break
			}
		}
		if !validPrefix {
			return fmt.Errorf("invalid %s mobile number: after +%s must start with %s",
				country, cc, strings.Join(rule.mobileStartDigits, ", "))
		}
	}

	return nil
}

// SanitizePhone cleans a single phone number to kwtSMS-accepted format and
// validates it against country-specific E.164 rules.
// Steps: convert Arabic-Indic digits, strip non-digits, strip leading zeros,
// check E.164 length (7-15), validate country format and mobile prefix.
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

	// Strip local trunk prefix (leading 0 after country code).
	// e.g. +966 055-XXX-XXXX → 9660559... → strip the 0 → 966559...
	// This happens when users include the local dialing prefix in international format.
	if cc := findCountryCode(s); cc != "" {
		local := s[len(cc):]
		if len(local) > 0 && local[0] == '0' {
			rule := phoneRules[cc]
			stripped := local[1:]
			for _, l := range rule.localLengths {
				if len(stripped) == l {
					s = cc + stripped
					break
				}
			}
		}
	}

	if len(s) < 7 {
		return "", fmt.Errorf("invalid phone number: %q is too short (minimum 7 digits)", input)
	}
	if len(s) > 15 {
		return "", fmt.Errorf("invalid phone number: %q is too long (maximum 15 digits)", input)
	}

	if err := validatePhoneFormat(s); err != nil {
		return "", err
	}

	return s, nil
}

// SanitizePhones splits a raw input on commas and newlines,
// applies SanitizePhone to each token, deduplicates, and returns the cleaned list.
// Spaces within a token are treated as part of the number and stripped by SanitizePhone.
// Returns an error if any token is invalid or the resulting list is empty.
func SanitizePhones(input string) ([]string, error) {
	// Split on commas and newlines only; spaces are stripped per-number by SanitizePhone
	tokens := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\t' || r == '\r'
	})
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no phone numbers provided")
	}
	seen := make(map[string]bool)
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
		if !seen[cleaned] {
			seen[cleaned] = true
			results = append(results, cleaned)
		}
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
