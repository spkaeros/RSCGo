package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["chatmsg"] = func(c *Client, p *packets.Packet) {
		//		for _, p1 := range c.player.NearbyPlayers() {
		//			if c1, ok := ClientList.Get(p1.Index).(*Client); c1 != nil && ok {
		//				c1.outgoingPackets <- packets.TeleBubble(diffX, diffY)
		//			}
		//		}

		for _, v := range c.player.LocalPlayers.List {
			if v, ok := v.(*entity.Player); ok {
				if c1, ok := ClientsIdx[v.Index]; ok {
					c1.outgoingPackets <- packets.PlayerChat(c.Index, string(strutil.PackChatMessage(strutil.FormatChatMessage(strutil.UnpackChatMessage(p.Payload)))))
				}
			}
		}
	}
	//	PacketHandlers[84] = func(c *Client, p *packets.Packet) {
	//		index, _ := p.ReadShort()
	//		c.player.Appearances = append(c.player.Appearances, int(index))
	//	}
}
