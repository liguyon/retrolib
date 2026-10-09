package info

import (
	"github.com/liguyon/retrolib/proto"
)

type GetMaps struct {
	proto.ClientSide
}

func (g *GetMaps) Opcode() proto.Opcode { return "IM" }

func (g *GetMaps) Serialize() (string, error) {
	return "", nil
}

func (g *GetMaps) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &GetMaps{} })
}
