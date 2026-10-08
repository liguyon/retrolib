package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type MapFrameObject2 struct {
	proto.ServerSide

	CellID        int
	FrameID       string
	IsInteractive bool
}

func (m *MapFrameObject2) Opcode() proto.Opcode { return "GDF" }

func (m *MapFrameObject2) Serialize() (string, error) {
	if m.IsInteractive {
		return fmt.Sprintf("|%d;%s;%d",
			m.CellID, m.FrameID, proto.BoolToInt(true)), nil
	}
	return fmt.Sprintf("|%d;%s", m.CellID, m.FrameID), nil
}

func (m *MapFrameObject2) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload[1:], ";")
	if len(sli) < 2 {
		return errors.New("invalid field count")
	}
	cid, err := strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	m.CellID = cid
	m.FrameID = sli[1]
	if len(sli) == 3 {
		b, err := proto.ParseBool(sli[2])
		if err != nil {
			return fmt.Errorf("%w: %q", err, sli[2])
		}
		m.IsInteractive = b
	}
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &MapFrameObject2{} })
}
