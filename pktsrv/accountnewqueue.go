package pktsrv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// AccountNewQueue represents a packet containing login queue status information.
type AccountNewQueue struct {
	// Position is the position of the client in the login queue
	Position int
	// NSubs is the number of clients that are subscribers in the login queue
	NSubs int
	// NNonSubs is the number of clients that are not subscribers in the login queue
	NNonSubs int
	// IsSub indicates whether the client is a subscriber
	IsSub bool
	// QueueID is the ID of the queue the client is currently in
	QueueID int
}

func (a *AccountNewQueue) TypeID() string {
	return "Af"
}

func (a *AccountNewQueue) Marshal() ([]byte, error) {
	var sub byte
	if a.IsSub {
		sub = '1'
	} else {
		sub = '0'
	}
	return []byte(fmt.Sprintf(
		"%s%d|%d|%d|%c|%d", a.TypeID(), a.Position, a.NSubs, a.NNonSubs, sub, a.QueueID)), nil
}

func (a *AccountNewQueue) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	tks := strings.Split(pl, "|")
	if len(tks) != 5 {
		return errors.New("invalid data")
	}
	var arr [5]int
	for i, v := range tks {
		var err error
		arr[i], err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	a.Position = arr[0]
	a.NSubs = arr[1]
	a.NNonSubs = arr[2]
	if arr[3] == 1 {
		a.IsSub = true
	} else {
		a.IsSub = false
	}
	a.QueueID = arr[4]
	return nil
}
