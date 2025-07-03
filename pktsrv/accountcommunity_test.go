package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccountCommunity_TypeID(t *testing.T) {
	a := AccountCommunity{}
	expected := "Ac"
	result := a.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountCommunity_Marshal(t *testing.T) {
	tests := []struct {
		name        string
		communityID int
		expected    string
	}{
		{"positive integer", 123, "Ac123"},
		{"zero", 0, "Ac0"},
		{"negative integer", -456, "Ac-456"},
		{"single digit", 7, "Ac7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountCommunity{CommunityID: tt.communityID}
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

func TestAccountCommunity_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountCommunity
		hasError bool
	}{
		{"positive integer", "Ac123", AccountCommunity{CommunityID: 123}, false},
		{"zero", "Ac0", AccountCommunity{CommunityID: 0}, false},
		{"negative integer", "Ac-456", AccountCommunity{CommunityID: -456}, false},
		{"single digit", "Ac7", AccountCommunity{CommunityID: 7}, false},
		{"empty payload", "Ac", AccountCommunity{}, true},
		{"non-numeric payload", "Acabc", AccountCommunity{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountCommunity
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

func TestAccountCommunity_MarshalUnmarshal_Roundtrip(t *testing.T) {
	tests := []struct {
		name        string
		communityID int
	}{
		{"positive integer", 123},
		{"zero", 0},
		{"negative integer", -456},
		{"single digit", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := AccountCommunity{CommunityID: tt.communityID}

			data, err := original.Marshal()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result AccountCommunity
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
