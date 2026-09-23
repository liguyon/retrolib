package aks

import (
	"github.com/liguyon/retrolib/proto"
)

type QuickPing struct {
}

func (q *QuickPing) Opcode() proto.Opcode { return "qping" }

func (q *QuickPing) Serialize() (string, error) {
	return "", nil
}

func (q *QuickPing) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterClientType("qping",
		func() proto.Deserializer { return &QuickPing{} })
}