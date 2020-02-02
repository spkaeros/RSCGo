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

	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	AddHandler("duelreq", func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if player.Busy() {
			return
		}
		target, ok := world.Players.FromIndex(index)
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
		if target.DuelBlocked() && !target.Friends(player.UsernameHash()) {
			player.Message("This player has duel requests blocked.")
			return
		}
		player.SetDuelTarget(target)
		if target.DuelTarget() != player {
			player.Message("Sending duel request")
			target.Message(player.Username() + " " + world.CombatPrefix(target.CombatDelta(player)) + "(level-" + strconv.Itoa(player.Skills().CombatLevel()) + ")@whi@ wishes to duel with you")
			return
		}
		if player.Busy() || target.Busy() {
			return
		}
		player.AddState(world.MSDueling)
		player.ResetPath()
		player.SendPacket(world.DuelOpen(target.Index))

		target.AddState(world.MSDueling)
		target.ResetPath()
		target.SendPacket(world.DuelOpen(player.Index))
	})
	AddHandler("duelupdate", func(player *world.Player, p *packet.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v attempted to update a duel it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.DuelTarget()
		if target == nil {
			log.Suspicious.Printf("%v attempted to update a duel with a non-existent target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.DuelTarget() != player {
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
		if (target.TransAttrs.VarBool("duel1accept", false) || target.TransAttrs.VarBool("duel2accept", false)) && (player.TransAttrs.VarBool("duel1accept", false) || player.TransAttrs.VarBool("duel2accept", false)) {
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
		itemCount := int(p.ReadByte())
		if itemCount < 0 || itemCount > 8 {
			log.Suspicious.Printf("%v attempted to offer an invalid amount[%v] of duel items!\n", player.String(), itemCount)
			return
		}
		if len(p.Payload) < 1+(itemCount*6) {
			log.Suspicious.Printf("%v attempted to send a duel offer update packet without enough data for the offer.\n", player.String())
			return
		}
		for i := 0; i < itemCount; i++ {
			player.DuelOffer.Add(p.ReadShort(), p.ReadInt())
		}
	})
	AddHandler("dueloptions", func(player *world.Player, p *packet.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v tried changing duel options in a duel that they are not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.DuelTarget()
		if target == nil {
			log.Suspicious.Printf("%v involved in duel with no target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.DuelTarget() != player {
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

		retreatsAllowed := p.ReadBool()
		magicAllowed := p.ReadBool()
		prayerAllowed := p.ReadBool()
		equipmentAllowed := p.ReadBool()

		player.TransAttrs.SetVar("duelCanRetreat", !retreatsAllowed)
		player.TransAttrs.SetVar("duelCanMagic", !magicAllowed)
		player.TransAttrs.SetVar("duelCanPrayer", !prayerAllowed)
		player.TransAttrs.SetVar("duelCanEquip", !equipmentAllowed)

		target.TransAttrs.SetVar("duelCanRetreat", !retreatsAllowed)
		target.TransAttrs.SetVar("duelCanMagic", !magicAllowed)
		target.TransAttrs.SetVar("duelCanPrayer", !prayerAllowed)
		target.TransAttrs.SetVar("duelCanEquip", !equipmentAllowed)
		player.SendPacket(world.DuelOptions(player))
		target.SendPacket(world.DuelOptions(target))
	})
	AddHandler("dueldecline", func(player *world.Player, p *packet.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v attempted to decline a duel it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.DuelTarget()
		if target == nil {
			log.Suspicious.Printf("%v attempted to decline a duel with a non-existent target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.DuelTarget() != player {
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
	AddHandler("duelaccept", func(player *world.Player, p *packet.Packet) {
		if !player.IsDueling() {
			log.Suspicious.Printf("%v attempted to decline a duel it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.DuelTarget()
		if target == nil {
			log.Suspicious.Printf("%v attempted to accept a duel with a non-existent target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.DuelTarget() != player {
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
		player.SetDuel1Accepted()
		if target.TransAttrs.VarBool("duel1accept", false) {
			player.SendPacket(world.DuelConfirmationOpen(player, target))
			target.SendPacket(world.DuelConfirmationOpen(target, player))
		} else {
			target.SendPacket(world.DuelTargetAccept(true))
		}
	})
	AddHandler("duelconfirmaccept", func(player *world.Player, p *packet.Packet) {
		if !player.IsDueling() || !player.TransAttrs.VarBool("duel1accept", false) {
			log.Suspicious.Printf("%v attempted to accept a duel confirmation it was not in!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		target := player.DuelTarget()
		if target == nil {
			log.Suspicious.Printf("%v involved in duel with no target!\n", player.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			return
		}
		if !target.IsDueling() || target.DuelTarget() != player || !target.TransAttrs.VarBool("duel1accept", false) {
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
		player.SetDuel2Accepted()
		if target.TransAttrs.VarBool("duel2accept", false) {
			player.ResetDuelAccepted()
			target.ResetDuelAccepted()
			if !player.TransAttrs.VarBool("duelCanPrayer", true) {
				for i := 0; i < 14; i++ {
					player.PrayerOff(i)
				}
				player.SendPrayers()
				player.Message("You cannot use prayer in this duel!")
			}
			if !player.TransAttrs.VarBool("duelCanEquip", true) {
				player.Inventory.Range(func(item *world.Item) bool {
					if item.Worn {
						player.DequipItem(item)
					}
					return true
				})
			}
			if !target.TransAttrs.VarBool("duelCanPrayer", true) {
				for i := 0; i < 14; i++ {
					target.PrayerOff(i)
				}
				target.SendPrayers()
				target.Message("You cannot use prayer in this duel!")
			}
			if !target.TransAttrs.VarBool("duelCanEquip", true) {
				target.Inventory.Range(func(item *world.Item) bool {
					if item.Worn {
						target.DequipItem(item)
					}
					return true
				})
			}
			player.StartCombat(target)
			player.SendPacket(world.DuelClose)
			target.SendPacket(world.DuelClose)
			player.Message("Commencing Duel!")
			target.Message("Commencing Duel!")
		}
	})
}
