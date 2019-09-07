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
		c.player.SetPrivacySettings(p.ReadBool(), p.ReadBool(), p.ReadBool(), p.ReadBool())
	}
}
