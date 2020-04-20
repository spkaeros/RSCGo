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
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	AddHandler("depositbank", func(player *world.Player, p *net.Packet) {
		if !player.HasState(world.StateBanking) {
			return
		}
		id := p.ReadUint16()
		amount := p.ReadUint32()
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
	})
	AddHandler("withdrawbank", func(player *world.Player, p *net.Packet) {
		if !player.HasState(world.StateBanking) {
			return
		}
		id := p.ReadUint16()
		amount := p.ReadUint32()
		//		botCheck := p.ReadUint32()
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
		if !player.Inventory.CanHold(id, amount) {
			player.Message("You don't have room to hold everything!")
			return
		}
		if player.Bank().RemoveByID(id, amount) > -1 {
			player.Inventory.Add(id, amount)
			player.SendInventory()
			player.SendPacket(world.BankUpdateItem(idx, id, cnt))
		}
	})
	AddHandler("closebank", func(player *world.Player, p *net.Packet) {
		if !player.HasState(world.StateBanking) {
			return
		}
		player.CloseBank()
	})
}
