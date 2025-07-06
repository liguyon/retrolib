package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/liguyon/retrolib/auth"
)

func TestAccountHosts_TypeID(t *testing.T) {
	h := AccountHosts{}
	expected := "AH"
	result := h.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountHosts_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		srvs     []auth.Server
		hasError bool
		expected string
	}{
		{
			name: "single server ok",
			srvs: []auth.Server{
				{
					ID:         0,
					State:      auth.ServerOnline,
					Completion: 100,
					CanLogIn:   true,
				},
			},
			expected: "AH0;1;100;1",
			hasError: false,
		},
		{
			name: "multiple servers ok",
			srvs: []auth.Server{
				{
					ID:         0,
					State:      auth.ServerOnline,
					Completion: 100,
					CanLogIn:   true,
				},
				{
					ID:         1,
					State:      auth.ServerStarting,
					Completion: 0,
					CanLogIn:   false,
				},
				{
					ID:         2,
					State:      auth.ServerOffline,
					Completion: 0,
					CanLogIn:   false,
				},
			},
			expected: "AH0;1;100;1|1;2;0;0|2;0;0;0",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountHosts{Servers: tt.srvs}
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

func TestAccountHosts_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountHosts
		hasError bool
	}{
		{
			name:  "single server ok",
			input: "AH0;1;100;1",
			expected: AccountHosts{Servers: []auth.Server{
				{
					ID:         0,
					State:      auth.ServerOnline,
					Completion: 100,
					CanLogIn:   true,
				},
			}},
			hasError: false,
		},
		{
			name:  "single server ok 1",
			input: "AH0;1;100;0",
			expected: AccountHosts{Servers: []auth.Server{
				{
					ID:         0,
					State:      auth.ServerOnline,
					Completion: 100,
					CanLogIn:   false,
				},
			}},
			hasError: false,
		},
		{
			name:  "multiple servers ok",
			input: "AH0;1;100;0|1;2;15;0|2;0;0;1",
			expected: AccountHosts{Servers: []auth.Server{
				{
					ID:         0,
					State:      auth.ServerOnline,
					Completion: 100,
					CanLogIn:   false,
				},
				{
					ID:         1,
					State:      auth.ServerStarting,
					Completion: 15,
					CanLogIn:   false,
				},
				{
					ID:         2,
					State:      auth.ServerOffline,
					Completion: 0,
					CanLogIn:   true,
				},
			}},
			hasError: false,
		},
		{
			name:     "invalid int",
			input:    "AHa;1;100;0",
			hasError: true,
		},
		{
			name:     "invalid int 1",
			input:    "AH1;a;100;0",
			hasError: true,
		},
		{
			name:     "invalid int 2",
			input:    "AH1;1;a;0",
			hasError: true,
		},
		{
			name:     "invalid int 2",
			input:    "AH1;1;1;a",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountHosts
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
