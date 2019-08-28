package server

import (
	"fmt"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"github.com/BurntSushi/toml"
)

//PacketHandlers A map with descriptive names for the keys, and functions to run for the value.
var PacketHandlers = make(map[string]func(*Client, *packets.Packet))

//packetHandlerTable Represents a mapping of descriptive names to packet opcodes.
var packetHandlerTable struct {
	Handlers []struct {
		Name   string `toml:"name"`
		Opcode int    `toml:"opcode"`
	} `toml:"packets"`
}

//InitPacketHandlerTable Deserializes the packet handler table into memory.
func InitPacketHandlerTable() {
	if _, err := toml.DecodeFile(TomlConfig.DataDir+TomlConfig.PacketHandlerFile, &packetHandlerTable); err != nil {
		LogError.Fatalln("Could not open packet handler table data file:", err)
		return
	}
}

//handlerFromOpcode Returns a descriptive name for a packet given its opcode.
func handlerFromOpcode(opcode int) func(*Client, *packets.Packet) {
	for _, v := range packetHandlerTable.Handlers {
		if v.Opcode == opcode {
			return PacketHandlers[v.Name]
		}
	}

	return nil
}

//HandlePacket Finds the mapped handler function for the specified packet, and calls it with the specified parameters.
func (c *Client) HandlePacket(p *packets.Packet) {
	handler := handlerFromOpcode(int(p.Opcode))
	if handler == nil {
		LogInfo.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, len(p.Payload))
		fmt.Printf("CONTENT: %v\n", p.Payload)
		return
	}

	handler(c, p)
}
