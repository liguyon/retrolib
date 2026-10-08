package specialization

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type Set struct {
	proto.ServerSide

	SpecializationID int
}

func (s *Set) Opcode() proto.Opcode { return "ZS" }

func (s *Set) Serialize() (string, error) {
	return fmt.Sprintf("%d", s.SpecializationID), nil
}

func (s *Set) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	id, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	s.SpecializationID = id
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Set{} })
}
