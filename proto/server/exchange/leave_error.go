package exchange

import (
	"github.com/liguyon/retrolib/proto"
)

type LeaveError struct {
	proto.ServerSide
}

func (l *LeaveError) Opcode() proto.Opcode { return "EVE" }

func (l *LeaveError) Serialize() (string, error) { return "", nil }

func (l *LeaveError) Deserialize(payload string) error { return nil }

func init() {
	proto.Register(func() proto.Deserializer { return &LeaveError{} })
}
