package proto

import (
	"bytes"
	"testing"
)

func TestKeyEncoding(t *testing.T) {
	tests := []struct {
		name              string
		id                byte
		material          []byte
		expectedErrEncode bool
		expectedErrDecode bool
	}{
		{"ok", 0x0f, []byte("abcdef"), false, false},
		{"invalid key id", 0xff, []byte("abc"), true, false},
		{"empty material", 0x0f, []byte{}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := Key{ID: tt.id, Material: tt.material}
			raw, err := key.Encode()
			switch {
			case err != nil && !tt.expectedErrEncode:
				t.Fatalf("expected no err; got: %v", err)
			case err == nil && tt.expectedErrEncode:
				t.Fatalf("expected err; got none")
			case err != nil:
				return
			}
			res, err := DecodeKey(raw)
			switch {
			case err != nil && !tt.expectedErrDecode:
				t.Fatalf("expected no err; got: %v", err)
			case err == nil && tt.expectedErrDecode:
				t.Fatalf("expected err; got none")
			case err != nil:
				return
			}
			if res.ID != key.ID || bytes.Compare(key.Material, res.Material) != 0 {
				t.Errorf("roundtrip failed: %v | %v", key, res)
			}
		})
	}
}

func TestMessageEncryption(t *testing.T) {
	tests := []struct {
		name               string
		key                Key
		payload            []byte
		expectedErrEncrypt bool
		expectedErrDecrypt bool
	}{
		{"ok", Key{ID: 1, Material: []byte("abcd")}, []byte("AlEf"), false, false},
		{"ok1", Key{ID: 15, Material: []byte("BEef")}, []byte("ihfadsf|0432+-"), false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, err := EncryptMessage(tt.payload, tt.key)
			switch {
			case err != nil && !tt.expectedErrEncrypt:
				t.Fatalf("expected no err; got: %v", err)
			case err == nil && tt.expectedErrEncrypt:
				t.Fatalf("expected err; got none")
			case err != nil:
				return
			}
			out, err := DecryptMessage(ct, tt.key)
			switch {
			case err != nil && !tt.expectedErrDecrypt:
				t.Fatalf("expected no err; got: %v", err)
			case err == nil && tt.expectedErrDecrypt:
				t.Fatalf("expected err; got none")
			case err != nil:
				return
			}
			if bytes.Compare(tt.payload, out) != 0 {
				t.Errorf("roundtrip failed: %v | %v", tt.payload, out)
			}
		})
	}
}
