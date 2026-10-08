package account

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/proto"
)

type QueuePosition struct {
	proto.ServerSide

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
		return errors.New("invalid field count")
	}

	var err error

	q.Position, err = strconv.Atoi(sli[0])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[0])
	}

	q.NSubs, err = strconv.Atoi(sli[1])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[1])
	}

	q.NNonSubs, err = strconv.Atoi(sli[2])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[2])
	}

	q.IsSub, err = proto.ParseBool(sli[3])
	if err != nil {
		return fmt.Errorf("%w: %q", err, sli[3])
	}

	q.QueueID, err = strconv.Atoi(sli[4])
	if err != nil {
		return fmt.Errorf("invalid number: %q", sli[4])
	}

	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &QueuePosition{} })
}
