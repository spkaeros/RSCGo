package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["tradereq"] = func(c *Client, p *packets.Packet) {
		index := p.ReadShort()
		c1, ok := Clients.FromIndex(index)
		if !ok {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to trade a player that does not exist.\n", c.player.Username, c.ip)
			return
		}
		if c.player.IsFighting() {
			return
		}
		if !c.player.WithinRange(c1.player.Location, 16) || c.player.Busy() {
			// TODO: Reasonably, 16 is really too far.  Visibly I think 5 or 6 tiles surrounding players is visible
			return
		}
		if c1.player.TradeBlocked() && !c1.player.Friends(c.player.UserBase37) {
			c.Message("This player has trade requests blocked.")
			return
		}
		c.player.SetTradeTarget(index)
		if c1.player.TradeTarget() == c.Index {
			if c1.player.IsFighting() || c.player.IsFighting() || c.player.Busy() || c1.player.Busy() {
				return
			}
			c.player.State = world.MSTrading
			c.player.ResetPath()
			c.TradeOpen()

			c1.player.State = world.MSTrading
			c1.player.ResetPath()
			c1.TradeOpen()
		} else {
			c.Message("Sending trade request.")
			c1.Message(c.player.Username + " wishes to trade with you.")
		}
	}
	PacketHandlers["tradeupdate"] = func(c *Client, p *packets.Packet) {
		if c.player.TradeTarget() == -1 || c.player.State != world.MSTrading {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to decline a non-existant trade!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to update a trade with a non-existent target!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		if c1.player.State != world.MSTrading || c1.player.TradeTarget() != c.Index || c.player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.player.Username, c.ip, c1.player.Username, c1.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c1.player.ResetTrade()
			c1.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			c1.outgoingPackets <- packets.TradeClose
			return
		}
		if (c1.player.TransAttrs.VarBool("trade1accept", false) || c1.player.TransAttrs.VarBool("trade2accept", false)) && (c.player.TransAttrs.VarBool("trade1accept", false) || c.player.TransAttrs.VarBool("trade2accept", false)) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade, player 1 attempted to alter offer after both players accepted!\n", c.player.Username, c.ip, c1.player.Username, c1.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c1.player.ResetTrade()
			c1.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			c1.outgoingPackets <- packets.TradeClose
			return
		}
		c.player.TransAttrs.UnsetVar("trade1accept")
		c1.player.TransAttrs.UnsetVar("trade1accept")
		c.player.TradeOffer.Clear()
		defer func() {
			c1.outgoingPackets <- packets.TradeUpdate(c.player)
		}()
		itemCount := int(p.ReadByte())
		if itemCount < 0 || itemCount > 12 {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to offer an invalid amount[%v] of trade items!\n", c.player.Username, c.ip, itemCount)
			return
		}
		if len(p.Payload) < 1 + (itemCount * 6) {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to send a trade offer update packet without enough data for the offer.\n", c.player.Username, c.ip)
			return
		}
		for i := 0; i < itemCount; i++ {
			c.player.TradeOffer.Put(p.ReadShort(), p.ReadInt())
		}
	}
	PacketHandlers["tradedecline"] = func(c *Client, p *packets.Packet) {
		if c.player.TradeTarget() == -1 || c.player.State != world.MSTrading {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to decline a trade it was not in!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to decline a trade with a non-existent target!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		if c1.player.State != world.MSTrading || c1.player.TradeTarget() != c.Index || c.player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.player.Username, c.ip, c1.player.Username, c1.ip)
		}
		c.player.ResetTrade()
		c.player.State = world.MSIdle
		c1.player.ResetTrade()
		c1.player.State = world.MSIdle
		c1.Message(c.player.Username + " has declined the trade.")
		c1.outgoingPackets <- packets.TradeClose
	}
	PacketHandlers["tradeaccept"] = func(c *Client, p *packets.Packet) {
		if c.player.TradeTarget() == -1 || c.player.State != world.MSTrading {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to accept a trade it was not in!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to accept a trade with a non-existent target!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		if c1.player.State != world.MSTrading || c1.player.TradeTarget() != c.Index || c.player.TradeTarget() != c1.Index {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.player.Username, c.ip, c1.player.Username, c1.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c1.player.ResetTrade()
			c1.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			c1.outgoingPackets <- packets.TradeClose
			return
		}
		c.player.TransAttrs.SetVar("trade1accept", true)
		if c1.player.TransAttrs.VarBool("trade1accept", false) {
			c.outgoingPackets <- packets.TradeConfirmationOpen(c.player, c1.player)
			c1.outgoingPackets <- packets.TradeConfirmationOpen(c1.player, c.player)
		} else {
			c1.outgoingPackets <- packets.TradeTargetAccept(true)
		}
	}
	PacketHandlers["tradeconfirmaccept"] = func(c *Client, p *packets.Packet) {
		if c.player.TradeTarget() == -1 || c.player.State != world.MSTrading || !c.player.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to accept a trade confirmation it was not in!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player['%v'@'%v'] attempted to accept a trade confirmation with a non-existent target!\n", c.player.Username, c.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			return
		}
		if c1.player.State != world.MSTrading || c1.player.TradeTarget() != c.Index || c.player.TradeTarget() != c1.Index || !c1.player.TransAttrs.VarBool("trade1accept", false) {
			log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in trade with apparently bad trade variables!\n", c.player.Username, c.ip, c1.player.Username, c1.ip)
			c.player.ResetTrade()
			c.player.State = world.MSIdle
			c1.player.ResetTrade()
			c1.player.State = world.MSIdle
			c.outgoingPackets <- packets.TradeClose
			c1.outgoingPackets <- packets.TradeClose
			return
		}
		c.player.TransAttrs.SetVar("trade2accept", true)
		if c1.player.TransAttrs.VarBool("trade2accept", false) {
			neededSlots := c1.player.TradeOffer.Size()
			availSlots := c.player.Items.Capacity - c.player.Items.Size() + c.player.TradeOffer.Size()
			theirNeededSlots := c.player.TradeOffer.Size()
			theirAvailSlots := c1.player.Items.Capacity - c1.player.Items.Size() + c1.player.TradeOffer.Size()
			if theirNeededSlots > theirAvailSlots {
				c.Message("The other player does not have room to accept your items.")
				c.player.ResetTrade()
				c.player.State = world.MSIdle
				c1.Message("You do not have room in your inventory to hold those items.")
				c1.player.ResetTrade()
				c1.player.State = world.MSIdle
				c.outgoingPackets <- packets.TradeClose
				c1.outgoingPackets <- packets.TradeClose
				return
			}
			if neededSlots > availSlots {
				c.Message("You do not have room in your inventory to hold those items.")
				c.player.ResetTrade()
				c.player.State = world.MSIdle
				c1.Message("The other player does not have room to accept your items.")
				c1.player.ResetTrade()
				c1.player.State = world.MSIdle
				c.outgoingPackets <- packets.TradeClose
				c1.outgoingPackets <- packets.TradeClose
				return
			}
			defer func() {
				c.outgoingPackets <- packets.InventoryItems(c.player)
				c.outgoingPackets <- packets.TradeClose
				c.player.ResetTrade()
				c.player.State = world.MSIdle
				c1.outgoingPackets <- packets.InventoryItems(c1.player)
				c1.outgoingPackets <- packets.TradeClose
				c1.player.ResetTrade()
				c1.player.State = world.MSIdle
			}()
			if c.player.Items.RemoveAll(c.player.TradeOffer) != c.player.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 1 did not have all items to give.", c.player.Username, c.ip, c1.player.Username, c1.ip)
				return
			}
			if c1.player.Items.RemoveAll(c1.player.TradeOffer) != c1.player.TradeOffer.Size() {
				log.Suspicious.Printf("Players{ 1:['%v'@'%v'];2:['%v'@'%v'] } involved in a trade, player 2 did not have all items to give.", c.player.Username, c.ip, c1.player.Username, c1.ip)
				return
			}
			for i := 0; i < c1.player.TradeOffer.Size(); i++ {
				item := c1.player.TradeOffer.Get(i)
				c.player.Items.Put(item.ID, item.Amount)
			}
			for i := 0; i < c.player.TradeOffer.Size(); i++ {
				item := c.player.TradeOffer.Get(i)
				c1.player.Items.Put(item.ID, item.Amount)
			}
			c.Message("Trade completed.")
			c1.Message("Trade completed.")
		}
	}
	PacketHandlers["duelreq"] = func(c *Client, p *packets.Packet) {
		c.Message("@que@@ora@Not yet implemented")
	}
}