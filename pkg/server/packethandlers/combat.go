/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"time"
)

func init() {
	PacketHandlers["attacknpc"] = func(player *world.Player, p *packet.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("player[%v] tried to attack nil NPC\n", player)
			return
		}
		if player.Busy() {
			return
		}
		if !world.NpcDefs[npc.ID].Attackable {
			log.Info.Println("Player attacked not attackable NPC!", world.NpcDefs[npc.ID])
			return
		}
		log.Info.Println(npc.ID)
		player.SetDistancedAction(func() bool {
			if player.NextTo(npc.Location) && player.WithinRange(npc.Location, 1) {
				for _, trigger := range script.NpcAtkTriggers {
					if trigger.Check(player, npc) {
						trigger.Action(player, npc)
						return true
					}
				}
				if time.Since(npc.TransAttrs.VarTime("lastFight")) <= time.Second*2 || npc.Busy() {
					return true
				}
				player.ResetPath()
				npc.ResetPath()
				player.StartCombat(npc)
				return true
			} else {
				player.SetPath(world.MakePath(player.Location, npc.Location))
			}
			return false
		})
	}
	PacketHandlers["attackplayer"] = func(player *world.Player, p *packet.Packet) {
		affectedPlayer, ok := world.Players.FromIndex(p.ReadShort())
		if affectedPlayer == nil || !ok {
			log.Suspicious.Printf("player[%v] tried to attack nil player\n", player)
			return
		}
		if player.Busy() {
			return
		}
		if affectedPlayer.Busy() {
			log.Info.Printf("Target player busy during attack request  State: %d\n", affectedPlayer.State)
			return
		}
		player.SetDistancedAction(func() bool {
			if player.NextTo(affectedPlayer.Location) && player.WithinRange(affectedPlayer.Location, 2) {
				player.ResetPath()
				if time.Since(affectedPlayer.TransAttrs.VarTime("lastRetreat")) <= time.Second*3 || affectedPlayer.IsFighting() {
					return true
				}
				player.ResetPath()
				affectedPlayer.ResetPath()
				player.StartCombat(affectedPlayer)
				return true
			}
			return player.FinishedPath()
		})
	}
	PacketHandlers["fightmode"] = func(player *world.Player, p *packet.Packet) {
		mode := p.ReadByte()
		if mode < 0 || mode > 3 {
			log.Suspicious.Printf("Invalid fightmode selected (%v) by %v", mode, player.String())
			return
		}
		player.SetFightMode(int(mode))
	}
}
