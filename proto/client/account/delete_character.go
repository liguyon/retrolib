package account

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type DeleteCharacter struct {
	proto.ClientSide

	CharacterID  int
	SecretAnswer string
}

func (d *DeleteCharacter) Opcode() proto.Opcode { return "AD" }

func (d *DeleteCharacter) Serialize() (string, error) {
	return fmt.Sprintf("%d|%s", d.CharacterID, d.SecretAnswer), nil
}

func (d *DeleteCharacter) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload, "|")
	if len(sli) != 2 {
		return errors.New("invalid field count")
	}
	id, err := strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[0])
	}
	d.CharacterID = id
	d.SecretAnswer = sli[0]
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &DeleteCharacter{} })
}
