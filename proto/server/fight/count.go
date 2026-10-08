package fight

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type Count struct {
	proto.ServerSide

	NFights int
}

func (c *Count) Opcode() proto.Opcode { return "fC" }

func (c *Count) Serialize() (string, error) {
	return fmt.Sprintf("%d", c.NFights), nil
}

func (c *Count) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	n, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	c.NFights = n
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Count{} })
}
