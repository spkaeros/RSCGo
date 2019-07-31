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


var handlers = make(map[byte]func(*Client, *packet))

func sessionRequest(c *Client, p *packet) {
	c.uID = p.payload[0]
	p1 := &packet{}
	p1.bare = true
	p1.addLong(rand.GetSecureRandomLong())
	c.writePacket(p1)
}

func loginRequest(c *Client, p *packet) {
	// TODO: RSA decryption, blabla.
	// Currently returns an invalid username or password response
	p1 := &packet{}
	p1.bare = true
	p1.addByte(3)
	c.writePacket(p1)
	c.kill <- struct{}{}
}

func init() {
	handlers[SessionRequest] = sessionRequest
	handlers[LoginRequest] = loginRequest
}

func (c *Client) handlePacket(p *packet) {
	handler, ok := handlers[p.opcode]
	if !ok {
		fmt.Printf("Unhandled packet: {opcode: %d; length: %d};\n", p.opcode, p.length)
		return
	}
	handler(c, p)
}
