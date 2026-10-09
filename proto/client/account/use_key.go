package account

import (
	"fmt"

	"github.com/liguyon/retrolib/proto"
)

type UseKey struct {
	proto.ClientSide

	KeyID byte
}

func (u *UseKey) Opcode() proto.Opcode { return "Ak" }

func (u *UseKey) Serialize() (string, error) {
	c, err := proto.HexDigit(u.KeyID)
	if err != nil {
		return "", fmt.Errorf("%w: %q", err, u.KeyID)
	}

	return fmt.Sprintf("%c", c), nil
}

func (u *UseKey) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	if len(payload) != 1 {
		return fmt.Errorf("not a char: %q", payload)
	}

	n, err := proto.HexNibble(payload[0])
	if err != nil {
		return fmt.Errorf("%w: %c", err, payload[0])
	}
	u.KeyID = n
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &UseKey{} })
}
