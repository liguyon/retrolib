package pktsrv

import (
	"strconv"
	"strings"
)

// AccountCommunity represents a packet containing the identifier for the community the account is part of.
// Type ID: Ac
// Wire format: Ac[CommunityID]
type AccountCommunity struct {
	// CommunityID is the identifier for the account's community
	CommunityID int
}

func (a *AccountCommunity) TypeID() string {
	return "Ac"
}

func (a *AccountCommunity) Marshal() ([]byte, error) {
	return []byte(a.TypeID() + strconv.Itoa(a.CommunityID)), nil
}

func (a *AccountCommunity) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), a.TypeID())
	var err error
	a.CommunityID, err = strconv.Atoi(pl)
	return err
}
