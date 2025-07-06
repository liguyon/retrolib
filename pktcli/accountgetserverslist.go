package pktcli

// AccountGetServersList represents a packet requesting the list of servers the account has characters on.
type AccountGetServersList struct{}

func (a *AccountGetServersList) TypeID() string {
	return "Ax"
}

func (a *AccountGetServersList) Marshal() ([]byte, error) {
	return []byte(a.TypeID()), nil
}

func (a *AccountGetServersList) Unmarshal(bytes []byte) error {
	return nil
}
