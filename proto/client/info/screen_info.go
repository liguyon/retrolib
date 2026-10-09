package info

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type ScreenInfo struct {
	proto.ClientSide

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
		return errors.New("invalid field count")
	}

	var err error
	s.Width, err = strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[0])
	}
	s.Height, err = strconv.Atoi(sli[1])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[1])
	}
	if len(sli[2]) != 1 {
		return fmt.Errorf("not a char: %q", sli[2])
	}
	if sli[2][0] < '0' || sli[2][0] > '9' {
		return fmt.Errorf("invalid digit: %c", sli[2][0])
	}
	s.State = sli[2][0]
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &ScreenInfo{} })
}
