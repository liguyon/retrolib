package proto

import (
	"errors"
)

// BoolToInt returns the digit representation of a bool.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ParseBool converts a bool literal to an actual bool.
func ParseBool(s string) (bool, error) {
	if s == "1" {
		return true, nil
	}
	if s == "0" {
		return false, nil
	}
	return false, errors.New("not a bool")
}

// HexNibble converts a hex digit char into the hex nibble it represents.
// '0'-'F' -> 0-15 (or \x00-\x0f)
func HexNibble(c byte) (byte, error) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', nil
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, nil
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, nil
	default:
		return 0, errors.New("not a hex digit")
	}
}

// HexDigit converts a hex nibble into its hex digit representation.
// 0-15 (or \x00-\x0f) -> '0'-'F'
func HexDigit(n byte) (byte, error) {
	switch {
	case n >= 0 && n <= 9:
		return n + '0', nil
	case n >= 10 && n <= 15:
		return n - 10 + 'A', nil
	default:
		return 0, errors.New("cannot be converted to a hex digit")
	}
}
