package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/liguyon/retrolib/auth"
)

func TestAccountServersList_TypeID(t *testing.T) {
	h := AccountServersList{}
	expected := "Ax"
	result := h.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountServersList_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		subTime  int64
		scs      []auth.ServerWithCharacters
		hasError bool
		expected string
	}{
		{
			name:    "ok",
			subTime: 6000000,
			scs: []auth.ServerWithCharacters{
				{
					ServerID:   0,
					NCharacter: 3,
				},
				{
					ServerID:   1,
					NCharacter: 2,
				},
			},
			expected: "AxK6000000|0,3|1,2",
			hasError: false,
		},
		{
			name:     "no server with character",
			subTime:  6000000,
			expected: "AxK6000000",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountServersList{RemainingSub: tt.subTime, Servers: tt.scs}
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

func TestAccountServersList_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountServersList
		hasError bool
	}{
		{
			name:     "no server ok",
			input:    "AxK60000000000",
			expected: AccountServersList{RemainingSub: 60_000_000_000},
			hasError: false,
		},
		{
			name:  "ok",
			input: "AxK60000000000|0,3|1,2",
			expected: AccountServersList{
				Servers: []auth.ServerWithCharacters{
					{
						ServerID:   0,
						NCharacter: 3,
					},
					{
						ServerID:   1,
						NCharacter: 2,
					},
				},
				RemainingSub: 60_000_000_000,
			},
			hasError: false,
		},
		{
			name:     "invalid data",
			input:    "AxKa60",
			hasError: true,
		},
		{
			name:     "invalid data 1",
			input:    "AxK|60",
			hasError: true,
		},
		{
			name:     "invalid data 2",
			input:    "AxK60|",
			hasError: true,
		},
		{
			name:     "invalid data 3",
			input:    "AxK60|1",
			hasError: true,
		},
		{
			name:     "invalid data 4",
			input:    "AxK60|1,",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountServersList
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
