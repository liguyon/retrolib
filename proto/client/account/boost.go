package account

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type Boost struct {
	proto.ClientSide

	CharacteristicID int
	Amount           int
}

func (b *Boost) Opcode() proto.Opcode { return "AB" }

func (b *Boost) Serialize() (string, error) {
	return fmt.Sprintf("%d|%d", b.CharacteristicID, b.Amount), nil
}

func (b *Boost) Deserialize(payload string) error {
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
	n, err := strconv.Atoi(sli[1])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[1])
	}
	b.CharacteristicID = id
	b.Amount = n
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Boost{} })
}
