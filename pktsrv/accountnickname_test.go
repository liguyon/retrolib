package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccountNickname_TypeID(t *testing.T) {
	a := AccountNickname{}
	expected := "Ad"
	result := a.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountNickname_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		nickname string
		hasError bool
		expected string
	}{
		{
			name:     "normal nickname",
			nickname: "TestUser",
			expected: "AdTestUser",
		},
		{
			name:     "empty nickname",
			nickname: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountNickname{Nickname: tt.nickname}
			result, err := a.Marshal()
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

func TestAccountNickname_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountNickname
		hasError bool
	}{
		{
			name:     "normal nickname",
			input:    "AdTestUser",
			expected: AccountNickname{Nickname: "TestUser"},
			hasError: false,
		},
		{
			name:     "empty nickname",
			input:    "Ad",
			expected: AccountNickname{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountNickname
			err := a.Unmarshal([]byte(tt.input))

			if tt.hasError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.hasError {
				if diff := cmp.Diff(tt.expected, a); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestAccountNickname_MarshalUnmarshal_Roundtrip(t *testing.T) {
	tests := []struct {
		name     string
		nickname string
	}{
		{"normal nickname", "TestUser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := AccountNickname{Nickname: tt.nickname}

			data, err := original.Marshal()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result AccountNickname
			err = result.Unmarshal(data)
			if err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if diff := cmp.Diff(original, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
