package account

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type SelectServer struct {
	ServerID int
}

func (s *SelectServer) Opcode() proto.Opcode { return "AX" }

func (s *SelectServer) Serialize() (string, error) {
	return fmt.Sprintf("%s", s.ServerID), nil
}

func (s *SelectServer) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	id, err := strconv.Atoi(payload)
	if err != nil {
		return proto.ErrMalformedPayload
	}
	s.ServerID = id
	return nil
}

func init() {
	proto.RegisterClientType("AX",
		func() proto.Deserializer { return &SelectServer{} })
}
