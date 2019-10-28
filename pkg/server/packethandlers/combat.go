package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["attacknpc"] = func(c clients.Client, p *packetbuilders.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("player[%v] tried to attack nil NPC\n", c)
			return
		}
		c.Player().QueueDistancedAction(func() bool {
			if c.Player().WithinRange(npc.Location, 1) {
				c.Player().ResetPath()
				npc.ResetPath()
				c.Player().Teleport(int(npc.Location.X.Load()), int(npc.Location.Y.Load()))
				c.Player().State = world.MSFighting
				npc.State = world.MSFighting
				c.Player().SetDirection(world.LeftFighting)
				npc.SetDirection(world.RightFighting)
				return true
			}
			return false
		})
	}
}
