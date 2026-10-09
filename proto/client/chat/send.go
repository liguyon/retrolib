package chat

import (
	"errors"
	"fmt"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type Send struct {
	proto.ClientSide

	Destination string
	Message     string
	Extra       string // TODO: item blobs: <item_id>!<compressed_effects>
}

func (s *Send) Opcode() proto.Opcode { return "BM" }

func (s *Send) Serialize() (string, error) {
	if s.Destination == "" {
		return "", errors.New("empty destination")
	}
	if s.Message == "" {
		return "", errors.New("empty message")
	}
	return fmt.Sprintf("%s|%s|%s", s.Destination, s.Message, s.Extra), nil
}

func (s *Send) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload, "|")
	if len(sli) != 3 {
		return errors.New("invalid field count")
	}

	s.Destination = sli[0]
	if s.Destination == "" {
		return errors.New("empty destination")
	}
	s.Message = sli[1]
	if s.Message == "" {
		return errors.New("empty message")
	}
	s.Extra = sli[2]
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Send{} })
}
