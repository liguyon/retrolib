package info

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type ScreenInfo struct {
	Width  int
	Height int
	State  byte
}

func (s *ScreenInfo) Opcode() proto.Opcode { return "Ir" }

func (s *ScreenInfo) Serialize() (string, error) {
	return fmt.Sprintf("%d;%d;%c", s.Width, s.Height, s.State), nil
}

func (s *ScreenInfo) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload, ";")
	if len(sli) != 3 {
		return proto.ErrMalformedPayload
	}

	var err error
	s.Width, err = strconv.Atoi(sli[0])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	s.Height, err = strconv.Atoi(sli[1])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	if len(sli[2]) != 1 {
		return proto.ErrMalformedPayload
	}
	if sli[2][0] < '0' || sli[2][0] > '9' {
		return proto.ErrMalformedPayload
	}
	s.State = sli[2][0]
	return nil
}

func init() {
	proto.RegisterClientType("Ir",
		func() proto.Deserializer { return &ScreenInfo{} })
}
