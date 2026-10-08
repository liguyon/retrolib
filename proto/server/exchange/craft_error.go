package exchange

import (
	"fmt"

	"github.com/liguyon/retrolib/proto"
)

type CraftError struct {
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
		return proto.ErrMalformedPayload
	}
	c.Reason = payload[0]
	return nil
}

func init() {
	proto.RegisterServerType("EcE",
		func() proto.Deserializer { return &CraftError{} })
}
