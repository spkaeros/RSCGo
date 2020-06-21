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

	"github.com/spkaeros/rscgo/pkg/game"

	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("shopbuy", func(player *world.Player, p *net.Packet) {
		if player.HasState(world.StateShopping) {
			id := p.ReadUint16()
			price := p.ReadUint32()

			shop := player.CurrentShop()
			if shop == nil || player.State()&world.StateShopping != world.StateShopping {
				log.Suspicious.Println(player, "tried purchasing from a shop but is not apparently accessing any shops.")
				return
			}
			if shop.Inventory.Count(id) < 1 {
				player.Message("The shop has ran out of stock")
				return
			}
			realPrice := int(shop.Inventory.Get(id).Price().Scale(shop.BaseSalePercent + shop.DeltaPercentModID(id)))
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}

			if !player.Inventory.CanHold(id, 1) && player.Inventory.CountID(10) != price {
				player.Message("You can't hold the objects you are trying to buy!")
				return
			}
			if player.Inventory.CountID(10) < price || player.Inventory.RemoveByID(10, price) == -1 {
				player.Message("You don't have enough coins")
				return
			}

			player.AddItem(id, 1)
			shop.Remove(id, 1)
			player.PlaySound("coins")
			shop.Players.RangePlayers(func(player *world.Player) bool {
				player.SendPacket(world.ShopOpen(shop))
				return false
			})
		}
	})
	game.AddHandler("shopsell", func(player *world.Player, p *net.Packet) {
		if player.HasState(world.StateShopping) {
			id := p.ReadUint16()
			price := p.ReadUint32()

			shop := player.CurrentShop()
			if shop == nil {
				log.Suspicious.Println(player, "tried selling to a shop but is not apparently accessing any shops.")
				return
			}

			if !shop.Stock.Contains(id) && !shop.BuysUnstocked {
				player.Message("This shop does not purchase foreign objects")
				return
			}

			realPrice := int(world.Price(definitions.Item(id).BasePrice).Scale(shop.BasePurchasePercent + shop.DeltaPercentModID(id)))
			if price != realPrice {
				log.Suspicious.Println(player, "tried buying item["+strconv.Itoa(id)+"] for ["+strconv.Itoa(price)+"gp] but actual price is currently ["+strconv.Itoa(realPrice)+"gp]")
				return
			}
			if player.Inventory.RemoveByID(id, 1) > -1 {
				player.PlaySound("coins")
				player.AddItem(10, price)
				shop.Inventory.Add(&world.Item{ID: id, Amount: 1})
				shop.Players.RangePlayers(func(player *world.Player) bool {
					if shop == player.CurrentShop() {
						player.SendPacket(world.ShopOpen(shop))
					}
					return false
				})
			}
		}
	})
	game.AddHandler("shopclose", func(player *world.Player, p *net.Packet) {
		if player.HasState(world.StateShopping) {
			player.CloseShop()
		}
	})
}
