package account

import (
	"github.com/liguyon/retrolib/proto"
)

type GetRegionalVersion struct{}

func (g *GetRegionalVersion) Opcode() proto.Opcode { return "AV" }

func (g *GetRegionalVersion) Serialize() (string, error) { return "", nil }

func (g *GetRegionalVersion) Deserialize(payload string) error { return nil }

func init() {
	proto.RegisterClientType("AV",
		func() proto.Deserializer { return &GetRegionalVersion{} })
}
