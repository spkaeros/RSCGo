package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["walkto"] = func(c *Client, p *packets.Packet) {
		startX := p.ReadShort()
		startY := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadSByte()))
			waypointsY = append(waypointsY, int(p.ReadSByte()))
		}
		if c.player.IsFollowing() {
			c.player.ResetFollowing()
		}
		c.player.SetPath(world.NewPathwayComplete(uint32(startX), uint32(startY), waypointsX, waypointsY))
	}
	PacketHandlers["walktoentity"] = func(c *Client, p *packets.Packet) {
		startX := p.ReadShort()
		startY := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadSByte()))
			waypointsY = append(waypointsY, int(p.ReadSByte()))
		}
		if c.player.IsFollowing() {
			c.player.ResetFollowing()
		}
		c.player.SetPath(world.NewPathwayComplete(uint32(startX), uint32(startY), waypointsX, waypointsY))
	}
	PacketHandlers["followreq"] = func(c *Client, p *packets.Packet) {
		playerID := p.ReadShort()
		affectedClient, ok := Clients.FromIndex(playerID)
		if !ok {
			c.Message("@que@Could not find the player you're looking for.")
			return
		}
		c.player.SetFollowing(playerID)
		c.Message("@que@Following " + affectedClient.player.Username)
		c.player.QueueDistancedAction(func() bool {
			if !c.player.IsFollowing() {
				// Following target no longer exists
				return true
			}
			if affectedClient == nil || !c.player.Location.WithinRange(affectedClient.player.Location, 16) {
				// We think we have a target, but they're miles away now or no longer exist
				c.player.ResetFollowing()
				return true
			}
			if !c.player.FinishedPath() && c.player.WithinRange(affectedClient.player.Location, 2) {
				// We're not done moving toward our target, but we're close enough that we should stop
				c.player.ResetPath()
			} else if c.player.FinishedPath() && !c.player.WithinRange(affectedClient.player.Location, 2) {
				// We're not moving, but our target is moving away, so we must try to get closer
				c.player.SetPath(world.NewPathwayFromLocation(&affectedClient.player.Location))
			}
			return false
		})
	}
	PacketHandlers["appearancerequest"] = func(c *Client, p *packets.Packet) {
		playerCount := p.ReadShort()
		for i := 0; i < playerCount; i++ {
			serverIndex := p.ReadShort()
			appearanceTicket := p.ReadShort()
			if ticket, ok := c.player.KnownAppearances[serverIndex]; !ok || ticket != appearanceTicket {
				if c1, ok := Clients.FromIndex(serverIndex); ok {
					c.player.AppearanceReqLock.Lock()
					c.player.AppearanceReq = append(c.player.AppearanceReq, c1.player)
					c.player.AppearanceReqLock.Unlock()
				}
			}
		}
	}
}
