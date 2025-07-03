package retrolib

import (
	"errors"
	"fmt"
)

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

// PacketHandler defines the contract for handling incoming packets in the network protocol.
type PacketHandler interface {
	// HandlePacket receives a deserialized packet and should:
	// - Type-assert the packet to determine its concrete type
	// - Execute the appropriate handler logic for that packet type
	// - Return an error if processing fails
	HandlePacket(pkt Packet) error
}

// InboundProcessor handles the conversion from raw wire data to deserialized packet instances.
type InboundProcessor interface {
	// ProcessWireData converts raw wire data into a deserialized packet. It should handle:
	// - EOM (End of Message) marker removal
	// - Data decryption if applicable
	// - Packet instantiation
	// - Deserialization via the Packet's Unmarshal method
	ProcessWireData(wireData []byte) (Packet, error)
}

// PacketDispatcher manages the routing of packets to their appropriate handlers based on packet type.
// It handles wire data processing and packet dispatching.
type PacketDispatcher struct {
	reg  map[string]PacketHandler // registry of packet type IDs mapped to their corresponding PacketHandler impl
	proc InboundProcessor         // inbound wire data processor
}

// NewPacketDispatcher returns a new PacketDispatcher with an initialized but empty handler registry.
func NewPacketDispatcher() *PacketDispatcher {
	return &PacketDispatcher{
		reg: make(map[string]PacketHandler),
	}
}

// RegisterHandler maps a packet type ID to a specific PacketHandler implementation.
// If a handler is already registered for the given typeID, it will be replaced with the new handler.
func (d *PacketDispatcher) RegisterHandler(typeID string, handler PacketHandler) {
	d.reg[typeID] = handler
}

// Process converts raw wire data to a packet and dispatches it to the appropriate handler.
func (d *PacketDispatcher) Process(wireData []byte) error {
	pkt, err := d.proc.ProcessWireData(wireData)
	if err != nil {
		return err
	}
	// TODO: handle packets with no typeID: version, credentials
	handler, ok := d.reg[pkt.TypeID()]
	if !ok {
		return errors.New(fmt.Sprintf("no registered handler for packet type: %s", pkt.TypeID()))
	}
	return handler.HandlePacket(pkt)
}
