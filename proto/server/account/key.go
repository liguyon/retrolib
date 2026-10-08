package account

import (
	"errors"

	"github.com/liguyon/retrolib/proto"
)

type Key struct {
	proto.ServerSide

	EncodedKey string
}

func (k *Key) Opcode() proto.Opcode { return "AK" }

func (k *Key) Serialize() (string, error) {
	if k.EncodedKey == "" {
		return "", errors.New("empty key")
	}
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
	proto.Register(func() proto.Deserializer { return &Key{} })
}
