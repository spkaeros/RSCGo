/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package handlers

import (
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
)

func init() {
	AddHandler("clientsetting", func(player *world.Player, p *net.Packet) {
		// 2 = mouse buttons
		// 0 = camera angle manual/auto
		// 3 = soundFX (false=on, wtf)
		player.SetClientSetting(int(p.ReadUint8()), p.ReadBoolean())
	})
	AddHandler("privacysettings", func(player *world.Player, p *net.Packet) {
		chatBlocked := p.ReadBoolean()
		friendBlocked := p.ReadBoolean()
		tradeBlocked := p.ReadBoolean()
		duelBlocked := p.ReadBoolean()
		if player.FriendBlocked() && !friendBlocked {
			// turning off private chat block
			world.Players.Range(func(c1 *world.Player) {
				if c1.FriendsWith(player.UsernameHash()) && !player.FriendsWith(c1.UsernameHash()) {
					c1.SendPacket(world.FriendUpdate(player.UsernameHash(), true))
				}
			})
		} else if !player.FriendBlocked() && friendBlocked {
			// turning on private chat block
			world.Players.Range(func(c1 *world.Player) {
				if c1.FriendsWith(player.UsernameHash()) && !player.FriendsWith(c1.UsernameHash()) {
					c1.SendPacket(world.FriendUpdate(player.UsernameHash(), false))
				}
			})
		}
		player.SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked)
	})
}
