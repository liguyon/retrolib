package account

import (
	"errors"

	"github.com/liguyon/retrolib/proto"
)

type SubmitTicket struct {
	proto.ClientSide

	Ticket string
}

func (t *SubmitTicket) Opcode() proto.Opcode { return "AT" }

func (t *SubmitTicket) Serialize() (string, error) {
	if t.Ticket == "" {
		return "", errors.New("empty ticket")
	}
	return t.Ticket, nil
}

func (t *SubmitTicket) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	t.Ticket = payload
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &SubmitTicket{} })
}
