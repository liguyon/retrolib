package crypto

import (
	"errors"
	"testing"
)

func TestKeyValidate(t *testing.T) {
	tests := []struct {
		name    string
		k       Key
		wantErr error
	}{
		{"valid", Key{ID: 1, Material: []byte("deadbeef")}, nil},
		{"valid", Key{ID: 15, Material: []byte("fsljfkdlsaf")}, nil},
		{"id = 0", Key{ID: 0, Material: []byte("deadbeef")}, ErrInvalidKey},
		{"id > 15", Key{ID: 16, Material: []byte("deadbeef")}, ErrInvalidKey},
		{"empty material", Key{ID: 1, Material: []byte{}}, ErrInvalidKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.k.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("err=%v; want %v", err, tt.wantErr)
			}
		})
	}
}
