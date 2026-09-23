package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type MapData struct {
	ID        int
	Timestamp string
	Key       string
}

func (m *MapData) Opcode() proto.Opcode { return "GDM" }

func (m *MapData) Serialize() (string, error) {
	return fmt.Sprintf("|%d|%s|%s", m.ID, m.Timestamp, m.Key), nil
}

func (m *MapData) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	sli := strings.Split(payload[1:], "|")
	if len(sli) != 3 {
		return proto.ErrMalformedPayload
	}

	id, err := strconv.Atoi(sli[0])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	m.ID = id
	m.Timestamp = sli[1]
	m.Key = sli[2]
	return nil
}

func init() {
	proto.RegisterServerType("GDM",
		func() proto.Deserializer { return &MapData{} })
}
