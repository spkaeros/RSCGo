package server

import "bitbucket.org/zlacki/rscgo/pkg/server/packets"

func init() {
	Handlers[32] = sessionRequest
	Handlers[0] = loginRequest
}

func sessionRequest(c *Client, p *packets.Packet) {
	c.uID, _ = p.ReadByte()
	seed := GenerateSessionID()
	c.isaacSeed[2] = uint32(seed >> 32)
	c.isaacSeed[3] = uint32(seed)
	c.WritePacket(packets.NewBarePacket(nil).AddLong(seed))
}

func loginRequest(c *Client, p *packets.Packet) {
	// TODO: Handle reconnect slightly different
	recon, _ := p.ReadByte()
	version, _ := p.ReadInt()
	if version != uint32(Version) {
		if len(Flags.Verbose) >= 1 {
			LogWarning.Printf("WARNING: Player tried logging in with invalid client version. Got %d, expected %d\n", version, Version)
		}
		c.sendLoginResponse(5)
		return
	}
	seed := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		seed[i], _ = p.ReadInt()
	}
	cipher := c.SeedISAAC(seed)
	if cipher == nil {
		c.sendLoginResponse(8)
		return
	}
	c.isaacStream = cipher
	username, _ := p.ReadString()
	password, _ := p.ReadString()
	LogInfo.Printf("Registered Player{idx:%v,ip:'%v'username:'%v',password:'%v',reconnecting:%v,version:%v}\n", c.index, c.ip, username, password, recon, version)
	c.sendLoginResponse(0)
}
