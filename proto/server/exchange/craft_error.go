package exchange

import (
	"fmt"

	"github.com/liguyon/retrolib/proto"
)

type CraftError struct {
	proto.ServerSide

	Reason byte
}

func (c *CraftError) Opcode() proto.Opcode { return "EcE" }

func (c *CraftError) Serialize() (string, error) {
	return fmt.Sprintf("%c", c.Reason), nil
}

func (c *CraftError) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	if len(payload) != 1 {
		return fmt.Errorf("expected a char: %q", payload)
	}
	c.Reason = payload[0]
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &CraftError{} })
}
