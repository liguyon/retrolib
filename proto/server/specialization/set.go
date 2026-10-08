package specialization

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type Set struct {
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
		return proto.ErrMalformedPayload
	}
	s.SpecializationID = id
	return nil
}

func init() {
	proto.RegisterServerType("ZS",
		func() proto.Deserializer { return &Set{} })
}
