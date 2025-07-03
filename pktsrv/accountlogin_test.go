package pktsrv

import (
	"testing"

	"github.com/liguyon/retrolib/auth"

	"github.com/google/go-cmp/cmp"
)

func TestAccountLogin_TypeID(t *testing.T) {
	h := AccountLogin{}
	expected := "Al"
	result := h.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountLogin_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		success  bool
		isGM     bool
		errID    auth.LoginErrorID
		extra    string
		hasError bool
		expected string
	}{
		{
			name:     "login success not gm",
			success:  true,
			isGM:     false,
			expected: "AlK0",
			hasError: false,
		},
		{
			name:     "login success gm",
			success:  true,
			isGM:     true,
			expected: "AlK1",
			hasError: false,
		},
		{
			name:     "login error any",
			success:  false,
			errID:    auth.ErrAccessDenied,
			expected: "AlEf",
			hasError: false,
		},
		{
			name:     "login error extra",
			success:  false,
			errID:    auth.ErrKicked,
			extra:    "whoknowswhy",
			expected: "AlEkwhoknowswhy",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountLogin{Success: tt.success, IsGM: tt.isGM, ErrID: tt.errID, Extra: tt.extra}
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

func TestAccountLogin_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountLogin
		hasError bool
	}{
		{
			name:     "login success not gm",
			input:    "AlK0",
			expected: AccountLogin{Success: true},
			hasError: false,
		},
		{
			name:     "login success gm",
			input:    "AlK1",
			expected: AccountLogin{Success: true, IsGM: true},
			hasError: false,
		},
		{
			name:     "invalid tokens",
			input:    "AlL",
			hasError: true,
		},
		{
			name:     "invalid packet",
			input:    "Al",
			hasError: true,
		},
		{
			name:     "login error access denied",
			input:    "AlEf",
			expected: AccountLogin{ErrID: auth.ErrAccessDenied},
			hasError: false,
		},
		{
			name:     "login error kicked extra",
			input:    "AlEkBaldloose",
			expected: AccountLogin{ErrID: auth.ErrKicked, Extra: "Baldloose"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountLogin
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
