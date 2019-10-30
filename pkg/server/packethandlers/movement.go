package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
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
		if c.Player().IsFollowing() {
			c.Player().ResetFollowing()
		}
		c.Player().ResetDistancedAction()
		c.Player().SetPath(world.NewPathwayComplete(uint32(startX), uint32(startY), waypointsX, waypointsY))
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
		if c.Player().IsFollowing() {
			c.Player().ResetFollowing()
		}
		c.Player().ResetDistancedAction()
		c.Player().SetPath(world.NewPathwayComplete(uint32(startX), uint32(startY), waypointsX, waypointsY))
	}
	PacketHandlers["followreq"] = func(c clients.Client, p *packetbuilders.Packet) {
		playerID := p.ReadShort()
		affectedClient, ok := clients.FromIndex(playerID)
		if !ok {
			c.Message("@que@Could not find the player you're looking for.")
			return
		}
		c.Player().SetFollowing(playerID)
		c.Message("@que@Following " + affectedClient.Player().Username)
		c.Player().SetDistancedAction(func() bool {
			if !c.Player().IsFollowing() {
				// Following target no longer exists
				return true
			}
			if affectedClient == nil || !c.Player().Location.WithinRange(affectedClient.Player().Location, 16) {
				// We think we have a target, but they're miles away now or no longer exist
				c.Player().ResetFollowing()
				return true
			}
			if !c.Player().FinishedPath() && c.Player().WithinRange(affectedClient.Player().Location, 2) {
				// We're not done moving toward our target, but we're close enough that we should stop
				c.Player().ResetPath()
			} else if c.Player().FinishedPath() && !c.Player().WithinRange(affectedClient.Player().Location, 2) {
				// We're not moving, but our target is moving away, so we must try to get closer
				c.Player().SetPath(world.NewPathwayFromLocation(affectedClient.Player().Location))
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
