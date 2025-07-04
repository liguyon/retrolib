package auth

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

// NewKey generates a 32-character key to be used by a client for password encryption.
func NewKey() (string, error) {
	var sb strings.Builder
	for i := 0; i < 32; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(alphabet[nBig.Int64()])
	}
	return sb.String(), nil
}

// EncryptPassword encrypts a password using Retro's custom encryption algorithm.
// This is the encryption method identified by #1 in the HC packet.
// key is the 32-character key contained in the HC packet.
// This function does not perform validation. The password plaintext and key are expected to be valid.
func EncryptPassword(pw, key string) string {
	l := byte(len(alphabet))
	var sb strings.Builder
	for i, v := range pw {
		// split rune into high and low nibble
		high := byte(v / 16)
		low := byte(v % 16)
		o := key[i] % l // offset
		// add offset to high and low, clamp, map to password alphabet
		sb.WriteByte(alphabet[(high+o)%l])
		sb.WriteByte(alphabet[(low+o)%l])
	}
	return sb.String()
}

// DecryptPassword decrypts a password that has been encrypted using Retro's encryption method #1.
// key is the 32-character key contained in the HC packet.
// This function does not perform validation. The ciphetext and key are expected to be valid.
func DecryptPassword(ct, key string) string {
	l := byte(len(alphabet))
	var sb strings.Builder
	for i := 0; i < len(ct); i += 2 {
		o := key[i/2] % l

		highPos := byte(strings.IndexByte(alphabet, ct[i]))
		if highPos+l < o {
			highPos += l
		}
		lowPos := byte(strings.IndexByte(alphabet, ct[i+1]))
		if lowPos+l < o {
			lowPos += l
		}

		high := (highPos - o + l) % l
		low := (lowPos - o + l) % l
		sb.WriteByte(high*16 + low)
	}
	return sb.String()
}
