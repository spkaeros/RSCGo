package server

import (
	"strings"

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
	c.Destroy()
}

func sessionRequest(c *Client, p *packets.Packet) {
	c.uID = p.ReadByte()
	c.player.SetServerSeed(GenerateSessionID())
	c.outgoingPackets <- packets.NewBarePacket(nil).AddLong(c.player.ServerSeed())
}

func loginRequest(c *Client, p *packets.Packet) {
	loginReply := make(chan byte)
	go c.HandleLogin(loginReply)
	// Login block encrypted with block cipher using shared secret, to send/recv credentials and stream cipher key securely
	// TODO: Re-enable RSA for 204 once JS implementation exists...
	/*
		p.Payload = DecryptRSABlock(p.Payload)
		if p.Payload == nil {
			LogWarning.Println("Could not decrypt RSA login block.")
			loginReply <- byte(9)
			return
		}
	*/
	c.player.SetReconnecting(p.ReadBool())
	// TODO: 204 sends version as short.
	/*
		if p.ReadInt() != TomlConfig.Version {
			loginReply <- byte(5)
			return
		}
	*/
	if ver := p.ReadShort(); ver != TomlConfig.Version {
		LogInfo.Printf("Invalid client version attempted to login: %d\n", ver)
		loginReply <- byte(5)
		return
	}

	// TODO: 204 sends what looks like a boolean named limit30.  WTF is it?
	p.ReadBool()

	// FIXME: 204 sends 0xA(10) as one byte before the ISAAC seeds.  It is always 10.  Why?
	p.ReadByte()

	// ISAAC seeds.
	p.ReadLong()
	p.ReadLong()

	// FIXME: 204 sends an int32 calling getLinkUID for its value.  Figure out what it is used for, seems to be more applet bullshit.
	p.ReadInt()

	usernameHash := strutil.Base37(strings.TrimSpace(p.ReadString(20)))
	if _, ok := Clients.FromUserHash(usernameHash); ok {
		loginReply <- byte(4)
		return
	}
	go c.LoadPlayer(usernameHash, HashPassword(strings.TrimSpace(p.ReadString(20))), loginReply)
}
