/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/crypto"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"strings"
	"time"

	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/log"
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
		if !db.ValidCredentials(player.UsernameHash(), crypto.Hash(oldPassword)) {
			player.Message("The old password you provided does not appear to be valid.  Try again.")
			return
		}
		db.UpdatePassword(player.UsernameHash(), crypto.Hash(newPassword))
		player.Message("Successfully updated your password to the new password you have provided.")
		return
	}
}

func closedConn(player *world.Player, p *packet.Packet) {
	logout(player, p)
}

func logout(player *world.Player, _ *packet.Packet) {
	if player.Busy() {
		player.SendPacket(world.CannotLogout)
		return
	}
	if player.Connected() {
		player.SendPacket(world.Logout)
		player.Destroy()
	}
}

//handleRegister This method will block until a byte is sent down the reply channel with the registration response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func handleRegister(player *world.Player, reply chan byte) {
	defer player.Destroy()
	defer close(reply)
	select {
	case r := <-reply:
		player.SendPacket(world.LoginResponse(int(r)))
		return
	case <-time.After(time.Second * 10):
		player.SendPacket(world.LoginResponse(0))
		return
	}
}

func newPlayer(player *world.Player, p *packet.Packet) {
	reply := make(chan byte)
	go handleRegister(player, reply)
	if version := p.ReadShort(); version != config.Version() {
		log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", player.CurrentIP(), version)
		reply <- 5
		return
	}
	username := strutil.Base37.Decode(strutil.Base37.Encode(strings.TrimSpace(p.ReadString(20))))
	password := strings.TrimSpace(p.ReadString(20))
	if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
		log.Suspicious.Printf("New player request contained invalid lengths: username:'%v'; password:'%v'\n", username, password)
		log.Info.Printf("New player denied: [ Reason:'username or password invalid length'; username='%s'; ip='%s'; passLen=%d ]\n", username, player.CurrentIP(), passLen)
		reply <- 0
		return
	}
	if db.UsernameExists(username) {
		log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, player.CurrentIP())
		reply <- 3
		return
	}

	if db.CreatePlayer(username, password) {
		log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, player.CurrentIP())
		reply <- 2
		return
	}
	log.Info.Printf("New player denied: [ Reason:'Most probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, player.CurrentIP())
	reply <- 0
	return
}

func sessionRequest(player *world.Player, p *packet.Packet) {
	player.SetConnected(true)
	p.ReadByte() // UID, useful?
	player.SetServerSeed(rand.Uint64())
	player.SendPacket(packet.NewBarePacket(nil).AddLong(player.ServerSeed()))
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
		player.OutgoingPackets <- world.LoginResponse(int(r))
		if isValid(r) {
			world.Players.Put(player)
			world.Players.BroadcastLogin(player, true)
			player.Initialize()
			for _, fn := range script.LoginTriggers {
				fn(player)
			}
			log.Info.Printf("Registered: %v\n", player)
			return
		}
		log.Info.Printf("Denied: %v (Response='%v')\n", player.String(), r)
		player.Destroy()
		return
	case <-time.After(time.Second * 10):
		player.SendPacket(world.LoginResponse(-1))
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

	// TODO: SetRegionRemoved all this bs from protocol...
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
	player.TransAttrs.SetVar("username", usernameHash)
	password := strings.TrimSpace(p.ReadString(20))
	if !db.UsernameExists(strutil.Base37.Decode(usernameHash)) {
		loginReply <- 3
		return
	}
	if _, ok := world.Players.FromUserHash(usernameHash); ok {
		loginReply <- byte(4)
		return
	}
	if !world.UpdateTime.IsZero() && time.Until(world.UpdateTime).Seconds() <= 0 {
		loginReply <- 8
		return
	}
	go db.LoadPlayer(player, usernameHash, password, loginReply)
}
