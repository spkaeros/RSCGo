package server

import (
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
				affectedClient.player.State = world.MSTrading
				affectedClient.player.ResetPath()
				affectedClient.TradeOpen()
			}
		}
	}
	PacketHandlers["tradeupdate"] = func(c *Client, p *packets.Packet) {
		c1, ok := Clients.FromIndex(c.player.TradeTarget())
		if ok {
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
	}
	PacketHandlers["duelreq"] = func(c *Client, p *packets.Packet) {
		c.Message("@que@@ora@Not yet implemented")
	}
}