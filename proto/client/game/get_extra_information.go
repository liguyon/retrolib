package game

import (
	"github.com/liguyon/retrolib/proto"
)

type GetExtraInformation struct {
	proto.ClientSide
}

func (g *GetExtraInformation) Opcode() proto.Opcode { return "GI" }

func (g *GetExtraInformation) Serialize() (string, error) { return "", nil }

func (g *GetExtraInformation) Deserialize(payload string) error { return nil }

func init() {
	proto.Register(func() proto.Deserializer { return &GetExtraInformation{} })
}
