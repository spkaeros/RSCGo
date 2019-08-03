package server

import (
	rscrand "bitbucket.org/zlacki/rscgo/pkg/rand"
	"fmt"
	"strings"
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
	p1.AddLong(rscrand.GetSecureRandomLong())
	c.WritePacket(p1)
}

func loginRequest(c *Client, p *Packet) {
	response := &Packet{bare: true}
	if !p.DecryptRSA() {
		response.AddByte(17)
		c.WritePacket(response)
		c.kill <- struct{}{}
	}
	recon := p.ReadByte() == 1
	version := p.ReadInt()
	c.decryptKey = p.ReadLong()
	c.encryptKey = p.ReadLong()
	username := strings.TrimSpace(p.ReadString(20))
	password := strings.TrimSpace(p.ReadString(20))
	fmt.Printf("reconnecting:%v,version:%v,clientSeed:%v,serverSeed:%v,username:%v,password:%v\n", recon, version, c.decryptKey, c.encryptKey, username, password)
	response.AddByte(0)
	c.WritePacket(response)
}

func ping(c *Client, p *Packet) {
	c.WritePacket(&Packet{opcode: 3, length: 0, payload: []byte{}})
}

func init() {
	handlers[32] = sessionRequest
	handlers[0] = loginRequest
	handlers[5] = ping
}

func (c *Client) HandlePacket(p *Packet) {
	handler, ok := handlers[p.opcode]
	if !ok {
		fmt.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.opcode, p.length)
		return
	}
	handler(c, p)
}
