package proto

import (
	"errors"
)

var (
	// ErrInvalidOpcode is returned when an opcode is structurally invalid.
	ErrInvalidOpcode = errors.New("invalid opcode")

	// ErrUnknownOpcode is returned when no message type is registered for an opcode.
	ErrUnknownOpcode = errors.New("unknown opcode")

	// ErrDuplicateOpcode is returned when an opcode is registered twice.
	ErrDuplicateOpcode = errors.New("opcode already registered")

	// ErrMissingPayload is returned when a message expected to carry a payload is
	// empty (has an opcode but no payload).
	ErrMissingPayload = errors.New("missing payload")

	// ErrMalformedPayload is returned when a message payload is structurally invalid.
	// It's an umbrella for payloads that cannot be deserialized correctly.
	ErrMalformedPayload = errors.New("malformed payload")

	// ErrNoDelimiter is returned when the end-of-message delimiter is missing from a
	// raw frame.
	ErrNoDelimiter = errors.New("no frame delimiter")

	// ErrDelimiterInBody is returned when trying to append the delimiter to a frame
	// that already has it.
	ErrDelimiterInBody = errors.New("delimiter already in body")

	// ErrInvalidMessage is returned when message fields cannot be serialized.
	ErrInvalidMessage = errors.New("invalid message")
)
