package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["walkto"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().State == world.MSMenuChoosing {
			c.Player().OptionMenuC <- -1
			c.Player().State = world.MSIdle
		}
		if c.Player().State != world.MSIdle {
			return
		}
		if c.Player().TransAttrs.VarBool("fighting", false) {
			curRound := c.Player().TransAttrs.VarInt("fightRound", 0)
			if curRound < 3 {
				c.Message("You can't retreat during the first 3 rounds of combat")
				return
			}
			if target := c.Player().TransAttrs.VarPlayer("fightTarget"); target != nil {
				target.SendPacket(packetbuilders.Sound("retreat"))
				target.SendPacket(packetbuilders.ServerMessage("Your opponent is retreating"))
			}
			c.Player().ResetFighting()
		}
		startX := p.ReadShort()
		startY := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadSByte()))
			waypointsY = append(waypointsY, int(p.ReadSByte()))
		}
		c.Player().ResetAll()
		c.Player().SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	}
	PacketHandlers["walktoentity"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().TransAttrs.VarBool("fighting", false) {
			c.Message("You can't do that whilst you are fighting.")
			return
		}
		if c.Player().State != world.MSIdle {
			return
		}
		startX := p.ReadShort()
		startY := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadSByte()))
			waypointsY = append(waypointsY, int(p.ReadSByte()))
		}
		c.Player().ResetAll()
		c.Player().SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	}
	PacketHandlers["followreq"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().State != world.MSIdle {
			return
		}
		playerID := p.ReadShort()
		affectedClient, ok := clients.FromIndex(playerID)
		if !ok {
			c.Message("@que@Could not find the player you're looking for.")
			return
		}
		c.Player().ResetAll()
		c.Player().StartFollowing(2)
		c.Message("@que@Following " + affectedClient.Player().Username)
		c.Player().SetDistancedAction(func() bool {
			if !c.Player().IsFollowing() {
				// Following vars have been reset.
				return true
			}
			if affectedClient == nil || !affectedClient.Player().TransAttrs.VarBool("connected", false) ||
				!c.Player().WithinRange(affectedClient.Player().Location, 16) {
				// We think we have a target, but they're miles away now or no longer exist
				c.Player().ResetFollowing()
				return true
			}
			if !c.Player().FinishedPath() && c.Player().WithinRange(affectedClient.Player().Location, c.Player().FollowRadius()) {
				// We're not done moving toward our target, but we're close enough that we should stop
				c.Player().ResetPath()
			} else if !c.Player().WithinRange(affectedClient.Player().Location, c.Player().FollowRadius()) {
				// We're not moving, but our target is moving away, so we must try to get closer
				c.Player().SetPath(world.MakePath(c.Player().Location, affectedClient.Player().Location))
			}
			return false
		})
	}
	PacketHandlers["appearancerequest"] = func(c clients.Client, p *packet.Packet) {
		playerCount := p.ReadShort()
		for i := 0; i < playerCount; i++ {
			serverIndex := p.ReadShort()
			appearanceTicket := p.ReadShort()
			c.Player().AppearanceLock.Lock()
			if ticket, ok := c.Player().KnownAppearances[serverIndex]; !ok || ticket != appearanceTicket {
				if c1, ok := clients.FromIndex(serverIndex); ok {
					c.Player().AppearanceReq = append(c.Player().AppearanceReq, c1.Player())
				}
			}
			c.Player().AppearanceLock.Unlock()
		}
	}
}
