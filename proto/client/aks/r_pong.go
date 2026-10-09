package aks

import (
	"errors"

	"github.com/liguyon/retrolib/proto"
)

type RPong struct {
	proto.ClientSide

	Payload string
}

func (r *RPong) Opcode() proto.Opcode { return "rpong" }

func (r *RPong) Serialize() (string, error) {
	if len(r.Payload) != 5 {
		return "", errors.New("invalid payload length")
	}
	return r.Payload, nil
}

func (r *RPong) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	if len(payload) != 5 {
		return errors.New("invalid payload length")
	}

	r.Payload = payload
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &RPong{} })
}
