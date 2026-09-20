package account

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type QueuePosition struct {
	Position int
	NSubs    int
	NNonSubs int
	IsSub    bool
	QueueID  int
}

func (q *QueuePosition) Opcode() proto.Opcode { return "Af" }

func (q *QueuePosition) Serialize() (string, error) {
	return fmt.Sprintf("%d|%d|%d|%d|%d", q.Position, q.NSubs,
		q.NNonSubs, proto.BoolToInt(q.IsSub), q.QueueID), nil
}

func (q *QueuePosition) Deserialize(payload string) error {
	if payload == "" {
		return proto.ErrMissingPayload
	}

	sli := strings.Split(payload, "|")
	if len(sli) != 5 {
		return proto.ErrMalformedPayload
	}

	var err error

	q.Position, err = strconv.Atoi(sli[0])
	if err != nil {
		return proto.ErrMalformedPayload
	}

	q.NSubs, err = strconv.Atoi(sli[1])
	if err != nil {
		return proto.ErrMalformedPayload
	}

	q.NNonSubs, err = strconv.Atoi(sli[2])
	if err != nil {
		return proto.ErrMalformedPayload
	}

	q.IsSub, err = proto.ParseBool(sli[3])
	if err != nil {
		return proto.ErrMalformedPayload
	}

	q.QueueID, err = strconv.Atoi(sli[4])
	if err != nil {
		return proto.ErrMalformedPayload
	}

	return nil
}

func init() {
	proto.RegisterServerType("Af",
		func() proto.Deserializer { return &QueuePosition{} })
}
