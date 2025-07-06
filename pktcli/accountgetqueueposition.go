package pktcli

// AccountGetQueuePosition represents a packet requesting the client's current position in the login queue.
// Type ID: Af
// Wire format: Af
type AccountGetQueuePosition struct {
}

func (a *AccountGetQueuePosition) TypeID() string {
	return "Af"
}

func (a *AccountGetQueuePosition) Marshal() ([]byte, error) {
	return []byte("Af"), nil
}

func (a *AccountGetQueuePosition) Unmarshal(bytes []byte) error {
	return nil
}
