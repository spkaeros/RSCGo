/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
)

// ShopPriceModifier represents the percentage to increase or decrease the price of an item by as players use it.
// For every item the shop buys, the asking price decreases by this much percentage of the items base price.
// However, for every item the shop sells, the asking price increases by this much percentage of the items base price.
//
// So if the shop starts with 10 of something, but has 5 right now, then that 5th item will cost you 15% MORE
// than the base asking price of that item.  However,  if the shop starts with 0 of something, and has 5 right now,
// then that 5th item will cost you 15% LESS than the base asking price of that item; so it does work both ways.
//
// There's a limit set that the shops prices for any item may not go below 10% of the items base price.
const ShopPriceModifier = 3

// shops contains all shop instances in the game world mapped to meaningful names.  Any calls to create a shop will
// automatically add them to this map.
var shops = map[string]*Shop{}

// shopItems is a type alias for a map to convert item IDs to item counts for shops.
type shopItems = map[int]int

// generalShop Default general shop definition
var generalShop = &Shop{true, 40, 130, shopItems{140: 2, 144: 2, 21: 2, 166: 2, 167: 2, 168: 5, 1263: 10}, shopItems{140: 2, 144: 2, 21: 2, 166: 2, 167: 2, 168: 5, 1263: 10}}

// NewShop returns a reference to a newly created shop with the specified parameters.
// See also:
// 	func NewGeneralShop(string) *Shop, func GetShop(string) *Shop, func AddShop(string), func RemoveShop(string)
//
// For every inventory count change, whether it be up or down from stock count, the shop changes its prices by a
// certain percentage according to whichever direction in which the count has gone.
// 	See also: func(*Shop) PurchasesPrice(int) int, func(*Shop) SalesPrice(int) int, func(*Shop) PriceCountMod(int) int
func NewShop(name string, percentPurchasesPrice, percentSalesPrice int, stock shopItems) *Shop {
	// clone map to be able to mutate the shop inventory while leaving its stock unchanged.
	inventory := shopItems{}
	for id, amt := range stock {
		inventory[id] = amt
	}
	return &Shop{PurchasesPricePercent: percentPurchasesPrice, SalesPricePercent: percentSalesPrice, Stock: stock, Inventory: inventory}
}

// NewGeneralShop returns a reference to a newly created general shop, and adds it to the worlds shop collection.
// 	See also: func NewShop(string) *Shop, func GetShop(string) *Shop, func AddShop(string), func RemoveShop(string)
//
//
//The distinction of what makes it a general shop is that it buys any tradeable items the player tries to sell it,
// where a normal shop only deals in its initial stock.
//
// Default stock:
//     2 jugs, 2 shears, 2 buckets, 2 tinderboxes, 2 chisels, 5 hammers, 10 sleeping bags
// Default percentage for sales price:
//     130% of item base price
// Default percentage for purchase price:
//     40% of item base price
//
// For every inventory count change, whether it be up or down from stock count, the shop changes its prices by a
// certain percentage according to whichever direction in which the count has gone.
// 	See also: func(*Shop) PurchasesPrice(int) int, func(*Shop) SalesPrice(int) int, func(*Shop) PriceCountMod(int) int
func NewGeneralShop(name string) *Shop {
	shops[name] = generalShop.Clone()
	return shops[name]
}

//GetShop attempts to fetch and return a shop by its name.  If no shop exists in the world with this name, it will
// return a brand new general shop instance and log the event
// See also: func NewShop(string) *Shop, func NewGeneralShop(string) *Shop, func AddShop(string), func RemoveShop(string)
func GetShop(name string) *Shop {
	s, ok := shops[name]
	if !ok || s == nil {
		log.Warning.Println("GetShop(\"" + name + "\") called but no such shop exists...")
		return nil
	}
	return s
}

//AddShop maps `name` to `s` in the world package-local shops collection.  If an existing entry is using that name,
// it will overwrite it.
func AddShop(name string, s *Shop) {
	shops[name] = s
}

//RemoveShop removes the shop registered with the given name if it exists, otherwise logs a warning to notify sysadmin.
func RemoveShop(name string) {
	_, ok := shops[name]
	if !ok {
		log.Warning.Println("Tried removing a shop that doesn't seem to exist: '" + name + "'")
		return
	}
	delete(shops, name)
}

//Shop represents an NPC owned shop in the game world.
type Shop struct {
	//BuysUnstocked is set true if this shop deals in unstocked items the player wants to sell, otherwise it's set false.
	BuysUnstocked bool
	//SalesPricePercent the percent of an items base price to ask for shop items.
	PurchasesPricePercent int
	//SalesPricePercent the percent of an items base price to offer for player items.
	SalesPricePercent int
	//Inventory represents the initial items in the shop.  This collection is immutable and should never be changed.
	Stock shopItems
	//Inventory represents the current items in the shop.  This collection is mutable and changes over time.
	Inventory shopItems
}

//Clone makes a clone of the receiver and returns it by reference.
// Primary usage was for easily copying the default instance of a general shop.
func (s *Shop) Clone() *Shop {
	return &Shop{s.BuysUnstocked, s.PurchasesPricePercent, s.SalesPricePercent, s.Stock, s.Inventory}
}

//SalesPrice will calculate and return the current sale price of the item with the provided id at this shop.
// It does so by taking the base sale price percentage (general shops are 130% of the items basePrice), and
// adding a modifier that is different for every item to that to figure how much percentage of the item base price we
// are currently asking for this item.  Then, we take that percentage figure and multiply it together with the item's
// base price, from its definition. Then, we divide the resulting product by 100 to get the final item sale price
//
// This is how Jagex calculated shop sale prices in RSClassic.
func (s *Shop) SalesPrice(id int) int {
	percentage := s.SalesPricePercent + s.PriceCountMod(id)
	if percentage < 10 {
		percentage = 10
	}
	return (percentage * ItemDefs[id].BasePrice) / 100
}

//PurchasesPrice will calculate and return the current purchase price of the item with the provided id at this shop.
// It does so by taking the base purchase price percentage (general shops are 40% of the items basePrice), and
// adding a modifier that is different for every item to that to get out how much percentage of the item base price we
// are currently paying for this item.  Then, we take that percentage figure and multiply it together with the item's
// base price, from its definition. Then, we divide the resulting product by 100 to get the final item purchase price.
//
// This is how Jagex calculated shop purchase prices in RSClassic.
func (s *Shop) PurchasesPrice(id int) int {
	percentage := s.PurchasesPricePercent + s.PriceCountMod(id)
	if percentage < 10 {
		percentage = 10
	}
	return (percentage * ItemDefs[id].BasePrice) / 100
}

//PriceCountMod returns the percentage to increase/decrease the price of the item with the given ID at this shop.
// This is calculated by subtracting the amount of said item the shop has in stock currently from the amount of
// said item that the shop's stock started with when it was created, and multiplying the difference by 3%.
func (s *Shop) PriceCountMod(id int) int {
	itemMod := (s.Stock[id] - s.Inventory[id]) * ShopPriceModifier
	if itemMod > 15 {
		itemMod = 15
	}
	return itemMod
}
