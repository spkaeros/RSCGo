package packethandlers

import (
	"github.com/BurntSushi/toml"
	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["pingreq"] = func(player *world.Player, p *packet.Packet) {
//		player.SendPacket(packetbuilders.ResponsePong)
	}
}

//handlerFunc Represents a function for handling incoming packetbuilders.
type handlerFunc func(*world.Player, *packet.Packet)

//PacketHandlers A map with descriptive names for the keys, and functions to run for the value.
var PacketHandlers = make(map[string]handlerFunc)

//packetHandler Definition of a packet handler.
type packetHandler struct {
	Opcode int    `toml:"opcode"`
	Name   string `toml:"name"`
	//	Handler handlerFunc
}

//packetHandlerTable Represents a mapping of descriptive names to packet opcodes.
type packetHandlerTable struct {
	Handlers []packetHandler `toml:"packets"`
}

var table packetHandlerTable

//Get Returns the packet handler function assigned to this opcode.  If it can't be found, returns nil.
func Get(opcode byte) handlerFunc {
	for _, handler := range table.Handlers {
		if byte(handler.Opcode) == opcode {
			return PacketHandlers[handler.Name]
		}
	}
	return nil
}

//Size returns the number of packet handlers currently defined.
func Size() int {
	return len(table.Handlers)
}

//Initialize Deserializes the packet handler table into memory.
func Initialize() {
	if _, err := toml.DecodeFile(config.DataDir()+config.PacketHandlers(), &table); err != nil {
		log.Error.Fatalln("Could not open packet handler table data file:", err)
		return
	}
}
