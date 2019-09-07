package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
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
	PacketHandlers["addfriend"] = func(c *Client, p *packets.Packet) {
		hash := p.ReadLong()
		if hash <= strutil.MaxBase37 && hash > 0 && !c.player.FriendsWith(hash) {
			c.player.FriendList = append(c.player.FriendList, hash)
		}
	}
	PacketHandlers["removefriend"] = func(c *Client, p *packets.Packet) {
		hash := p.ReadLong()
		if hash <= strutil.MaxBase37 && hash > 0 && c.player.FriendsWith(hash) {
			for i, v := range c.player.FriendList {
				if v == hash {
					newSize := len(c.player.FriendList) - 1
					c.player.FriendList[i] = c.player.FriendList[newSize]
					c.player.FriendList = c.player.FriendList[:newSize]
					return
				}
			}
		}
	}
	//	PacketHandlers[84] = func(c *Client, p *packets.Packet) {
	//		index, _ := p.ReadShort()
	//		c.player.Appearances = append(c.player.Appearances, int(index))
	//	}
}
