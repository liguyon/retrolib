package account

import (
	"errors"

	"github.com/liguyon/retrolib/proto"
)

type Rescue struct {
	proto.ClientSide

	Ticket string
}

func (r *Rescue) Opcode() proto.Opcode { return "Ar" }

func (r *Rescue) Serialize() (string, error) {
	if r.Ticket == "" {
		return "", errors.New("empty ticket")
	}
	return r.Ticket, nil
}

func (r *Rescue) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	r.Ticket = payload
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Rescue{} })
}
