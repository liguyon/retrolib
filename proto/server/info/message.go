package info

import (
	"fmt"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type InfoMessage struct {
	ID   string
	Args string
}

func (i InfoMessage) Serialize() string {
	if i.Args == "" {
		return i.ID
	}
	return fmt.Sprintf("%s;%s", i.ID, i.Args)
}

type Message struct {
	Channel  byte
	Messages []InfoMessage
}

func (m *Message) Opcode() proto.Opcode { return "Im" }

func (m *Message) Serialize() (string, error) {
	var msgs []string
	for _, entry := range m.Messages {
		msgs = append(msgs, entry.Serialize())
	}
	return fmt.Sprintf("%c%s", m.Channel, strings.Join(msgs, "|")), nil
}

func (m *Message) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	m.Channel = payload[0]

	sli := strings.Split(payload[1:], "|")
	for _, entry := range sli {
		toks := strings.Split(entry, ";")
		msg := InfoMessage{ID: toks[0]}
		if len(toks) == 2 {
			msg.Args = toks[1]
		}
		m.Messages = append(m.Messages, msg)
	}
	return nil
}

func init() {
	proto.RegisterServerType("Im",
		func() proto.Deserializer { return &Message{} })
}
