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
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

var reasons = []string{
	"Offensive Language",                // 1
	"Item scamming",                     // 2
	"Password scamming",                 // 3
	"Bug abuse",                         // 4
	"Staff impersonation",               // 5
	"Account sharing",                   // 6
	"Macroing",                          // 7
	"Multiple logging in",               // 8
	"Encouraging others to break rules", // 9
	"Misuse of customer support",        // 10
	"Advertising/website",               // 11
	"Real world item trading",           // 12
}

var actions = []string{
	"reported",
	"muted",
	"banned",
}

func init() {
	game.AddHandler("reportabuse", func(player *world.Player, p *net.Packet) {
		userHash := p.ReadUint64()
		reasonIndex := int(p.ReadUint8() - 1)
		actionIndex := int(p.ReadUint8())

		if userHash == player.UsernameHash() {
			player.Message("You can't report yourself!!")
			return
		}

		// validate reason for report
		if reasonIndex < 0 || reasonIndex > len(reasons)-1 {
			log.Suspicious.Printf("Report had invalid reason:\n[\n\taction:%d ('%s'),\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d\n];\n", actionIndex, actions[actionIndex], player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1)
			log.Info.Printf("Report had invalid reason:\n[\n\taction:%d ('%s'),\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d\n];\n", actionIndex, actions[actionIndex], player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1)
			return
		}
		// validate action report results in
		if actionIndex < 0 || actionIndex > len(actions)-1 {
			log.Suspicious.Printf("Report had invalid action:\n[\n\taction:%d,\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d ('%s')\n];\n", actionIndex, player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1, reasons[reasonIndex])
			log.Info.Printf("Report had invalid action:\n[\n\taction:%d,\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d ('%s')\n];\n", actionIndex, player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1, reasons[reasonIndex])
			return
		}
		// validate username provided for report is a real player
		if !db.DefaultPlayerService.PlayerNameExists(strutil.Base37.Decode(userHash)) {
			player.Message("Invalid player name.")
			return
		}

		log.Info.Printf("Report:\n[\n\taction:%d ('%s'),\n\tsender:'%s',\n\ttarget:'%s',\n\treason:%d ('%s')\n];\n", actionIndex, actions[actionIndex], player.Username(), strutil.Base37.Decode(userHash), reasonIndex+1, reasons[reasonIndex])
		log.Info.Println(player.Username(), actions[actionIndex], strutil.Base37.Decode(userHash), "for breaking rule", reasonIndex+1, "('"+reasons[reasonIndex]+"')")

		log.Suspicious.Println(player.Username(), actions[actionIndex], strutil.Base37.Decode(userHash), "for breaking rule", reasonIndex+1, "('"+reasons[reasonIndex]+"')")
		player.Message("Thank-you, your abuse report has been received.")
	})
}
