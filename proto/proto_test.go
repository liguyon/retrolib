package proto

import (
	"bytes"
	"errors"
	"testing"
)

func TestTrimMessageDelim(t *testing.T) {
	tests := []struct {
		name     string
		msg      []byte
		dir      Direction
		expected []byte
	}{
		{"svr", []byte("AlK0\x00"), ServerToClient, []byte("AlK0")},
		{"cli", []byte("Af\n\x00"), ClientToServer, []byte("Af")},
		{"no delim", []byte("Af"), ClientToServer, []byte("Af")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := TrimMessageDelim(tt.msg, tt.dir)
			if bytes.Compare(tt.expected, res) != 0 {
				t.Errorf("want %q; got %q", string(tt.expected), string(res))
			}
		})
	}
}

func TestAppendMessageDelim(t *testing.T) {
	tests := []struct {
		name     string
		msg      []byte
		dir      Direction
		expected []byte
	}{
		{"svr", []byte("AlK0"), ServerToClient, []byte("AlK0\x00")},
		{"cli", []byte("Af"), ClientToServer, []byte("Af\n\x00")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := AppendMessageDelim(tt.msg, tt.dir)
			if bytes.Compare(res, tt.expected) != 0 {
				t.Errorf("want %q; got %q", string(tt.expected), string(res))
			}
		})
	}
}

func TestRegisterType(t *testing.T) {
	clientRegistry = map[Opcode]func() Deserializer{}
	maxClientOpcodeLen = 0

	tests := []struct {
		name        string
		op          Opcode
		dir         Direction
		expectedLen int
	}{
		{"cli", "Ax", ClientToServer, 2},
		{"longer", "ABCDEF", ClientToServer, 6},
		{"len unchanged", "AH", ClientToServer, 6},
		// test panic?
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterType(tt.op, nil, tt.dir)
			_, exists := clientRegistry[tt.op]
			if !exists {
				t.Fatal("not registered")
			}
			if maxClientOpcodeLen != tt.expectedLen {
				t.Errorf("want maxLen %d; got %d", tt.expectedLen,
					maxClientOpcodeLen)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	clientRegistry = map[Opcode]func() Deserializer{}
	maxClientOpcodeLen = 0
	RegisterClientType("AT", nil)
	RegisterClientType("Af", nil)

	tests := []struct {
		name            string
		msg             []byte
		expectedOpcode  Opcode
		expectedPayload string
		expectedErr     error
	}{
		{"no payload", []byte("Af"), "Af", "", nil},
		{"with payload", []byte("ATblabla"), "AT", "blabla", nil},
		{"empty msg", []byte{}, "", "", ErrInvalidOpcode},
		{"unknown opcode", []byte("HCblabla"), "", "", ErrUnknownOpcode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, pl, err := ParseClientMessage(tt.msg)
			if err == nil && tt.expectedErr != nil {
				t.Fatalf("want err; got none")
			}
			if err != nil && tt.expectedErr == nil {
				t.Fatalf("want no err; got %v", err)
			}

			if err != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("want err %v; got %v", tt.expectedErr, err)
				}
			}

			if op != tt.expectedOpcode {
				t.Errorf("want opcode %q; got %q", tt.expectedOpcode, op)
			}

			if pl != tt.expectedPayload {
				t.Errorf("want payload %q; got %q", tt.expectedPayload, pl)
			}
		})
	}
}

type fakeTicket struct {
	Ticket string
}

func (f *fakeTicket) Opcode() Opcode { return "AT" }

func (f *fakeTicket) Serialize() (string, error) { return f.Ticket, nil }

func (f *fakeTicket) Deserialize(payload string) error {
	f.Ticket = payload
	return nil
}

type fakeDisconnect struct{}

func (f *fakeDisconnect) Opcode() Opcode { return "BYE" }

func (f *fakeDisconnect) Serialize() (string, error) { return "", nil }

func (f *fakeDisconnect) Deserialize(payload string) error { return nil }

func TestSerializeMessage(t *testing.T) {
	tests := []struct {
		name        string
		msg         Serializer
		expected    []byte
		expectedErr error
	}{
		{"with payload", &fakeTicket{Ticket: "ABCdef123"}, []byte("ATABCdef123"), nil},
		{"empty payload", &fakeDisconnect{}, []byte("BYE"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := SerializeMessage(tt.msg)
			if err == nil && tt.expectedErr != nil {
				t.Fatalf("want err; got none")
			}
			if err != nil && tt.expectedErr == nil {
				t.Fatalf("want no err; got %v", err)
			}
			if err != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("want err %v; got %v", tt.expectedErr, err)
				}
			}

			if bytes.Compare(res, tt.expected) != 0 {
				t.Errorf("want %q; got %q", string(tt.expected), string(res))
			}
		})
	}
}

func TestDeserializeMessage(t *testing.T) {
	clientRegistry = map[Opcode]func() Deserializer{}
	maxClientOpcodeLen = 0
	RegisterClientType("AT", func() Deserializer { return &fakeTicket{} })
	RegisterClientType("BYE", func() Deserializer { return &fakeDisconnect{} })

	msg, err := DeserializeClientMessage("AT", "123abc")
	if err != nil {
		t.Fatalf("want no err; got %v", err)
	}
	switch m := msg.(type) {
	case *fakeTicket:
		if m.Ticket != "123abc" {
			t.Fatalf("want ticket=\"123abc\"; got %q", m.Ticket)
		}
	default:
		t.Fatalf("could not type-switch")
	}

	msg, err = DeserializeClientMessage("BYE", "")
	if err != nil {
		t.Fatalf("want no err; got %v", err)
	}
	switch msg.(type) {
	case *fakeDisconnect:
	default:
		t.Fatalf("could not type-switch")
	}
}
