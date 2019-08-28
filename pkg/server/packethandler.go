package server

import (
	"fmt"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"github.com/BurntSushi/toml"
)

type packethandler struct {
	Name   string
	Opcode int
}

type packethandlerTable struct {
	handlers []packethandler
}

var packethandlers packethandlerTable

//PacketHandlers A map with descriptive names for the keys, and functions to run for the value.
var PacketHandlers = make(map[string]func(*Client, *packets.Packet))

//LoadPacketHandlerTable Deserializes the packet handler table into memory.
func LoadPacketHandlerTable(file string) {
	if _, err := toml.DecodeFile(TomlConfig.DataDir+file, &packethandlers); err != nil {
		LogError.Fatalln("Could not open packet handler table data file:", err)
		return
	}
}

//HandlerName Returns a descriptive name for a packet given its opcode.
func HandlerName(opcode int) string {
	for _, v := range packethandlers.handlers {
		if v.Opcode == opcode {
			return v.Name
		}
	}
	return ""
}

//HandlePacket Finds the mapped handler function for the specified packet, and calls it with the specified parameters.
func (c *Client) HandlePacket(p *packets.Packet) {
	name := HandlerName(int(p.Opcode))
	if len(name) <= 0 {
		LogInfo.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, len(p.Payload))
		fmt.Printf("CONTENT: %v\n", p.Payload)
		return
	}
	// If the opcode maps to any name at all, it should exist here.
	PacketHandlers[name](c, p)
}
