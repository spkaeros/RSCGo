/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package handlers

import (
	"strings"
	"time"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/crypto"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/login"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/tasks"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	AddHandler("forgotpass", func(player *world.Player, p *net.Packet) {
		// TODO: These non-login handlers must be isolated and rewrote
		go func() {
			usernameHash := p.ReadLong()
			if !db.DefaultPlayerService.PlayerHasRecoverys(usernameHash) {
				player.SendPacket(net.NewBarePacket([]byte{0}))
				player.Destroy()
				return
			}
			player.SendPacket(net.NewBarePacket([]byte{1}))
			for _, question := range db.DefaultPlayerService.PlayerLoadRecoverys(usernameHash) {
				player.SendPacket(net.NewBarePacket([]byte{byte(len(question))}).AddBytes([]byte(question)))
			}
		}()
	})
	AddHandler("loginreq", func(player *world.Player, p *net.Packet) {
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
		loginReply := login.NewLoginListener(player).ResponseListener()
		if login.Throttler.Recent(player.CurrentIP(), time.Minute*5) >= 5 {
			loginReply <- login.ResponseSpamTimeout
			return
		}
		player.SetReconnecting(p.ReadBool())
		if ver := p.ReadShort(); ver != config.Version() {
			log.Info.Printf("Invalid client version attempted to login: %d\n", ver)
			loginReply <- login.ResponseUpdated
			return
		}

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
		player.FriendList.Owner = player.Username()
		password := strings.TrimSpace(p.ReadString(20))
		if _, ok := world.Players.FromUserHash(usernameHash); ok {
			loginReply <- login.ResponseLoggedIn
			return
		}
		if !world.UpdateTime.IsZero() && time.Until(world.UpdateTime).Seconds() <= 0 {
			loginReply <- login.ResponseLoginServerRejection
			return
		}

		db.DefaultPlayerService.PlayerLoad(player, usernameHash, password, loginReply)
	})
	AddHandler("logoutreq", func(player *world.Player, p *net.Packet) {
		tasks.TickerList.Add("playerDestroy", func() bool {
			if player.Busy() {
				player.SendPacket(world.CannotLogout)
				return true
			}
			if player.Connected() {
				player.Destroy()
			}
			return true
		})
	})
	AddHandler("closeconn", func(player *world.Player, p *net.Packet) {
		if player.Busy() {
			log.Suspicious.Println("CLOSECONN!!", player, p.String(), player.State())
			log.Info.Println("CLOSECONN!!", player, p.String(), player.State())
			player.SendPacket(world.CannotLogout)
			return
		}
		if player.Connected() {
			player.Destroy()
		}
	})
	AddHandler("cancelpq", func(player *world.Player, p *net.Packet) {
		// empty net
	})
	AddHandler("setpq", func(player *world.Player, p *net.Packet) {
		var questions []string
		var answers []uint64
		for i := 0; i < 5; i++ {
			length := p.ReadByte()
			questions = append(questions, p.ReadString(int(length)))
			answers = append(answers, p.ReadLong())
		}
		log.Info.Println(questions, answers)
	})
	AddHandler("changepq", func(player *world.Player, p *net.Packet) {
		player.SendPacket(net.NewOutgoingPacket(224))
	})
	AddHandler("changepass", func(player *world.Player, p *net.Packet) {
		oldPassword := strings.TrimSpace(p.ReadString(20))
		newPassword := strings.TrimSpace(p.ReadString(20))
		go func() {
			if !db.DefaultPlayerService.PlayerValidLogin(player.UsernameHash(), crypto.Hash(oldPassword)) {
				player.Message("The old password you provided does not appear to be valid.  Try again.")
				return
			}
			db.DefaultPlayerService.PlayerChangePassword(player.UsernameHash(), crypto.Hash(newPassword))
			player.Message("Successfully updated your password to the new password you have provided.")
		}()
	})
	AddHandler("newplayer", func(player *world.Player, p *net.Packet) {
		reply := login.NewRegistrationListener(player).ResponseListener()
		if login.RegisterThrottler.Recent(player.CurrentIP(), time.Minute*5) >= 5 {
			reply <- login.ResponseSpamTimeout
			return
		}
		if version := p.ReadShort(); version != config.Version() {
			log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", player.CurrentIP(), version)
			reply <- login.ResponseUpdated
			return
		}
		username := strutil.Base37.Decode(strutil.Base37.Encode(strings.TrimSpace(p.ReadString(20))))
		password := strings.TrimSpace(p.ReadString(20))
		if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
			log.Suspicious.Printf("New player request contained invalid lengths: username:'%v'; password:'%v'\n", username, password)
			log.Info.Printf("New player denied: [ Reason:'username or password invalid length'; username='%s'; ip='%s'; passLen=%d ]\n", username, player.CurrentIP(), passLen)
			reply <- 17
			return
		}
		go func() {
			if db.DefaultPlayerService.PlayerNameTaken(username) {
				log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, player.CurrentIP())
				reply <- login.ResponseUsernameTaken
				return
			}

			if db.DefaultPlayerService.PlayerCreate(username, password) {
				log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, player.CurrentIP())
				reply <- login.ResponseRegisterSuccess
				return
			}
			log.Info.Printf("New player denied: [ Reason:'Most probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, player.CurrentIP())
			reply <- -1
		}()
	})
}
