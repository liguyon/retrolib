package chat

import (
	"github.com/liguyon/retrolib/proto"
)

type ServerMessage struct {
	Message string
}

func (s *ServerMessage) Opcode() proto.Opcode { return "cs" }

func (s *ServerMessage) Serialize() (string, error) {
	if s.Message == "" {
		return "", proto.ErrMissingPayload
	}
	return s.Message, nil
}

func (s *ServerMessage) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	s.Message = payload
}

func init() {
	proto.RegisterServerType("cs",
		func() proto.Deserializer { return &ServerMessage{} })
}
