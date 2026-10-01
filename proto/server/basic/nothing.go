package basic

import (
	"github.com/liguyon/retrolib/proto"
)

type Nothing struct {
}

func (n *Nothing) Opcode() proto.Opcode { return "BN" }

func (n *Nothing) Serialize() (string, error) {
	return "", nil
}

func (n *Nothing) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterServerType("BN",
		func() proto.Deserializer { return &Nothing{} })
}
