package pktcli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCredentials_TypeID(t *testing.T) {
	c := Credentials{}
	expected := ""
	result := c.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestCredentials_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		username string
		enc      int
		pwct     string
		hasError bool
		expected string
	}{
		{
			name:     "ok",
			username: "test",
			enc:      1,
			pwct:     "aaa",
			expected: "test\n#1aaa",
			hasError: false,
		},
		{
			name:     "empty username",
			username: "",
			enc:      1,
			pwct:     "aaa",
			hasError: true,
		},
		{
			name:     "empty password ciphertext",
			username: "le",
			enc:      1,
			pwct:     "",
			hasError: true,
		},
		{
			name:     "enc",
			username: "el",
			enc:      10000,
			pwct:     "le",
			hasError: false,
			expected: "el\n#10000le",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Credentials{Username: tt.username, EncID: tt.enc, PasswordCT: tt.pwct}
			result, err := c.Marshal()
			if tt.hasError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.hasError {
				if string(result) != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, string(result))
				}
			}
		})
	}
}

func TestCredentials_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Credentials
		hasError bool
	}{
		{
			name:     "ok",
			input:    "test\n#1aaa",
			expected: Credentials{Username: "test", EncID: 1, PasswordCT: "aaa"},
			hasError: false,
		},
		{
			name:     "empty username",
			input:    "\n#1aaa",
			hasError: true,
		},
		{
			name:     "empty password",
			input:    "test\n#1",
			hasError: true,
		},
		{
			name:     "invalid encID",
			input:    "test\n#abcd",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Credentials
			err := c.Unmarshal([]byte(tt.input))

			if tt.hasError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.hasError {
				if diff := cmp.Diff(tt.expected, c); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
