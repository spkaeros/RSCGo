package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

func init() {
	PacketHandlers["chatmsg"] = func(c *Client, p *packets.Packet) {
		for _, v := range c.player.LocalPlayers.List {
			if v, ok := v.(*entity.Player); ok && (!v.ChatBlocked() || v.FriendsWith(c.player.UserBase37)) {
				c1 := ClientFromIndex(v.Index)
				if c1 != nil {
					c1.outgoingPackets <- packets.PlayerChat(c.Index, string(p.Payload))
				}
			}
		}
	}
	//	PacketHandlers[84] = func(c *Client, p *packets.Packet) {
	//		index, _ := p.ReadShort()
	//		c.player.Appearances = append(c.player.Appearances, int(index))
	//	}
}
