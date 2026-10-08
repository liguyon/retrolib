package info

import (
	"github.com/liguyon/retrolib/proto"
)

type GetMaps struct {
}

func (g *GetMaps) Opcode() proto.Opcode { return "IM" }

func (g *GetMaps) Serialize() (string, error) {
	return "", nil
}

func (g *GetMaps) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterClientType("IM",
		func() proto.Deserializer { return &GetMaps{} })
}
