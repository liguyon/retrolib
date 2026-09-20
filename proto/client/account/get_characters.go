package account

import (
	"github.com/liguyon/retrolib/proto"
)

type GetCharacters struct{}

func (g *GetCharacters) Opcode() proto.Opcode { return "AL" }

func (g *GetCharacters) Serialize() (string, error) { return "", nil }

func (g *GetCharacters) Deserialize(payload string) error { return nil }

func init() {
	proto.RegisterClientType("AL",
		func() proto.Deserializer { return &GetCharacters{} })
}
