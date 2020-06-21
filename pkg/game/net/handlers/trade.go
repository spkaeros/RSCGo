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
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("tradereq", func(player *world.Player, p *net.Packet) {
		if player.Busy() {
			return
		}
		index := p.ReadUint16()
		p1, ok := world.Players.FindIndex(index)
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
		if p1.TradeBlocked() && !p1.FriendsWith(player.UsernameHash()) {
			player.Message("This player has trade requests blocked.")
			return
		}
		player.SetTradeTarget(index)
		if p1.TradeTarget() == player.Index {
			if player.Busy() || p1.Busy() {
				return
			}
			player.AddState(world.StateTrading)
			player.ResetPath()
			player.SendPacket(world.TradeOpen(p1.Index))

			p1.AddState(world.StateTrading)
			p1.ResetPath()
			p1.SendPacket(world.TradeOpen(player.Index))
		} else {
			player.Message("Sending trade request.")
			p1.Message(player.Username() + " wishes to trade with you.")
		}
	})
	game.AddHandler("tradeupdate", func(player *world.Player, p *net.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("%v attempted to decline a non-existent trade!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FindIndex(player.TradeTarget())
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
		if (c1.VarBool("trade1accept", false) || c1.VarBool("trade2accept", false)) && (player.VarBool("trade1accept", false) || player.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ %v;2:%v } involved in trade, player 1 attempted to alter offer after both players accepted!\n", player.String(), c1.String())
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(world.TradeClose)
			c1.SendPacket(world.TradeClose)
			return
		}
		player.UnsetVar("trade1accept")
		c1.UnsetVar("trade1accept")
		player.TradeOffer.Clear()
		defer func() {
			c1.SendPacket(world.TradeUpdate(player))
		}()
		itemCount := int(p.ReadUint8())
		if itemCount < 0 || itemCount > 12 {
			log.Suspicious.Printf("%v attempted to offer an invalid amount[%v] of trade items!\n", player.String(), itemCount)
			return
		}
		if p.Length() < 1+(itemCount*6) {
			log.Suspicious.Printf("%v attempted to send a trade offer update net without enough data for the offer.\n", player.String())
			return
		}
		for i := 0; i < itemCount; i++ {
			player.TradeOffer.Add(p.ReadUint16(), p.ReadUint32())
		}
	})
	game.AddHandler("tradedecline", func(player *world.Player, p *net.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("%v attempted to decline a trade it was not in!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FindIndex(player.TradeTarget())
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
	})
	game.AddHandler("tradeaccept", func(player *world.Player, p *net.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("%v attempted to accept a trade it was not in!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := world.Players.FindIndex(player.TradeTarget())
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
		player.SetVar("trade1accept", true)
		if c1.VarBool("trade1accept", false) {
			player.SendPacket(world.TradeConfirmationOpen(player, c1))
			c1.SendPacket(world.TradeConfirmationOpen(c1, player))
		} else {
			c1.SendPacket(world.TradeTargetAccept(true))
		}
	})
	game.AddHandler("tradeconfirmaccept", func(player *world.Player, p *net.Packet) {
		if !player.IsTrading() || !player.VarBool("trade1accept", false) {
			log.Suspicious.Printf("%v attempted to accept a trade confirmation it was not in!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		target, ok := world.Players.FindIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("%v attempted to accept a trade confirmation with a non-existent target!\n", player.String())
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !target.IsTrading() || target.TradeTarget() != player.Index || player.TradeTarget() != target.Index || !target.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:%v; 2:%v } involved in trade with apparently bad trade variables!\n", player.String(), target.String())
			player.ResetTrade()
			target.ResetTrade()
			player.SendPacket(world.TradeClose)
			target.SendPacket(world.TradeClose)
			return
		}
		player.SetVar("trade2accept", true)
		if target.VarBool("trade2accept", false) {
			neededSlots := target.TradeOffer.Size()
			availSlots := player.Inventory.Capacity - player.Inventory.Size() + player.TradeOffer.Size()
			theirNeededSlots := player.TradeOffer.Size()
			theirAvailSlots := target.Inventory.Capacity - target.Inventory.Size() + target.TradeOffer.Size()
			if theirNeededSlots > theirAvailSlots {
				player.Message("The other player does not have room to accept your items.")
				player.ResetTrade()
				target.Message("You do not have room in your inventory to hold those items.")
				target.ResetTrade()
				player.SendPacket(world.TradeClose)
				target.SendPacket(world.TradeClose)
				return
			}
			if neededSlots > availSlots {
				player.Message("You do not have room in your inventory to hold those items.")
				player.ResetTrade()
				target.Message("The other player does not have room to accept your items.")
				target.ResetTrade()
				player.SendPacket(world.TradeClose)
				target.SendPacket(world.TradeClose)
				return
			}
			defer func() {
				player.SendPacket(world.InventoryItems(player))
				player.SendPacket(world.TradeClose)
				player.ResetTrade()
				target.SendPacket(world.InventoryItems(target))
				target.SendPacket(world.TradeClose)
				target.ResetTrade()
			}()
			if player.Inventory.RemoveAll(player.TradeOffer) != player.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ %v;2:%v } involved in a trade, player 1 did not have all items to give.", player.String(), target.String())
				return
			}
			if target.Inventory.RemoveAll(target.TradeOffer) != target.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ %v;2:%v } involved in a trade, player 2 did not have all items to give.", player.String(), target.String())
				return
			}
			for i := 0; i < target.TradeOffer.Size(); i++ {
				item := target.TradeOffer.Get(i)
				player.Inventory.Add(item.ID, item.Amount)
			}
			for i := 0; i < player.TradeOffer.Size(); i++ {
				item := player.TradeOffer.Get(i)
				target.Inventory.Add(item.ID, item.Amount)
			}
			player.Message("Trade completed.")
			target.Message("Trade completed.")
		}
	})
}
