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
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/game"
)

var dataService = db.DefaultPlayerService

func init() {
	game.AddHandler("forgotpass", func(player *world.Player, p *net.Packet) {
		usernameHash := p.ReadUint64()
		//dataService is a db.PlayerService that all login-related functions should use to access or change player profile data.
		go func() {
			if !dataService.PlayerHasRecoverys(usernameHash) {
				player.Destroy()
				return
			}
			player.SendPacket(net.NewReplyPacket([]byte{1}))
			for _, question := range dataService.PlayerLoadRecoverys(usernameHash) {
				player.SendPacket(net.NewReplyPacket([]byte{byte(len(question))}).AddBytes([]byte(question)))
			}
		}()
	})
}
