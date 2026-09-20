package account

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type RegionalVersion struct {
	Version int
}

func (r *RegionalVersion) Opcode() proto.Opcode { return "AV" }

func (r *RegionalVersion) Serialize() (string, error) {
	return fmt.Sprintf("%d", r.Version), nil
}

func (r *RegionalVersion) Deserialize(payload string) error {
	if len(payload) == 0 {
		return proto.ErrMissingPayload
	}

	v, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("%w: %v", proto.ErrMalformedPayload, err)
	}

	r.Version = v
	return nil
}

func init() {
	proto.RegisterServerType("AV",
		func() proto.Deserializer { return &RegionalVersion{} })
}
