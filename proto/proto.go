// Package proto implements the wire format for Retro's network protocol.
// Each frame contains a single delimited message, optionally obfuscated
// (see retrolib/crypto), carrying an opcode-prefixed payload.
//
// The two directions of traffic are handled as separate and composable stages.
//
// Inbound traffic generally goes through the following pipeline:
// raw frame []byte -> TrimDelim -> Decrypt-> ParseFrame -> DeserializeMessage -> Message
//
// Outbound traffic:
// Message (concrete type) -> SerializeMessage -> Encrypt-> AppendDelim -> raw frame []byte
package proto

import (
	"bytes"
	"fmt"
	"strings"
)

// Direction identifies which way a message travels.
type Direction int

const (
	ClientToServer Direction = iota
	ServerToClient
)

func (d Direction) String() string {
	if d == ServerToClient {
		return "s2c"
	}
	return "c2s"
}

// Delim returns the end-of-message delimiter for a traffic direction.
func (d Direction) Delim() []byte {
	if d == ServerToClient {
		return []byte("\x00")
	}
	return []byte("\n\x00")
}

// AppendDelim appends the direction's full delimiter to a serialized and typically encrypted // message, producing a frame ready to write to the connection.
// Returns ErrDelimiterInBody if body already contains the delimiter, since that would split
// the frame on the wire.
func AppendDelim(body []byte, dir Direction) ([]byte, error) {
	if bytes.HasSuffix(body, []byte("\x00")) {
		return nil, ErrDelimiterInBody
	}
	frame := make([]byte, 0, len(body)+len(dir.Delim()))
	frame = append(frame, body...)
	frame = append(frame, dir.Delim()...)
	return frame, nil
}

// CutDelim removes the direction's full delimiter from a raw frame read off the wire.
// Returns ErrNoDelimiter if the frame doesn't end with it (e.g. a client frame missing '\n').
func CutDelim(frame []byte, dir Direction) ([]byte, error) {
	if !bytes.HasSuffix(frame, dir.Delim()) {
		return nil, ErrNoDelimiter
	}
	return bytes.TrimSuffix(frame, dir.Delim()), nil
}

// Opcode is the prefix identifying a message type.
type Opcode string

// Validate verifies that opcode isn't empty and isn't prefixed or suffixed with spaces.
// Returns ErrInvalidOpcode on error.
func (o Opcode) Validate() error {
	trimmed := strings.TrimSpace(string(o))
	if o == "" || trimmed == "" {
		return fmt.Errorf("%w: empty", ErrInvalidOpcode)
	}
	if trimmed != string(o) {
		return fmt.Errorf("%w: leading or trailing spaces", ErrInvalidOpcode)
	}
	return nil
}

// Message is implemented by concrete message types. It's the interface callers type-switch
// on after deserializing.
type Message interface {
	Opcode() Opcode
	Direction() Direction
}

type ClientMessage interface {
	Message
	clientSide()
}

type ClientSide struct{}

func (ClientSide) Direction() Direction { return ClientToServer }
func (ClientSide) clientSide()          {}

type ServerMessage interface {
	Message
	serverSide()
}

type ServerSide struct{}

func (ServerSide) Direction() Direction { return ServerToClient }
func (ServerSide) serverSide()          {}

// Serializer is implemented by message types that can serialize their payload (everything
// that comes after the opcode prefix).
type Serializer interface {
	Message
	Serialize() (string, error)
}

// Deserializer is implemented by message types that can be populated from deserializing
// an inbound payload (everything that comes after the opcode prefix).
type Deserializer interface {
	Message
	Deserialize(payload string) error
}

// TypeRegistry maps opcodes to constructors for concrete message types.
// A default registry is declared in this package (see var. "registry" below). It can be
// populated with the types implemented under retrolib/proto/{client,server}/* by importing
// those packages. Importing retrolib/proto/{client,server}/all blank-imports all namespaces
// for convenience.
// Custom type implementations can be registered in the default registry by calling
// "RegisterDefault".
// Users can also create their TypeRegistry to separate their own type implementations.
type TypeRegistry struct {
	reg       map[Direction]map[Opcode]func() Deserializer
	maxLenCli int
	maxLenSvr int
}

// NewTypeRegistry constructs an empty message type registry.
func NewTypeRegistry() *TypeRegistry {
	cli := make(map[Opcode]func() Deserializer)
	svr := make(map[Opcode]func() Deserializer)
	return &TypeRegistry{
		reg: map[Direction]map[Opcode]func() Deserializer{
			ClientToServer: cli, ServerToClient: svr},
		maxLenCli: 0,
		maxLenSvr: 0,
	}
}

// registry is the default registry where concrete message types under proto/{client,server}
// register their deserializer.
var registry = NewTypeRegistry()

// Register registers a concrete message type that can be deserialized, so that frames
// with the associated opcode can later be parsed and deserialized.
func (r *TypeRegistry) Register(factory func() Deserializer) error {
	msg := factory()

	if err := msg.Opcode().Validate(); err != nil {
		return fmt.Errorf("%w: %q", err, msg.Opcode())
	}

	reg := r.reg[msg.Direction()]
	_, exists := reg[msg.Opcode()]
	if exists {
		return fmt.Errorf("%w: %s: %q",
			ErrDuplicateOpcode, msg.Direction(), msg.Opcode())
	}

	reg[msg.Opcode()] = factory
	switch msg.Direction() {
	case ClientToServer:
		if len(msg.Opcode()) > r.maxLenCli {
			r.maxLenCli = len(msg.Opcode())
		}
	case ServerToClient:
		if len(msg.Opcode()) > r.maxLenSvr {
			r.maxLenSvr = len(msg.Opcode())
		}
	}
	return nil
}

// Register registers in the default registry a concrete message type that can be
// deserialized.
func Register(factory func() Deserializer) {
	err := registry.Register(factory)
	if err != nil {
		panic(err)
	}
}

// ParseFrame extracts the opcode and the raw payload from a decrypted frame.
func (r *TypeRegistry) ParseFrame(
	frame []byte, dir Direction) (op Opcode, payload string, err error) {
	if len(frame) == 0 {
		return "", "", fmt.Errorf("%w: empty message", ErrInvalidOpcode)
	}

	reg := r.reg[dir]
	var maxLen int

	switch dir {
	case ClientToServer:
		maxLen = r.maxLenCli
	case ServerToClient:
		maxLen = r.maxLenSvr
	}

	for n := maxLen; n >= 1; n-- {
		if len(frame) < n {
			continue
		}
		_, exists := reg[Opcode(frame[:n])]
		if exists {
			return Opcode(frame[:n]), string(frame[n:]), nil
		}
	}
	return "", "", fmt.Errorf("%w: %q", ErrUnknownOpcode, string(frame))
}

// ParseFrame extracts the opcode and the payload from a message using the default registry.
func ParseFrame(frame []byte, dir Direction) (op Opcode, payload string, err error) {
	return registry.ParseFrame(frame, dir)
}

// DeserializeMessage instantiates the message type associated with op and populates it from
// payload. The returned Message is meant to be type-switched on to get the concrete type.
func (r *TypeRegistry) DeserializeMessage(
	op Opcode, payload string, dir Direction) (Message, error) {
	reg := r.reg[dir]

	if err := op.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %q", err, op)
	}

	factory, exists := reg[op]
	if !exists {
		return nil, fmt.Errorf("%w: %q", ErrUnknownOpcode, op)
	}
	msg := factory()

	err := msg.Deserialize(payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMalformedPayload, err)
	}
	return msg, nil
}

// DeserializeMessage instantiates the message type registered in the default registry
// for op and populates it from payload. The returned Message is meant to be type-switched on
// to get the concrete type.
func DeserializeMessage(op Opcode, payload string, dir Direction) (Message, error) {
	return registry.DeserializeMessage(op, payload, dir)
}

// SerializeMessage serializes a message into its pre-processed wire format: the opcode prefix
// followed by the payload string.
func SerializeMessage(msg Serializer) ([]byte, error) {
	body, err := msg.Serialize()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidMessage, err)
	}
	return []byte(fmt.Sprintf("%s%s", msg.Opcode(), body)), nil
}
