package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["clientsetting"] = func(player *world.Player, p *packet.Packet) {
		// 2 = mouse buttons
		// 0 = camera angle manual/auto
		// 3 = soundFX (false=on, wtf)
		player.SetClientSetting(int(p.ReadByte()), p.ReadBool())
	}
	PacketHandlers["privacysettings"] = func(player *world.Player, p *packet.Packet) {
		chatBlocked := p.ReadBool()
		friendBlocked := p.ReadBool()
		tradeBlocked := p.ReadBool()
		duelBlocked := p.ReadBool()
		if player.FriendBlocked() && !friendBlocked {
			// turning off private chat block
			players.Range(func(c1 *world.Player) {
				if c1.Friends(player.UsernameHash()) && !player.Friends(c1.UsernameHash()) {
					c1.SendPacket(world.FriendUpdate(player.UsernameHash(), true))
				}
			})
		} else if !player.FriendBlocked() && friendBlocked {
			// turning on private chat block
			players.Range(func(c1 *world.Player) {
				if c1.Friends(player.UsernameHash()) && !player.Friends(c1.UsernameHash()) {
					c1.SendPacket(world.FriendUpdate(player.UsernameHash(), false))
				}
			})
		}
		player.SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked)
	}
}
