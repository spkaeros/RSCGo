package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["tradereq"] = func(c *Client, p *packets.Packet) {
		index := p.ReadShort()
		if affectedClient, ok := Clients.FromIndex(index); ok {
			if !c.player.WithinRange(&affectedClient.player.Location, 16) || c.player.Busy() {
				return
			}
			c.player.SetTradeTarget(index)
			c.Message("Sending trade request.")
			affectedClient.Message(c.player.Username + " wishes to trade with you.")
			if affectedClient.player.TradeTarget() == c.Index {
				c.player.State = world.MSTrading
				c.player.ResetPath()
				c.TradeOpen()
				c.player.TransAttrs.SetVar("trading", true)
				affectedClient.player.State = world.MSTrading
				affectedClient.player.ResetPath()
				affectedClient.TradeOpen()
				affectedClient.player.TransAttrs.SetVar("trading", true)
			}
		}
	}
	PacketHandlers["tradeupdate"] = func(c *Client, p *packets.Packet) {
		if c.player.TradeTarget() == -1 || c.player.State != world.MSTrading {
			log.Suspicious.Printf("Player[%v] attempted to decline a trade it was not in!\n", c)
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player[%v] attempted to update a trade with a non-existent target!\n", c)
			return
		}
		c.player.TransAttrs.UnsetVar("trade1accept")
		c1.player.TransAttrs.UnsetVar("trade1accept")
		c.player.TradeOffer.Clear()
		offerLen := int(p.ReadByte())
		var ids []int
		var amts []int
		for i := 0; i < offerLen; i++ {
			id := p.ReadShort()
			amt := p.ReadInt()
			ids = append(ids, id)
			amts = append(amts, amt)
			c.player.TradeOffer.Put(id, amt)
		}
		c1.outgoingPackets <- packets.TradeUpdate(offerLen, ids, amts)
	}
	PacketHandlers["tradedecline"] = func(c *Client, p *packets.Packet) {
		if c.player.TradeTarget() == -1 || c.player.State != world.MSTrading {
			log.Suspicious.Printf("Player[%v] attempted to decline a trade it was not in!\n", c)
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player[%v] attempted to decline a trade with a non-existent target!\n", c)
			return
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
			log.Suspicious.Printf("Player[%v] attempted to decline a trade it was not in!\n", c)
			return
		}
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if !ok {
			log.Suspicious.Printf("Player[%v] attempted to decline a trade with a non-existent target!\n", c)
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
	PacketHandlers["duelreq"] = func(c *Client, p *packets.Packet) {
		c.Message("@que@@ora@Not yet implemented")
	}
}