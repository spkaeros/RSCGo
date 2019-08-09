package server

import "fmt"

func init() {
	handlers[32] = sessionRequest
	handlers[0] = loginRequest
}

func sessionRequest(c *Client, p *Packet) {
	c.uID = p.ReadByte()
	seed := GenerateSessionID()
	c.isaacSeed[2] = uint32(seed >> 32)
	c.isaacSeed[3] = uint32(seed)
	c.WritePacket(NewBarePacket(nil).AddLong(seed))
}

func loginRequest(c *Client, p *Packet) {
	// TODO: Handle reconnect slightly different
	fmt.Println(p.Payload)
	recon, version := p.ReadByte() == 1, int(p.ReadInt())
	if version != Version {
		LogDebug(1, "WARNING: Player tried logging in with invalid client version. Got %d, expected %d\n", version, Version)
		c.sendLoginResponse(5)
		return
	}
	seed := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		seed[i] = p.ReadInt()
	}
	cipher := c.SeedISAAC(seed)
	if cipher == nil {
		c.sendLoginResponse(8)
		return
	}
	c.isaacStream = cipher
	username, password := p.ReadString(), p.ReadString()
	LogDebug(0, "Registered Player{idx:%v,ip:'%v'username:'%v',password:'%v',reconnecting:%v,version:%v}\n", c.index, c.ip, username, password, recon, version)
	c.sendLoginResponse(0)
}

func (c *Client) sendLoginResponse(i byte) {
	c.WritePacket(NewBarePacket([]byte{i}))
	if i != 0 {
		c.kill <- struct{}{}
	}
}