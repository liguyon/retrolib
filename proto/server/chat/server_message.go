package chat

import (
	"errors"

	"github.com/liguyon/retrolib/proto"
)

type ServerMessage struct {
	proto.ServerSide

	Message string
}

func (s *ServerMessage) Opcode() proto.Opcode { return "cs" }

func (s *ServerMessage) Serialize() (string, error) {
	if s.Message == "" {
		return "", errors.New("empty message")
	}
	return s.Message, nil
}

func (s *ServerMessage) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	s.Message = payload
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &ServerMessage{} })
}
