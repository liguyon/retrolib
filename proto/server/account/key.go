package account

import (
	"github.com/liguyon/retrolib/proto"
)

type Key struct {
	EncodedKey string
}

func (k *Key) Opcode() proto.Opcode { return "AK" }

func (k *Key) Serialize() (string, error) {
	return k.EncodedKey, nil
}

func (k *Key) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	k.EncodedKey = payload
	return nil
}

func init() {
	proto.RegisterServerType("AK", func() proto.Deserializer { return &Key{} })
}
