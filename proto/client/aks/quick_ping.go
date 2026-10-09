package aks

import (
	"github.com/liguyon/retrolib/proto"
)

type QuickPing struct {
	proto.ClientSide
}

func (q *QuickPing) Opcode() proto.Opcode { return "qping" }

func (q *QuickPing) Serialize() (string, error) {
	return "", nil
}

func (q *QuickPing) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &QuickPing{} })
}
