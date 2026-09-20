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

// AppendMessageDelim appends the end-of-message delimiter to an encoded
// (and typically encrypted) message, producing a frame ready to write to the connection.
func AppendMessageDelim(msg []byte, dir Direction) []byte {
	res := make([]byte, 0, len(msg)+len(dir.Delim()))
	res = append(res, msg...)
	res = append(res, dir.Delim()...)
	return res
}

// TrimMessageDelim removes the end-of-message delimiter from a raw frame read off the wire.
func TrimMessageDelim(raw []byte, dir Direction) []byte {
	return bytes.TrimSuffix(raw, dir.Delim())
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

var (
	clientRegistry     = map[Opcode]func() Deserializer{}
	maxClientOpcodeLen int

	serverRegistry     = map[Opcode]func() Deserializer{}
	maxServerOpcodeLen int
)

// RegisterType associates an opcode to a constructor for the Deserializer it identifies.
// It's meant to be called from an init() or before any message parsing or deserialization
// happens.
func RegisterType(op Opcode, factory func() Deserializer, dir Direction) {
	if op == "" {
		panic(fmt.Errorf("%w: empty opcode", ErrInvalidOpcode))
	}

	var dirStr string
	var reg map[Opcode]func() Deserializer
	var maxLen *int

	switch dir {
	case ServerToClient:
		dirStr = "server"
		reg = serverRegistry
		maxLen = &maxServerOpcodeLen
	case ClientToServer:
		dirStr = "client"
		reg = clientRegistry
		maxLen = &maxClientOpcodeLen
	}

	_, exists := reg[op]
	if exists {
		panic(fmt.Errorf("%s: %w: %q", dirStr, ErrDuplicateOpcode, op))
	}

	reg[op] = factory
	if *maxLen < len(op) {
		*maxLen = len(op)
	}
}

// RegisterClientType is the same as calling RegisterType with dir=ClientToServer.
func RegisterClientType(op Opcode, factory func() Deserializer) {
	RegisterType(op, factory, ClientToServer)
}

// RegisterServerType is a the same as calling RegisterType with dir=ServerToClient.
func RegisterServerType(op Opcode, factory func() Deserializer) {
	RegisterType(op, factory, ServerToClient)
}

// ParseMessage extracts the opcode and the raw payload from a message.
func ParseMessage(msg []byte, dir Direction) (op Opcode, payload string, err error) {
	if len(msg) == 0 {
		return "", "", fmt.Errorf("%w: empty message", ErrInvalidOpcode)
	}

	var reg map[Opcode]func() Deserializer
	var maxLen int

	switch dir {
	case ClientToServer:
		reg = clientRegistry
		maxLen = maxClientOpcodeLen
	case ServerToClient:
		reg = serverRegistry
		maxLen = maxServerOpcodeLen
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

// ParseClientMessage is the same as calling ParseMessage with dir=ClientToServer.
func ParseClientMessage(msg []byte) (op Opcode, payload string, err error) {
	return ParseMessage(msg, ClientToServer)
}

// ParseServerMessage is the same as calling ParseMessage with dir=ServerToClient.
func ParseServerMessage(msg []byte) (op Opcode, payload string, err error) {
	return ParseMessage(msg, ServerToClient)
}

// DeserializeMessage instantiates the message type registered for op and populates it from
// payload. The returned Message is meant to be type-switched on to get the concrete type.
func DeserializeMessage(op Opcode, payload string, dir Direction) (Message, error) {
	var reg map[Opcode]func() Deserializer
	switch dir {
	case ServerToClient:
		reg = serverRegistry
	case ClientToServer:
		reg = clientRegistry
	}

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

// DeserializeClientMessage is the same a calling DeserializeMessage with dir=ClientToServer.
func DeserializeClientMessage(op Opcode, payload string) (Message, error) {
	return DeserializeMessage(op, payload, ClientToServer)
}

// DeserializeServerMessage is the same as calling DeserializeMessage with dir=ServerToClient.
func DeserializeServerMessage(op Opcode, payload string) (Message, error) {
	return DeserializeMessage(op, payload, ServerToClient)
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
