package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type CreateSuccess struct {
	proto.ServerSide

	Type int
}

func (c *CreateSuccess) Opcode() proto.Opcode { return "GCK" }

func (c *CreateSuccess) Serialize() (string, error) {
	return fmt.Sprintf("%d", c.Type), nil
}

func (c *CreateSuccess) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload[1:], "|")
	t, err := strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[0])
	}
	c.Type = t
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &CreateSuccess{} })
}
