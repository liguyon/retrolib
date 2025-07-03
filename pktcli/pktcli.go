package pktcli

import "github.com/liguyon/retrolib"

// CreatePacketByID creates a new packet instance based on the provided type ID.
// Returns a pointer to the appropriate packet struct, or nil if the type ID is not recognized.
// This factory function is used after the packet type is extracted from the wire data and the concrete struct
// needs to be instantiated.
func CreatePacketByID(typeID string) retrolib.Packet {
	switch typeID {
	case "Af":
		return &AccountGetQueuePosition{}
	}
	return nil
}
