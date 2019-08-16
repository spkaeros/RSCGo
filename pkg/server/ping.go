package server

import "bitbucket.org/zlacki/rscgo/pkg/server/packets"

func init() {
	Handlers[5] = ping
}

func ping(c *Client, p *packets.Packet) {
	c.WritePacket(packets.ResponsePong)
}
