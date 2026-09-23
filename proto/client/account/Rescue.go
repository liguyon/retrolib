package account

import (
	"github.com/liguyon/retrolib/proto"
)

type Rescue struct {
	Ticket string
}

func (r *Rescue) Opcode() proto.Opcode { return "Ar" }

func (r *Rescue) Serialize() (string, error) {
	if r.Ticket == "" {
		return "", proto.ErrMissingPayload
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
	proto.RegisterClientType("Ar",
		func() proto.Deserializer { return &Rescue{} })
}
