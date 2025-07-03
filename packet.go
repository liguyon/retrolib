package retrolib

// Packet is the interface that defines the contract for all packet types in Retro's custom network protocol.
type Packet interface {

	// TypeID returns a string identifier for the packet type.
	TypeID() string

	// Marshal serializes the packet into raw byte data for network transmission.
	// This method handles only the packet's data serialization and should NOT:
	// - Apply packet encryption
	// - Add the EOM (End Of Message) marker
	Marshal() ([]byte, error)

	// Unmarshal deserializes raw byte data into the packet structure.
	// This method handles only the packet's data deserialization and shoud NOT:
	// - Apply packet decryption
	// - Remove the EOM (End Of Message) marker
	Unmarshal([]byte) error
}
