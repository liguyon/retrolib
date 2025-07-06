package pktsrv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liguyon/retrolib/auth"
)

// AccountHosts represents a packet containing info about all game servers.
// Type ID: AH
// Wire format: AH[Srv0ID][Srv0State][Srv0Completion][Srv0CanLogIn]|[Srv1ID][Srv1State][Srv1Completion][Srv1CanLogIn]...
type AccountHosts struct {
	Servers []auth.Server
}

func (a *AccountHosts) TypeID() string {
	return "AH"
}

func (a *AccountHosts) Marshal() ([]byte, error) {
	var sb strings.Builder
	sb.WriteString(a.TypeID())
	for i, v := range a.Servers {
		var cl int
		if v.CanLogIn {
			cl = 1
		}
		sb.WriteString(fmt.Sprintf("%d;%d;%d;%d", v.ID, v.State, v.Completion, cl))
		if i < len(a.Servers)-1 {
			sb.WriteByte('|')
		}
	}
	return []byte(sb.String()), nil
}

func (a *AccountHosts) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	srvs := strings.Split(pl, "|")
	if len(srvs) == 0 {
		return errors.New("insufficient data")
	}
	for _, s := range srvs {
		tks := strings.Split(s, ";")
		if len(tks) != 4 {
			return errors.New("invalid data")
		}
		var arr [4]int
		var err error
		for i, t := range tks {
			arr[i], err = strconv.Atoi(t)
			if err != nil {
				return errors.New("invalid data")
			}
		}
		var srv auth.Server
		if arr[3] == 1 {
			srv = auth.Server{ID: arr[0], State: auth.ServerState(arr[1]), Completion: arr[2], CanLogIn: true}
		} else {
			srv = auth.Server{ID: arr[0], State: auth.ServerState(arr[1]), Completion: arr[2], CanLogIn: false}

		}
		a.Servers = append(a.Servers, srv)
	}
	return nil
}
