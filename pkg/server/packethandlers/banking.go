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
		count := player.Inventory.CountID(id)
		if count < amount {
			log.Suspicious.Println("Attempted to deposit more than owned:", player.String())
			return
		}

		if player.Inventory.RemoveByID(id, amount) > -1 {
			player.Bank().Add(id, amount)
			player.SendPacket(world.BankUpdateItem(player.Bank().GetIndex(id), id, player.Bank().GetByID(id).Amount))
		}
	}
	PacketHandlers["withdrawbank"] = func(player *world.Player, p *packet.Packet) {
		if !player.HasState(world.MSBanking) {
			return
		}
		id := p.ReadShort()
		amount := p.ReadShort()
		//		botCheck := p.ReadInt()
		idx := player.Bank().GetIndex(id)
		if idx == -1 {
			log.Suspicious.Println("Attempted withdraw of item they do not have:", player.String(), id, amount)
			return
		}
		item := player.Bank().Get(idx)
		cnt := item.Amount - amount
		if item == nil || item.Amount < amount {
			log.Suspicious.Println("Attempted withdraw of items they do not have:", player.String(), id, amount)
			return
		}
		if player.Bank().RemoveByID(id, amount) > -1 {
			player.Inventory.Add(id, amount)
			player.SendInventory()
			player.SendPacket(world.BankUpdateItem(idx, id, cnt))
		}
	}
	PacketHandlers["closebank"] = func(player *world.Player, p *packet.Packet) {
		if !player.HasState(world.MSBanking) {
			return
		}
		player.CloseBank()
	}
}
