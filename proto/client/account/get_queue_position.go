package account

import (
	"github.com/liguyon/retrolib/proto"
)

type GetQueuePosition struct{}

func (g *GetQueuePosition) Opcode() proto.Opcode             { return "Af" }
func (g *GetQueuePosition) Serialize() (string, error)       { return "", nil }
func (g *GetQueuePosition) Deserialize(payload string) error { return nil }

func init() {
	proto.RegisterClientType("Af",
		func() proto.Deserializer { return &GetQueuePosition{} })
}
