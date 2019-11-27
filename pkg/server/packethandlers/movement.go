package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"time"
)

func init() {
	PacketHandlers["walkto"] = func(c *world.Player, p *packet.Packet) {
		if c.HasState(world.MSMenuChoosing) {
			c.OptionMenuC <- -1
			c.RemoveState(world.MSMenuChoosing)
		}
		if c.IsFighting() {
			target := c.FightTarget()
			if target == nil {
				c.ResetFighting()
				return
			}
			curRound := target.FightRound()
			if curRound < 3 {
				c.SendPacket(packetbuilders.ServerMessage("You can't retreat during the first 3 rounds of combat"))
				return
			}
			if target, ok := target.(*world.Player); ok {
				target.SendPacket(packetbuilders.Sound("retreat"))
				target.SendPacket(packetbuilders.ServerMessage("Your opponent is retreating"))
			}
			c.TransAttrs.SetVar("lastRetreat", time.Now())
			c.UpdateLastRetreat()
			c.ResetFighting()
		}
		if c.Busy() {
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
		c.ResetAll()
		c.SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	}
	PacketHandlers["walktoentity"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
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
		c.ResetAll()
		c.SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	}
	PacketHandlers["followreq"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
			return
		}
		playerID := p.ReadShort()
		affectedClient, ok := players.FromIndex(playerID)
		if !ok {
			c.SendPacket(packetbuilders.ServerMessage("@que@Could not find the player you're looking for."))
			return
		}
		c.ResetAll()
		c.StartFollowing(2)
		c.SendPacket(packetbuilders.ServerMessage("@que@Following " + affectedClient.Username))
		c.SetDistancedAction(func() bool {
			if !c.IsFollowing() {
				// Following vars have been reset.
				return true
			}
			if affectedClient == nil || !affectedClient.Connected() ||
				!c.WithinRange(affectedClient.Location, 16) {
				// We think we have a target, but they're miles away now or no longer exist
				c.ResetFollowing()
				return true
			}
			if !c.FinishedPath() && c.WithinRange(affectedClient.Location, c.FollowRadius()) {
				// We're not done moving toward our target, but we're close enough that we should stop
				c.ResetPath()
			} else if !c.WithinRange(affectedClient.Location, c.FollowRadius()) {
				// We're not moving, but our target is moving away, so we must try to get closer
				c.SetPath(world.MakePath(c.Location, affectedClient.Location))
			}
			return false
		})
	}
	PacketHandlers["appearancerequest"] = func(c *world.Player, p *packet.Packet) {
		playerCount := p.ReadShort()
		for i := 0; i < playerCount; i++ {
			serverIndex := p.ReadShort()
			appearanceTicket := p.ReadShort()
			c.AppearanceLock.Lock()
			if ticket, ok := c.KnownAppearances[serverIndex]; !ok || ticket != appearanceTicket {
				if c1, ok := players.FromIndex(serverIndex); ok {
					c.AppearanceReq = append(c.AppearanceReq, c1)
				}
			}
			c.AppearanceLock.Unlock()
		}
	}
}
