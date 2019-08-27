/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-20-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-27-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	Handlers[174] = func(c *Client, p *packets.Packet) {
		//		for _, p1 := range c.player.NearbyPlayers() {
		//			if c1, ok := ClientList.Get(p1.Index).(*Client); c1 != nil && ok {
		//				c1.outgoingPackets <- packets.TeleBubble(diffX, diffY)
		//			}
		//		}
		for _, v := range c.player.LocalPlayers.List {
			v, ok := v.(*entity.Player)
			if ok {
				c1, ok := ClientList.Get(v.Index).(*Client)
				if ok {
					c1.outgoingPackets <- packets.PlayerChat(c.index, string(strutil.PackChatMessage(strutil.FormatChatMessage(strutil.UnpackChatMessage(p.Payload)))))
				}
			}
		}
	}
	Handlers[84] = func(c *Client, p *packets.Packet) {
		index, _ := p.ReadShort()
		c.player.Appearances = append(c.player.Appearances, int(index))
	}
}
