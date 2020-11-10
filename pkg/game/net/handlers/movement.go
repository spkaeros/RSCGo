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

	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
)

func init() {
	game.AddHandler("followreq", func(player *world.Player, p *net.Packet) {
		if player.IsFighting() {
			return
		}
		if !player.CanWalk() {
			return
		}
		playerID := p.ReadUint16()
		target, ok := world.Players.FindIndex(playerID)
		if !ok {
			player.Message("@que@Could not find the player you're looking for.")
			return
		}
		player.ResetAll()
		player.Message("@que@Following " + target.Username())
		rad := 2
		player.SetVar("following", true)
		player.SetTickAction(func() bool {
			// if !player.VarBool("following", false) {
				// Following vars have been reset.
				// return true
			// }
			if target == nil || !target.Connected() ||
				!player.Near(target, player.ViewRadius()) {
				// We think we have a target, but they're miles away now or no longer exist
				player.UnsetVar("following")
				return true
			}
			if !player.FinishedPath() && player.WithinRange(target.Point(), rad) {
				// We're not done moving toward our target, but we're close enough that we should stop
				player.ResetPath()
			} else if !player.WithinRange(target.Point(), rad) {
				// We're not moving, but our target is moving away, so we must try to get closer
				if !player.WalkTo(target.Point()) {
					return true
				}
			}
			return !player.VarBool("following", false)
		})
	})
}
