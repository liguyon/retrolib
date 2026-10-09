package account

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type ResetCharacter struct {
	proto.ClientSide

	CharacterID int
}

func (r *ResetCharacter) Opcode() proto.Opcode { return "AR" }

func (r *ResetCharacter) Serialize() (string, error) {
	return fmt.Sprintf("%d", r.CharacterID), nil
}

func (r *ResetCharacter) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	id, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	r.CharacterID = id
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &ResetCharacter{} })
}
