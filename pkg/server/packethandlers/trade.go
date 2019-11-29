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
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["tradereq"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		c1, ok := players.FromIndex(index)
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to trade a player that does not exist.\n", player.Username, player.IP)
			return
		}
		if !player.WithinRange(c1.Location, 16) || player.Busy() {
			// TODO: Reasonably, 16 is really too far.  Visibly I think 5 or 6 tiles surrounding players is visible
			return
		}
		if c1.TradeBlocked() && !c1.Friends(player.UserBase37) {
			player.Message("This player has trade requests blocked.")
			return
		}
		player.SetTradeTarget(index)
		if c1.TradeTarget() == player.Index {
			if player.Busy() || c1.Busy() {
				return
			}
			player.AddState(world.MSTrading)
			player.ResetPath()
			player.SendPacket(world.TradeOpen(player))

			c1.AddState(world.MSTrading)
			c1.ResetPath()
			c1.SendPacket(world.TradeOpen(c1))
		} else {
			player.Message("Sending trade request.")
			c1.Message(player.Username + " wishes to trade with you.")
		}
	}
	PacketHandlers["tradeupdate"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a non-existant trade!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to update a trade with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(world.TradeClose)
			c1.SendPacket(world.TradeClose)
			return
		}
		if (c1.TransAttrs.VarBool("trade1accept", false) || c1.TransAttrs.VarBool("trade2accept", false)) && (player.TransAttrs.VarBool("trade1accept", false) || player.TransAttrs.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade, player 1 attempted to alter offer after both players accepted!\n", player.Username, player.IP, c1.Username, c1.IP)
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
			log.Suspicious.Printf("player['%v'@'%v'] attempted to offer an invalid amount[%v] of trade items!\n", player.Username, player.IP, itemCount)
			return
		}
		if len(p.Payload) < 1+(itemCount*6) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to send a trade offer update packet without enough data for the offer.\n", player.Username, player.IP)
			return
		}
		for i := 0; i < itemCount; i++ {
			player.TradeOffer.Add(p.ReadShort(), p.ReadInt())
		}
	}
	PacketHandlers["tradedecline"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade it was not in!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
		}
		player.ResetTrade()
		c1.ResetTrade()
		c1.Message(player.Username + " has declined the trade.")
		c1.SendPacket(world.TradeClose)
	}
	PacketHandlers["tradeaccept"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade it was not in!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
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
	PacketHandlers["tradeconfirmaccept"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() || !player.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation it was not in!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(world.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index || !c1.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
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
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 1 did not have all items to give.", player.Username, player.IP, c1.Username, c1.IP)
				return
			}
			if c1.Inventory.RemoveAll(c1.TradeOffer) != c1.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 2 did not have all items to give.", player.Username, player.IP, c1.Username, c1.IP)
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
	PacketHandlers["duelreq"] = func(player *world.Player, p *packet.Packet) {
		player.Message("@que@@ora@Not yet implemented")
	}
}
