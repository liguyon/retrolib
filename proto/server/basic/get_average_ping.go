package basic

import (
	"github.com/liguyon/retrolib/proto"
)

type GetAveragePing struct {
	proto.ServerSide
}

func (g *GetAveragePing) Opcode() proto.Opcode { return "Bp" }

func (g *GetAveragePing) Serialize() (string, error) {
	return "", nil
}

func (g *GetAveragePing) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &GetAveragePing{} })
}
