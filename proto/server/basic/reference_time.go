package basic

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type ReferenceTime struct {
	Timestamp int64
}

func (r *ReferenceTime) Opcode() proto.Opcode { return "BT" }

func (r *ReferenceTime) Serialize() (string, error) {
	return fmt.Sprintf("%d", r.Timestamp), nil
}

func (r *ReferenceTime) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	t, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return proto.ErrMalformedPayload
	}
	r.Timestamp = t
	return nil
}

func init() {
	proto.RegisterServerType("BT",
		func() proto.Deserializer { return &ReferenceTime{} })
}
