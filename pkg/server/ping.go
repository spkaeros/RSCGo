package server

func init() {
	handlers[5] = ping
}

func ping(c *Client, p *Packet) {
	c.WritePacket(&Packet{Opcode: 3, Length: 0, Payload: []byte{}})
}