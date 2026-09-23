package aks

import (
	"github.com/liguyon/retrolib/proto"
)

type RPong struct {
	Payload string
}

func (r *RPong) Opcode() proto.Opcode { return "rpong" }

func (r *RPong) Serialize() (string, error) {
	if len(r.Payload) != 5 {
		return "", proto.ErrMalformedPayload
	}
	return r.Payload, nil
}

func (r *RPong) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	if len(payload) != 5 {
		return proto.ErrMalformedPayload
	}

	r.Payload = payload
	return nil
}

func init() {
	proto.RegisterClientType("rpong",
		func() proto.Deserializer { return &RPong{} })
}
