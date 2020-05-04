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
	"strconv"
	"strings"
	"time"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/crypto"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/net"
	`github.com/spkaeros/rscgo/pkg/game/net/handlers`
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	handlers.AddHandler("newplayer", func(player *world.Player, p *net.Packet) {
		player.SetConnected(true)

		reply := NewRegistrationListener(player).ResponseListener()
		username := "nil"
		sendReply := func(code ResponseCode, reason string) {
			log.Info.Printf("New player denied: [ Reason:'%s'; username='%s'; ip='%s'; response=%d; ]\n", reason, username, player.CurrentIP(), code)
			reply <- ResponseCode(code)
		}
		if RegisterThrottle.Recent(player.CurrentIP(), time.Hour) >= 2 {
			sendReply(ResponseSpamTimeout, "Recently registered too many other characters")
			return
		}
		if version := p.ReadUint16(); version != config.Version() {
			sendReply(ResponseUpdated, "Client version (" + strconv.Itoa(version) + ") out of date")
			return
		}
		username = strutil.Base37.Decode(p.ReadUint64())
		password := strings.TrimSpace(p.ReadString())
		player.Transients().SetVar("username", username)
		if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
			sendReply(ResponseBadInputLength, "Password and/or username too long and/or too short.")
			return
		}
		go func() {
			dataService := db.DefaultPlayerService
			if dataService.PlayerNameExists(username) {
				sendReply(ResponseUsernameTaken, "Username is taken")
				return
			}

			if !dataService.PlayerCreate(username, crypto.Hash(password), player.CurrentIP()) {
				sendReply(-1, "Unknown issue during PlayerCreate call")
				return
			}
			log.Info.Printf("New player created: [ username='%s'; ip='%s' ]", username, player.CurrentIP())
			reply <- ResponseRegisterSuccess
		}()
	})
}
