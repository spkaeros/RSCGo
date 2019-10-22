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
		log.Suspicious.Printf("New player request contained invalid lengths: username:'%v'; password:'%v'\n", username, password)
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

	// TODO: Remove all this bs from protocol...
	p.ReadBool()
	p.ReadByte()

	// ISAAC seeds.
	p.ReadLong()
	p.ReadLong()

	// TODO: Remove all this bs from protocol...
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
