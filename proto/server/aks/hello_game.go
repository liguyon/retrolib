package aks

import (
	"github.com/liguyon/retrolib/proto"
)

type HelloGame struct {
	proto.ServerSide
}

func (h *HelloGame) Opcode() proto.Opcode { return "HG" }

func (h *HelloGame) Serialize() (string, error) { return "", nil }

func (h *HelloGame) Deserialize(payload string) error { return nil }

func init() {
	proto.Register(func() proto.Deserializer { return &HelloGame{} })
}
