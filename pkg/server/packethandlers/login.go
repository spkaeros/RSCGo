package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/crypto"
	"strings"

	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["sessionreq"] = sessionRequest
	PacketHandlers["loginreq"] = loginRequest
	PacketHandlers["logoutreq"] = logout
	PacketHandlers["closeconn"] = closedConn
	PacketHandlers["newplayer"] = newPlayer
	PacketHandlers["forgotpass"] = func(c clients.Client, p *packetbuilders.Packet) {
		usernameHash := p.ReadLong()
		if !db.HasRecoveryQuestions(usernameHash) {
			c.SendPacket(packetbuilders.NewBarePacket([]byte{0}))
			c.Destroy()
			return
		}
		c.SendPacket(packetbuilders.NewBarePacket([]byte{1}))
		for _, question := range db.GetRecoveryQuestions(usernameHash) {
			c.SendPacket(packetbuilders.NewBarePacket([]byte{byte(len(question))}).AddBytes([]byte(question)))
		}
	}
	PacketHandlers["cancelpq"] = func(c clients.Client, p *packetbuilders.Packet) {
		// empty packet
	}
	PacketHandlers["setpq"] = func(c clients.Client, p *packetbuilders.Packet) {
		var questions []string
		var answers []uint64
		for i := 0; i < 5; i++ {
			length := p.ReadByte()
			questions = append(questions, p.ReadString(int(length)))
			answers = append(answers, p.ReadLong())
		}
		log.Info.Println(questions, answers)
	}
	PacketHandlers["changepq"] = func(c clients.Client, p *packetbuilders.Packet) {
		c.SendPacket(packetbuilders.NewOutgoingPacket(224))
	}
	PacketHandlers["changepass"] = func(c clients.Client, p *packetbuilders.Packet) {
		oldPassword := strings.TrimSpace(p.ReadString(20))
		newPassword := strings.TrimSpace(p.ReadString(20))
		if !db.ValidCredentials(c.Player().UserBase37, crypto.Hash(oldPassword)) {
			c.Message("The old password you provided does not appear to be valid.  Try again.")
			return
		}
		db.UpdatePassword(c.Player().UserBase37, crypto.Hash(newPassword))
		c.Message("Successfully updated your password to the new password you have provided.")
		return
	}
}

func closedConn(c clients.Client, p *packetbuilders.Packet) {
	logout(c, p)
}

func logout(c clients.Client, _ *packetbuilders.Packet) {
	if c.Player().TransAttrs.VarBool("connected", false) {
		c.SendPacket(packetbuilders.Logout)
		c.Destroy()
	}
}

func newPlayer(c clients.Client, p *packetbuilders.Packet) {
	reply := make(chan byte)
	go c.HandleRegister(reply)
	if version := p.ReadShort(); version != config.Version() {
		log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", c.Player().IP, version)
		reply <- 5
		return
	}
	username := strutil.Base37.Decode(strutil.Base37.Encode(strings.TrimSpace(p.ReadString(20))))
	password := strings.TrimSpace(p.ReadString(20))
	if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
		log.Suspicious.Printf("New player request contained invalid lengths: username:'%v'; password:'%v'\n", username, password)
		log.Info.Printf("New player denied: [ Reason:'username or password invalid length'; username='%s'; ip='%s'; passLen=%d ]\n", username, c.Player().IP, passLen)
		reply <- 0
		return
	}
	if db.UsernameExists(username) {
		log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, c.Player().IP)
		reply <- 3
		return
	}

	if db.CreatePlayer(username, password) {
		log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, c.Player().IP)
		reply <- 2
		return
	}
	log.Info.Printf("New player denied: [ Reason:'Most probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, c.Player().IP)
	reply <- 0
	return
}

func sessionRequest(c clients.Client, p *packetbuilders.Packet) {
	c.Player().UID = p.ReadByte()
	c.Player().SetServerSeed(rand.Uint64())
	c.SendPacket(packetbuilders.NewBarePacket(nil).AddLong(c.Player().ServerSeed()))
}

func loginRequest(c clients.Client, p *packetbuilders.Packet) {
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
	c.Player().SetReconnecting(p.ReadBool())
	if ver := p.ReadShort(); ver != config.Version() {
		log.Info.Printf("Invalid client version attempted to login: %d\n", ver)
		loginReply <- byte(5)
		return
	}

	// TODO: Remove all this bs from protocol...
	p.ReadBool() // limit30
	p.ReadByte() // 0xA.  Some sort of separator I think?

	// ISAAC seeds.
	p.ReadLong()
	p.ReadLong()

	// TODO: Remove all this bs from protocol...
	//  getLinkUID--Jagex used this as a means of identification
	//  it was a random var read from the RS cache to help identify individuals and assist in cheat detection
	//  My understanding is that this is exactly what they used to trigger the too many accounts logged in reply,
	//  hence why running unsigned client back in the day, with its own temp RS cache, allowed you to login anyways
	p.ReadInt()

	usernameHash := strutil.Base37.Encode(strings.TrimSpace(p.ReadString(20)))
	c.Player().Username = strutil.Base37.Decode(usernameHash)
	password := strings.TrimSpace(p.ReadString(20))
	if !db.UsernameExists(strutil.Base37.Decode(usernameHash)) {
		loginReply <- 3
		return
	}
	if _, ok := clients.FromUserHash(usernameHash); ok {
		loginReply <- byte(4)
		return
	}
	go db.LoadPlayer(c.Player(), usernameHash, password, loginReply)
}
