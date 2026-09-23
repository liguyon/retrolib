package account

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type Boost struct {
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
		return proto.ErrMalformedPayload
	}
	id, err := strconv.Atoi(sli[0])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	n, err := strconv.Atoi(sli[1])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	b.CharacteristicID = id
	b.Amount = n
	return nil
}

func init() {
	proto.RegisterClientType("AB",
		func() proto.Deserializer { return &Boost{} })
}
