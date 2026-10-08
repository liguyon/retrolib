package game

import (
	"github.com/liguyon/retrolib/proto"
)

type CreateError struct {
	proto.ServerSide
}

func (c *CreateError) Opcode() proto.Opcode { return "GCE" }

func (c *CreateError) Serialize() (string, error) { return "", nil }

func (c *CreateError) Deserialize(payload string) error { return nil }

func init() {
	proto.Register(func() proto.Deserializer { return &CreateError{} })
}
