package exchange

import (
	"github.com/liguyon/retrolib/proto"
)

// Actually, the client (1.43.7) doesn't read 'K' or 'E' from the ExchangeLeave
// response. It reads the extra data coming after and tests extra == 'a'.
// It seems that 'a' is sent when the exchange is ended due to both parties
// accepting the exchange instead of cancelling it.
type LeaveSuccess struct {
	OnAccept bool
}

func (l *LeaveSuccess) Opcode() proto.Opcode { return "EVK" }

func (l *LeaveSuccess) Serialize() (string, error) {
	if l.OnAccept {
		return "a", nil
	}
	return "", nil
}

func (l *LeaveSuccess) Deserialize(payload string) error {
	if payload == "a" {
		l.OnAccept = true
	}
	return nil
}

func init() {
	proto.RegisterServerType("EVK",
		func() proto.Deserializer { return &LeaveSuccess{} })
}
