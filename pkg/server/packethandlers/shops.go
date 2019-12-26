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
			id := p.ReadShort()
			price := p.ReadInt()

			shop := player.CurrentShop()
			if shop == nil {
				log.Suspicious.Println(player, "tried purchasing from a shop but is not apparently accessing any shops.")
				return
			}

			if !shop.Inventory.Contains(id) {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but the ("+shop.Name+") shop is out of that item.")
				// TODO: og rsc msg here
				player.Message("There is no more of those in stock right now")
				return
			}

			item := shop.Inventory.Get(id)
			realPrice := int(item.Price().Scale(shop.BaseSalePercent + shop.StockDeltaPercentage(item)))
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}
			if player.Inventory.RemoveByID(10, price) > -1 {
				player.AddItem(id, 1)
				item.Amount--
				if item.Amount == 0 && shop.Stock.Count(id) == 0 {
					shop.Inventory.Remove(item)
				}
				world.Players.Range(func(player *world.Player) {
					if shop == player.CurrentShop() {
						player.SendPacket(world.ShopOpen(shop))
					}
				})
			} else {

			}
		}
	}
	PacketHandlers["shopsell"] = func(player *world.Player, p *packet.Packet) {
		if player.HasState(world.MSShopping) {
			id := p.ReadShort()
			price := p.ReadInt()

			shop := player.CurrentShop()
			if shop == nil {
				log.Suspicious.Println(player, "tried selling to a shop but is not apparently accessing any shops.")
				return
			}

			if !shop.Stock.Contains(id) && !shop.BuysUnstocked {
				//log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but the (" + shop.Name + ") shop is out of that item.")
				// TODO: og rsc msg here
				player.Message("This shop does not purchase foreign objects")
				return
			}

			realPrice := int(world.Price(world.ItemDefs[id].BasePrice).Scale(shop.BasePurchasePercent + shop.StockDeltaPercentID(id)))
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}
			if player.Inventory.RemoveByID(id, 1) > -1 {
				player.AddItem(10, price)
				shop.Inventory.Add(&world.Item{ID: id, Amount: 1})
				world.Players.Range(func(player *world.Player) {
					if shop == player.CurrentShop() {
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
