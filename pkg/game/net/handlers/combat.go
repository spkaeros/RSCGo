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
	"time"

	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("attacknpc", func(player *world.Player, p *net.Packet) {
		npc := world.AsNpc( world.Npcs.Get(p.ReadUint16()) )
		if npc == nil || !npc.Attackable() {
			log.Suspicious.Printf("%v tried to attack nil NPC\n", player)
			player.Message("The character does not appear interested in fighting")
			player.ResetPath()
			return
		}
		if player.IsFighting() {
			player.Message("You're already fighting!")
			player.ResetPath()
			return
		}
		if npc.IsFighting() {
			player.Message("Your opponent is busy!")
			player.ResetPath()
			return
		}
		if player.Busy() {
			return
		}
		player.WalkingArrivalAction(npc, 1, func() {
			if player.IsFighting() {
				player.Message("You're already fighting!")
				return
			}
			if npc.IsFighting() {
				player.Message("Your opponent is busy!")
				return
			}
			if player.Busy() {
				return
			}
			player.ResetPath()
			if time.Since(npc.VarTime("lastFight")) <= time.Second*2 || npc.Busy() {
				return
			}
			npc.ResetPath()
			for _, trigger := range world.NpcAtkTriggers {
				if trigger.Check(player, npc) {
					trigger.Action(player, npc)
					return
				}
			}
			player.StartCombat(npc)
		})
	})
	game.AddHandler("attackplayer", func(player *world.Player, p *net.Packet) {
		affectedPlayer, ok := world.Players.FindIndex(p.ReadUint16())
		if affectedPlayer == nil || !ok {
			log.Suspicious.Printf("player[%v] tried to attack nil player\n", player)
			return
		}
		if player.IsFighting() {
			player.Message("You're already fighting!")
			return
		}
		if affectedPlayer.IsFighting() {
			player.Message("Your opponent is busy!")
			return
		}
		if player.Busy() {
			return
		}
		player.WalkingArrivalAction(affectedPlayer, 2, func() {
			if player.IsFighting() {
				player.Message("You're already fighting!")
				return
			}
			if affectedPlayer.IsFighting() {
				player.Message("Your opponent is busy!")
				return
			}
			if player.Busy() || !player.CanAttack(affectedPlayer) {
				return
			}
			player.ResetPath()
			if time.Since(affectedPlayer.VarTime("lastRetreat")) <= time.Second*3 {
				return
			}
			affectedPlayer.ResetPath()
			affectedPlayer.Message("You are under attack!")
			player.StartCombat(affectedPlayer)
		})
	})
	game.AddHandler("fightmode", func(player *world.Player, p *net.Packet) {
		mode := p.ReadUint8()
		if mode < 0 || mode > 3 {
			log.Suspicious.Printf("Invalid fightmode(%v) selected by %s", mode, player.String())
			return
		}
		player.SetFightMode(int(mode))
	})
}
