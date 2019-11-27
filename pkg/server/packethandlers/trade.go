package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["tradereq"] = func(c clients.Client, p *packet.Packet) {
		index := p.ReadShort()
		c1, ok := clients.FromIndex(index)
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to trade a player that does not exist.\n", c.Player().Username, c.Player().IP)
			return
		}
		if !c.Player().WithinRange(c1.Player().Location, 16) || c.Player().Busy() {
			// TODO: Reasonably, 16 is really too far.  Visibly I think 5 or 6 tiles surrounding players is visible
			return
		}
		if c1.Player().TradeBlocked() && !c1.Player().Friends(c.Player().UserBase37) {
			c.Message("This player has trade requests blocked.")
			return
		}
		c.Player().SetTradeTarget(index)
		if c1.Player().TradeTarget() == c.Player().Index {
			if c.Player().Busy() || c1.Player().Busy() {
				return
			}
			c.Player().AddState(world.MSTrading)
			c.Player().ResetPath()
			c.TradeOpen()

			c1.Player().AddState(world.MSTrading)
			c1.Player().ResetPath()
			c1.TradeOpen()
		} else {
			c.Message("Sending trade request.")
			c1.Message(c.Player().Username + " wishes to trade with you.")
		}
	}
	PacketHandlers["tradeupdate"] = func(c clients.Client, p *packet.Packet) {
		if !c.Player().IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a non-existant trade!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := clients.FromIndex(c.Player().TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to update a trade with a non-existent target!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.Player().IsTrading() || c1.Player().TradeTarget() != c.Player().Index || c.Player().TradeTarget() != c1.Player().Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
			c.Player().ResetTrade()
			c1.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		if (c1.Player().TransAttrs.VarBool("trade1accept", false) || c1.Player().TransAttrs.VarBool("trade2accept", false)) && (c.Player().TransAttrs.VarBool("trade1accept", false) || c.Player().TransAttrs.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade, player 1 attempted to alter offer after both players accepted!\n", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
			c.Player().ResetTrade()
			c1.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		c.Player().TransAttrs.UnsetVar("trade1accept")
		c1.Player().TransAttrs.UnsetVar("trade1accept")
		c.Player().TradeOffer.Clear()
		defer func() {
			c1.SendPacket(packetbuilders.TradeUpdate(c.Player()))
		}()
		itemCount := int(p.ReadByte())
		if itemCount < 0 || itemCount > 12 {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to offer an invalid amount[%v] of trade items!\n", c.Player().Username, c.Player().IP, itemCount)
			return
		}
		if len(p.Payload) < 1+(itemCount*6) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to send a trade offer update packet without enough data for the offer.\n", c.Player().Username, c.Player().IP)
			return
		}
		for i := 0; i < itemCount; i++ {
			c.Player().TradeOffer.Add(p.ReadShort(), p.ReadInt())
		}
	}
	PacketHandlers["tradedecline"] = func(c clients.Client, p *packet.Packet) {
		if !c.Player().IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade it was not in!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := clients.FromIndex(c.Player().TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to decline a trade with a non-existent target!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.Player().IsTrading() || c1.Player().TradeTarget() != c.Player().Index || c.Player().TradeTarget() != c1.Player().Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
		}
		c.Player().ResetTrade()
		c1.Player().ResetTrade()
		c1.Message(c.Player().Username + " has declined the trade.")
		c1.SendPacket(packetbuilders.TradeClose)
	}
	PacketHandlers["tradeaccept"] = func(c clients.Client, p *packet.Packet) {
		if !c.Player().IsTrading() {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade it was not in!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := clients.FromIndex(c.Player().TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade with a non-existent target!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.Player().IsTrading() || c1.Player().TradeTarget() != c.Player().Index || c.Player().TradeTarget() != c1.Player().Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
			c.Player().ResetTrade()
			c1.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		c.Player().TransAttrs.SetVar("trade1accept", true)
		if c1.Player().TransAttrs.VarBool("trade1accept", false) {
			c.SendPacket(packetbuilders.TradeConfirmationOpen(c.Player(), c1.Player()))
			c1.SendPacket(packetbuilders.TradeConfirmationOpen(c1.Player(), c.Player()))
		} else {
			c1.SendPacket(packetbuilders.TradeTargetAccept(true))
		}
	}
	PacketHandlers["tradeconfirmaccept"] = func(c clients.Client, p *packet.Packet) {
		if !c.Player().IsTrading() || !c.Player().TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation it was not in!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		c1, ok := clients.FromIndex(c.Player().TradeTarget())
		if !ok {
			log.Suspicious.Printf("player['%v'@'%v'] attempted to accept a trade confirmation with a non-existent target!\n", c.Player().Username, c.Player().IP)
			c.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			return
		}
		if !c1.Player().IsTrading() || c1.Player().TradeTarget() != c.Player().Index || c.Player().TradeTarget() != c1.Player().Index || !c1.Player().TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
			c.Player().ResetTrade()
			c1.Player().ResetTrade()
			c.SendPacket(packetbuilders.TradeClose)
			c1.SendPacket(packetbuilders.TradeClose)
			return
		}
		c.Player().TransAttrs.SetVar("trade2accept", true)
		if c1.Player().TransAttrs.VarBool("trade2accept", false) {
			neededSlots := c1.Player().TradeOffer.Size()
			availSlots := c.Player().Items.Capacity - c.Player().Items.Size() + c.Player().TradeOffer.Size()
			theirNeededSlots := c.Player().TradeOffer.Size()
			theirAvailSlots := c1.Player().Items.Capacity - c1.Player().Items.Size() + c1.Player().TradeOffer.Size()
			if theirNeededSlots > theirAvailSlots {
				c.Message("The other player does not have room to accept your items.")
				c.Player().ResetTrade()
				c1.Message("You do not have room in your inventory to hold those items.")
				c1.Player().ResetTrade()
				c.SendPacket(packetbuilders.TradeClose)
				c1.SendPacket(packetbuilders.TradeClose)
				return
			}
			if neededSlots > availSlots {
				c.Message("You do not have room in your inventory to hold those items.")
				c.Player().ResetTrade()
				c1.Message("The other player does not have room to accept your items.")
				c1.Player().ResetTrade()
				c.SendPacket(packetbuilders.TradeClose)
				c1.SendPacket(packetbuilders.TradeClose)
				return
			}
			defer func() {
				c.SendPacket(packetbuilders.InventoryItems(c.Player()))
				c.SendPacket(packetbuilders.TradeClose)
				c.Player().ResetTrade()
				c1.SendPacket(packetbuilders.InventoryItems(c1.Player()))
				c1.SendPacket(packetbuilders.TradeClose)
				c1.Player().ResetTrade()
			}()
			if c.Player().Items.RemoveAll(c.Player().TradeOffer) != c.Player().TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 1 did not have all items to give.", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
				return
			}
			if c1.Player().Items.RemoveAll(c1.Player().TradeOffer) != c1.Player().TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 2 did not have all items to give.", c.Player().Username, c.Player().IP, c1.Player().Username, c1.Player().IP)
				return
			}
			for i := 0; i < c1.Player().TradeOffer.Size(); i++ {
				item := c1.Player().TradeOffer.Get(i)
				c.Player().Items.Add(item.ID, item.Amount)
			}
			for i := 0; i < c.Player().TradeOffer.Size(); i++ {
				item := c.Player().TradeOffer.Get(i)
				c1.Player().Items.Add(item.ID, item.Amount)
			}
			c.Message("Trade completed.")
			c1.Message("Trade completed.")
		}
	}
	PacketHandlers["duelreq"] = func(c clients.Client, p *packet.Packet) {
		c.Message("@que@@ora@Not yet implemented")
	}
}
