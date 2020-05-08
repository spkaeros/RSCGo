/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package handshake

import (
	"strings"
	"strconv"
	"time"

	
	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/crypto"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/net"
	`github.com/spkaeros/rscgo/pkg/game/net/handlers`
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	handlers.AddHandler("loginreq", func(player *world.Player, p *net.Packet) {
		loginReply := NewLoginListener(player).attachPlayer(player)
		if !world.UpdateTime.IsZero() {
			loginReply <- response{ResponseLoginServerRejection, "System update in progress"}
			return
		}
		if world.Players.Size() >= config.MaxPlayers() {
			loginReply <- response{ResponseWorldFull, "Out of usable player slots"}
			return
		}
		if loginThrottle.Recent(player.CurrentIP(), time.Minute*5) >= 5 {
			loginReply <- response{ResponseSpamTimeout, "Too many recent invalid login attempts (5 in 5 minutes)"}
			return
		}

		player.SetReconnecting(p.ReadBoolean())
		if ver := p.ReadUint16(); ver != config.Version() {
			loginReply <- response{ResponseUpdated, "Invalid client version (" + strconv.Itoa(ver) + ")"}
			return
		}
		go func() {
			// RSA encryption is expensive; in practice I am uncertain if this is really needed to be isolated
			//  to its own goroutine, but why not since we already launch one for the disk I/O ahead.
			rsaSize := p.ReadUint16()
			data := make([]byte, rsaSize)
			rsaRead := p.Read(data)
			if rsaRead < rsaSize {
				loginReply <- response{ResponseLoginServerRejection, "Invalid RSA block"}
				return
			}
			decryptedP := net.NewPacket(0x0, crypto.DecryptRSA(data))
	//		decryptedP.Skip(1)
	
			// Note: Classic's encryption scheme consisting of an RSA login block to exchange symmetric cipher keys,
			// and after that has been communicated, shifting the packet opcodes using the symmetric ciphers output stream,
			// is deprecated in favor of TLS.  It is simpler and easier to use, while also being an incredibly more secure choice.
			// I am leaving it setup as-is because of 235 protocol compatibility concerns.
			ourSeeds := []int{int(decryptedP.ReadUint32()), int(decryptedP.ReadUint32())}
			theirSeeds := []int{int(decryptedP.ReadUint32()), int(decryptedP.ReadUint32())}
			player.SetVar("ourSeeds", ourSeeds)
			player.SetVar("theirSeeds", theirSeeds)
	
			// this was named linkUID by jagex; it identifys a unique user agent I think
			//p.ReadUint32()
	
			player.Transients().SetVar("username", strutil.Base37.Encode(strings.TrimSpace(decryptedP.ReadString())))
			password := strings.TrimSpace(decryptedP.ReadString())
	//		xteaSize := p.ReadUint16()
	//		data = make([]byte, xteaSize)
	//		xteaRead := p.Read(data)
	//		if xteaRead < xteaSize {
	//			log.Info.Println("Invalid xtea block; the buffer contains:", p.FrameBuffer)
	//			loginReply <- ResponseLoginServerRejection
	//			return
	//		}
	//		keyBuf := make([]byte, 4*4)
	//		for i, v := range append(ourSeeds, theirSeeds...) {
	//			binary.BigEndian.PutUint32(keyBuf[4*i:], v)
	//		}
	//		c, err := xtea.NewCipher(keyBuf)
	//		if err != nil {
	//			log.Info.Println(err)
	//		}
	//		out := make([]byte, xteaSize)
	//		c.Decrypt(out, data)
	//		p = net.NewPacket(0, crypto.DecryptXTEA(xteaSize, data, append(ourSeeds, theirSeeds...)...))
	//		log.Info.Println(password)
	//		player.Transients().SetVar("username", strutil.Base37.Encode(strings.TrimSpace(p.ReadString())))
	
	
			if world.Players.ContainsHash(player.UsernameHash()) {
				loginReply <- response{ResponseLoggedIn, "Player with same username is already logged in"}
				return
			}
			//dataService is a db.PlayerService that all login-related functions should use to access or change player profile data.
			var dataService = db.DefaultPlayerService
			if !dataService.PlayerNameExists(player.Username()) || !dataService.PlayerValidLogin(player.UsernameHash(), crypto.Hash(password)) {
				loginReply <- response{ResponseBadPassword, "Invalid credentials"}
				return
			}
			if !dataService.PlayerLoad(player) {
				loginReply <- response{ResponseDecodeFailure, "Could not load player profile; is the dataService setup properly?"}
				return
			}

			if player.Reconnecting() {
				loginReply <- response{ResponseReconnected, ""}
				return
			}
			switch player.Rank() {
			case 2:
				loginReply <- response{ResponseAdministrator, ""}
			case 1:
				loginReply <- response{ResponseModerator, ""}
			default:
				loginReply <- response{ResponseLoginSuccess, ""}
			}
		}()
	})
	handlers.AddHandler("forgotpass", func(player *world.Player, p *net.Packet) {
		// TODO: These non-login handlers must be isolated and rewrote
		go func() {
			//dataService is a db.PlayerService that all login-related functions should use to access or change player profile data.
			var dataService = db.DefaultPlayerService
			usernameHash := p.ReadUint64()
			if !dataService.PlayerHasRecoverys(usernameHash) {
				player.SendPacket(net.NewReplyPacket([]byte{0}))
				player.Destroy()
				return
			}
			player.SendPacket(net.NewReplyPacket([]byte{1}))
			for _, question := range dataService.PlayerLoadRecoverys(usernameHash) {
				player.SendPacket(net.NewReplyPacket([]byte{byte(len(question))}).AddBytes([]byte(question)))
			}
		}()
	})
}
