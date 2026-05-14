package middleware

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "user@example.com", true},
		{"valid with subdomain", "user@mail.example.com", true},
		{"valid with plus", "user+tag@example.com", true},
		{"invalid no @", "userexample.com", false},
		{"invalid no domain", "user@", false},
		{"invalid no TLD", "user@example", false},
		{"invalid spaces", "user @example.com", false},
		{"empty string", "", false},
		{"too long", string(make([]byte, 256)) + "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantOk  bool
		wantMsg string
	}{
		{
			name:     "valid password",
			password: "SecurePass123!",
			wantOk:   true,
			wantMsg:  "",
		},
		{
			name:     "too short",
			password: "Short1!",
			wantOk:   false,
			wantMsg:  "Password must be at least 12 characters long",
		},
		{
			name:     "no uppercase",
			password: "securepass123!",
			wantOk:   false,
			wantMsg:  "Password must contain at least one uppercase letter",
		},
		{
			name:     "no lowercase",
			password: "SECUREPASS123!",
			wantOk:   false,
			wantMsg:  "Password must contain at least one lowercase letter",
		},
		{
			name:     "no number",
			password: "SecurePassword!",
			wantOk:   false,
			wantMsg:  "Password must contain at least one number",
		},
		{
			name:     "no special char",
			password: "SecurePass123",
			wantOk:   false,
			wantMsg:  "Password must contain at least one special character",
		},
		{
			name:     "too long",
			password: string(make([]byte, 130)),
			wantOk:   false,
			wantMsg:  "Password must not exceed 128 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotMsg := ValidatePassword(tt.password)
			if gotOk != tt.wantOk {
				t.Errorf("ValidatePassword() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("ValidatePassword() msg = %q, want %q", gotMsg, tt.wantMsg)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no special chars", "Hello World", "Hello World"},
		{"with HTML tags", "<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"with quotes", `He said "hello"`, "He said &quot;hello&quot;"},
		{"with single quotes", "It's working", "It&#39;s working"},
		{"with spaces", "  trimmed  ", "trimmed"},
		{"mixed", `<div class="test">Content</div>`, "&lt;div class=&quot;test&quot;&gt;Content&lt;/div&gt;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeString(tt.input); got != tt.want {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
