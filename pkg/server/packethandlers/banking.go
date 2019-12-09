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
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["depositbank"] = func(player *world.Player, p *packet.Packet) {
		if !player.HasState(world.MSBanking) {
			return
		}
		id := p.ReadShort()
		amount := p.ReadShort()
		if amount < 1 {
			log.Suspicious.Println("Attempted to deposit less than 1:", player.String())
			return
		}
		//		botCheck := p.ReadInt()
		count := player.Inventory.CountID(id)
		if count < amount {
			log.Suspicious.Println("Attempted to deposit more than owned:", player.String())
			return
		}

		if world.ItemDefs[id].Stackable {
			if item := player.Bank.GetByID(id); item != nil && player.Inventory.RemoveByID(id, amount) > -1 {
				item.Amount += amount
			} else if player.Bank.Size() < player.Bank.Capacity-1 && player.Inventory.RemoveByID(id, amount) > -1 {
				player.Bank.Add(id, amount)
			}
		} else {
			if item := player.Inventory.GetByID(id); item.Worn {
				player.DequipItem(item)
			}
			for j := 0; j < amount; j++ {
				if item := player.Bank.GetByID(id); item != nil && player.Inventory.RemoveByID(id, 1) > -1 {
					item.Amount += 1
				} else if player.Bank.Size() < player.Bank.Capacity-1 && player.Inventory.RemoveByID(id, 1) > -1 {
					player.Bank.Add(id, 1)
				}
			}
		}

		if deposited := player.Bank.GetByID(id); deposited != nil && deposited.Index > -1 {
			player.SendInventory()
			player.SendPacket(world.BankUpdateItem(deposited))
		}
	}
	PacketHandlers["withdrawbank"] = func(player *world.Player, p *packet.Packet) {
		if !player.HasState(world.MSBanking) {
			return
		}
		id := p.ReadShort()
		amount := p.ReadShort()
		//		botCheck := p.ReadInt()
		item := player.Bank.GetByID(id)
		if item == nil || item.Amount < amount {
			log.Suspicious.Println("Attempted withdraw of items they do not have:", player.String(), id, amount)
			return
		}
		if item.Stackable() {
			if invItem := player.Inventory.GetByID(item.ID); invItem != nil && player.Bank.RemoveByID(id, amount) > -1 {
				player.Inventory.Add(item.ID, item.Amount)
			} else if player.Inventory.Size() < player.Inventory.Capacity-1 && player.Bank.RemoveByID(id, amount) > -1 {
				player.Inventory.Add(item.ID, item.Amount)
			}
		} else {
			for j := 0; j < amount; j++ {
				if player.Bank.RemoveByID(item.ID, 1) > -1 {
					if player.Inventory.Add(item.ID, 1) < 0 {
						break
					}
				}
			}
		}
		player.SendInventory()
		player.SendPacket(world.BankUpdateItem(item))
	}
	PacketHandlers["closebank"] = func(player *world.Player, p *packet.Packet) {
		if !player.HasState(world.MSBanking) {
			return
		}
		player.CloseBank()
	}
}
