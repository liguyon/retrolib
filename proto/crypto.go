package proto

import (
	"encoding/hex"
	"errors"
	"fmt"
)

const hexDigits = "0123456789ABCDEF"

// Key is used in the protocol obfuscation cipher. The client keeps a collection of keys
// indexed by key ID.
type Key struct {
	ID       byte
	Material []byte
}

// Encode encodes a key into its wire format: key ID encoded as a hex digit in the first byte,
// followed by the hex encoded key material.
func (k Key) Encode() (string, error) {
	if len(k.Material) == 0 {
		return "", fmt.Errorf("%w: empty key material", ErrMalformedKey)
	}

	out := make([]byte, 1+hex.EncodedLen(len(k.Material)))

	id, err := HexDigit(k.ID)
	if err != nil {
		return "", fmt.Errorf("%w: invalid id", ErrMalformedKey)
	}
	out[0] = id

	_ = hex.Encode(out[1:], k.Material)
	return string(out), nil
}

// Decode parses a key that's been encoded using the encoding method implemented by Key.Encode.
func DecodeKey(raw string) (Key, error) {
	if len(raw)%2 == 0 {
		return Key{}, ErrMalformedKey
	}
	b := []byte(raw)

	id, err := HexNibble(b[0])
	if err != nil {
		return Key{}, fmt.Errorf("%w: invalid id", ErrMalformedKey)
	}

	material := make([]byte, hex.DecodedLen(len(b[1:])))
	_, err = hex.Decode(material, b[1:])
	if err != nil {
		return Key{}, fmt.Errorf("%w: %v", ErrMalformedKey, err)
	}
	if len(material) == 0 {
		return Key{}, fmt.Errorf("%w: empty material", ErrMalformedKey)
	}
	return Key{ID: id, Material: material}, nil
}

// EncryptMessage implements Retro's symmetric cipher used to obfuscate client<->server traffic.
// It encodes key ID in the first byte of the output, a checksum of the payload in the second,
// and hex encodes the ciphertext after.
// The cipher is a XOR stream, it XORs each byte of the payload with a byte from key offset by
// a counter.
func EncryptMessage(plaintext []byte, key Key) ([]byte, error) {
	if len(key.Material) == 0 {
		return nil, errors.New("empty key")
	}

	sum := checksum(plaintext)
	cipherBytes := xorStream(plaintext, key.Material, int(sum)*2)

	out := make([]byte, 2+hex.EncodedLen(len(cipherBytes)))
	out[0] = hexDigits[key.ID]
	out[1] = hexDigits[sum]
	hex.Encode(out[2:], cipherBytes)
	return out, nil
}

// DecryptMessage is used to decrypt a message encrypted with the algorithm implemented by
// "EncryptMessage".
func DecryptMessage(message []byte, key Key) ([]byte, error) {
	if len(message) < 2 {
		return nil, errors.New("message too short")
	}
	if len(key.Material) == 0 {
		return nil, errors.New("empty key")
	}

	id, err := HexNibble(message[0])
	if err != nil {
		return nil, fmt.Errorf("invalid key id: %v", err)
	}
	if id != key.ID {
		return nil, errors.New("key id mismatch")
	}

	sum, err := HexNibble(message[1])
	if err != nil {
		return nil, fmt.Errorf("invalid checksum: %v", err)
	}

	cipherBytes := make([]byte, hex.DecodedLen(len(message[2:])))
	_, err = hex.Decode(cipherBytes, message[2:])
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext: %w", err)
	}

	plainBytes := xorStream(cipherBytes, key.Material, int(sum)*2)

	if checksum(plainBytes) != sum {
		return nil, errors.New("checksum mismatch")
	}
	return plainBytes, nil
}

func checksum(data []byte) byte {
	var sum int
	for _, b := range data {
		sum += int(b % 16)
	}
	return byte(sum % 16)
}

func xorStream(data, key []byte, offset int) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ key[(i+offset)%len(key)]
	}
	return out
}
