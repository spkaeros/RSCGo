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

	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

// thin wrapper around Anko spell handling subroutines.
// Utilizing the scripting environments runtime we can gain cool benefits, e.g
// fully dynamic spell dispatching, since the scripts and their env are all loaded at runtime.
// This also means that it is easier to shoot oneself in the foot (without realizing it, at times).
// For this reason, the API for scripts to use is implemented in stable fast Go, while scripts using said API
// can get away with suffering a minor performance decrease without causing any issues to the game.

func init() {
	game.AddHandler("spellnpc", func(player *world.Player, p *net.Packet) {
		targetIndex := p.ReadUint16()
		target := world.GetNpc(targetIndex)
		if target == nil {
			return
		}
		spellIndex := p.ReadUint16()
		//		log.Info.Println("cast on npc:", targetIndex, target.ID, spellIndex)
		dispatchSpellAction(player, spellIndex, target)
	})
	game.AddHandler("spellplayer", func(player *world.Player, p *net.Packet) {
		targetIndex := p.ReadUint16()
		target, ok := world.Players.FindIndex(targetIndex)
		if !ok {
			return
		}
		spellIndex := p.ReadUint16()
		//		log.Info.Println("cast on player:", targetIndex, target.String(), spellIndex)
		dispatchSpellAction(player, spellIndex, target)
	})
	game.AddHandler("spellself", func(player *world.Player, p *net.Packet) {
		idx := p.ReadUint16()

		//		log.Info.Println("Cast on self:", idx)
		dispatchSpellAction(player, idx, nil)
	})
	game.AddHandler("spellinvitem", func(player *world.Player, p *net.Packet) {
		p.Skip(2) // short uint itemIndex
		spellIndex := p.ReadUint16()
		//		log.Info.Println("Cast on invitem:", spellIndex, "on", itemIndex)
		dispatchSpellAction(player, spellIndex, nil)
	})
	game.AddHandler("spellgrounditem", func(player *world.Player, p *net.Packet) {
		//		itemX := p.ReadUint16()
		//		itemY := p.ReadUint16()
		//		itemID := p.ReadUint16()
		p.Skip(6)
		spellIndex := p.ReadUint16()
		//		log.Info.Println(itemX, itemY, itemID, "cast on grounditem:", spellIndex, "on", strconv.Itoa(itemID), "at", strconv.Itoa(itemX)+","+strconv.Itoa(itemY))
		dispatchSpellAction(player, spellIndex, nil)
	})
}

func dispatchSpellAction(player *world.Player, idx int, target entity.MobileEntity) {
	s, ok := world.SpellTriggers[idx]
	if !ok || s == nil {
		log.Info.Printf("Couldn't find spell handler ID: %v, status=`%v`\n", idx, ok)
		return
	}

	s(player, map[string]interface{}{"idx": idx, "target": target})
}
