package account

import (
	"github.com/liguyon/retrolib/proto"
)

type GetServers struct {
}

func (g *GetServers) Opcode() proto.Opcode { return "Ax" }

func (g *GetServers) Serialize() (string, error) {
	return "", nil
}

func (g *GetServers) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterClientType("Ax",
		func() proto.Deserializer { return &GetServers{} })
}