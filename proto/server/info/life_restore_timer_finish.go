package info

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type LifeRestoreTimerFinish struct {
	TotalHealed int
}

func (l *LifeRestoreTimerFinish) Opcode() proto.Opcode { return "ILF" }

func (l *LifeRestoreTimerFinish) Serialize() (string, error) {
	return fmt.Sprintf("%d", l.TotalHealed), nil
}

func (l *LifeRestoreTimerFinish) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	n, err := strconv.Atoi(payload)
	if err != nil {
		return proto.ErrMalformedPayload
	}
	l.TotalHealed = n
	return nil
}

func init() {
	proto.RegisterServerType("ILF",
		func() proto.Deserializer { return &LifeRestoreTimerFinish{} })
}
