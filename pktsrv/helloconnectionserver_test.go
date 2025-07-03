package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHelloConnectionServer_TypeID(t *testing.T) {
	h := HelloConnectionServer{}
	expected := "HC"
	result := h.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestHelloConnectionServer_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		hasError bool
		expected string
	}{
		{
			name:     "valid 32-character key",
			key:      "abcdefghijklmnopqrstuvwxyz123456",
			expected: "HCabcdefghijklmnopqrstuvwxyz123456",
		},
		{
			name:     "empty key",
			key:      "",
			hasError: true,
		},
		{
			name:     "short key",
			key:      "shortkey",
			hasError: true,
		},
		{
			name:     "long key",
			key:      "abcdefghijklmnopqrstuvwxyz1234567890",
			hasError: true,
		},
		{
			name:     "31-character key",
			key:      "abcdefghijklmnopqrstuvwxyz12345",
			hasError: true,
		},
		{
			name:     "33-character key",
			key:      "abcdefghijklmnopqrstuvwxyz1234567",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HelloConnectionServer{Key: tt.key}
			result, err := h.Marshal()
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

func TestHelloConnectionServer_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected HelloConnectionServer
		hasError bool
	}{
		{
			name:     "valid 32-character key",
			input:    "HCabcdefghijklmnopqrstuvwxyz123456",
			expected: HelloConnectionServer{Key: "abcdefghijklmnopqrstuvwxyz123456"},
			hasError: false,
		},
		{
			name:     "empty key",
			input:    "HC",
			expected: HelloConnectionServer{},
			hasError: true,
		},
		{
			name:     "short key",
			input:    "HCshortkey",
			expected: HelloConnectionServer{},
			hasError: true,
		},
		{
			name:     "long key",
			input:    "HCabcdefghijklmnopqrstuvwxyz1234567890",
			expected: HelloConnectionServer{},
			hasError: true,
		},
		{
			name:     "31-character key",
			input:    "HCabcdefghijklmnopqrstuvwxyz12345",
			expected: HelloConnectionServer{},
			hasError: true,
		},
		{
			name:     "33-character key",
			input:    "HCabcdefghijklmnopqrstuvwxyz1234567",
			expected: HelloConnectionServer{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h HelloConnectionServer
			err := h.Unmarshal([]byte(tt.input))

			if tt.hasError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.hasError {
				if diff := cmp.Diff(tt.expected, h); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestHelloConnectionServer_MarshalUnmarshal_Roundtrip(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"valid 32-character key", "abcdefghijklmnopqrstuvwxyz123456"},
		{"numeric key", "01234567890123456789012345678901"},
		{"mixed key", "abc123def456ghi789jkl012mno345pq"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := HelloConnectionServer{Key: tt.key}

			data, err := original.Marshal()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result HelloConnectionServer
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
