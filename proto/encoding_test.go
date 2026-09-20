package proto

import (
	"testing"
)

func TestHexNibble(t *testing.T) {
	tests := []struct {
		name        string
		in          byte
		expected    byte
		expectedErr bool
	}{
		{"0", '0', 0, false},
		{"9", '9', 9, false},
		{"10 lower", 'a', 10, false},
		{"10 upper", 'A', 10, false},
		{"15 lower", 'f', 15, false},
		{"10 upper", 'F', 15, false},
		{"err", 'g', 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := HexNibble(tt.in)

			if err != nil && !tt.expectedErr {
				t.Fatalf("want no err; got %v", err)
			}
			if err == nil && tt.expectedErr {
				t.Fatal("want err; got not")
			}

			if tt.expected != n {
				t.Errorf("want %d; got %d", tt.expected, n)
			}
		})
	}
}

func TestHexDigit(t *testing.T) {
	tests := []struct {
		name        string
		in          byte
		expected    byte
		expectedErr bool
	}{
		{"0", 0, '0', false},
		{"9", 9, '9', false},
		{"10", 0x0a, 'A', false},
		{"15", 0x0f, 'F', false},
		{"err", 16, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := HexDigit(tt.in)

			if err != nil && !tt.expectedErr {
				t.Fatalf("want no err; got %v", err)
			}
			if err == nil && tt.expectedErr {
				t.Fatal("want err; got not")
			}

			if tt.expected != n {
				t.Errorf("want %d; got %d", tt.expected, n)
			}
		})
	}
}
