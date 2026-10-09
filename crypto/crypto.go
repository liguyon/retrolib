package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	// ErrInvalidKey is returned when key ID is outside 1..15 or material is empty.
	ErrInvalidKey = errors.New("malformed key")

	// ErrKeyMismatch is returned when the frame's key ID differs from Key.ID.
	ErrKeyMismatch = errors.New("key id mismatch")

	// ErrNotEncrypted is returned when the input data is not valid encrypted data.
	// It's an umbrella for data that cannot be decrypted correctly.
	ErrNotEncrypted = errors.New("not a valid encrypted frame")
)

const hexDigits = "0123456789ABCDEF"

// Key is used in the protocol obfuscation cipher and for mapdata encryption.
// The client keeps a collection of up to 15 keys indexed by key ID.
type Key struct {
	ID       byte
	Material []byte
}

// Validate checks whether a key is structurally valid. Returns ErrInvalidKey on error.
// ID must be in 1..15 (0 means no key), Material must not be empty.
func (k Key) Validate() error {
	if k.ID < 1 || k.ID > 15 || len(k.Material) == 0 {
		return ErrInvalidKey
	}
	return nil
}

// Encode encodes a key into its wire format: key ID encoded as a hex digit in the first byte,
// followed by the hex encoded key material.
func (k Key) Encode() (string, error) {
	if err := k.Validate(); err != nil {
		return "", err
	}
	out := make([]byte, 1+hex.EncodedLen(len(k.Material)))
	out[0] = hexDigits[k.ID]
	_ = hex.Encode(out[1:], k.Material)
	return string(out), nil
}

// Decode parses a key encoded with the encoding method implemented by Key.Encode.
func DecodeKey(raw string) (Key, error) {
	if len(raw)%2 == 0 {
		return Key{}, ErrInvalidKey
	}
	b := []byte(raw)

	id, err := hexNibble(b[0])
	if err != nil {
		return Key{}, fmt.Errorf("%w: invalid id", ErrInvalidKey)
	}

	material := make([]byte, hex.DecodedLen(len(b[1:])))
	_, err = hex.Decode(material, b[1:])
	if err != nil {
		return Key{}, fmt.Errorf("%w: material is not a valid hex", ErrInvalidKey)
	}
	if len(material) == 0 {
		return Key{}, fmt.Errorf("%w: empty material", ErrInvalidKey)
	}
	return Key{ID: id, Material: material}, nil
}

// PeekKeyID reports whether the frame's first character is a valid hex digit.
// Returns its value as a hex nibble if true.
func PeekKeyID(frame []byte) (byte, bool) {
	id, err := hexNibble(frame[0])
	if err != nil {
		return 0, false
	}
	return id, true
}

// Encrypt implements Retro's cipher used to encrypt client-server traffic and MapData.
func Encrypt(plaintext []byte, key Key) ([]byte, error) {
	if err := key.Validate(); err != nil {
		return nil, err
	}

	sum := checksum(plaintext)
	cipherBytes := xorStream(escape(plaintext), key.Material, int(sum)*2)

	out := make([]byte, 2+hex.EncodedLen(len(cipherBytes)))
	out[0] = hexDigits[key.ID]
	out[1] = hexDigits[sum]
	hex.Encode(out[2:], cipherBytes)
	return out, nil
}

// Decrypt decrypts encrypted frames and MadData.
func Decrypt(data []byte, key Key) ([]byte, error) {
	if err := key.Validate(); err != nil {
		return nil, err
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("%w: too short", ErrNotEncrypted)
	}

	id, ok := PeekKeyID(data)
	if !ok || id == 0 {
		// the encrypted data has an invalid key id prefix, not the key itself
		return nil, fmt.Errorf("%w: invalid key id", ErrNotEncrypted)
	}
	if id != key.ID {
		return nil, ErrKeyMismatch
	}

	sum, err := hexNibble(data[1])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid checksum", ErrNotEncrypted)
	}

	cipherBytes := make([]byte, hex.DecodedLen(len(data[2:])))
	_, err = hex.Decode(cipherBytes, data[2:])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid ciphertext: %w", ErrNotEncrypted, err)
	}

	plainBytes := xorStream(cipherBytes, key.Material, int(sum)*2)
	unescaped, err := unescape(plainBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid ciphertext: %w", ErrNotEncrypted, err)
	}

	if checksum(unescaped) != sum {
		return nil, fmt.Errorf("%w: checksum mismatch", ErrNotEncrypted)
	}
	return unescaped, nil
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

func hexNibble(c byte) (byte, error) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', nil
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, nil
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, nil
	default:
		return 0, fmt.Errorf("not a hex digit: %q", c)
	}
}

// escape implements Retro's custom URL encoding.
func escape(s []byte) []byte {
	out := make([]byte, 0, len(s))
	for _, c := range s {
		if c < 32 || c > 127 || c == '%' || c == '+' {
			out = append(out, '%', hexDigits[c>>4], hexDigits[c&0x0F])
		} else {
			out = append(out, c)
		}
	}
	return out
}

func unescape(s []byte) ([]byte, error) {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			hi, err := hexNibble(s[i+1])
			if err != nil {
				return nil, err
			}
			lo, err := hexNibble(s[i+2])
			if err != nil {
				return nil, err
			}
			out = append(out, hi<<4|lo)
			i += 2
			continue
		}
		out = append(out, s[i])
	}
	return out, nil
}
