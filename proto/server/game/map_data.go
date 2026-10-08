package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type MapData struct {
	proto.ServerSide

	ID   int
	Date string
	Key  string
}

func (m *MapData) Opcode() proto.Opcode { return "GDM" }

func (m *MapData) Serialize() (string, error) {
	return fmt.Sprintf("|%d|%s|%s", m.ID, m.Date, m.Key), nil
}

func (m *MapData) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	sli := strings.Split(payload[1:], "|")
	if len(sli) != 3 {
		return errors.New("invalid field count")
	}

	id, err := strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[0])
	}
	m.ID = id
	m.Date = sli[1]
	m.Key = sli[2]
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &MapData{} })
}
