package job

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type Options struct {
	JobIndex int
	Options  int
	MinSlots int
}

func (o *Options) Opcode() proto.Opcode { return "JO" }

func (o *Options) Serialize() (string, error) {
	return fmt.Sprintf("%d|%d|%d", o.JobIndex, o.Options, o.MinSlots), nil
}

func (o *Options) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}
	sli := strings.Split(payload, "|")
	if len(sli) != 3 {
		return proto.ErrMalformedPayload
	}
	var err error
	o.JobIndex, err = strconv.Atoi(sli[0])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	o.Options, err = strconv.Atoi(sli[1])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	o.MinSlots, err = strconv.Atoi(sli[2])
	if err != nil {
		return proto.ErrMalformedPayload
	}
	return nil
}

func init() {
	proto.RegisterServerType("JO",
		func() proto.Deserializer { return &Options{} })
}
