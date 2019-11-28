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
	PacketHandlers["forgotpass"] = func(player *world.Player, p *packet.Packet) {
		usernameHash := p.ReadLong()
		if !db.HasRecoveryQuestions(usernameHash) {
			player.SendPacket(packet.NewBarePacket([]byte{0}))
			player.Destroy()
			return
		}
		player.SendPacket(packet.NewBarePacket([]byte{1}))
		for _, question := range db.GetRecoveryQuestions(usernameHash) {
			player.SendPacket(packet.NewBarePacket([]byte{byte(len(question))}).AddBytes([]byte(question)))
		}
	}
	PacketHandlers["cancelpq"] = func(player *world.Player, p *packet.Packet) {
		// empty packet
	}
	PacketHandlers["setpq"] = func(player *world.Player, p *packet.Packet) {
		var questions []string
		var answers []uint64
		for i := 0; i < 5; i++ {
			length := p.ReadByte()
			questions = append(questions, p.ReadString(int(length)))
			answers = append(answers, p.ReadLong())
		}
		log.Info.Println(questions, answers)
	}
	PacketHandlers["changepq"] = func(player *world.Player, p *packet.Packet) {
		player.SendPacket(packet.NewOutgoingPacket(224))
	}
	PacketHandlers["changepass"] = func(player *world.Player, p *packet.Packet) {
		oldPassword := strings.TrimSpace(p.ReadString(20))
		newPassword := strings.TrimSpace(p.ReadString(20))
		if !db.ValidCredentials(player.UserBase37, crypto.Hash(oldPassword)) {
			player.SendPacket(packetbuilders.ServerMessage("The old password you provided does not appear to be valid.  Try again."))
			return
		}
		db.UpdatePassword(player.UserBase37, crypto.Hash(newPassword))
		player.SendPacket(packetbuilders.ServerMessage("Successfully updated your password to the new password you have provided."))
		return
	}
}

func closedConn(player *world.Player, p *packet.Packet) {
	logout(player, p)
}

func logout(player *world.Player, _ *packet.Packet) {
	if player.Busy() {
		player.SendPacket(packetbuilders.CannotLogout)
		return
	}
	if player.Connected() {
		player.SendPacket(packetbuilders.Logout)
		player.Destroy()
	}
}

//handleRegister This method will block until a byte is sent down the reply channel with the registration response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func handleRegister(player *world.Player, reply chan byte) {
	defer player.Destroy()
	defer close(reply)
	select {
	case r := <-reply:
		player.SendPacket(packetbuilders.LoginResponse(int(r)))
		return
	case <-time.After(time.Second * 10):
		player.SendPacket(packetbuilders.LoginResponse(0))
		return
	}
}

func newPlayer(player *world.Player, p *packet.Packet) {
	reply := make(chan byte)
	go handleRegister(player, reply)
	if version := p.ReadShort(); version != config.Version() {
		log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", player.IP, version)
		reply <- 5
		return
	}
	username := strutil.Base37.Decode(strutil.Base37.Encode(strings.TrimSpace(p.ReadString(20))))
	password := strings.TrimSpace(p.ReadString(20))
	if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
		log.Suspicious.Printf("New player request contained invalid lengths: username:'%v'; password:'%v'\n", username, password)
		log.Info.Printf("New player denied: [ Reason:'username or password invalid length'; username='%s'; ip='%s'; passLen=%d ]\n", username, player.IP, passLen)
		reply <- 0
		return
	}
	if db.UsernameExists(username) {
		log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, player.IP)
		reply <- 3
		return
	}

	if db.CreatePlayer(username, password) {
		log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, player.IP)
		reply <- 2
		return
	}
	log.Info.Printf("New player denied: [ Reason:'Most probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, player.IP)
	reply <- 0
	return
}

func sessionRequest(player *world.Player, p *packet.Packet) {
	player.UID = p.ReadByte()
	player.SetServerSeed(rand.Uint64())
	player.SendPacket(packet.NewBarePacket(nil).AddLong(player.ServerSeed()))
}

//initialize
func initialize(player *world.Player) {
	for user := range player.FriendList {
		if players.ContainsHash(user) {
			player.FriendList[user] = true
		}
	}
	world.AddPlayer(player)
	player.Change()
	player.SetConnected(true)
	if player.Skills().Experience(world.StatHits) < 10 {
		for i := 0; i < 18; i++ {
			level := 1
			exp := 0
			if i == 3 {
				level = 10
				exp = 1154
			}
			player.Skills().SetCur(i, level)
			player.Skills().SetMax(i, level)
			player.Skills().SetExp(i, exp)
		}
	}
	if s := time.Until(script.UpdateTime).Seconds(); s > 0 {
		player.SendPacket(packetbuilders.SystemUpdate(int(s)))
	}
	player.SendPacket(packetbuilders.PlaneInfo(player))
	player.SendPacket(packetbuilders.FriendList(player))
	player.SendPacket(packetbuilders.IgnoreList(player))
	if !player.Reconnecting() {
		// Reconnecting implies that the client has all of this data already, so as an optimization, we don't send it again
		player.SendPacket(packetbuilders.PlayerStats(player))
		player.SendPacket(packetbuilders.EquipmentStats(player))
		player.SendPacket(packetbuilders.Fatigue(player))
		player.SendPacket(packetbuilders.InventoryItems(player))
		// TODO: Not canonical RSC, but definitely good QoL update...
		//  player.SendPacket(packetbuilders.FightMode(player)
		player.SendPacket(packetbuilders.ClientSettings(player))
		player.SendPacket(packetbuilders.PrivacySettings(player))
		player.SendPacket(packetbuilders.WelcomeMessage)
		t, err := time.Parse(time.ANSIC, player.Attributes.VarString("lastLogin", time.Time{}.Format(time.ANSIC)))
		if err != nil {
			log.Info.Println(err)
			return
		}

		days := int(time.Since(t).Hours() / 24)
		if t.IsZero() {
			days = 0
		}
		player.Attributes.SetVar("lastLogin", time.Now().Format(time.ANSIC))
		player.SendPacket(packetbuilders.LoginBox(days, player.Attributes.VarString("lastIP", "127.0.0.1")))
	}
	players.BroadcastLogin(player, true)
	if player.FirstLogin() {
		player.SetFirstLogin(false)
		player.AddState(world.MSChangingAppearance)
		player.SendPacket(packetbuilders.ChangeAppearance)
	}
}

//handleLogin This method will block until a byte is sent down the reply channel with the login response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func handleLogin(player *world.Player, reply chan byte) {
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
		player.SendPacket(packetbuilders.LoginResponse(int(r)))
		if isValid(r) {
			players.Put(player)
			log.Info.Printf("Registered: %v\n", player)
			initialize(player)
			return
		}
		log.Info.Printf("Denied: %v (Response='%v')\n", player.String(), r)
		player.Destroy()
		return
	case <-time.After(time.Second * 10):
		player.SendPacket(packetbuilders.LoginResponse(-1))
		return
	}
}

func loginRequest(player *world.Player, p *packet.Packet) {
	loginReply := make(chan byte)
	go handleLogin(player, loginReply)
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
	player.SetReconnecting(p.ReadBool())
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
	player.Username = strutil.Base37.Decode(usernameHash)
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
	go db.LoadPlayer(player, usernameHash, password, loginReply)
}
