package chat

import (
	"fmt"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type Send struct {
	Destination string
	Message     string
	Extra       string // TODO: items blob: <item_id>!<compressed_effects>
}

func (s *Send) Opcode() proto.Opcode { return "BM" }

func (s *Send) Serialize() (string, error) {
	return fmt.Sprintf("%s|%s|%s", s.Destination, s.Message, s.Extra), nil
}

func (s *Send) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload, "|")
	if len(sli) != 3 {
		return proto.ErrMalformedPayload
	}

	s.Destination = sli[0]
	s.Message = sli[1]
	s.Extra = sli[2]
	return nil
}

func init() {
	proto.RegisterClientType("BM",
		func() proto.Deserializer { return &Send{} })
}
