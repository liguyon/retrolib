package info

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type LifeRestoreTimerStart struct {
	Interval int
}

func (l *LifeRestoreTimerStart) Opcode() proto.Opcode { return "ILS" }

func (l *LifeRestoreTimerStart) Serialize() (string, error) {
	return fmt.Sprintf("%d", l.Interval), nil
}

func (l *LifeRestoreTimerStart) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	n, err := strconv.Atoi(payload)
	if err != nil {
		return proto.ErrMalformedPayload
	}
	l.Interval = n
	return nil
}

func init() {
	proto.RegisterServerType("ILS",
		func() proto.Deserializer { return &LifeRestoreTimerStart{} })
}
