/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-20-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

func init() {
	Handlers[186] = func(c *Client, p *packets.Packet) {
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
