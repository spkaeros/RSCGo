package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

func init() {
	PacketHandlers["walkto"] = func(c *Client, p *packets.Packet) {
		startX, _ := p.ReadShort()
		startY, _ := p.ReadShort()
		numWaypoints := (len(p.Payload) - 4) / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			nextX, _ := p.ReadSByte()
			nextY, _ := p.ReadSByte()
			waypointsX = append(waypointsX, int(nextX))
			waypointsY = append(waypointsY, int(nextY))
		}
		c.player.Path = &entity.Pathway{StartX: int(startX), StartY: int(startY), WaypointsX: waypointsX, WaypointsY: waypointsY, CurrentWaypoint: -1}
	}
	// 157 is gay.  Walk to other entity with distanced action
	//	Handlers[157] = func(c *Client, p *packets.Packet) {
	//	}
}
