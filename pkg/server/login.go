package server

import (
	"crypto/rand"
	"crypto/rsa"

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
	close(c.kill)
	//c.kill <- struct{}{}
}

func sessionRequest(c *Client, p *packets.Packet) {
	c.uID, _ = p.ReadByte()
	seed := GenerateSessionID()
	c.isaacSeed[1] = seed
	c.outgoingPackets <- packets.NewBarePacket(nil).AddLong(seed)
}

func loginRequest(c *Client, p *packets.Packet) {
	buf, err := rsa.DecryptPKCS1v15(rand.Reader, RsaKey, p.Payload)
	if err != nil {
		LogWarning.Printf("Could not decrypt RSA login block: `%v`\n", err.Error())
		c.sendLoginResponse(9)
		return
	}
	p.Payload = buf
	player := c.player
	// Login block encrypted with block cipher using shared secret, to send/recv credentials and stream cipher key
	// TODO: Handle reconnect slightly different
	c.reconnecting, err = p.ReadBool()
	if err != nil {
		c.sendLoginResponse(6)
		return
	}
	version, err := p.ReadInt()
	if err != nil {
		c.sendLoginResponse(6)
		return
	}
	if int(version) != TomlConfig.Version {
		c.sendLoginResponse(5)
		return
	}
	clientSeed, err := p.ReadLong()
	if err != nil {
		c.sendLoginResponse(6)
		return
	}
	serverSeed, err := p.ReadLong()
	if err != nil {
		c.sendLoginResponse(6)
		return
	}
	c.isaacStream = c.SeedISAAC(clientSeed, serverSeed)
	username, err := p.ReadString()
	if err != nil {
		c.sendLoginResponse(6)
		return
	}
	player.UserBase37 = strutil.Base37(username)
	player.Username = strutil.DecodeBase37(player.UserBase37)
	if _, ok := Clients[player.UserBase37]; ok {
		c.sendLoginResponse(4)
		return
	}
	password, err := p.ReadString()
	if err != nil {
		c.sendLoginResponse(6)
		return
	}
	c.sendLoginResponse(byte(c.LoadPlayer(player.UserBase37, HashPassword(password))))
}
