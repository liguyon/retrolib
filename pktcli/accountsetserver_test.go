package pktcli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccountSetServer_TypeID(t *testing.T) {
	a := AccountSetServer{}
	expected := "AX"
	result := a.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountSetServer_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		serverID int
		expected string
	}{
		{"ok", 123, "AX123"},
		{"ok 1", 0, "AX0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountSetServer{ServerID: tt.serverID}
			result, err := a.Marshal()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestAccountSetServer_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountSetServer
		hasError bool
	}{
		{"ok", "AX123", AccountSetServer{ServerID: 123}, false},
		{name: "invalid server ID", input: "AXQ12", hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountSetServer
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
