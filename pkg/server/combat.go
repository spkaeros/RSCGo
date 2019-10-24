package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["attacknpc"] = func(c *Client, p *packets.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("Player[%v] tried to attack nil NPC\n", c)
			return
		}
		c.player.QueueDistancedAction(func() bool {
			if c.player.WithinRange(&npc.Location, 1) {
				c.player.ResetPath()
				npc.ResetPath()
				c.player.Teleport(int(npc.Location.X.Load()), int(npc.Location.Y.Load()))
				c.player.State = world.MSFighting
				npc.State = world.MSFighting
				c.player.SetDirection(world.LeftFighting)
				npc.SetDirection(world.RightFighting)
				return true
			}
			return false
		})
	}
}