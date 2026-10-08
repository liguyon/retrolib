package info

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type LifeRestoreTimerStart struct {
	proto.ServerSide

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
		return fmt.Errorf("invalid number: %q", payload)
	}
	l.Interval = n
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &LifeRestoreTimerStart{} })
}
