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
	"strconv"

	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["shopbuy"] = func(player *world.Player, p *packet.Packet) {
		if player.HasState(world.MSShopping) {
			shop := player.CurrentShop()
			if shop == nil {
				log.Suspicious.Println(player, "tried purchasing from a shop but is not apparently accessing any shops.")
				return
			}

			id := p.ReadShort()
			price := p.ReadInt()
			curAmt := shop.Inventory[id]
			if curAmt < 1 {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but the shop is apparently out of that item.")
				return
			}
			realPrice := shop.SalesPrice(id)
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}
			if player.Inventory.RemoveByID(10, price) > -1 {
				log.Info.Println(id, price)
				if player.Inventory.Add(id, 1) > -1 {
					player.SendInventory()
				}
				shop.Inventory[id] = curAmt - 1
				world.Players.Range(func(player *world.Player) {
					if player.CurrentShop() == shop {
						player.SendPacket(world.ShopOpen(shop))
					}
				})
			}
		}
	}
	PacketHandlers["shopsell"] = func(player *world.Player, p *packet.Packet) {
		if player.HasState(world.MSShopping) {
			shop := player.CurrentShop()
			if shop == nil {
				log.Suspicious.Println(player, "tried selling to a shop but is not apparently accessing any shops.")
				return
			}

			id := p.ReadShort()
			price := p.ReadInt()
			realPrice := shop.PurchasesPrice(id)
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}
			if player.Inventory.RemoveByID(id, 1) > -1 {
				player.Inventory.Add(10, price)
				shop.Inventory[id] = shop.Inventory[id] + 1
				log.Info.Println(id, price)
				world.Players.Range(func(player *world.Player) {
					if player.CurrentShop() == shop {
						player.SendPacket(world.ShopOpen(shop))
					}
				})
			}
		}
	}
	PacketHandlers["shopclose"] = func(player *world.Player, p *packet.Packet) {
		if player.HasState(world.MSShopping) {
			player.CloseShop()
		}
	}
}
