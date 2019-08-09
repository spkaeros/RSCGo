package server

func init() {
	handlers[5] = ping
}

func ping(c *Client, p *Packet) {
	c.WritePacket(NewOutgoingPacket(3))
}