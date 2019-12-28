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
	"strconv"

	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["tradereq"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		p1, ok := world.Players.FromIndex(index)
		if !ok {
			log.Suspicious.Printf("%v attempted to trade a player that does not exist.\n", player.String())
			return
		}
		if !player.WithinRange(p1.Location, 16) || player.Busy() {
			return
		}
		if !player.WithinRange(p1.Location, 5) {
			player.Message("You are too far away to do that")
			return
		}
		if p1.TradeBlocked() && !p1.Friends(player.UsernameHash()) {
			player.Message("This player has trade requests blocked.")
			return
		}
		player.SetTradeTarget(index)
		if p1.TradeTarget() == player.Index {
			if player.Busy() || p1.Busy() {
				return
			}
			player.AddState(world.MSTrading)
			player.ResetPath()
			player.SendPacket(world.TradeOpen(p1.Index))

			p1.AddState(world.MSTrading)
			p1.ResetPath()
			p1.SendPacket(world.TradeOpen(player.Index))
		} else {
			player.Message("Sending trade request.")
			p1.Message(player.Username() + " wishes to trade with you.")
		}
	}
	PacketHandlers["duelreq"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		p1, ok := world.Players.FromIndex(index)
		if !ok {
			log.Suspicious.Printf("%v attempted to duel a player that does not exist.\n", player.String())
			return
		}
		if !player.WithinRange(p1.Location, 16) || player.Busy() {
			return
		}
		if !player.WithinRange(p1.Location, 5) {
			player.Message("You are too far away to do that")
			return
		}
		if p1.DuelBlocked() && !p1.Friends(player.UsernameHash()) {
			player.Message("This player has duel requests blocked.")
			return
		}
		player.SetDuelTarget(p1)
		if p1.DuelTarget() != player {
			player.Message("Sending duel request")
			p1.Message(player.Username() + " " + world.CombatPrefix(p1.CombatDelta(player)) + "(level-" + strconv.Itoa(player.Skills().CombatLevel()) + ")@whi@ wishes to duel with you")
			return
		}
		if player.Busy() || p1.Busy() {
			return
		}
		player.AddState(world.MSDueling)
		player.ResetPath()
		player.SendPacket(world.DuelOpen(p1.Index))

		p1.AddState(world.MSDueling)
		p1.ResetPath()
		p1.SendPacket(world.DuelOpen(player.Index))
	}
	PacketHandlers["tradeupdate"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("%v attempted to decline a non-existant trade!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("%v attempted to update a trade with a non-existent target!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in trade with apparently bad trade variables!\n", player.String(), c1.String())
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(world.TradeClose)
			c1.SendPacket(world.TradeClose)
			return
		}
		if (c1.TransAttrs.VarBool("trade1accept", false) || c1.TransAttrs.VarBool("trade2accept", false)) && (player.TransAttrs.VarBool("trade1accept", false) || player.TransAttrs.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in trade, player 1 attempted to alter offer after both players accepted!\n", player.String(), c1.String())
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(world.TradeClose)
			c1.SendPacket(world.TradeClose)
			return
		}
		player.TransAttrs.UnsetVar("trade1accept")
		c1.TransAttrs.UnsetVar("trade1accept")
		player.TradeOffer.Clear()
		defer func() {
			c1.SendPacket(world.TradeUpdate(player))
		}()
		itemCount := int(p.ReadByte())
		if itemCount < 0 || itemCount > 12 {
			log.Suspicious.Printf("%v attempted to offer an invalid amount[%v] of trade items!\n", player.String(), itemCount)
			return
		}
		if len(p.Payload) < 1+(itemCount*6) {
			log.Suspicious.Printf("%v attempted to send a trade offer update packet without enough data for the offer.\n", player.String())
			return
		}
		for i := 0; i < itemCount; i++ {
			player.TradeOffer.Add(p.ReadShort(), p.ReadInt())
		}
	}
	PacketHandlers["duelupdate"] = func(player *world.Player, p *packet.Packet) {
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
	}
	PacketHandlers["tradedecline"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("%v attempted to decline a trade it was not in!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("%v attempted to decline a trade with a non-existent target!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in trade with apparently bad trade variables!\n", player.String(), c1.String())
		}
		player.ResetTrade()
		c1.ResetTrade()
		c1.Message(player.Username() + " has declined the trade.")
		c1.SendPacket(world.TradeClose)
	}
	PacketHandlers["dueldecline"] = func(player *world.Player, p *packet.Packet) {
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
		player.ResetDuel()
		player.SendPacket(world.DuelClose)
		target.ResetDuel()
		target.Message(player.Username() + " has declined the duel")
		target.SendPacket(world.DuelClose)
	}
	PacketHandlers["tradeaccept"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("%v attempted to accept a trade it was not in!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("%v attempted to accept a trade with a non-existent target!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in trade with apparently bad trade variables!\n", player.String(), c1.String())
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(world.TradeClose)
			c1.SendPacket(world.TradeClose)
			return
		}
		player.TransAttrs.SetVar("trade1accept", true)
		if c1.TransAttrs.VarBool("trade1accept", false) {
			player.SendPacket(world.TradeConfirmationOpen(player, c1))
			c1.SendPacket(world.TradeConfirmationOpen(c1, player))
		} else {
			c1.SendPacket(world.TradeTargetAccept(true))
		}
	}
	PacketHandlers["duelaccept"] = func(player *world.Player, p *packet.Packet) {
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
		player.SetDuel1Accepted()
		if target.TransAttrs.VarBool("duel1accept", false) {
			player.SendPacket(world.DuelConfirmationOpen(player, target))
			target.SendPacket(world.DuelConfirmationOpen(target, player))
		} else {
			target.SendPacket(world.DuelTargetAccept(true))
		}
	}
	PacketHandlers["tradeconfirmaccept"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() || !player.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("%v attempted to accept a trade confirmation it was not in!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("%v attempted to accept a trade confirmation with a non-existent target!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index || !c1.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:%v; 2:%v } involved in trade with apparently bad trade variables!\n", player.String(), c1.String())
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(world.TradeClose)
			c1.SendPacket(world.TradeClose)
			return
		}
		player.TransAttrs.SetVar("trade2accept", true)
		if c1.TransAttrs.VarBool("trade2accept", false) {
			neededSlots := c1.TradeOffer.Size()
			availSlots := player.Inventory.Capacity - player.Inventory.Size() + player.TradeOffer.Size()
			theirNeededSlots := player.TradeOffer.Size()
			theirAvailSlots := c1.Inventory.Capacity - c1.Inventory.Size() + c1.TradeOffer.Size()
			if theirNeededSlots > theirAvailSlots {
				player.Message("The other player does not have room to accept your items.")
				player.ResetTrade()
				c1.Message("You do not have room in your inventory to hold those items.")
				c1.ResetTrade()
				player.SendPacket(world.TradeClose)
				c1.SendPacket(world.TradeClose)
				return
			}
			if neededSlots > availSlots {
				player.Message("You do not have room in your inventory to hold those items.")
				player.ResetTrade()
				c1.Message("The other player does not have room to accept your items.")
				c1.ResetTrade()
				player.SendPacket(world.TradeClose)
				c1.SendPacket(world.TradeClose)
				return
			}
			defer func() {
				player.SendPacket(world.InventoryItems(player))
				player.SendPacket(world.TradeClose)
				player.ResetTrade()
				c1.SendPacket(world.InventoryItems(c1))
				c1.SendPacket(world.TradeClose)
				c1.ResetTrade()
			}()
			if player.Inventory.RemoveAll(player.TradeOffer) != player.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ %v;2:%v } involved in a trade, player 1 did not have all items to give.", player.String(), c1.String())
				return
			}
			if c1.Inventory.RemoveAll(c1.TradeOffer) != c1.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ %v;2:%v } involved in a trade, player 2 did not have all items to give.", player.String(), c1.String())
				return
			}
			for i := 0; i < c1.TradeOffer.Size(); i++ {
				item := c1.TradeOffer.Get(i)
				player.Inventory.Add(item.ID, item.Amount)
			}
			for i := 0; i < player.TradeOffer.Size(); i++ {
				item := player.TradeOffer.Get(i)
				c1.Inventory.Add(item.ID, item.Amount)
			}
			player.Message("Trade completed.")
			c1.Message("Trade completed.")
		}
	}

	PacketHandlers["dueloptions"] = func(player *world.Player, p *packet.Packet) {
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
		if target.DuelTarget() != player {
			log.Suspicious.Printf("Players{ 1:%v; 2:%v } involved in duel with apparently bad state!\n", player.String(), target.String())
			player.ResetDuel()
			player.SendPacket(world.DuelClose)
			target.ResetDuel()
			target.SendPacket(world.DuelClose)
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
	}
	PacketHandlers["duelconfirmaccept"] = func(player *world.Player, p *packet.Packet) {
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
	}
}
