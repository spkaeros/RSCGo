package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/rand"
	"fmt"
)

// TODO: Maybe load this from some sort of persistent storage medium, e.g YAML/TOML/JSON file
const (
	LoginRequest   = 0
	SessionRequest = 32
)

var handlers = make(map[byte]func(*Client, *Packet))

func sessionRequest(c *Client, p *Packet) {
	c.uID = p.payload[0]
	p1 := &Packet{bare: true}
	p1.AddLong(rand.GetSecureRandomLong())
	c.WritePacket(p1)
}

func loginRequest(c *Client, p *Packet) {
	// TODO: RSA decryption, blabla.
	// Currently returns an invalid username or password response
	p1 := &Packet{bare: true}
	p1.AddByte(3)
	c.WritePacket(p1)
	c.kill <- struct{}{}
}

func init() {
	handlers[32] = sessionRequest
	handlers[0] = loginRequest
}

func (c *Client) HandlePacket(p *Packet) {
	handler, ok := handlers[p.opcode]
	if !ok {
		fmt.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.opcode, p.length)
		return
	}
	handler(c, p)
}
