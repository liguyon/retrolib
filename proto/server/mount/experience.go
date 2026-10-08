package mount

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type Experience struct {
	proto.ServerSide

	Percent int
}

func (e *Experience) Opcode() proto.Opcode { return "Rx" }

func (e *Experience) Serialize() (string, error) {
	return fmt.Sprintf("%d", e.Percent), nil
}

func (e *Experience) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	var err error
	e.Percent, err = strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Experience{} })
}
