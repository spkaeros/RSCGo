package server

import "bitbucket.org/zlacki/rscgo/pkg/server/packets"

func init() {
	Handlers[5] = func(c *Client, p *packets.Packet) {
		c.outgoingPackets <- packets.ResponsePong
	}
}
