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
	"sync"
	"time"
)

const (
	// The percentage to increase or decrease the price of an item by as players use it.
	//
	// For every item the shop buys, the asking price decreases by itemDeltaCost% of the items base price.
	// However, for every item the shop sells, the asking price increases by itemDeltaCost% of the items base price.
	itemDeltaCost = 3
)

type (
	// Collection to hold all of the items that a shop is carrying.  Items may have 0 amount and continue to live here.
	// However, They will only be sent to the client with 0 amount if the shop has it in its default stock.
	//
	// The main reason for populating this collection with shop items that have 0 amount is to easily call Item related
	// methods when we need price information.
	ShopItems []*Item

	Shop struct {
		// Shop represents a shop in the game world, typically managed by an NPC, but not necessarily (for RSC preservation,
		// it is necessary I believe...no shops other than ones accessed by NPCs there)
		//
		// Each time a shop buys or sells a shop item, that specific item type's price is scaled to the following percent:
		//	basePercent+(item.DeltaAmount(shop.GetStockItem(item.ID))*itemDeltaCost)
		// where basePercent can be distinct for selling or buying items, or they can be identical it doesn't matter. Only
		// limitation as far as this goes, really, seems to be that you may not scale to percentages lower than 10%, it is
		// the absolute lowest percentage of base item value that we can scale shop items to

		// True if this shop deals in unstocked items the player wants to sell, otherwise false.
		BuysUnstocked bool
		// The shop's base percentage at which it buys items from players.
		BasePurchasePercent int
		// The shop's base percentage at which it sells items to players.
		BaseSalePercent int
		// Contains all of the initial shop items that are always restocked and available, and their initial amounts.
		// Never change the contents of this unless you want to add permanently restocking items to a shop or something like that.
		Stock ShopItems
		// Contains all of the active shop items and may differ greatly from Stock, but it can never remove things that are
		// in Stock entirely from itself(they will remain with Amount=0), and occasionally it normalizes itself, referencing
		// Stock for the normal amounts to replenish toward, and the default IDs of what items to replenish.
		Inventory ShopItems
	}
	ShopContainer struct {
		// This is a concurrency-friendly collection set to simplify containing world-scoped shops without introducing any
		// data race conditions.
		//
		// The structure exposes a lot of easy accessor methods to manipulate the underlying set using a simple API in a
		// concurrency friendly manner.  All accesses via these methods are locked using a RWMutex structure to provide
		// safe concurrent access from any number of goroutines.
		//
		// Some examples from the exported API are:
		//  Add(name string, shop Shop) or Remove(name string) shops, storing them using name as their key in its set
		//  Get(name string) a shop by its identifying name.  If it does not find any shop with that name, it will return the Shop zero value
		//  Contains(name string) to check if a shop with a specific name exists or not in this collection set.
		//  Range(fn func(Shop)) to range over the collection set and perform an action for every entry.
		set map[string]*Shop
		sync.RWMutex
	}
)

func (s *ShopContainer) Add(name string, shop Shop) {
	s.Lock()
	s.set[name] = &shop
	s.Unlock()
}

func (s *ShopContainer) Contains(name string) bool {
	s.RLock()
	_, ok := s.set[name]
	s.RUnlock()
	return ok
}

func (s *ShopContainer) Get(name string) *Shop {
	s.RLock()
	shop, ok := s.set[name]
	s.RUnlock()
	if !ok {
		return &Shop{}
	}
	return shop
}

func (s *ShopContainer) Range(fn func(*Shop)) {
	s.RLock()
	for _, value := range s.set {
		fn(value)
	}
	s.RUnlock()
}

func (s *ShopContainer) Remove(name string) {
	s.Lock()
	delete(s.set, name)
	s.Unlock()
}

var (
	//Shops contains all shop instances in the game world mapped to identifying names.
	//
	// No Shops need to be loaded until a player requests one, in my opinion.  From that point on, Shops are active for
	// the entire runtime of the server, unless a system admin does something to stop them.  The structures use so
	// little memory in practice, that I can't imagine a shop ever needing unset for performance or resource usage issues.
	Shops = &ShopContainer{
		set: make(map[string]*Shop),
	}
	//generalStock defines all default items that any general shop carries in their stock, and amounts.
	generalStock = ShopItems{
		{ID: 140, Amount: 2},   // 2 jugs
		{ID: 144, Amount: 2},   // 2 shears
		{ID: 21, Amount: 2},    // 2 buckets
		{ID: 166, Amount: 2},   // 2 tinderboxes
		{ID: 167, Amount: 2},   // 2 chisels
		{ID: 168, Amount: 5},   // 5 hammers
		{ID: 1263, Amount: 10}, // 10 sleeping bags
	}
	// The distinction of what makes it a general shop is that it buys any tradeable items the player tries to sell it,
	// where a normal shop only deals in its initial stock.
	// For every inventory count change, whether it be up or down from stock count, the shop changes its prices by a
	// certain percentage according to whichever direction in which the count has gone.
	//
	// This is the basic default definition of a general shop:
	// Default stock:
	// 	2 jugs
	//	2 shears
	//	2 buckets
	//	2 tinderboxes
	//	2 chisels
	//	5 hammers
	//	10 sleeping bags
	// Base sale percent:
	// 	130% of basePrice
	// Base purchase percent:
	// 	40% of basePrice
	//
	// Defaults for general stores from RSClassic are:
	//	Shop{
	//		BuysUnstocked: true, // obvious general stores are known for this
	//		BasePurchasePercent: 40, // this is used when the shop owner buys an item from you.
	//		BaseSalePercent: 130, // this is used when the shop owner sells an item to you.
	//		Stock: generalStock.Clone(),
	//		Inventory: generalStock.Clone(),
	//	}
	generalShop = Shop{
		true, // General Shops are known for buying all items, whether or not they normally stock it
		40, // the shop pays 40% of base price to purchase any item, less for each one it buys
		130, // the shop asks 130% of base price to purchase any item, more for each one it sells
		generalStock.Clone(), // We want to avoid ending up with identical references and so we clone
		generalStock.Clone(), // see above
	}
)

//Clone will make a new ShopItems collection populated with clones of the values in the existing collection set.
// Very useful to easily define an initial inventory and stock when defining a shop.
func (s ShopItems) Clone() ShopItems {
	clone := ShopItems(nil)
	for _, v := range s {
		clone = append(clone, &Item{ID: v.ID, Amount: v.Amount})
	}
	return clone
}

//NewShop creates a new Shop instance using the arguments provided, and returns it.
//
// Returns: a new Shop instance, made with the given arguments
func NewShop(name string, percentPurchasesPrice, percentSalesPrice int, stock ShopItems) Shop {
	return Shop{BasePurchasePercent: percentPurchasesPrice, BaseSalePercent: percentSalesPrice, Stock: stock, Inventory: append(ShopItems(nil), stock...)}
}

// Creates a new general shop, and adds it automatically to the world-local ShopContainer instance before returning it
//
// Returns: Shops.get(name), after building and adding a new general shop to it, using a generic general shop definition.
func NewGeneralShop(name string) *Shop {
	Shops.set[name] = (&generalShop).Clone()
	go func() {
		for {
			time.Sleep(12400*time.Millisecond)
			shop := Shops.Get(name)
			if shop == nil {
				return
			}
			changed := false
			for idx, item := range shop.Inventory {
				stocked := shop.GetStockItem(item.ID)
				if stocked == nil || stocked.Amount == 0 {
					if item.Amount > 0 {
						item.Amount--
						changed = true
					}
					if item.Amount == 0 {
						changed = true
						if idx < len(shop.Inventory) {
							shop.Inventory = append(shop.Inventory[:idx], shop.Inventory[idx+1:]...)
						} else {
							shop.Inventory = shop.Inventory[:idx]
						}
					}
				} else if stocked.ID == item.ID {
					if stocked.Amount < item.Amount {
						changed = true
						item.Amount--
					} else if stocked.Amount > item.Amount {
						changed = true
						item.Amount++
					}
				}
			}
			if changed {
				Players.Range(func(player *Player) {
					if player.HasState(MSShopping) && player.CurrentShop() == shop {
						player.SendPacket(ShopOpen(*shop))
					}
				})
			}
		}
	}()
	return Shops.set[name]
}

//Get returns the shop item entry with the given ID, or if it can't find it, adds a new one then returns it.
func (s *ShopItems) Get(id int) *Item {
	for _, item := range *s {
		if item.ID == id {
			return item
		}
	}
	*s = append(*s, &Item{ID: id})
	return (*s)[len(*s)-1]
}

//Clone makes a clone of the receiver shop and returns it.
func (s *Shop) Clone() *Shop {
	return &Shop{s.BuysUnstocked, s.BasePurchasePercent, s.BaseSalePercent, s.Stock.Clone(), s.Inventory.Clone()}
}

//GetStockItem returns the shop default stock item with the specified ID.
// If there is no stock item with this ID, instantiates an otherwise zero-value item with the ID changed to id, and
// returns a reference to it.
func (s *Shop) GetStockItem(id int) *Item {
	for _, item := range s.Stock {
		if item.ID == id {
			return item
		}
	}

	return &Item{ID: id}
}

//GetItem returns the shop active inventory item with the specified ID.
// If there is no shop inventory item with this ID, instantiates an otherwise zero-value item with the ID changed to id,
// and returns a reference to it.
func (s *Shop) GetItem(id int) *Item {
	return s.Inventory.Get(id)
}

//StockDeltaPercentage calculates the percentage to scale the item's price up or down from its respective base percentage.
// The formula simply subtracts the shop's base stocked amount of item from item's amount, and multiplies the difference
// by itemDeltaCost.
func (s *Shop) StockDeltaPercentage(item *Item) int {
	return item.DeltaAmount(s.GetStockItem(item.ID)) * itemDeltaCost
}
