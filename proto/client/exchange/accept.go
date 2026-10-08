package exchange

import (
	"github.com/liguyon/retrolib/proto"
)

type Accept struct {
}

func (a *Accept) Opcode() proto.Opcode { return "EK" }

func (a *Accept) Serialize() (string, error) {
	return "", nil
}

func (a *Accept) Deserialize(payload string) error {
	return nil
}

func init() {
	proto.RegisterClientType("EK",
		func() proto.Deserializer { return &Accept{} })
}
