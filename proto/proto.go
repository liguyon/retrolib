// Package proto implements the wire format for Retro's network protocol.
// Each frame contains a single delimited message, optionally obfuscated (see crypto.go),
// carrying an opcode-prefixed payload.
//
// The two directions of traffic are handled as separate and composable stages.
//
// Inbound traffic generally goes through the following pipeline:
// raw []byte (read from wire) -> TrimDelim -> Decrypt-> ParseOpcode -> Deserialize -> Message
//
// Outbound traffic:
// Message (concrete, typed) -> Serialize -> Encrypt-> AppendDelim -> raw []byte
package proto

import (
	"bytes"
	"fmt"
)

// Direction identifies which way a message travels.
type Direction int

const (
	ClientToServer Direction = iota
	ServerToClient
)

// Delim returns the end-of-message delimiter for a traffic direction.
func (d Direction) Delim() []byte {
	if d == ServerToClient {
		return []byte("\x00")
	}
	return []byte("\n\x00")
}

// AppendDelim appends the end-of-message delimiter to an encoded
// (and typically encrypted) message, producing a frame ready to write to the connection.
func AppendDelim(msg []byte, dir Direction) []byte {
	frame := make([]byte, 0, len(msg)+len(dir.Delim()))
	frame = append(frame, msg...)
	frame = append(frame, dir.Delim()...)
	return frame
}

// TrimDelim removes the end-of-message delimiter from a raw frame read off the wire.
func TrimDelim(frame []byte, dir Direction) []byte {
	return bytes.TrimSuffix(frame, dir.Delim())
}

// Opcode is the prefix identifying a message type.
type Opcode string

// Message is implemented by concrete packet types. It's the interface callers type-switch
// on after deserializing.
type Message interface {
	Opcode() Opcode
}

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

// TypeRegistry maps opcodes to constructors/factories of the associated message type.
// A default global registry is declared in this package (see var. "registry" below). It can be
// populated with the types implemented under retrolib/proto/{client,server}/* by importing
// those packages. Importing retrolib/proto/{client,server}/all blank-imports all namespaces
// for convenience.
// Custom type implementations can be registered in the default registry by calling the static
// "RegisterClientType" and "RegisterServerType".
// Users can also construct their own TypeRegistry to separate their own type implementations.
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

func (r *TypeRegistry) register(op Opcode, factory func() Deserializer, dir Direction) error {
	if op == "" {
		return ErrInvalidOpcode
	}
	reg := r.reg[dir]
	_, exists := reg[op]
	if exists {
		return ErrDuplicateOpcode
	}
	reg[op] = factory
	switch dir {
	case ClientToServer:
		if len(op) > r.maxLenCli {
			r.maxLenCli = len(op)
		}
	case ServerToClient:
		if len(op) > r.maxLenSvr {
			r.maxLenSvr = len(op)
		}
	}
	return nil
}

// RegisterServerType associates an opcode to a constructor for the server Deserializer it
// identifies.
func (r *TypeRegistry) RegisterServerType(op Opcode, factory func() Deserializer) error {
	return r.register(op, factory, ServerToClient)
}

// RegisterClientType associates an opcode to a constructor for the client Deserializer it
// identifies.
func (r *TypeRegistry) RegisterClientType(op Opcode, factory func() Deserializer) error {
	return r.register(op, factory, ClientToServer)
}

func (r *TypeRegistry) parse(
	msg []byte, dir Direction) (op Opcode, payload string, err error) {
	if len(msg) == 0 {
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
		if len(msg) < n {
			continue
		}
		_, exists := reg[Opcode(msg[:n])]
		if exists {
			return Opcode(msg[:n]), string(msg[n:]), nil
		}
	}
	return "", "", ErrUnknownOpcode
}

// ParseClientMessage extracts the opcode and the raw payload from a message going from a
// client to a server.
func (r *TypeRegistry) ParseClientMessage(msg []byte) (op Opcode, payload string, err error) {
	return r.parse(msg, ClientToServer)
}

// ParseServerMessage extracts the opcode and the raw payload from a message going from a
// server to a client.
func (r *TypeRegistry) ParseServerMessage(msg []byte) (op Opcode, payload string, err error) {
	return r.parse(msg, ServerToClient)
}

// DeserializeMessage instantiates the message type registered for op and populates it from
// payload. The returned Message is meant to be type-switched on to get the concrete type.
func (r *TypeRegistry) DeserializeMessage(
	op Opcode, payload string, dir Direction) (Message, error) {
	reg := r.reg[dir]

	factory, exists := reg[op]
	if !exists {
		return nil, fmt.Errorf("%w: %q", ErrUnknownOpcode, op)
	}
	msg := factory()

	err := msg.Deserialize(payload)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

var registry = NewTypeRegistry()

// RegisterClientType
func RegisterClientType(op Opcode, factory func() Deserializer) {
	err := registry.RegisterClientType(op, factory)
	if err != nil {
		panic(fmt.Errorf("client: %w: %q", ErrInvalidOpcode, op))
	}
}

// RegisterServerType
func RegisterServerType(op Opcode, factory func() Deserializer) {
	err := registry.RegisterServerType(op, factory)
	if err != nil {
		panic(fmt.Errorf("server: %w: %q", ErrInvalidOpcode, op))
	}
}

// ParseClientMessage extracts the opcode and the payload from a message by using the default
// registry.
func ParseClientMessage(msg []byte) (op Opcode, payload string, err error) {
	return registry.ParseClientMessage(msg)
}

// ParseServerMessage extracts the opcode and the payload from a message by using the default
// registry.
func ParseServerMessage(msg []byte) (op Opcode, payload string, err error) {
	return registry.ParseServerMessage(msg)
}

// DeserializeClientMessage instantiates the message type registered in the default registry
// for op and populates it from payload. The returned Message is meant to be type-switched on
// to get the concrete type.
func DeserializeClientMessage(op Opcode, payload string) (Message, error) {
	return registry.DeserializeMessage(op, payload, ClientToServer)
}

// DeserializeServerMessage instantiates the message type registered in the default registry
// for op and populates it from payload. The returned Message is meant to be type-switched on
// to get the concrete type.
func DeserializeServerMessage(op Opcode, payload string) (Message, error) {
	return registry.DeserializeMessage(op, payload, ServerToClient)
}

// SerializeMessage serializes a message into its pre-processed wire format: the opcode prefix
// followed by the payload string.
func SerializeMessage(msg Serializer) ([]byte, error) {
	body, err := msg.Serialize()
	if err != nil {
		return nil, fmt.Errorf("serialize %q: %w", msg.Opcode(), err)
	}
	return []byte(fmt.Sprintf("%s%s", msg.Opcode(), body)), nil
}
