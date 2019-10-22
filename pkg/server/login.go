package server

import (
	"strings"

	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["sessionreq"] = sessionRequest
	PacketHandlers["loginreq"] = loginRequest
	PacketHandlers["logoutreq"] = logout
	PacketHandlers["newplayer"] = newPlayer
}

func logout(c *Client, _ *packets.Packet) {
	c.outgoingPackets <- packets.Logout
	c.Destroy()
}

func newPlayer(c *Client, p *packets.Packet) {
	reply := make(chan byte)
	go c.HandleRegister(reply)
	if version := p.ReadShort(); version != config.Version() {
		log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", c.ip, version)
		reply <- 5
		return
	}
	username := strutil.DecodeBase37(strutil.Base37(strings.TrimSpace(p.ReadString(20))))
	password := strings.TrimSpace(p.ReadString(20))
	if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
		// TODO: log it, it's suspicious.  Client should prevent this under normal circumstances.
		log.Info.Printf("New player denied: [ Reason:'username or password invalid length'; username='%s'; ip='%s'; passLen=%d ]\n", username, c.ip, passLen)
		reply <- 0
		return
	}
	if db.UsernameExists(username) {
		log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, c.ip)
		reply <- 3
		return
	}

	if db.CreatePlayer(username, HashPassword(password)) {
		log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, c.ip)
		reply <- 2
		return
	}
	log.Info.Printf("New player denied: [ Reason:'Most probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, c.ip)
	reply <- 0
	return
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
	if ver := p.ReadShort(); ver != config.Version() {
		log.Info.Printf("Invalid client version attempted to login: %d\n", ver)
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
	if !db.UsernameExists(strutil.DecodeBase37(usernameHash)) {
		loginReply <- 3
		return
	}
	if _, ok := Clients.FromUserHash(usernameHash); ok {
		loginReply <- byte(4)
		return
	}
	go db.LoadPlayer(c.player, usernameHash, HashPassword(strings.TrimSpace(p.ReadString(20))), loginReply)
}
