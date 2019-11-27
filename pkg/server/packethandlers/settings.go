package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["clientsetting"] = func(c *world.Player, p *packet.Packet) {
		// 2 = mouse buttons
		// 0 = camera angle manual/auto
		// 3 = soundFX (false=on, wtf)
		c.SetClientSetting(int(p.ReadByte()), p.ReadBool())
	}
	PacketHandlers["privacysettings"] = func(c *world.Player, p *packet.Packet) {
		chatBlocked := p.ReadBool()
		friendBlocked := p.ReadBool()
		tradeBlocked := p.ReadBool()
		duelBlocked := p.ReadBool()
		if c.FriendBlocked() && !friendBlocked {
			// turning off private chat block
			players.Range(func(c1 *world.Player) {
				if c1.Friends(c.UserBase37) && !c.Friends(c1.UserBase37) {
					c1.SendPacket(packetbuilders.FriendUpdate(c.UserBase37, true))
				}
			})
		} else if !c.FriendBlocked() && friendBlocked {
			// turning on private chat block
			players.Range(func(c1 *world.Player) {
				if c1.Friends(c.UserBase37) && !c.Friends(c1.UserBase37) {
					c1.SendPacket(packetbuilders.FriendUpdate(c.UserBase37, false))
				}
			})
		}
		c.SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked)
	}
}
