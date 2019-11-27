package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["tradereq"] = func(c *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		c1, ok := players.FromIndex(index)
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to trade a player that does not exist.\n", c.Username, c.IP)
			return
		}
		if !c.WithinRange(c1.Location, 16) || c.Busy() {
			// TODO: Reasonably, 16 is really too far.  Visibly I think 5 or 6 tiles surrounding players is visible
			return
		}
		if c1.TradeBlocked() && !c1.Friends(c.UserBase37) {
			c.SendPacket(packetbuilders.ServerMessage("This player has trade requests blocked."))
			return
		}
		c.SetTradeTarget(index)
		if c1.TradeTarget() == c.Index {
			if c.Busy() || c1.Busy() {
				return
			}
			c.AddState(world.MSTrading)
			c.ResetPath()
			c.SendPacket(packetbuilders.TradeOpen(c))

			c1.AddState(world.MSTrading)
			c1.ResetPath()
			c1.SendPacket(packetbuilders.TradeOpen(c1))
		} else {
			c.SendPacket(packetbuilders.ServerMessage("Sending trade request."))
			c1.SendPacket(packetbuilders.ServerMessage(c.Username + " wishes to trade with you."))
		}
	}
	PacketHandlers["tradeupdate"] = func(c *world.Player, p *packet.Packet) {
		if !c.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a non-existant trade!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(c.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to update a trade with a non-existent target!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != c.Index || c.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Username, c.IP, c1.Username, c1.IP)
			c.ResetTrade()
			c1.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		if (c1.TransAttrs.VarBool("trade1accept", false) || c1.TransAttrs.VarBool("trade2accept", false)) && (c.TransAttrs.VarBool("trade1accept", false) || c.TransAttrs.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade, player 1 attempted to alter offer after both players accepted!\n", c.Username, c.IP, c1.Username, c1.IP)
			c.ResetTrade()
			c1.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		c.TransAttrs.UnsetVar("trade1accept")
		c1.TransAttrs.UnsetVar("trade1accept")
		c.TradeOffer.Clear()
		defer func() {
			c1.SendPacket(packetbuilders.TradeUpdate(c))
		}()
		itemCount := int(p.ReadByte())
		if itemCount < 0 || itemCount > 12 {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to offer an invalid amount[%v] of trade items!\n", c.Username, c.IP, itemCount)
			return
		}
		if len(p.Payload) < 1+(itemCount*6) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to send a trade offer update packet without enough data for the offer.\n", c.Username, c.IP)
			return
		}
		for i := 0; i < itemCount; i++ {
			c.TradeOffer.Add(p.ReadShort(), p.ReadInt())
		}
	}
	PacketHandlers["tradedecline"] = func(c *world.Player, p *packet.Packet) {
		if !c.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade it was not in!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(c.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade with a non-existent target!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != c.Index || c.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Username, c.IP, c1.Username, c1.IP)
		}
		c.ResetTrade()
		c1.ResetTrade()
		c1.SendPacket(packetbuilders.ServerMessage(c.Username + " has declined the trade."))
		c1.SendPacket(packetbuilders.TradeClose)
	}
	PacketHandlers["tradeaccept"] = func(c *world.Player, p *packet.Packet) {
		if !c.IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade it was not in!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(c.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade with a non-existent target!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != c.Index || c.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Username, c.IP, c1.Username, c1.IP)
			c.ResetTrade()
			c1.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		c.TransAttrs.SetVar("trade1accept", true)
		if c1.TransAttrs.VarBool("trade1accept", false) {
			c.SendPacket(packetbuilders.TradeConfirmationOpen(c, c1))
			c1.SendPacket(packetbuilders.TradeConfirmationOpen(c1, c))
		} else {
			c1.SendPacket(packetbuilders.TradeTargetAccept(true))
		}
	}
	PacketHandlers["tradeconfirmaccept"] = func(c *world.Player, p *packet.Packet) {
		if !c.IsTrading() || !c.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation it was not in!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := players.FromIndex(c.TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation with a non-existent target!\n", c.Username, c.IP)
			c.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.IsTrading() || c1.TradeTarget() != c.Index || c.TradeTarget() != c1.Index || !c1.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Username, c.IP, c1.Username, c1.IP)
			c.ResetTrade()
			c1.ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		c.TransAttrs.SetVar("trade2accept", true)
		if c1.TransAttrs.VarBool("trade2accept", false) {
			neededSlots := c1.TradeOffer.Size()
			availSlots := c.Items.Capacity - c.Items.Size() + c.TradeOffer.Size()
			theirNeededSlots := c.TradeOffer.Size()
			theirAvailSlots := c1.Items.Capacity - c1.Items.Size() + c1.TradeOffer.Size()
			if theirNeededSlots > theirAvailSlots {
				c.SendPacket(packetbuilders.ServerMessage("The other player does not have room to accept your items."))
				c.ResetTrade()
				c1.SendPacket(packetbuilders.ServerMessage("You do not have room in your inventory to hold those items."))
				c1.ResetTrade()
				c.SendPacket(packetbuilders.TradeClose)
				c1.SendPacket(packetbuilders.TradeClose)
				return
			}
			if neededSlots > availSlots {
				c.SendPacket(packetbuilders.ServerMessage("You do not have room in your inventory to hold those items."))
				c.ResetTrade()
				c1.SendPacket(packetbuilders.ServerMessage("The other player does not have room to accept your items."))
				c1.ResetTrade()
				c.SendPacket(packetbuilders.TradeClose)
				c1.SendPacket(packetbuilders.TradeClose)
				return
			}
			defer func() {
				c.SendPacket(packetbuilders.InventoryItems(c))
				c.SendPacket(packetbuilders.TradeClose)
				c.ResetTrade()
				c1.SendPacket(packetbuilders.InventoryItems(c1))
				c1.SendPacket(packetbuilders.TradeClose)
				c1.ResetTrade()
			}()
			if c.Items.RemoveAll(c.TradeOffer) != c.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 1 did not have all items to give.", c.Username, c.IP, c1.Username, c1.IP)
				return
			}
			if c1.Items.RemoveAll(c1.TradeOffer) != c1.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 2 did not have all items to give.", c.Username, c.IP, c1.Username, c1.IP)
				return
			}
			for i := 0; i < c1.TradeOffer.Size(); i++ {
				item := c1.TradeOffer.Get(i)
				c.Items.Add(item.ID, item.Amount)
			}
			for i := 0; i < c.TradeOffer.Size(); i++ {
				item := c.TradeOffer.Get(i)
				c1.Items.Add(item.ID, item.Amount)
			}
			c.SendPacket(packetbuilders.ServerMessage("Trade completed."))
			c1.SendPacket(packetbuilders.ServerMessage("Trade completed."))
		}
	}
	PacketHandlers["duelreq"] = func(c *world.Player, p *packet.Packet) {
		c.SendPacket(packetbuilders.ServerMessage("@que@@ora@Not yet implemented"))
	}
}
