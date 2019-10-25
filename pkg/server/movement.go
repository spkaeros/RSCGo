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
		c.player.QueueDistancedAction(func() bool {
			if !c.player.IsFollowing() {
				return true
			}
			followingClient, ok := Clients.FromIndex(c.player.FollowIndex())
			if followingClient == nil || !ok || !c.player.Location.WithinRange(followingClient.player.Location, 15) {
				c.player.ResetFollowing()
				return true
			} else if !c.player.FinishedPath() && c.player.WithinRange(followingClient.player.Location, 2) {
				c.player.ResetPath()
			} else if c.player.FinishedPath() && !c.player.WithinRange(followingClient.player.Location, 2) {
				c.player.SetPath(world.NewPathwayFromLocation(&followingClient.player.Location))
			}
			return false
		})
		c.Message("@que@Following " + affectedClient.player.Username)
	}
}
