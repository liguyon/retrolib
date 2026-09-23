package game

import (
	"github.com/liguyon/retrolib/proto"
)

type MapLoaded struct{}

func (m *MapLoaded) Opcode() proto.Opcode { return "GDK" }

func (m *MapLoaded) Serialize() (string, error) { return "", nil }

func (m *MapLoaded) Deserialize(payload string) error { return nil }

func init() {
	proto.RegisterServerType("GDK",
		func() proto.Deserializer { return &MapLoaded{} })
}
