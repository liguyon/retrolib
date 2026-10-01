package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
)

var ErrMalformedKey = errors.New("malformed key")

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
	if k.ID > 0x0f {
		return "", fmt.Errorf("%w: invalid id", ErrMalformedKey)
	}
	out[0] = hexDigits[k.ID]
	_ = hex.Encode(out[1:], k.Material)
	return string(out), nil
}

// Decode parses a key encoded with the encoding method implemented by Key.Encode.
func DecodeKey(raw string) (Key, error) {
	if len(raw)%2 == 0 {
		return Key{}, ErrMalformedKey
	}
	b := []byte(raw)

	id, err := hexNibble(b[0])
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

// Encrypt implements Retro's cipher used to encrypt client-server traffic and MapData.
func Encrypt(plaintext []byte, key Key) ([]byte, error) {
	if len(key.Material) == 0 {
		return nil, fmt.Errorf("%w: empty", ErrMalformedKey)
	}

	sum := checksum(plaintext)
	cipherBytes := xorStream(escape(plaintext), key.Material, int(sum)*2)

	out := make([]byte, 2+hex.EncodedLen(len(cipherBytes)))
	out[0] = hexDigits[key.ID]
	out[1] = hexDigits[sum]
	hex.Encode(out[2:], cipherBytes)
	return out, nil
}

// Decrypt is used to decrypt data that's been encrypted with the algorithm implemented by
// "Encrypt".
func Decrypt(message []byte, key Key) ([]byte, error) {
	if len(message) < 2 {
		return nil, errors.New("message too short")
	}
	if len(key.Material) == 0 {
		return nil, fmt.Errorf("%w: empty", ErrMalformedKey)
	}

	id, err := hexNibble(message[0])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid id", ErrMalformedKey)
	}
	if id != key.ID {
		return nil, errors.New("key id mismatch")
	}

	sum, err := hexNibble(message[1])
	if err != nil {
		return nil, fmt.Errorf("invalid checksum: %v", err)
	}

	cipherBytes := make([]byte, hex.DecodedLen(len(message[2:])))
	_, err = hex.Decode(cipherBytes, message[2:])
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext: %v", err)
	}

	plainBytes := xorStream(cipherBytes, key.Material, int(sum)*2)
	unescaped, err := unescape(plainBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext: %v", err)
	}

	if checksum(unescaped) != sum {
		return nil, errors.New("checksum mismatch")
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
