package pktcli

// Version represents a packet containing the client's version.
// Type ID: none
// Wire format: [V]
type Version struct {
	// V is the client's version formatted as MAJOR.MINOR.PATCH
	V string
}

func (v *Version) TypeID() string {
	return ""
}

func (v *Version) Marshal() ([]byte, error) {
	return []byte(v.V), nil
}

func (v *Version) Unmarshal(bytes []byte) error {
	v.V = string(bytes)
	return nil
}
