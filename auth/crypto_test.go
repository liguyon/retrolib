package auth

import (
	"testing"
)

func TestNewKey(t *testing.T) {
	for i := 0; i < 100; i++ {
		k, err := NewKey()
		if err != nil {
			t.Error(err.Error())
		}
		if len(k) != 32 {
			t.Errorf("invalid key length: got %d", len(k))
		}
	}
}

func TestEncryptPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		key      string
		expected string
	}{
		{
			name:     "case 0",
			password: "a",
			key:      "lochhixmjyemvksmkojqvahgrhkhnpfg",
			expected: "YT",
		},
		{
			name:     "case 1",
			password: "test123",
			key:      "kujczhnrntlhudhntfakqoiyzzalfgxd",
			expected: "YV76XTQN97RQXX",
		},
		{
			name:     "case 2",
			password: "--------------------",
			key:      "fxdmrtjdzkvwozmwvbemezsrtrlrsbcr",
			expected: "OZ6fMXV60_2bS3MX8hT44d5eX88hV65e4dKVNYV6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncryptPassword(tt.password, tt.key)
			if result != tt.expected {
				t.Errorf("EncryptPassword(%q, %q) = %q, expected %q", tt.password, tt.key, result, tt.expected)
			}
		})
	}
}

func TestDecryptPassword(t *testing.T) {
	tests := []struct {
		name       string
		ciphertext string
		key        string
		expected   string
	}{
		{
			name:       "case 0",
			ciphertext: "YT",
			key:        "lochhixmjyemvksmkojqvahgrhkhnpfg",
			expected:   "a",
		},
		{
			name:       "case 1",
			ciphertext: "YV76XTQN97RQXX",
			key:        "kujczhnrntlhudhntfakqoiyzzalfgxd",
			expected:   "test123",
		},
		{
			name:       "case 2",
			ciphertext: "OZ6fMXV60_2bS3MX8hT44d5eX88hV65e4dKVNYV6",
			key:        "fxdmrtjdzkvwozmwvbemezsrtrlrsbcr",
			expected:   "--------------------",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecryptPassword(tt.ciphertext, tt.key)
			if result != tt.expected {
				t.Errorf("DecryptPassword(%q, %q) = %q, expected %q", tt.ciphertext, tt.key, result, tt.expected)
			}
		})
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		password string
		key      string
	}{
		{
			name:     "case 0",
			password: "iamhungryrightnow",
			key:      "lochhixmjyemvksmkojqvahgrhkhnpfg",
		},
		{
			name:     "case 1",
			password: "butitsgettingtoolate",
			key:      "kujczhnrntlhudhntfakqoiyzzalfgxd",
		},
		{
			name:     "case 2",
			password: "okilljustwaittomorrow",
			key:      "fxdmrtjdzkvwozmwvbemezsrtrlrsbcr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted := EncryptPassword(tt.password, tt.key)
			decrypted := DecryptPassword(encrypted, tt.key)
			if decrypted != tt.password {
				t.Errorf("Round trip failed for password %q with key %q: got %q, expected %q",
					tt.password, tt.key, decrypted, tt.password)
			}
		})
	}
}
