package pktsrv

import (
	"errors"
	"strings"
)

// HelloConnectionServer represents a login server handshake packet containing an encryption key.
// Type ID: HC
// Wire format: HC[Key]
type HelloConnectionServer struct {
	// Key is a 32-character string used for password encryption
	Key string
}

func (h *HelloConnectionServer) TypeID() string {
	return "HC"
}

func (h *HelloConnectionServer) Marshal() ([]byte, error) {
	if len(h.Key) != 32 {
		return nil, errors.New("invalid key")
	}
	return []byte(h.TypeID() + h.Key), nil
}

func (h *HelloConnectionServer) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), h.TypeID())
	h.Key = pl
	if len(h.Key) != 32 {
		return errors.New("invalid key")
	}
	return nil
}
