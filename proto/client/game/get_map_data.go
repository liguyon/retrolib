package game

import (
	"fmt"
	"strconv"

	"github.com/liguyon/retrolib/proto"
)

type GetMapData struct {
	proto.ClientSide

	MapID *int
}

func (g *GetMapData) Opcode() proto.Opcode { return "GD" }

func (g *GetMapData) Serialize() (string, error) {
	if g.MapID != nil {
		return fmt.Sprintf("%d", g.MapID), nil
	}
	return "", nil
}

func (g *GetMapData) Deserialize(payload string) error {
	if payload == "" {
		return nil
	}
	id, err := strconv.Atoi(payload)
	if err != nil {
		return fmt.Errorf("invalid number: %q", payload)
	}
	*g.MapID = id
	return nil
}

func init() {
	proto.Register(func() proto.Deserializer { return &GetMapData{} })
}
