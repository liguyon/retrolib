package basic

import (
	"github.com/liguyon/retrolib/proto"
)

type GetDate struct {
	proto.ClientSide
}

func (g *GetDate) Opcode() proto.Opcode { return "BD" }

func (g *GetDate) Serialize() (string, error) {
	return "", nil
}

func (g *GetDate) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &GetDate{} })
}
