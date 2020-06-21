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
	"time"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/log"
)

var (
	crewHead           = 1
	metalHead          = 4
	downsHead          = 6
	beardHead          = 7
	baldHead           = 8
	validHeads         = []int{crewHead, metalHead, downsHead, beardHead, baldHead}
	validFemaleHeads   = []int{crewHead, metalHead, downsHead, baldHead}
	maleBody           = 2
	femaleBody         = 5
	validBodys         = []int{maleBody, femaleBody}
	validSkinColors    = []int{0xecded0, 0xccb366, 0xb38c40, 0x997326, 0x906020}
	validHeadColors    = []int{0xffc030, 0xffa040, 0x805030, 0x604020, 0x303030, 0xff6020, 0xff4000, 0xffffff, 65280, 65535}
	validBodyLegColors = []int{0xff0000, 0xff8000, 0xffe000, 0xa0e000, 57344, 32768, 41088, 45311, 33023, 12528, 0xe000e0, 0x303030, 0x604000, 0x805000, 0xffffff}
)

func inArray(a []int, i int) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func init() {
	game.AddHandler("changeappearance", func(player *world.Player, p *net.Packet) {
		if !player.HasState(world.StateChangingLooks) {
			// Make sure the player either has never logged in before, or talked to the makeover mage to get here.
			return
		}
		isMale := p.ReadBoolean()
		headType := int(p.ReadUint8() + 1)
		bodyType := int(p.ReadUint8() + 1)
		legType := int(p.ReadUint8() + 1) // appearance2Colour, seems to be a client const, value seems to remain 2.  ofc, legs never change
		hairColor := int(p.ReadUint8())
		topColor := int(p.ReadUint8())
		legColor := int(p.ReadUint8())
		skinColor := int(p.ReadUint8())
		/*		if !inArray(validHeads, int(headType)) || !inArray(validBodys, int(bodyType)) || !inArray(validBodyLegColors, int(topColor)) ||
				!inArray(validBodyLegColors, int(legColor)) || !inArray(validHeadColors, int(hairColor))  ||
				!inArray(validSkinColors, int(skinColor)) || legType != 2 {*/
		if hairColor >= len(validHeadColors) || !inArray(validHeads, headType) || topColor >= len(validBodyLegColors) ||
			legColor >= len(validBodyLegColors) || skinColor >= len(validSkinColors) || !inArray(validBodys, bodyType) || legType != 3 || legColor >= len(validBodyLegColors) {
			log.Warnf("Invalid appearance data provided by %v: (headType:%v, bodyType:%v, legType:%v, hairColor:%v, topColor:%v, legColor:%v, skinColor:%v, gender:%v)\n", player.String(), headType, bodyType, legType, hairColor, topColor, legColor, skinColor, isMale)
			return
		}
		if config.Verbosity >= 2 {
			log.Debugf("(headType:%v, bodyType:%v, legType:%v, hairColor:%v, topColor:%v, legColor:%v, skinColor:%v, gender:%v)\n", headType, bodyType, legType, hairColor, topColor, legColor, skinColor, isMale)
		}
		if !isMale {
			if bodyType != femaleBody {
				log.Cheat("Correcting invalid packet data: female asked for male body type; setting to female body type, packet from", player)
				bodyType = femaleBody
			}
			if headType == beardHead {
				log.Cheat("Correcting invalid packet data: female asked for male head type; setting to female head type, packet from", player)
				headType = metalHead
			}
		}
		{
			sprites := player.Equips()
			if sprites[0] == player.Appearance.Head {
				sprites[0] = headType
			}
			if sprites[1] == player.Appearance.Body {
				sprites[1] = bodyType
			}
			player.Appearance.Body = bodyType
			player.Appearance.Head = headType
			player.Appearance.Male = isMale
			player.Appearance.HeadColor = hairColor
			player.Appearance.SkinColor = skinColor
			player.Appearance.BodyColor = topColor
			player.Appearance.LegsColor = legColor
			player.UpdateAppearance()
		}
		player.RemoveState(world.StateChangingLooks)
		if !player.Attributes.Contains("madeAvatar") {
			player.SendPacket(world.WelcomeMessage)
			player.Attributes.SetVar("madeAvatar", time.Now())
			player.Attributes.SetVar("lastLogin", time.Now())
		}
	})
}
