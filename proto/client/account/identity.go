package account

import (
	"errors"

	"github.com/liguyon/retrolib/proto"
)

type Identity struct {
	proto.ClientSide

	ID string
}

func (i *Identity) Opcode() proto.Opcode { return "Ai" }

func (i *Identity) Serialize() (string, error) {
	if i.ID == "" {
		return "", errors.New("empty identity")
	}
	return i.ID, nil
}

func (i *Identity) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	i.ID = payload
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Identity{} })
}
