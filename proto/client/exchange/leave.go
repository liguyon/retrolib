package exchange

import (
	"github.com/liguyon/retrolib/proto"
)

type Leave struct {
	proto.ClientSide
}

func (l *Leave) Opcode() proto.Opcode { return "EV" }

func (l *Leave) Serialize() (string, error) { return "", nil }

func (l *Leave) Deserialize(payload string) error { return nil }

func init() {
	proto.Register(func() proto.Deserializer { return &Leave{} })
}
