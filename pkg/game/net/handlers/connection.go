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
	//	"github.com/spkaeros/rscgo/pkg/engine/tasks"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("logoutreq", func(player *world.Player, p *net.Packet) {
		//		tasks.Tickers.Add("playerDestroy", func() bool {
		if player.Busy() {
			player.SendPacket(world.CannotLogout)
			return
		}
		if player.Connected() {
			//				player.SendPacket(world.Logout)
			player.Destroy()
		}
		//			return true
		//		})
	})
	game.AddHandler("closeconn", func(player *world.Player, p *net.Packet) {
		if player.Busy() {
			log.Suspicious.Println("CLOSECONN!!", player, p.String(), player.State())
			log.Info.Println("CLOSECONN!!", player, p.String(), player.State())
			player.SendPacket(world.CannotLogout)
			return
		}
		if player.Connected() {
			player.Destroy()
		}
	})
}
