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

	// ErrMissingPayload is returned when a message payload expected to carry data is empty.
	ErrMissingPayload = errors.New("missing payload")

	// ErrMalformedPayload is returned when a message payload is malformed.
	ErrMalformedPayload = errors.New("malformed payload")
)
