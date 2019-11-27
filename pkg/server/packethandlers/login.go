package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/crypto"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"strings"
	"time"

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
	PacketHandlers["forgotpass"] = func(c *world.Player, p *packet.Packet) {
		usernameHash := p.ReadLong()
		if !db.HasRecoveryQuestions(usernameHash) {
			c.SendPacket(packet.NewBarePacket([]byte{0}))
			c.Destroy()
			return
		}
		c.SendPacket(packet.NewBarePacket([]byte{1}))
		for _, question := range db.GetRecoveryQuestions(usernameHash) {
			c.SendPacket(packet.NewBarePacket([]byte{byte(len(question))}).AddBytes([]byte(question)))
		}
	}
	PacketHandlers["cancelpq"] = func(c *world.Player, p *packet.Packet) {
		// empty packet
	}
	PacketHandlers["setpq"] = func(c *world.Player, p *packet.Packet) {
		var questions []string
		var answers []uint64
		for i := 0; i < 5; i++ {
			length := p.ReadByte()
			questions = append(questions, p.ReadString(int(length)))
			answers = append(answers, p.ReadLong())
		}
		log.Info.Println(questions, answers)
	}
	PacketHandlers["changepq"] = func(c *world.Player, p *packet.Packet) {
		c.SendPacket(packet.NewOutgoingPacket(224))
	}
	PacketHandlers["changepass"] = func(c *world.Player, p *packet.Packet) {
		oldPassword := strings.TrimSpace(p.ReadString(20))
		newPassword := strings.TrimSpace(p.ReadString(20))
		if !db.ValidCredentials(c.UserBase37, crypto.Hash(oldPassword)) {
			c.SendPacket(packetbuilders.ServerMessage("The old password you provided does not appear to be valid.  Try again."))
			return
		}
		db.UpdatePassword(c.UserBase37, crypto.Hash(newPassword))
		c.SendPacket(packetbuilders.ServerMessage("Successfully updated your password to the new password you have provided."))
		return
	}
}

func closedConn(c *world.Player, p *packet.Packet) {
	logout(c, p)
}

func logout(c *world.Player, _ *packet.Packet) {
	if c.Busy() {
		c.SendPacket(packetbuilders.CannotLogout)
		return
	}
	if c.Connected() {
		c.SendPacket(packetbuilders.Logout)
		c.Destroy()
	}
}

//handleRegister This method will block until a byte is sent down the reply channel with the registration response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func handleRegister(c *world.Player, reply chan byte) {
	defer c.Destroy()
	defer close(reply)
	select {
	case r := <-reply:
		c.SendPacket(packetbuilders.LoginResponse(int(r)))
		return
	case <-time.After(time.Second * 10):
		c.SendPacket(packetbuilders.LoginResponse(0))
		return
	}
}

func newPlayer(c *world.Player, p *packet.Packet) {
	reply := make(chan byte)
	go handleRegister(c, reply)
	if version := p.ReadShort(); version != config.Version() {
		log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", c.IP, version)
		reply <- 5
		return
	}
	username := strutil.Base37.Decode(strutil.Base37.Encode(strings.TrimSpace(p.ReadString(20))))
	password := strings.TrimSpace(p.ReadString(20))
	if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
		log.Suspicious.Printf("New player request contained invalid lengths: username:'%v'; password:'%v'\n", username, password)
		log.Info.Printf("New player denied: [ Reason:'username or password invalid length'; username='%s'; ip='%s'; passLen=%d ]\n", username, c.IP, passLen)
		reply <- 0
		return
	}
	if db.UsernameExists(username) {
		log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, c.IP)
		reply <- 3
		return
	}

	if db.CreatePlayer(username, password) {
		log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, c.IP)
		reply <- 2
		return
	}
	log.Info.Printf("New player denied: [ Reason:'Most probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, c.IP)
	reply <- 0
	return
}

func sessionRequest(c *world.Player, p *packet.Packet) {
	c.UID = p.ReadByte()
	c.SetServerSeed(rand.Uint64())
	c.SendPacket(packet.NewBarePacket(nil).AddLong(c.ServerSeed()))
}

//initialize
func initialize(c *world.Player) {
	for user := range c.FriendList {
		if players.ContainsHash(user) {
			c.FriendList[user] = true
		}
	}
	world.AddPlayer(c)
	c.Change()
	c.SetConnected(true)
	if c.Skills().Experience(world.StatHits) < 10 {
		for i := 0; i < 18; i++ {
			level := 1
			exp := 0
			if i == 3 {
				level = 10
				exp = 1154
			}
			c.Skills().SetCur(i, level)
			c.Skills().SetMax(i, level)
			c.Skills().SetExp(i, exp)
		}
	}
	if s := time.Until(script.UpdateTime).Seconds(); s > 0 {
		c.SendPacket(packetbuilders.SystemUpdate(int(s)))
	}
	c.SendPacket(packetbuilders.PlaneInfo(c))
	c.SendPacket(packetbuilders.FriendList(c))
	c.SendPacket(packetbuilders.IgnoreList(c))
	if !c.Reconnecting() {
		// Reconnecting implies that the client has all of this data already, so as an optimization, we don't send it again
		c.SendPacket(packetbuilders.PlayerStats(c))
		c.SendPacket(packetbuilders.EquipmentStats(c))
		c.SendPacket(packetbuilders.Fatigue(c))
		c.SendPacket(packetbuilders.InventoryItems(c))
		// TODO: Not canonical RSC, but definitely good QoL update...
		//  c.SendPacket(packetbuilders.FightMode(c)
		c.SendPacket(packetbuilders.ClientSettings(c))
		c.SendPacket(packetbuilders.PrivacySettings(c))
		c.SendPacket(packetbuilders.WelcomeMessage)
		t, err := time.Parse(time.ANSIC, c.Attributes.VarString("lastLogin", time.Time{}.Format(time.ANSIC)))
		if err != nil {
			log.Info.Println(err)
			return
		}

		days := int(time.Since(t).Hours() / 24)
		if t.IsZero() {
			days = 0
		}
		c.Attributes.SetVar("lastLogin", time.Now().Format(time.ANSIC))
		c.SendPacket(packetbuilders.LoginBox(days, c.Attributes.VarString("lastIP", "127.0.0.1")))
	}
	players.BroadcastLogin(c, true)
	if c.FirstLogin() {
		c.SetFirstLogin(false)
		c.AddState(world.MSChangingAppearance)
		c.SendPacket(packetbuilders.ChangeAppearance)
	}
}

//handleLogin This method will block until a byte is sent down the reply channel with the login response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func handleLogin(c *world.Player, reply chan byte) {
	isValid := func(r byte) bool {
		valid := [...]byte{0, 1, 24, 25}
		for _, i := range valid {
			if i == r {
				return true
			}
		}
		return false
	}
	defer close(reply)
	select {
	case r := <-reply:
		c.SendPacket(packetbuilders.LoginResponse(int(r)))
		if isValid(r) {
			players.Put(c)
			log.Info.Printf("Registered: %v\n", c)
			initialize(c)
			return
		}
		log.Info.Printf("Denied Client: {IP:'%v', username:'%v', Response='%v'}\n", c.IP, c.Username, r)
		c.Destroy()
		return
	case <-time.After(time.Second * 10):
		c.SendPacket(packetbuilders.LoginResponse(-1))
		return
	}
}

func loginRequest(c *world.Player, p *packet.Packet) {
	loginReply := make(chan byte)
	go handleLogin(c, loginReply)
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
	c.SetReconnecting(p.ReadBool())
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
	c.Username = strutil.Base37.Decode(usernameHash)
	password := strings.TrimSpace(p.ReadString(20))
	if !db.UsernameExists(strutil.Base37.Decode(usernameHash)) {
		loginReply <- 3
		return
	}
	if _, ok := players.FromUserHash(usernameHash); ok {
		loginReply <- byte(4)
		return
	}
	if !script.UpdateTime.IsZero() && time.Until(script.UpdateTime).Seconds() <= 0 {
		loginReply <- 8
		return
	}
	go db.LoadPlayer(c, usernameHash, password, loginReply)
}
