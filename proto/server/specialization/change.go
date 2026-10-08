package specialization

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type Change struct {
	SpecializationID int
}

func (c *Change) Opcode() proto.Opcode { return "ZC" }

func (c *Change) Serialize() (string, error) {
	return fmt.Sprintf("%d", c.SpecializationID), nil
}

func (c *Change) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	id, err := strconv.Atoi(payload)
	if err != nil {
		return proto.ErrMalformedPayload
	}
	c.SpecializationID = id
	return nil
}

func init() {
	proto.RegisterServerType("ZC",
		func() proto.Deserializer { return &Change{} })
}
