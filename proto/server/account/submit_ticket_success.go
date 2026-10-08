package account

import (
	"fmt"

	"github.com/liguyon/retrolib/proto"
)

type SubmitTicketSuccess struct {
	proto.ServerSide

	KeyID byte
}

func (s *SubmitTicketSuccess) Opcode() proto.Opcode { return "ATK" }

func (s *SubmitTicketSuccess) Serialize() (string, error) {
	c, err := proto.HexDigit(s.KeyID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%c", c), nil
}

func (s *SubmitTicketSuccess) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	if len(payload) != 1 {
		return fmt.Errorf("not a byte: %q", payload)
	}

	n, err := proto.HexNibble(payload[0])
	if err != nil {
		return fmt.Errorf("%w: %c", err, payload[0])
	}
	s.KeyID = n

	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &SubmitTicketSuccess{} })
}
