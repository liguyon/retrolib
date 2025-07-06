package pktsrv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/login"
)

// AccountSelectServer represents a packet sent in response to a client select a server.
// Type ID: AY
// Wire format: AY[Success:K][Addr][Ticket] or AE[Success:E][ErrID]
type AccountSelectServer struct {
	Success bool
	Addr    string
	Ticket  int
	ErrID   login.SelectServerErrorID
}

func (a *AccountSelectServer) TypeID() string {
	return "AY"
}

func (a *AccountSelectServer) Marshal() ([]byte, error) {
	if a.Success {
		return []byte(fmt.Sprintf("%sK%s;%d", a.TypeID(), a.Addr, a.Ticket)), nil
	}
	return []byte(fmt.Sprintf("%sE%c", a.TypeID(), a.ErrID)), nil
}

func (a *AccountSelectServer) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	if len(pl) < 2 {
		return errors.New("insufficient data")
	}
	if pl[0] == 'E' {
		a.ErrID = login.SelectServerErrorID(pl[1])
		return nil
	}
	if pl[0] != 'K' {
		return errors.New("invalid data")
	}
	a.Success = true
	tks := strings.Split(pl[1:], ";")
	if len(tks) != 2 {
		return errors.New("invalid data")
	}
	a.Addr = tks[0]
	t, err := strconv.Atoi(tks[1])
	if err != nil {
		return errors.New("invalid data")
	}
	a.Ticket = t
	return nil
}
