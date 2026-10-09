package basic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type AveragePing struct {
	proto.ClientSide

	Ping       int
	NSamples   int
	BufferSize int
}

func (a *AveragePing) Opcode() proto.Opcode { return "Bp" }

func (a *AveragePing) Serialize() (string, error) {
	return fmt.Sprintf("%d|%d|%d", a.Ping, a.NSamples, a.BufferSize), nil
}

func (a *AveragePing) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload, "|")
	if len(sli) != 3 {
		return errors.New("invalid field count")
	}
	ping, err := strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[0])
	}
	samples, err := strconv.Atoi(sli[1])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[1])
	}
	bufsize, err := strconv.Atoi(sli[2])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[2])
	}
	a.Ping = ping
	a.NSamples = samples
	a.BufferSize = bufsize
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &AveragePing{} })
}
