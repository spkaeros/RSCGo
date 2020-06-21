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
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("invwield", func(player *world.Player, p *net.Packet) {
		if player.IsDueling() && player.IsFighting() && !player.DuelEquipment() {
			player.Message("You can not use equipment in this duel")
			return
		}

		index := p.ReadUint16()
		if index < 0 || index > player.Inventory.Size() {
			log.Cheatf("Player[%v] tried to wield an item with an out-of-bounds inventory index: %d\n", player, index)
			return
		}

		item := player.Inventory.Get(index)
		if item == nil || item.Worn {
			return
		}

		player.EquipItem(item)
	})
	game.AddHandler("removeitem", func(player *world.Player, p *net.Packet) {
		index := p.ReadUint16()
		if index < 0 || index > player.Inventory.Size() {
			log.Cheatf("Player[%v] tried to unwield an item with an out-of-bounds inventory index: %d\n", player, index)
			return
		}

		item := player.Inventory.Get(index)
		if item == nil || !item.Worn {
			return
		}

		player.DequipItem(item)
		player.PlaySound("click")
	})
	game.AddHandler("takeitem", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		x := p.ReadUint16()
		y := p.ReadUint16()
		if x < 0 || x >= world.MaxX || y < 0 || y >= world.MaxY {
			log.Cheatf("%v attempted to pick up an item at an invalid location: [%d,%d]\n", player, x, y)
			return
		}

		id := p.ReadUint16()
		if id < 0 || id > len(definitions.Items)-1 {
			log.Cheatf("%v attempted to pick up an item with an out-of-bounds ID: %d\n", player, id)
			return
		}

		p.ReadUint16() // Unused in this opcode this variable is what affect-type we are applying to the ground item, e.g casting, using item with, etc...but we are not using any affects

		player.SetTickAction(func() bool {
			if player.Busy() {
				return false
			}

			item := world.GetItem(x, y, id)
			if item == nil || !item.VisibleTo(player) {
				log.Cheatf("%v attempted to pick up an item that doesn't exist: %s@{%d,%d}\n", player, definitions.Items[id].Name, x, y)
				return false
			}

			maxDelta := 0
			if world.IsTileBlocking(x, y, 0x40, false) {
				maxDelta++
			}
			if delta := player.Delta(item.Location); delta > maxDelta || delta == 1 && !player.ReachableCoords(item.X(), item.Y()) {
				return player.FinishedPath()
			}

			if !player.Inventory.CanHold(item.ID, item.Amount) {
				player.Message("You do not have room for that item in your inventory.")
				return false
			}

			player.ResetPath()
			item.Remove()
			player.Inventory.Add(item.ID, item.Amount)
			player.SendInventory()
			player.PlaySound("takeobject")
			return false
		})
	})
	game.AddHandler("dropitem", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		index := p.ReadUint16()
		// Just to prevent drops mid-path, and perform drop on path completion
		player.SetTickAction(func() bool {
			if player.Busy() {
				return false
			}
			if !player.FinishedPath() {
				return true
			}

			if player.Inventory.Size() < index {
				return false
			}

			item := player.Inventory.Get(index)
			if !player.Inventory.Remove(index) {
				return false
			}
			world.AddItem(world.NewGroundItemFor(player.UsernameHash(), item.ID, item.Amount, player.X(), player.Y()))
			player.PlaySound("dropobject")
			player.SendInventory()
			return true
		})
	})
	game.AddHandler("invaction1", func(player *world.Player, p *net.Packet) {
		index := p.ReadUint16()
		item := player.Inventory.Get(index)
		if item == nil || player.Busy() || player.IsFighting() {
			return
		}
		player.AddState(world.MSItemAction)
		go func() {
			defer func() {
				player.RemoveState(world.MSItemAction)
			}()
			for _, triggerDef := range world.ItemTriggers {
				if triggerDef.Check(item) {
					triggerDef.Action(player, item)
					return
				}
			}
			player.SendPacket(world.DefaultActionMessage)
		}()
	})
}
