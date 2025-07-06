package pktsrv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/login"
)

// AccountServersList represents a packet containing the list of servers used by the account.
// Type ID: Ax
// Wire format: AxK[RemainingSub](|[Server0ID],[Server0NChars]...)
type AccountServersList struct {
	//Success bool // always Ok
	RemainingSub int64 // Remaining subscription time in milliseconds
	Servers      []login.ServerWithCharacters
}

func (a *AccountServersList) TypeID() string {
	return "Ax"
}

func (a *AccountServersList) Marshal() ([]byte, error) {
	var srvstr []string
	for _, v := range a.Servers {
		srvstr = append(srvstr, fmt.Sprintf("%d,%d", v.ServerID, v.NCharacter))
	}
	if len(a.Servers) == 0 {
		return []byte(
			fmt.Sprintf("%sK%d", a.TypeID(), a.RemainingSub)), nil
	}
	return []byte(
		fmt.Sprintf("%sK%d|%s", a.TypeID(), a.RemainingSub, strings.Join(srvstr, "|"))), nil
}

func (a *AccountServersList) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	if len(pl) < 2 {
		return errors.New("insufficient data")
	}
	s := strings.Split(pl[1:], "|")
	t, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		return errors.New("invalid data")
	}
	a.RemainingSub = t
	if len(s) == 1 {
		return nil
	}
	for _, v := range s[1:] {
		tks := strings.Split(v, ",")
		if len(tks) != 2 {
			return errors.New("invalid data")
		}
		id, err := strconv.Atoi(tks[0])
		if err != nil {
			return errors.New("invalid data")
		}
		n, err := strconv.Atoi(tks[1])
		if err != nil {
			return errors.New("invalid data")
		}
		a.Servers = append(a.Servers, login.ServerWithCharacters{id, n})
	}
	return nil
}
