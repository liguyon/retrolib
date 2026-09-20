package account

import (
	"github.com/liguyon/retrolib/proto"
)

type SubmitTicket struct {
	Ticket string
}

func (t *SubmitTicket) Opcode() proto.Opcode { return "AT" }

func (t *SubmitTicket) Serialize() (string, error) {
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
	proto.RegisterClientType("AT", func() proto.Deserializer { return &SubmitTicket{} })
}
