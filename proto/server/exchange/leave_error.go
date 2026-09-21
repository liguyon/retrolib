package exchange

import (
	"github.com/liguyon/retrolib/proto"
)

type LeaveError struct{}

func (l *LeaveError) Opcode() proto.Opcode { return "EVE" }

func (l *LeaveError) Serialize() (string, error) { return "", nil }

func (l *LeaveError) Deserialize(payload string) error { return nil }

func init() {
	proto.RegisterServerType("EVE",
		func() proto.Deserializer { return &LeaveError{} })
}
