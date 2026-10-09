package game

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type Create struct {
	proto.ClientSide

	Type int
}

func (c *Create) Opcode() proto.Opcode { return "GC" }

func (c *Create) Serialize() (string, error) {
	return fmt.Sprintf("%d", c.Type), nil
}

func (c *Create) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	t, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	c.Type = t
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Create{} })
}
