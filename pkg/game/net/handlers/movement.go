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
	AddHandler("walkto", func(player *world.Player, p *net.Packet) {
		if player.IsFighting() {
			target := player.FightTarget()
			if target == nil {
				player.ResetFighting()
				return
			}
			if player.IsFighting() && player.IsDueling() && !player.TransAttrs.VarBool("duelCanRetreat", true) {
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
		numWaypoints := p.Available() / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadInt8()))
			waypointsY = append(waypointsY, int(p.ReadInt8()))
		}
		player.ResetAll()
		player.SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	})
	AddHandler("walktoentity", func(player *world.Player, p *net.Packet) {
		if player.IsFighting() {
			return
		}
		if !player.CanWalk() {
			return
		}
		startX := p.ReadUint16()
		startY := p.ReadUint16()
		numWaypoints := p.Available() / 2
		var waypointsX, waypointsY []int
		for i := 0; i < numWaypoints; i++ {
			waypointsX = append(waypointsX, int(p.ReadInt8()))
			waypointsY = append(waypointsY, int(p.ReadInt8()))
		}
		player.ResetAll()
		player.SetPath(world.NewPathway(startX, startY, waypointsX, waypointsY))
	})
	AddHandler("followreq", func(player *world.Player, p *net.Packet) {
		if player.IsFighting() {
			return
		}
		if !player.CanWalk() {
			return
		}
		playerID := p.ReadUint16()
		affectedClient, ok := world.Players.FromIndex(playerID)
		if !ok {
			player.Message("@que@Could not find the player you're looking for.")
			return
		}
		player.ResetAll()
		player.StartFollowing(2)
		player.Message("@que@Following " + affectedClient.Username())
		player.SetDistancedAction(func() bool {
			if !player.IsFollowing() {
				// Following vars have been reset.
				return true
			}
			if affectedClient == nil || !affectedClient.Connected() ||
				!player.WithinRange(affectedClient.Location, 16) {
				// We think we have a target, but they're miles away now or no longer exist
				player.ResetFollowing()
				return true
			}
			if !player.FinishedPath() && player.WithinRange(affectedClient.Location, player.FollowRadius()) {
				// We're not done moving toward our target, but we're close enough that we should stop
				player.ResetPath()
			} else if !player.WithinRange(affectedClient.Location, player.FollowRadius()) {
				// We're not moving, but our target is moving away, so we must try to get closer
				player.SetPath(world.MakePath(player.Location, affectedClient.Location))
			}
			return false
		})
	})
	AddHandler("appearancerequest", func(player *world.Player, p *net.Packet) {
		playerCount := p.ReadUint16()
		for i := 0; i < playerCount; i++ {
			serverIndex := p.ReadUint16()
			appearanceTicket := p.ReadUint16()
			player.AppearanceLock.Lock()
			if ticket, ok := player.KnownAppearances[serverIndex]; !ok || ticket != appearanceTicket {
				if c1, ok := world.Players.FromIndex(serverIndex); ok {
					player.AppearanceReq = append(player.AppearanceReq, c1)
				}
			}
			player.AppearanceLock.Unlock()
		}
	})
}
