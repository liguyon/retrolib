package account

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type SelectCharacter struct {
	proto.ClientSide

	CharacterID int
}

func (s *SelectCharacter) Opcode() proto.Opcode { return "AS" }

func (s *SelectCharacter) Serialize() (string, error) {
	return fmt.Sprintf("%d", s.CharacterID), nil
}

func (s *SelectCharacter) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	id, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	s.CharacterID = id
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &SelectCharacter{} })
}
