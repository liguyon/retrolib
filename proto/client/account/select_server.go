package account

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type SelectServer struct {
	proto.ClientSide

	ServerID int
}

func (s *SelectServer) Opcode() proto.Opcode { return "AX" }

func (s *SelectServer) Serialize() (string, error) {
	return fmt.Sprintf("%d", s.ServerID), nil
}

func (s *SelectServer) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	id, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	s.ServerID = id
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &SelectServer{} })
}
