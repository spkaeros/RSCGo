package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["attacknpc"] = func(c clients.Client, p *packetbuilders.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("player[%v] tried to attack nil NPC\n", c)
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().WithinRange(npc.Location, 1) {
				c.Player().ResetPath()
				npc.ResetPath()
				c.Player().Teleport(npc.CurX(), npc.CurY())
				c.Player().State = world.MSFighting
				npc.State = world.MSFighting
				c.Player().SetDirection(world.LeftFighting)
				npc.SetDirection(world.RightFighting)
				c.Player().TransAttrs.SetVar("fighting", true)
				c.Player().TransAttrs.SetVar("fightTarget", npc)
				npc.TransAttrs.SetVar("fighting", true)
				npc.TransAttrs.SetVar("fightTarget", c.Player())
				return true
			} else {
				c.Player().SetPath(world.MakePath(c.Player().Location, npc.Location))
			}
			return false
		})
	}
}
