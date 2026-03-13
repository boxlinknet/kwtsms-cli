package sanitize

import (
	"testing"
)

func TestSanitizePhone(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"+96598765432", "96598765432", false},
		{"0096598765432", "96598765432", false},
		{"965 9876 5432", "96598765432", false},
		{"965-9876-5432", "96598765432", false},
		{"(965) 9876-5432", "96598765432", false},
		{"٩٦٥٩٨٧٦٥٤٣٢", "96598765432", false},
		{"۹۶۵۹۸۷۶۵۴۳۲", "96598765432", false},
		{"96598765432", "96598765432", false},
		{"abc", "", true},
		{"", "", true},
		{"   ", "", true},
	}
	for _, tt := range tests {
		got, err := SanitizePhone(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("SanitizePhone(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("SanitizePhone(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSanitizePhones(t *testing.T) {
	tests := []struct {
		input   string
		want    []string
		wantErr bool
	}{
		{"9651234,9659876", []string{"9651234", "9659876"}, false},
		// spaces within a token are stripped per-number, not treated as delimiters
		{"965 9876 5432", []string{"96598765432"}, false},
		{"+9651234,+9659876", []string{"9651234", "9659876"}, false},
		{"9651234\n9659876", []string{"9651234", "9659876"}, false},
		{"9651234, 9659876", []string{"9651234", "9659876"}, false},
		// deduplication
		{"9651234,9651234,9659876", []string{"9651234", "9659876"}, false},
		{"", nil, true},
		{"abc", nil, true},
	}
	for _, tt := range tests {
		got, err := SanitizePhones(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("SanitizePhones(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if len(got) != len(tt.want) {
			t.Errorf("SanitizePhones(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("SanitizePhones(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestSanitizeMessage(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"Hello world", "Hello world", false},
		{"<b>Hello</b>", "Hello", false},
		{"Hello<script>alert(1)</script>world", "Helloalert(1)world", false},
		{"Hello&#x3C;world", "Helloworld", false},
		{"Hello\x00world", "Helloworld", false},
		{"Hello\nworld", "Hello\nworld", false},
		{"Hello\tworld", "Hello\tworld", false},
		// SQL meta stripped
		{"Hello; DROP TABLE", "Hello DROP TABLE", false},
		{"Hello--comment", "Hellocomment", false},
		// Zero-width space stripped
		{"Hello\u200Bworld", "Helloworld", false},
		// BOM stripped
		{"\uFEFFHello", "Hello", false},
		// Empty after sanitization
		{"<script></script>", "", true},
		{"", "", true},
		{"   ", "", true},
	}
	for _, tt := range tests {
		got, err := SanitizeMessage(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("SanitizeMessage(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("SanitizeMessage(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSanitizeSenderID(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"MY-SENDER", "MY-SENDER", false},
		{"App.Name", "App.Name", false},
		{"My Brand", "My Brand", false},
		{"KWT-SMS", "KWT-SMS", false},
		{"ABC123", "ABC123", false},
		{"  MY-SENDER  ", "MY-SENDER", false},
		// Invalid characters
		{"<script>", "", true},
		{"sender!", "", true},
		{"sender@id", "", true},
		// Too long (>11 chars)
		{"TOOLONGNAME12", "", true},
		// Empty
		{"", "", true},
		{"   ", "", true},
	}
	for _, tt := range tests {
		got, err := SanitizeSenderID(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("SanitizeSenderID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("SanitizeSenderID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSanitizeConfigValue(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"myvalue", "myvalue", false},
		{"  myvalue  ", "myvalue", false},
		{`"myvalue"`, "myvalue", false},
		{"'myvalue'", "myvalue", false},
		{`"  myvalue  "`, "myvalue", false},
		// Empty
		{"", "", true},
		{"   ", "", true},
		{`""`, "", true},
		{"''", "", true},
	}
	for _, tt := range tests {
		got, err := SanitizeConfigValue(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("SanitizeConfigValue(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("SanitizeConfigValue(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
