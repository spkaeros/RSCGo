package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

var ObjectCommands = make(map[string]func(p *entity.Player, o *entity.Object))

func init() {
	ObjectCommands["open"] = func(p *entity.Player, o *entity.Object) {
		LogInfo.Printf("Object with command 'open' interacted with at %v...\n", o.Location())
	}
	PacketHandlers["objectaction"] = func(c *Client, p *packets.Packet) {
		targetObject := c.player.LocalObjects.GetObject(p.ReadShort(), p.ReadShort())
		if targetObject != nil {
			if handler, ok := ObjectCommands[ObjectDefinitions[targetObject.ID].Commands[0]]; ok {
				handler(c.player, targetObject)
			}
		}
	}
}
