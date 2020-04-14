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
	"time"

	"github.com/spkaeros/rscgo/pkg/config"
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
		if RegisterThrottle.Recent(player.CurrentIP(), time.Hour) >= 2 {
			reply <- ResponseSpamTimeout
			return
		}
		if version := p.ReadUint16(); version != config.Version() {
			log.Info.Printf("New player denied: [ Reason:'Wrong client version'; ip='%s'; version=%d ]\n", player.CurrentIP(), version)
			reply <- ResponseUpdated
			return
		}
		username := strutil.Base37.Decode(p.ReadUint64())
		password := strings.TrimSpace(p.ReadString())
		player.Transients().SetVar("username", username)
		if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
			log.Suspicious.Printf("New player request contained invalid lengths: %v username=%v; password:'%v'\n", player.CurrentIP(), username, password)
			reply <- ResponseBadInputLength
			return
		}
		go func() {
			dataService := db.DefaultPlayerService
			if dataService.PlayerNameExists(username) {
				log.Info.Printf("New player denied: [ Reason:'Username is taken'; username='%s'; ip='%s' ]\n", username, player.CurrentIP())
				reply <- ResponseUsernameTaken
				return
			}

			if dataService.PlayerCreate(username, password, player.CurrentIP()) {
				log.Info.Printf("New player accepted: [ username='%s'; ip='%s' ]", username, player.CurrentIP())
				reply <- ResponseRegisterSuccess
				return
			}
			log.Info.Printf("New player denied: [ Reason:'unknown; probably database related.  Debug required'; username='%s'; ip='%s' ]\n", username, player.CurrentIP())
			reply <- -1
		}()
	})
}
