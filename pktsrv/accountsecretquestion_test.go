package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccountSecretQuestion_TypeID(t *testing.T) {
	h := AccountSecretQuestion{}
	expected := "AQ"
	result := h.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountSecretQuestion_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		question string
		hasError bool
		expected string
	}{
		{
			name:     "ok",
			question: "delete?",
			expected: "AQdelete?",
			hasError: false,
		},
		{
			name:     "empty question",
			question: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountSecretQuestion{Question: tt.question}
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

func TestAccountSecretQuestion_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountSecretQuestion
		hasError bool
	}{
		{
			name:     "ok",
			input:    "AQdelete?",
			expected: AccountSecretQuestion{Question: "delete?"},
			hasError: false,
		},
		{
			name:     "empty question",
			input:    "AQ",
			expected: AccountSecretQuestion{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountSecretQuestion
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

func TestAccountSecretQuestion_MarshalUnmarshal_Roundtrip(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"ok", "AQdelete?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := AccountSecretQuestion{Question: tt.key}

			data, err := original.Marshal()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result AccountSecretQuestion
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
