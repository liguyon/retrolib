package proto

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCutDelim(t *testing.T) {
	tests := []struct {
		name string
		msg  []byte
		dir  Direction
		want []byte
	}{
		{"svr", []byte("AlK0\x00"), ServerToClient, []byte("AlK0")},
		{"cli", []byte("Af\n\x00"), ClientToServer, []byte("Af")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CutDelim(tt.msg, tt.dir)
			if err != nil {
				t.Fatalf("unexpected err=%v", err)
			}
			if !bytes.Equal(tt.want, res) {
				t.Errorf("want %q; got %q", string(tt.want), string(res))
			}
		})
	}
}

func TestCutDelim_Err(t *testing.T) {
	tests := []struct {
		name  string
		frame []byte
		dir   Direction
	}{
		{"no delim", []byte("AlK0"), ServerToClient},
		{"client missing newline", []byte("Af\x00"), ClientToServer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CutDelim(tt.frame, tt.dir)
			if !errors.Is(err, ErrNoDelimiter) {
				t.Fatalf("want ErrNoDelimiter; got err=%v", err)
			}
		})
	}
}

func TestAppendDelim(t *testing.T) {
	tests := []struct {
		name  string
		frame []byte
		dir   Direction
		want  []byte
	}{
		{"svr", []byte("AlK0"), ServerToClient, []byte("AlK0\x00")},
		{"cli", []byte("Af"), ClientToServer, []byte("Af\n\x00")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := AppendDelim(tt.frame, tt.dir)
			if err != nil {
				t.Fatalf("unexpected err=%v", err)
			}
			if !bytes.Equal(res, tt.want) {
				t.Errorf("want %q; got %q", string(tt.want), string(res))
			}
		})
	}
}

func TestAppendDelim_Err(t *testing.T) {
	tests := []struct {
		name  string
		frame []byte
		dir   Direction
	}{
		{"svr", []byte("AlK0\x00"), ServerToClient},
		{"cli", []byte("Af\n\x00"), ClientToServer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AppendDelim(tt.frame, tt.dir)
			if !errors.Is(err, ErrDelimiterInBody) {
				t.Fatalf("want ErrDelimiterInBody; got err=%v", err)
			}
		})
	}
}

type fakeHello struct {
	ServerSide
	Rest string
}

func (f *fakeHello) Opcode() Opcode { return "hello" }
func (f *fakeHello) Deserialize(payload string) error {
	if payload == "" {
		return ErrMissingPayload
	}
	f.Rest = payload
	return nil
}
func (f *fakeHello) Serialize() (string, error) {
	if f.Rest == "" {
		return "", errors.New("empty")
	}
	return f.Rest, nil
}

type fakeHelloCli struct {
	ClientSide
	Rest string
}

func (f *fakeHelloCli) Opcode() Opcode { return "hello" }
func (f *fakeHelloCli) Deserialize(payload string) error {
	if payload == "" {
		return ErrMissingPayload
	}
	f.Rest = payload
	return nil
}
func (f *fakeHelloCli) Serialize() (string, error) {
	if f.Rest == "" {
		return "", errors.New("empty")
	}
	return f.Rest, nil
}

type fakeEmpty struct{ ServerSide }

func (f *fakeEmpty) Opcode() Opcode           { return "" }
func (f *fakeEmpty) Deserialize(string) error { return nil }

func newTestRegistry(t *testing.T) *TypeRegistry {
	t.Helper()

	reg := NewTypeRegistry()

	err := reg.Register(func() Deserializer { return &fakeHello{} })
	if err != nil {
		t.Fatalf("unexpected err=%v", err)
	}

	err = reg.Register(func() Deserializer { return &fakeHello{} })
	if !errors.Is(err, ErrDuplicateOpcode) {
		t.Fatalf("want ErrDuplicateOpcode; got err=%v", err)
	}

	err = reg.Register(func() Deserializer { return &fakeHelloCli{} })
	if err != nil {
		t.Fatalf("unexpected err=%v", err)
	}

	err = reg.Register(func() Deserializer { return &fakeHelloCli{} })
	if !errors.Is(err, ErrDuplicateOpcode) {
		t.Fatalf("want ErrDuplicateOpcode; got err=%v", err)
	}

	err = reg.Register(func() Deserializer { return &fakeEmpty{} })
	if !errors.Is(err, ErrInvalidOpcode) {
		t.Fatalf("want ErrInvalidOpcode; got err=%v", err)
	}
	return reg
}

func TestParseFrame(t *testing.T) {
	reg := newTestRegistry(t)

	tests := []struct {
		name        string
		frame       string
		dir         Direction
		wantOp      Opcode
		wantPayload string
		wantErr     error
	}{
		{"valid svr", "helloworld", ServerToClient, "hello", "world", nil},
		{"valid cli", "helloworld", ClientToServer, "hello", "world", nil},
		{"err empty", "", ClientToServer, "", "", ErrInvalidOpcode},
		{"err unknown", "holazawarudo", ServerToClient, "", "", ErrUnknownOpcode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, pl, err := reg.ParseFrame([]byte(tt.frame), tt.dir)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v; want %v", err, tt.wantErr)
			}
			if op != tt.wantOp {
				t.Fatalf("op=%q; want %q", op, tt.wantOp)
			}
			if pl != tt.wantPayload {
				t.Fatalf("payload=%q; want %q", pl, tt.wantPayload)
			}
		})
	}
}

func TestDeserializeMessage(t *testing.T) {
	reg := newTestRegistry(t)

	tests := []struct {
		name    string
		op      Opcode
		pl      string
		dir     Direction
		wantMsg Message
		wantErr error
	}{
		{"valid svr", "hello", "World", ServerToClient,
			&fakeHello{Rest: "World"}, nil},
		{"valid cli", "hello", "lesgens", ClientToServer,
			&fakeHelloCli{Rest: "lesgens"}, nil},
		{"unknown op", "xinchao", "moinguoi", ClientToServer,
			nil, ErrUnknownOpcode},
		{"empty op", "   ", "", ClientToServer,
			nil, ErrInvalidOpcode},
		{"deserialize error propagates", "hello", "", ServerToClient,
			nil, ErrMissingPayload},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := reg.DeserializeMessage(tt.op, tt.pl, tt.dir)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v; want %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantMsg, msg); diff != "" {
				t.Errorf("-want +got:\n%s", diff)
			}
		})
	}
}

func TestSerializeMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     Serializer
		want    []byte
		wantErr error
	}{
		{"valid svr", &fakeHello{Rest: "world"}, []byte("helloworld"), nil},
		{"valid cli", &fakeHelloCli{Rest: "olleh"}, []byte("helloolleh"), nil},
		{"serialize error propagates", &fakeHello{}, nil, ErrInvalidMessage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := SerializeMessage(tt.msg)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v; want %v", err, tt.wantErr)
			}
			if !bytes.Equal(body, tt.want) {
				t.Fatalf("body=%q; want %q", body, tt.want)
			}
		})
	}
}
