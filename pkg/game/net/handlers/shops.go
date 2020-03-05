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
	"strconv"

	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	AddHandler("shopbuy", func(player *world.Player, p *net.Packet) {
		if player.HasState(world.MSShopping) {
			id := p.ReadUint16()
			price := p.ReadUint32()

			shop := player.CurrentShop()
			if shop == nil {
				log.Suspicious.Println(player, "tried purchasing from a shop but is not apparently accessing any shops.")
				return
			}

			if shop.Inventory.Count(id) < 1 {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but the ("+shop.Name+") shop is out of that item.")
				// TODO: og rsc msg here
				player.Message("There is no more of those in stock right now")
				return
			}

			item := shop.Inventory.Get(id)
			realPrice := int(item.Price().Scale(shop.BaseSalePercent + shop.DeltaPercentMod(item)))
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}
			if player.Inventory.RemoveByID(10, price) > -1 {
				player.AddItem(id, 1)
				if shop.Stock.Count(id) == 0 {
					if item.Amount == 1 {
						shop.Inventory.Remove(item)
					} else {
						item.Amount--
					}
				} else {
					item.Amount--
				}
				world.Players.Range(func(player *world.Player) {
					if shop == player.CurrentShop() {
						player.SendPacket(world.ShopOpen(shop))
					}
				})
			}
		}
	})
	AddHandler("shopsell", func(player *world.Player, p *net.Packet) {
		if player.HasState(world.MSShopping) {
			id := p.ReadUint16()
			price := p.ReadUint32()

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

			realPrice := int(world.Price(world.ItemDefs[id].BasePrice).Scale(shop.BasePurchasePercent + shop.DeltaPercentModID(id)))
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
	})
	AddHandler("shopclose", func(player *world.Player, p *net.Packet) {
		if player.HasState(world.MSShopping) {
			player.CloseShop()
		}
	})
}
