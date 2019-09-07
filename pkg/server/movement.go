package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
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
		c.player.SetPath(entity.NewPathwayComplete(startX, startY, waypointsX, waypointsY))
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
		c.player.SetPath(entity.NewPathwayComplete(startX, startY, waypointsX, waypointsY))
	}
	PacketHandlers["followreq"] = func(c *Client, p *packets.Packet) {
		playerID := p.ReadShort()
		affectedClient := ClientFromIndex(playerID)
		if affectedClient == nil {
			c.outgoingPackets <- packets.ServerMessage("@que@Could not find the player you're looking for.")
			return
		}
		c.player.SetFollowing(playerID)
		c.outgoingPackets <- packets.ServerMessage("@que@Following " + affectedClient.player.Username)
	}
	// 157 is gay.  Walk to other entity with distanced action
	//	Handlers[157] = func(c *Client, p *packets.Packet) {
	//	}
}
