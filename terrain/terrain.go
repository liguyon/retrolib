package terrain

import (
	"fmt"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

var charValue = func() map[rune]uint8 {
	m := make(map[rune]uint8, len(alphabet))
	for i, r := range alphabet {
		m[r] = uint8(i)
	}
	return m
}()

type Cell struct {
	Active      bool
	LineOfSight bool

	GroundLevel int // 4 bits
	GroundRot   int // 2 bits
	GroundSlope int // 4 bits
	GroundFlip  bool
	GroundNum   int // 11 bits

	Movement int // 3 bitr

	Object1Num  int // 14 bits
	Object1Rot  int // 2 bits
	Object1Flip bool

	Object2Num         int // 14 bits
	Object2Flip        bool
	Object2Interactive bool
}

func (c Cell) Encode() string {
	var raw [10]int

	raw[0] = flag(c.Active)<<5 |
		flag(c.LineOfSight) |
		(c.GroundNum&0x600)>>6 |
		(c.Object1Num&0x2000)>>11 |
		(c.Object2Num&0x2000)>>12

	raw[1] = (c.GroundRot&0x3)<<4 | c.GroundLevel&0xF

	raw[2] = (c.Movement&0x7)<<3 | (c.GroundNum>>6)&0x7

	raw[3] = c.GroundNum & 0x3F

	raw[4] = (c.GroundSlope&0xF)<<2 | flag(c.GroundFlip)<<1 | (c.Object1Num>>12)&0x1

	raw[5] = (c.Object1Num >> 6) & 0x3F
	raw[6] = c.Object1Num & 0x3F

	raw[7] = (c.Object1Rot&0x3)<<4 |
		flag(c.Object1Flip)<<3 |
		flag(c.Object2Flip)<<2 |
		flag(c.Object2Interactive)<<1 |
		(c.Object2Num>>12)&0x1

	raw[8] = (c.Object2Num >> 6) & 0x3F
	raw[9] = c.Object2Num & 0x3F

	var sb strings.Builder
	sb.Grow(len(raw))
	for _, v := range raw {
		sb.WriteByte(alphabet[v])
	}
	return sb.String()
}

// Decode unpacks a 10-character wire string into a Cell.
//
// If the cell's Active bit is unset, Decode normally returns a zero-value
// (inactive) Cell without decoding the remaining fields, matching the
// original codec's behavior of leaving unused cells sparse. Passing
// forced=true decodes all fields regardless of the Active bit.
func Decode(data string, forced bool) (Cell, error) {
	if len(data) != 10 {
		return Cell{}, fmt.Errorf("encoded data must be 10 characters, got %d", len(data))
	}

	var raw [10]uint8
	for i, r := range data {
		v, ok := charValue[r]
		if !ok {
			return Cell{}, fmt.Errorf("invalid character %q at position %d", r, i)
		}
		raw[i] = v
	}

	c := Cell{Active: raw[0]&0x20 != 0}
	if !c.Active && !forced {
		return c, nil
	}

	c.LineOfSight = raw[0]&0x01 != 0

	c.GroundRot = int(raw[1]&0x30) >> 4
	c.GroundLevel = int(raw[1] & 0x0F)

	c.Movement = int(raw[2]&0x38) >> 3
	c.GroundNum = int(raw[0]&0x18)<<6 | int(raw[2]&0x07)<<6 | int(raw[3])

	c.GroundSlope = int(raw[4]&0x3C) >> 2
	c.GroundFlip = raw[4]&0x02 != 0

	c.Object1Num = int(raw[0]&0x04)<<11 | int(raw[4]&0x01)<<12 | int(raw[5])<<6 | int(raw[6])
	c.Object1Rot = int(raw[7]&0x30) >> 4
	c.Object1Flip = raw[7]&0x08 != 0

	c.Object2Flip = raw[7]&0x04 != 0
	c.Object2Interactive = raw[7]&0x02 != 0
	c.Object2Num = int(raw[0]&0x02)<<12 | int(raw[7]&0x01)<<12 | int(raw[8])<<6 | int(raw[9])

	return c, nil
}

func flag(v bool) int {
	if v {
		return 1
	}
	return 0
}
