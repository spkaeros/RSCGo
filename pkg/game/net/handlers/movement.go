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
	"math"
)

func init() {
	game.AddHandler("walkto", func(player *world.Player, p *net.Packet) {
		if player.IsFighting() {
			target := player.FightTarget()
			if target == nil {
				player.ResetFighting()
				return
			}
			if player.IsDueling() && player.IsFighting() && !player.DuelRetreating() {
				player.Message("You cannot retreat during this duel!")
				return
			}
			curRound := target.FightRound()
			if curRound < 3 {
				player.Message("You can't retreat during the first 3 rounds of combat")
				return
			}
			if target, ok := target.(*world.Player); ok {
				target.PlaySound("retreat")
				target.Message("Your opponent is retreating")
			}
			player.UpdateLastRetreat()
			player.ResetFighting()
		}
		if !player.CanWalk() {
			return
		}
		startX := p.ReadUint16()
		startY := p.ReadUint16()
		numWaypoints := math.Ceil(float64(p.Length()-5) / 2)
		var waypointsX, waypointsY []int
		for i := float64(0); i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadInt8()))
			waypointsY = append(waypointsY, int(p.ReadInt8()))
		}
		player.ResetAll()
		player.SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	})
	game.AddHandler("walktoentity", func(player *world.Player, p *net.Packet) {
		if player.IsFighting() {
			return
		}
		if !player.CanWalk() {
			return
		}
		startX := p.ReadUint16()
		startY := p.ReadUint16()
		//		numWaypoints := (p.Length()-5) / 2
		numWaypoints := math.Ceil(float64(p.Length()-5) / 2)
		var waypointsX, waypointsY []int
		for i := float64(0); i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadInt8()))
			waypointsY = append(waypointsY, int(p.ReadInt8()))
		}
		player.ResetAll()
		player.SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	})
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
			if !player.VarBool("following", false) {
				// Following vars have been reset.
				return true
			}
			if target == nil || !target.Connected() ||
				!player.WithinRange(target.Location, 16) {
				// We think we have a target, but they're miles away now or no longer exist
				player.UnsetVar("following")
				return true
			}
			if !player.FinishedPath() && player.WithinRange(target.Location, rad) {
				// We're not done moving toward our target, but we're close enough that we should stop
				player.ResetPath()
			} else if !player.WithinRange(target.Location, rad) {
				// We're not moving, but our target is moving away, so we must try to get closer
				if !player.WalkTo(target.Location) {
					return true
				}
			}
			return false
		})
	})
	game.AddHandler("appearancerequest", func(player *world.Player, p *net.Packet) {
		playerCount := p.ReadUint16()
		for i := 0; i < playerCount; i++ {
			serverIndex := p.ReadUint16()
			appearanceTicket := p.ReadUint16()
			if ticket, ok := player.KnownAppearances[serverIndex]; !ok || ticket != appearanceTicket {
				if c1, ok := world.Players.FindIndex(serverIndex); ok {
					player.AppearanceReq = append(player.AppearanceReq, c1)
				}
			}
		}
	})
}
