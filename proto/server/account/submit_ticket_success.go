package account

import (
	"fmt"

	"github.com/liguyon/retrolib/proto"
)

type SubmitTicketSuccess struct {
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
		return proto.ErrMalformedPayload
	}

	n, err := proto.HexNibble(payload[0])
	if err != nil {
		return fmt.Errorf("%w: %v", proto.ErrMalformedKey, err)
	}
	s.KeyID = n

	return nil
}

func init() {
	proto.RegisterServerType("ATK",
		func() proto.Deserializer { return &SubmitTicketSuccess{} })
}
