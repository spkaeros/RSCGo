package server

import "bitbucket.org/zlacki/rscgo/pkg/server/packets"

func init() {
	PacketHandlers["pingreq"] = func(c *Client, p *packets.Packet) {
		c.outgoingPackets <- packets.ResponsePong
	}
}
