package info

import (
	"fmt"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type MessageDTO struct {
	ID   string
	Args string
}

func (m MessageDTO) Serialize() string {
	if m.Args == "" {
		return m.ID
	}
	return fmt.Sprintf("%s;%s", m.ID, m.Args)
}

func parseMessageDTO(s string) MessageDTO {
	sli := strings.Split(s, ";")
	msg := MessageDTO{ID: sli[0]}
	if len(sli) == 2 {
		msg.Args = sli[1]
	}
	return msg
}

type Message struct {
	proto.ServerSide

	Channel  byte
	Messages []MessageDTO
}

func (m *Message) Opcode() proto.Opcode { return "Im" }

func (m *Message) Serialize() (string, error) {
	var ser []string
	for _, v := range m.Messages {
		ser = append(ser, v.Serialize())
	}
	return fmt.Sprintf("%c%s", m.Channel, strings.Join(ser, "|")), nil
}

func (m *Message) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	m.Channel = payload[0]

	for part := range strings.SplitSeq(payload[1:], "|") {
		m.Messages = append(m.Messages, parseMessageDTO(part))
	}
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &Message{} })
}
