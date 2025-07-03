package pktsrv

import (
	"errors"
	"strings"
)

// AccountNickname represents a packet containing the account's nickname(pseudo).
// Type ID: Ad
// Wire format: Ad[Nickname]
type AccountNickname struct {
	// Nickname is the account's nickname(pseudo)
	Nickname string
}

func (a *AccountNickname) TypeID() string {
	return "Ad"
}

func (a *AccountNickname) Marshal() ([]byte, error) {
	if len(a.Nickname) == 0 {
		return nil, errors.New("insufficient data")
	}
	return []byte(a.TypeID() + a.Nickname), nil
}

func (a *AccountNickname) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	if len(pl) == 0 {
		return errors.New("insufficient data")
	}
	a.Nickname = pl
	return nil
}
