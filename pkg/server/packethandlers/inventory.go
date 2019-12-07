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
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["invwield"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", player, index)
			return
		}
		if item := player.Inventory.Get(index); item != nil {
			if item.Worn {
				return
			}
			player.PlaySound("click")
			player.EquipItem(item)
			player.SendPacket(world.EquipmentStats(player))
			player.SendPacket(world.InventoryItems(player))
		}
	}
	PacketHandlers["removeitem"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", player, index)
			return
		}
		// TODO: Wielding
		if item := player.Inventory.Get(index); item != nil {
			if !item.Worn {
				return
			}
			player.PlaySound("click")
			player.DequipItem(item)
			player.SendPacket(world.EquipmentStats(player))
			player.SendPacket(world.InventoryItems(player))
		}
	}
	PacketHandlers["takeitem"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		id := p.ReadShort()
		p.ReadShort() // Useless, this variable is for what affect we are applying to the ground item, e.g casting, using item with
		if x < 0 || x >= world.MaxX || y < 0 || y >= world.MaxY {
			log.Suspicious.Printf("Player[%v] attempted to pick up an item at an invalid location: [%d,%d]\n", player, x, y)
			return
		}
		if id < 0 || id > 1289 {
			log.Suspicious.Printf("Player[%v] attempted to pick up an item with an invalid ID: %d\n", player, id)
			return
		}

		player.SetDistancedAction(func() bool {
			item := world.GetItem(x, y, id)
			if item == nil || !item.VisibleTo(player) {
				log.Suspicious.Printf("Player[%v] attempted to pick up an item that doesn't exist: %d,%d,%d\n", player, id, x, y)
				return true
			}
			if !player.WithinRange(item.Location, 0) {
				return false
			}
			player.PlaySound("takeobject")
			player.ResetPath()
			if player.Inventory.Size() >= 30 {
				player.Message("You do not have room for that item in your inventory.")
				return true
			}
			item.Remove()
			player.Inventory.Add(item.ID, item.Amount)
			player.SendPacket(world.InventoryItems(player))
			return true
		})
	}
	PacketHandlers["dropitem"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		index := p.ReadShort()
		item := player.Inventory.Get(index)
		if item != nil {
			player.SetDistancedAction(func() bool {
				if player.FinishedPath() {
					groundItem := world.NewGroundItemFor(player.UserBase37, item.ID, item.Amount, player.X(), player.Y())
					if player.Inventory.Remove(index, item.Amount) {
						world.AddItem(groundItem)
						player.SendPacket(world.InventoryItems(player))
						player.PlaySound("dropobject")
					}
					return true
				}
				return false
			})
		}
	}
	PacketHandlers["invaction1"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		index := p.ReadShort()
		item := player.Inventory.Get(index)
		if item != nil {
			player.AddState(world.MSBusy)
			go func() {
				defer func() {
					player.RemoveState(world.MSBusy)
				}()
				for _, triggerDef := range script.ItemTriggers {
					if triggerDef.Check(item) {
						triggerDef.Action(player, item)
						return
					}
				}
				player.SendPacket(world.DefaultActionMessage)
			}()
		}
	}
}
