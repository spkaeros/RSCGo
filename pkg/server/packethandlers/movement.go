package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["walkto"] = func(c clients.Client, p *packetbuilders.Packet) {
		startX := p.ReadShort()
		startY := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadSByte()))
			waypointsY = append(waypointsY, int(p.ReadSByte()))
		}
		c.Player().ResetAll()
		c.Player().SetPath(world.NewPathway(uint32(startX), uint32(startY), waypointsX, waypointsY))
	}
	PacketHandlers["walktoentity"] = func(c clients.Client, p *packetbuilders.Packet) {
		startX := p.ReadShort()
		startY := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadSByte()))
			waypointsY = append(waypointsY, int(p.ReadSByte()))
		}
		c.Player().ResetAll()
		c.Player().SetPath(world.NewPathway(uint32(startX), uint32(startY), waypointsX, waypointsY))
	}
	PacketHandlers["followreq"] = func(c clients.Client, p *packetbuilders.Packet) {
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
			if affectedClient == nil || affectedClient.Player().IsFalsy() ||
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
				if dest := c.Player().NextTileToward(affectedClient.Player().Location); !dest.Equals(c.Player().Location) {
					c.Player().SetLocation(dest)
					c.Player().Move()
				} else {
					log.Info.Printf("Could not traverse the world to follow a client: from %v to %v, %v was following %v\n", c.Player().Location, affectedClient.Player().Location, c, affectedClient)
					c.Player().ResetFollowing()
					return true
				}
			}
			return false
		})
	}
	PacketHandlers["appearancerequest"] = func(c clients.Client, p *packetbuilders.Packet) {
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
