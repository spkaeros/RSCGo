package server

import (
	"fmt"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"github.com/BurntSushi/toml"
)

//PacketHandlers A map with descriptive names for the keys, and functions to run for the value.
var PacketHandlers = make(map[string]func(*Client, *packets.Packet))

//packethandler Definition of a packet handler.
type packetHandler struct {
	Opcode uint8  `toml:"opcode"`
	Name   string `toml:"name"`
	//	Handle func(c *Client, p *packets.Packet)
}

//packetHandlerTable Represents a mapping of descriptive names to packet opcodes.
type packetHandlerTable struct {
	Handlers []packetHandler `toml:"packets"`
}

var table packetHandlerTable

func (p packetHandlerTable) Get(opcode byte) func(*Client, *packets.Packet) {
	for _, handler := range p.Handlers {
		if handler.Opcode == opcode {
			return PacketHandlers[handler.Name]
		}
	}
	return nil
}

//InitPacketHandlerTable Deserializes the packet handler table into memory.
func InitPacketHandlerTable() {
	if _, err := toml.DecodeFile(TomlConfig.DataDir+TomlConfig.PacketHandlerFile, &table); err != nil {
		LogError.Fatalln("Could not open packet handler table data file:", err)
		return
	}
}

//HandlePacket Finds the mapped handler function for the specified packet, and calls it with the specified parameters.
func (c *Client) HandlePacket(p *packets.Packet) {
	handler := table.Get(p.Opcode)
	if handler == nil {
		LogInfo.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, len(p.Payload))
		fmt.Printf("CONTENT: %v\n", p.Payload)
		return
	}

	handler(c, p)
}
