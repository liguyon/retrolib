package basic

import (
	"github.com/liguyon/retrolib/proto"
)

type GetAveragePing struct {
}

func (g *GetAveragePing) Opcode() proto.Opcode { return "Bp" }

func (g *GetAveragePing) Serialize() (string, error) {
	return "", nil
}

func (g *GetAveragePing) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterServerType("Bp",
		func() proto.Deserializer { return &GetAveragePing{} })
}