package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

func init() {
	PacketHandlers["clientsetting"] = func(c *Client, p *packets.Packet) {
		// 2 = mouse buttons
		// 0 = camera angle manual/auto
		// 3 = soundFX (false=on, wtf)
		c.player.SetClientSetting(int(p.ReadByte()), p.ReadBool())
	}
	PacketHandlers["privacysettings"] = func(c *Client, p *packets.Packet) {
		chatBlocked := p.ReadBool()
		friendBlocked := p.ReadBool()
		tradeBlocked := p.ReadBool()
		duelBlocked := p.ReadBool()
		if c.player.FriendBlocked() && !friendBlocked {
			// turning off private chat block
			Clients.Broadcast(func(c1 *Client) {
				if c1.player.Friends(c.player.UserBase37) && !c.player.Friends(c1.player.UserBase37) {
					c1.outgoingPackets <- packets.FriendUpdate(c.player.UserBase37, true)
				}
			})
		} else if !c.player.FriendBlocked() && friendBlocked {
			// turning on private chat block
			Clients.Broadcast(func(c1 *Client) {
				if c1.player.Friends(c.player.UserBase37) && !c.player.Friends(c1.player.UserBase37) {
					c1.outgoingPackets <- packets.FriendUpdate(c.player.UserBase37, false)
				}
			})
		}
		c.player.SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked)
	}
}
