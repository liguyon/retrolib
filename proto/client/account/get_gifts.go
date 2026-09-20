package account

import (
	"github.com/liguyon/retrolib/proto"
)

type GetGifts struct {
	LanguageCode string
}

func (g *GetGifts) Opcode() proto.Opcode { return "Ag" }

func (g *GetGifts) Serialize() (string, error) {
	if g.LanguageCode == "" {
		return "", proto.ErrMissingPayload
	}
	return g.LanguageCode, nil
}

func (g *GetGifts) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	g.LanguageCode = payload
	return nil
}

func init() {
	proto.RegisterClientType("Ag",
		func() proto.Deserializer { return &GetGifts{} })
}
