package aks

import (
	"github.com/liguyon/retrolib/proto"
)

type Ping struct {
	proto.ClientSide
}

func (p *Ping) Opcode() proto.Opcode { return "ping" }

func (p *Ping) Serialize() (string, error) {
	return "", nil
}

func (p *Ping) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Ping{} })
}
