package pktsrv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccountNewQueue_TypeID(t *testing.T) {
	a := AccountNewQueue{}
	expected := "Af"
	result := a.TypeID()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAccountNewQueue_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		packet   AccountNewQueue
		expected string
	}{
		{"subscriber true", AccountNewQueue{Position: 1, NSubs: 5, NNonSubs: 10, IsSub: true, QueueID: 100}, "Af1|5|10|1|100"},
		{"subscriber false", AccountNewQueue{Position: 2, NSubs: 3, NNonSubs: 7, IsSub: false, QueueID: 200}, "Af2|3|7|0|200"},
		{"zero values", AccountNewQueue{Position: 0, NSubs: 0, NNonSubs: 0, IsSub: false, QueueID: 0}, "Af0|0|0|0|0"},
		{"negative values", AccountNewQueue{Position: -1, NSubs: -2, NNonSubs: -3, IsSub: true, QueueID: -4}, "Af-1|-2|-3|1|-4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.packet.Marshal()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestAccountNewQueue_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AccountNewQueue
		hasError bool
	}{
		{"subscriber true", "Af1|5|10|1|100", AccountNewQueue{Position: 1, NSubs: 5, NNonSubs: 10, IsSub: true, QueueID: 100}, false},
		{"subscriber false", "Af2|3|7|0|200", AccountNewQueue{Position: 2, NSubs: 3, NNonSubs: 7, IsSub: false, QueueID: 200}, false},
		{"zero values", "Af0|0|0|0|0", AccountNewQueue{Position: 0, NSubs: 0, NNonSubs: 0, IsSub: false, QueueID: 0}, false},
		{"negative values", "Af-1|-2|-3|1|-4", AccountNewQueue{Position: -1, NSubs: -2, NNonSubs: -3, IsSub: true, QueueID: -4}, false},
		{"non-zero non-one subscriber", "Af1|2|3|5|6", AccountNewQueue{Position: 1, NSubs: 2, NNonSubs: 3, IsSub: false, QueueID: 6}, false},
		{"too few fields", "Af1|2|3|4", AccountNewQueue{}, true},
		{"too many fields", "Af1|2|3|4|5|6", AccountNewQueue{}, true},
		{"empty payload", "Af", AccountNewQueue{}, true},
		{"non-numeric field", "Afabc|2|3|4|5", AccountNewQueue{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AccountNewQueue
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

func TestAccountNewQueue_MarshalUnmarshal_Roundtrip(t *testing.T) {
	tests := []struct {
		name   string
		packet AccountNewQueue
	}{
		{"subscriber true", AccountNewQueue{Position: 1, NSubs: 5, NNonSubs: 10, IsSub: true, QueueID: 100}},
		{"subscriber false", AccountNewQueue{Position: 2, NSubs: 3, NNonSubs: 7, IsSub: false, QueueID: 200}},
		{"zero values", AccountNewQueue{Position: 0, NSubs: 0, NNonSubs: 0, IsSub: false, QueueID: 0}},
		{"negative values", AccountNewQueue{Position: -1, NSubs: -2, NNonSubs: -3, IsSub: true, QueueID: -4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.packet

			data, err := original.Marshal()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result AccountNewQueue
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
