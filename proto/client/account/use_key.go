package account

import (
	"fmt"

	"github.com/liguyon/retrolib/proto"
)

type UseKey struct {
	KeyID byte
}

func (u *UseKey) Opcode() proto.Opcode { return "Ak" }

func (u *UseKey) Serialize() (string, error) {
	c, err := proto.HexDigit(u.KeyID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%c", c), nil
}

func (u *UseKey) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	if len(payload) != 1 {
		return proto.ErrMalformedPayload
	}

	n, err := proto.HexNibble(payload[0])
	if err != nil {
		return fmt.Errorf("%w: %v", proto.ErrMalformedKey, err)
	}
	u.KeyID = n
	return nil
}

func init() {
	proto.RegisterClientType("Ak", func() proto.Deserializer { return &UseKey{} })
}
