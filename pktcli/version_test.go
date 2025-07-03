package pktcli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVersion_TypeID(t *testing.T) {
	v := Version{}
	expected := ""
	result := v.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestVersion_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"semver", "1.2.3", "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{V: tt.version}
			result, err := v.Marshal()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestVersion_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Version
	}{
		{"semver", "1.2.3", Version{V: "1.2.3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Version
			err := v.Unmarshal([]byte(tt.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, v); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestVersion_MarshalUnmarshal_Roundtrip(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"semver", "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := Version{V: tt.version}

			data, err := original.Marshal()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result Version
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
