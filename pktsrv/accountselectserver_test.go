package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/liguyon/retrolib/login"
)

func TestAccountSelectServer_TypeID(t *testing.T) {
	h := AccountSelectServer{}
	expected := "AY"
	result := h.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountSelectServer_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		success  bool
		addr     string
		ticket   int
		errID    login.SelectServerErrorID
		hasError bool
		expected string
	}{
		{
			name:     "selection success ok",
			success:  true,
			addr:     "localhost:5555",
			ticket:   12345,
			expected: "AYKlocalhost:5555;12345",
			hasError: false,
		},
		{
			name:     "selection error ok",
			success:  false,
			errID:    login.SelectServerFull,
			expected: "AYEF",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountSelectServer{Success: tt.success, Addr: tt.addr, Ticket: tt.ticket, ErrID: tt.errID}
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

func TestAccountSelectServer_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountSelectServer
		hasError bool
	}{
		{
			name:  "selection success ok",
			input: "AYKlocalhost:5555;12345",
			expected: AccountSelectServer{
				Success: true,
				Addr:    "localhost:5555",
				Ticket:  12345,
			},
			hasError: false,
		},
		{
			name:  "selection error ok",
			input: "AYEF",
			expected: AccountSelectServer{
				ErrID: login.SelectServerFull,
			},
			hasError: false,
		},
		{
			name:     "invalid data",
			input:    "AYT",
			hasError: true,
		},
		{
			name:     "invalid data",
			input:    "AYKlocalhost:5555",
			hasError: true,
		},
		{
			name:     "invalid data",
			input:    "AY",
			hasError: true,
		},
		{
			name:     "invalid data",
			input:    "AYKlocalhost:5555;",
			hasError: true,
		},
		{
			name:     "invalid data",
			input:    "AYKlocalhost:5555;a",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountSelectServer
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
