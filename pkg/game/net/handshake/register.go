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
	// "time"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/crypto"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/net/handlers"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	handlers.AddHandler("newplayer", func(player *world.Player, p *net.Packet) {
		reply := NewRegistrationListener(player).attachPlayer(player)
		// if registerThrottle.Recent(player.CurrentIP(), time.Hour) >= 2 {
			// reply <- Response{ResponseSpamTimeout, "Recently registered too many other characters"}
			// return
		// }
		if version := p.ReadUint16(); version != config.Version() {
			reply <- Response{ResponseUpdated, "Invalid client version (" + strconv.Itoa(version) + ")"}
			return
		}
		username := strutil.Base37.Decode(p.ReadUint64())
		password := strings.TrimSpace(p.ReadString())
		player.SetVar("username", username)
		go func() {
			if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
				reply <- Response{ResponseBadInputLength, "Password and/or username too long and/or too short."}
				return
			}
			dataService := db.DefaultPlayerService
			if dataService.PlayerNameExists(username) {
				reply <- Response{ResponseUsernameTaken, "Username is taken by another player already."}
				return
			}

			if !dataService.PlayerCreate(username, crypto.Hash(password), player.CurrentIP()) {
				reply <- Response{-1, ""}
				return
			}
			reply <- Response{ResponseRegisterSuccess, ""}
		}()
	})
}
