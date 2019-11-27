package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
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
			player.SendPacket(packetbuilders.ServerMessage("This player has trade requests blocked."))
			return
		}
		player.SetTradeTarget(index)
		if c1.TradeTarget() == player.Index {
			if player.Busy() || c1.Busy() {
				return
			}
			player.AddState(world.MSTrading)
			player.ResetPath()
			player.SendPacket(packetbuilders.TradeOpen(player))

			c1.AddState(world.MSTrading)
			c1.ResetPath()
			c1.SendPacket(packetbuilders.TradeOpen(c1))
		} else {
			player.SendPacket(packetbuilders.ServerMessage("Sending trade request."))
			c1.SendPacket(packetbuilders.ServerMessage(player.Username + " wishes to trade with you."))
		}
	}
	PacketHandlers["tradeupdate"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a non-existant trade!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to update a trade with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		if (c1.TransAttrs.VarBool("trade1accept", false) || c1.TransAttrs.VarBool("trade2accept", false)) && (player.TransAttrs.VarBool("trade1accept", false) || player.TransAttrs.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade, player 1 attempted to alter offer after both players accepted!\n", player.Username, player.IP, c1.Username, c1.IP)
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		player.TransAttrs.UnsetVar("trade1accept")
		c1.TransAttrs.UnsetVar("trade1accept")
		player.TradeOffer.Clear()
		defer func() {
			c1.SendPacket(packetbuilders.TradeUpdate(player))
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
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
		}
		player.ResetTrade()
		c1.ResetTrade()
		c1.SendPacket(packetbuilders.ServerMessage(player.Username + " has declined the trade."))
		c1.SendPacket(packetbuilders.TradeClose)
	}
	PacketHandlers["tradeaccept"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade it was not in!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		player.TransAttrs.SetVar("trade1accept", true)
		if c1.TransAttrs.VarBool("trade1accept", false) {
			player.SendPacket(packetbuilders.TradeConfirmationOpen(player, c1))
			c1.SendPacket(packetbuilders.TradeConfirmationOpen(c1, player))
		} else {
			c1.SendPacket(packetbuilders.TradeTargetAccept(true))
		}
	}
	PacketHandlers["tradeconfirmaccept"] = func(player *world.Player, p *packet.Packet) {
		if !player.IsTrading() || !player.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation it was not in!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation with a non-existent target!\n", player.Username, player.IP)
			player.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != player.Index || player.TradeTarget() != c1.Index || !c1.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", player.Username, player.IP, c1.Username, c1.IP)
			player.ResetTrade()
			c1.ResetTrade()
			player.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		player.TransAttrs.SetVar("trade2accept", true)
		if c1.TransAttrs.VarBool("trade2accept", false) {
			neededSlots := c1.TradeOffer.Size()
			availSlots := player.Items.Capacity - player.Items.Size() + player.TradeOffer.Size()
			theirNeededSlots := player.TradeOffer.Size()
			theirAvailSlots := c1.Items.Capacity - c1.Items.Size() + c1.TradeOffer.Size()
			if theirNeededSlots > theirAvailSlots {
				player.SendPacket(packetbuilders.ServerMessage("The other player does not have room to accept your items."))
				player.ResetTrade()
				c1.SendPacket(packetbuilders.ServerMessage("You do not have room in your inventory to hold those items."))
				c1.ResetTrade()
				player.SendPacket(packetbuilders.TradeClose)
				c1.SendPacket(packetbuilders.TradeClose)
				return
			}
			if neededSlots > availSlots {
				player.SendPacket(packetbuilders.ServerMessage("You do not have room in your inventory to hold those items."))
				player.ResetTrade()
				c1.SendPacket(packetbuilders.ServerMessage("The other player does not have room to accept your items."))
				c1.ResetTrade()
				player.SendPacket(packetbuilders.TradeClose)
				c1.SendPacket(packetbuilders.TradeClose)
				return
			}
			defer func() {
				player.SendPacket(packetbuilders.InventoryItems(player))
				player.SendPacket(packetbuilders.TradeClose)
				player.ResetTrade()
				c1.SendPacket(packetbuilders.InventoryItems(c1))
				c1.SendPacket(packetbuilders.TradeClose)
				c1.ResetTrade()
			}()
			if player.Items.RemoveAll(player.TradeOffer) != player.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 1 did not have all items to give.", player.Username, player.IP, c1.Username, c1.IP)
				return
			}
			if c1.Items.RemoveAll(c1.TradeOffer) != c1.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 2 did not have all items to give.", player.Username, player.IP, c1.Username, c1.IP)
				return
			}
			for i := 0; i < c1.TradeOffer.Size(); i++ {
				item := c1.TradeOffer.Get(i)
				player.Items.Add(item.ID, item.Amount)
			}
			for i := 0; i < player.TradeOffer.Size(); i++ {
				item := player.TradeOffer.Get(i)
				c1.Items.Add(item.ID, item.Amount)
			}
			player.SendPacket(packetbuilders.ServerMessage("Trade completed."))
			c1.SendPacket(packetbuilders.ServerMessage("Trade completed."))
		}
	}
	PacketHandlers["duelreq"] = func(player *world.Player, p *packet.Packet) {
		player.SendPacket(packetbuilders.ServerMessage("@que@@ora@Not yet implemented"))
	}
}
