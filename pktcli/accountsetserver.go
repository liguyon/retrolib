package pktcli

import (
	"errors"
	"strconv"
	"strings"
)

type AccountSetServer struct {
	ServerID int
}

func (a *AccountSetServer) TypeID() string {
	return "AX"
}

func (a *AccountSetServer) Marshal() ([]byte, error) {
	return []byte(a.TypeID() + strconv.Itoa(a.ServerID)), nil
}

func (a *AccountSetServer) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	var err error
	a.ServerID, err = strconv.Atoi(pl)
	if err != nil {
		return errors.New("invalid data")
	}
	return nil
}
