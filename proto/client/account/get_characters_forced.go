package account

import (
	"github.com/liguyon/retrolib/proto"
)

type GetCharactersForced struct {
}

func (g *GetCharactersForced) Opcode() proto.Opcode { return "ALf" }

func (g *GetCharactersForced) Serialize() (string, error) {
	return "", nil
}

func (g *GetCharactersForced) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterClientType("ALf",
		func() proto.Deserializer { return &GetCharactersForced{} })
}