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
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("invaction1", func(player *world.Player, p *net.Packet) {
		index := p.ReadUint16()
		item := player.Inventory.Get(index)
		if item == nil || player.Busy() || player.IsFighting() {
			return
		}
		player.AddState(world.MSItemAction)
		go func() {
			defer func() {
				player.RemoveState(world.MSItemAction)
			}()
			for _, triggerDef := range world.ItemTriggers {
				if triggerDef.Check(item) {
					triggerDef.Action(player, item)
					return
				}
			}
			player.SendPacket(world.DefaultActionMessage)
		}()
	})
}
