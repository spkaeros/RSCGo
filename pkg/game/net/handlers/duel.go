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
	"strconv"

	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	game.AddHandler("duelreq", func(player *world.Player, p *net.Packet) {
		if player.Busy() {
			return
		}
		index := p.ReadUint16()
		target, ok := world.Players.FindIndex(index)
		if !ok {
			log.Suspicious.Printf("%v attempted to duel a player that does not exist.\n", player.String())
			return
		}
		if !player.WithinRange(target.Location, 16) || player.Busy() {
			return
		}
		if !player.WithinRange(target.Location, 5) {
			player.Message("You are too far away to do that")
			return
		}
		if target.DuelBlocked() && !target.FriendsWith(player.UsernameHash()) {
			player.Message("This player has duel requests blocked.")
			return
		}
		player.SetDuelTarget(target)
		player.Duel.Target = target
		if target.Duel.Target != player {
			player.Message("Sending duel request")
			target.Message(player.Username() + " " + strutil.CombatPrefix(target.CombatDelta(player)) + "(level-" + strconv.Itoa(player.Skills().CombatLevel()) + ")@whi@ wishes to duel with you")
			return
		}
		if player.Busy() || target.Busy() {
			return
		}
		player.AddState(world.StateDueling)
		player.ResetPath()
		player.SendPacket(world.DuelOpen(target.Index))

		target.AddState(world.StateDueling)
		target.ResetPath()
		target.SendPacket(world.DuelOpen(player.Index))
	})
	game.AddHandler("duelupdate", func(player *world.Player, p *net.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v attempted to update a duel it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.Duel.Target
		if target == nil {
			log.Suspicious.Printf("%v attempted to update a duel with a non-existent target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.Duel.Target != player {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in duel with apparently bad state!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			target.ResetDuel()
			target.SendPacket(world.DuelClose)
			return
		}
		if player.IsFighting() || target.IsFighting() {
			log.Suspicious.Println(player, "attempted modifying duel state (with", target, ") during the duels fight!!")
			return
		}
		if (target.DuelAccepted(1) && player.DuelAccepted(1)) || (target.DuelAccepted(2) && player.DuelAccepted(2)) {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in duel, player 1 attempted to alter offer after both players accepted!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			target.ResetDuel()
			target.SendPacket(world.DuelClose)
			return
		}
		player.ResetDuelAccepted()
		target.ResetDuelAccepted()

		player.DuelOffer.Clear()
		defer func() {
			target.SendPacket(world.DuelUpdate(player))
		}()
		itemCount := int(p.ReadUint8())
		if itemCount < 0 || itemCount > 8 {
			log.Suspicious.Printf("%v attempted to offer an invalid amount[%v] of duel items!\n", player.String(), itemCount)
			return
		}
		if p.Length() < 1+(itemCount*6) {
			log.Suspicious.Printf("%v attempted to send a duel offer update without enough data for the offer.\n", player.String())
			log.Suspicious.Println(p.FrameBuffer)
			return
		}
		for i := 0; i < itemCount; i++ {
			player.DuelOffer.Add(p.ReadUint16(), p.ReadUint32())
		}
	})
	game.AddHandler("dueloptions", func(player *world.Player, p *net.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v tried changing duel options in a duel that they are not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.Duel.Target
		if target == nil {
			log.Suspicious.Printf("%v involved in duel with no target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.Duel.Target != player {
			log.Suspicious.Printf("Players{ 1:%v; 2:%v } involved in duel with apparently bad state!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			target.ResetDuel()
			target.SendPacket(world.DuelClose)
			return
		}
		if player.IsFighting() || target.IsFighting() {
			log.Suspicious.Println(player, "attempted modifying duel state (with", target, ") during the duels fight!!")
			return
		}
		player.ResetDuelAccepted()
		target.ResetDuelAccepted()

		//retreatsAllowed := p.ReadBoolean()
		//magicAllowed := p.ReadBoolean()
		//prayerAllowed := p.ReadBoolean()
		//equipmentAllowed := p.ReadBoolean()

		for i := 0; i < p.Length(); i++ {
			flag := p.ReadBoolean()
			player.SetDuelRule(i, flag)
			target.SetDuelRule(i, flag)
		}
		//player.SetVar("duelCanRetreat", !retreatsAllowed)
		//player.SetVar("duelCanMagic", !magicAllowed)
		//player.SetVar("duelCanPrayer", !prayerAllowed)
		//player.SetVar("duelCanEquip", !equipmentAllowed)
		player.SendPacket(world.DuelOptions(player))

		//target.SetVar("duelCanRetreat", !retreatsAllowed)
		//target.SetVar("duelCanMagic", !magicAllowed)
		//target.SetVar("duelCanPrayer", !prayerAllowed)
		//target.SetVar("duelCanEquip", !equipmentAllowed)
		target.SendPacket(world.DuelOptions(target))
	})
	game.AddHandler("dueldecline", func(player *world.Player, p *net.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v attempted to decline a duel it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.Duel.Target
		if target == nil {
			log.Suspicious.Printf("%v attempted to decline a duel with a non-existent target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.Duel.Target != player {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in duel with apparently bad state!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if player.IsFighting() || target.IsFighting() {
			log.Suspicious.Println(player, "attempted modifying duel state (with", target, ") during the duels fight!!")
			return
		}
		player.ResetDuel()
		player.SendPacket(world.DuelClose)
		target.ResetDuel()
		target.Message(player.Username() + " has declined the duel")
		target.SendPacket(world.DuelClose)
	})
	game.AddHandler("duelaccept", func(player *world.Player, p *net.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v attempted to decline a duel it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.Duel.Target
		if target == nil {
			log.Suspicious.Printf("%v attempted to accept a duel with a non-existent target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.Duel.Target != player {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in duel with apparently bad state!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			target.ResetDuel()
			target.SendPacket(world.DuelClose)
			return
		}
		if player.IsFighting() || target.IsFighting() {
			log.Suspicious.Println(player, "attempted modifying duel state (with", target, ") during the duels fight!!")
			return
		}
		player.SetDuelAccepted(1, true)
		if target.DuelAccepted(1) {
			player.SendPacket(world.DuelConfirmationOpen(player, target))
			target.SendPacket(world.DuelConfirmationOpen(target, player))
		} else {
			target.SendPacket(world.DuelTargetAccept(true))
		}
	})
	game.AddHandler("duelconfirmaccept", func(player *world.Player, p *net.Packet) {
		if !player.IsDueling() || !player.DuelAccepted(1) {
			log.Suspicious.Printf("%v attempted to accept a duel confirmation it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.Duel.Target
		if target == nil {
			log.Suspicious.Printf("%v involved in duel with no target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.Duel.Target != player || !target.DuelAccepted(1) {
			log.Suspicious.Printf("Players{ 1:%v; 2:%v } involved in duel with apparently bad state!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			target.ResetDuel()
			target.SendPacket(world.DuelClose)
			return
		}
		if player.IsFighting() || target.IsFighting() {
			log.Suspicious.Println(player, "attempted modifying duel state (with", target, ") during the duels fight!!")
			return
		}
		player.SetDuelAccepted(2, true)
		if target.DuelAccepted(2) {
			player.ResetDuelAccepted()
			target.ResetDuelAccepted()
			if !player.VarBool("duelCanPrayer", true) {
				for i := 0; i < 14; i++ {
					player.DeactivatePrayer(i)
				}
				player.SendPrayers()
				player.Message("You cannot use prayer in this duel!")
			}
			if !player.VarBool("duelCanEquip", true) {
				player.Inventory.Range(func(item *world.Item) bool {
					if item.Worn {
						player.DequipItem(item)
					}
					return true
				})
				player.SendInventory()
			}
			if !target.VarBool("duelCanPrayer", true) {
				for i := 0; i < 14; i++ {
					target.DeactivatePrayer(i)
				}
				target.SendPrayers()
				target.Message("You cannot use prayer in this duel!")
			}
			if !target.VarBool("duelCanEquip", true) {
				target.Inventory.Range(func(item *world.Item) bool {
					if item.Worn {
						target.DequipItem(item)
					}
					return true
				})
				target.SendInventory()
			}
			player.StartCombat(target)
			player.SendPacket(world.DuelClose)
			target.SendPacket(world.DuelClose)
			player.Message("Commencing Duel!")
			target.Message("Commencing Duel!")
		}
	})
}
