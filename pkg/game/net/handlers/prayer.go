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
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	//requiredLevels contains prayer level requirements for each prayer in order from prayer 0 to prayer 13
	requiredLevels := []int{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34, 37, 40}
	game.AddHandler("prayeron", func(player *world.Player, p *net.Packet) {
		idx := p.ReadUint8()
		if idx < 0 || idx > 13 {
			log.Suspicious.Printf("%v turned on a prayer that doesn't exist: %d\n", player, idx)
			return
		}
		if requiredLevels[idx] > player.Skills().Maximum(entity.StatPrayer) {
			log.Suspicious.Printf("%v turned on a prayer that he is too low level for: %d\n", player, idx)
			return
		}
		player.ActivatePrayer(int(idx))
		player.PrayerOn(int(idx))
		player.SendPrayers()
	})
	game.AddHandler("prayeroff", func(player *world.Player, p *net.Packet) {
		idx := p.ReadUint8()
		if idx < 0 || idx > 13 {
			log.Suspicious.Printf("%v turned off a prayer that doesn't exist: %d\n", player, idx)
			return
		}
		if requiredLevels[idx] > player.Skills().Maximum(entity.StatPrayer) {
			log.Suspicious.Printf("%v turned off a prayer that he is too low level for: %d\n", player, idx)
			return
		}
		if player.PrayerActivated(int(idx)) {
			player.DeactivatePrayer(int(idx))
		}
		player.SendPrayers()
	})
}
