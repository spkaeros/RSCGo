package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["sessionreq"] = sessionRequest
	PacketHandlers["loginreq"] = loginRequest
	PacketHandlers["logoutreq"] = logout
}

func logout(c *Client, p *packets.Packet) {
	c.outgoingPackets <- packets.Logout

	if !c.destroying {
		close(c.Kill)
	}
}

func sessionRequest(c *Client, p *packets.Packet) {
	c.uID = p.ReadByte()
	c.serverSeed = GenerateSessionID()
	c.outgoingPackets <- packets.NewBarePacket(nil).AddLong(c.serverSeed)
}

func loginRequest(c *Client, p *packets.Packet) {
	loginReply := make(chan byte)
	go c.HandleLogin(loginReply)
	// Login block encrypted with block cipher using shared secret, to send/recv credentials and stream cipher key securely
	p.Payload = DecryptRSABlock(p.Payload)
	if p.Payload == nil {
		LogWarning.Println("Could not decrypt RSA login block.")
		loginReply <- byte(9)
		return
	}
	// TODO: Handle reconnect slightly different
	c.reconnecting = p.ReadBool()
	if p.ReadInt() != TomlConfig.Version {
		loginReply <- byte(5)
		return
	}
	c.isaacStream = c.SeedISAAC(p.ReadLong(), p.ReadLong())
	usernameHash := strutil.Base37(p.ReadString())
	if _, ok := Clients[usernameHash]; ok {
		loginReply <- byte(4)
		return
	}
	go c.LoadPlayer(usernameHash, HashPassword(p.ReadString()), loginReply)
}
