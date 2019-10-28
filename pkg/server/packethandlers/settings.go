package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
)

func init() {
	PacketHandlers["clientsetting"] = func(c clients.Client, p *packetbuilders.Packet) {
		// 2 = mouse buttons
		// 0 = camera angle manual/auto
		// 3 = soundFX (false=on, wtf)
		c.Player().SetClientSetting(int(p.ReadByte()), p.ReadBool())
	}
	PacketHandlers["privacysettings"] = func(c clients.Client, p *packetbuilders.Packet) {
		chatBlocked := p.ReadBool()
		friendBlocked := p.ReadBool()
		tradeBlocked := p.ReadBool()
		duelBlocked := p.ReadBool()
		if c.Player().FriendBlocked() && !friendBlocked {
			// turning off private chat block
			clients.Range(func(c1 clients.Client) {
				if c1.Player().Friends(c.Player().UserBase37) && !c.Player().Friends(c1.Player().UserBase37) {
					c1.SendPacket(packetbuilders.FriendUpdate(c.Player().UserBase37, true))
				}
			})
		} else if !c.Player().FriendBlocked() && friendBlocked {
			// turning on private chat block
			clients.Range(func(c1 clients.Client) {
				if c1.Player().Friends(c.Player().UserBase37) && !c.Player().Friends(c1.Player().UserBase37) {
					c1.SendPacket(packetbuilders.FriendUpdate(c.Player().UserBase37, false))
				}
			})
		}
		c.Player().SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked)
	}
}
