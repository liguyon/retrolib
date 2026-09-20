package account

import (
	"github.com/liguyon/retrolib/proto"
)

type Identity struct {
	ID string
}

func (i *Identity) Opcode() proto.Opcode { return "Ai" }

func (i *Identity) Serialize() (string, error) { return i.ID, nil }

func (i *Identity) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	i.ID = payload
	return nil
}

func init() {
	proto.RegisterClientType("Ai",
		func() proto.Deserializer { return &Identity{} })
}
