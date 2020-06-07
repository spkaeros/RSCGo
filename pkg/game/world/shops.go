/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"sync"

	"github.com/spkaeros/rscgo/pkg/tasks"
)

const (
	// The percentage to increase or decrease the price of an item by as players use it.
	//
	// For every item the shop buys, the asking price decreases by ShopNormalDeltaRate% of the items base price.
	// However, for every item the shop sells, the asking price increases by ShopNormalDeltaRate% of the items base price.
	ShopNormalDeltaRate = 3
	// Defines the base selling price for buying items from players in the general store.
	ShopBuyPriceBasePercent = 40
	// Defines the base asking price for selling items to players in the general store.
	ShopSellPriceBasePercent = 130
	// Defines how often to normalize most general stores inventorys.  Might be an exception to this rule later.
	ShopGeneralRespawnTime = 50
)

type (
	Shop struct {
		// Shop represents a shop in the game world, typically managed by an NPC, but not necessarily (for RSC preservation,
		// it is necessary I believe...no shops other than ones accessed by NPCs there)
		//
		// Each time a shop buys or sells a shop item, that specific item type's price is scaled to the following percent:
		//	basePercent+(item.DeltaAmount(shop.GetStockItem(item.ID))*ShopNormalDeltaRate)
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
		Stock *ShopItems
		// Contains all of the active shop items and may differ greatly from Stock, but it can never remove things that are
		// in Stock entirely from itself(they will remain with Amount=0), and occasionally it normalizes itself, referencing
		// Stock for the normal amounts to replenish toward, and the default IDs of what items to replenish.
		Inventory *ShopItems
		// Descriptive name for this shop.
		Name string
		// List of players actively using the shop
		Players *MobList
	}
	ShopItems struct {
		// This is a concurrency-friendly collection set to simplify containing shop-scoped item lists without introducing any
		// data race conditions.
		set shopItemSet
		sync.RWMutex
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

type shopItemSet []*Item

func (s shopItemSet) Clone() shopItemSet {
	clone := make([]*Item, len(s))
	for i, v := range s {
		clone[i] = &Item{ID: v.ID, Amount: v.Amount}
	}
	return clone
}

func (s *ShopContainer) Add(name string, shop *Shop) {
	s.Lock()
	s.set[name] = shop
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
	// the entire runtime of the game, unless a system admin does something to stop them.  The structures use so
	// little memory in practice, that I can't imagine a shop ever needing unset for performance or resource usage issues.
	Shops = &ShopContainer{
		set: make(map[string]*Shop),
	}
	// Represents a default general shop's stock.
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
	generalStock = &ShopItems{
		set: shopItemSet{
			{ID: 140, Amount: 2},   // 2 jugs
			{ID: 144, Amount: 2},   // 2 shears
			{ID: 21, Amount: 2},    // 2 buckets
			{ID: 166, Amount: 2},   // 2 tinderboxes
			{ID: 167, Amount: 2},   // 2 chisels
			{ID: 168, Amount: 5},   // 5 hammers
			{ID: 1263, Amount: 10}, // 10 sleeping bags
		},
	}
)

//NewShop creates a new Shop instance using the arguments provided, and returns it.
//
// Returns: a new Shop instance, made with the given arguments
func NewShop(percentPurchasesPrice, percentSalesPrice int, stock shopItemSet, name string) *Shop {
	s := &ShopItems{set: stock}
	return &Shop{BasePurchasePercent: percentPurchasesPrice, BaseSalePercent: percentSalesPrice, Stock: s, Inventory: s.Clone(), Name: name}
}

// Creates a new general shop, and adds it automatically to the world-local ShopContainer instance before returning it
//
// Returns: Shops.get(name), after building and adding a new general shop to it, using a generic general shop definition.
func NewGeneralShop(name string) *Shop {
	shop := &Shop{true, 40, 130, generalStock.Clone(), generalStock.Clone(), name, NewMobList()}
	Shops.Add(name, shop)
	shopTicker := 0
	tasks.TickList.Add(func() bool {
		shopTicker++
		if shopTicker == 20 { // 12.8s
			shopTicker = 0
			changed := false
			shop.Inventory.Range(func(item *Item) bool {
				if shop.Stock.Count(item.ID) == item.Amount {
					return false
				}
				changed = true
				//				if stockedAmount <= 0 {
				//					changed = true

				//					if item.Amount > 0 {
				//						item.Amount--
				//					}
				//					return item.Amount == 0
				//				} else {
				// We always have these items around and work on changing the amount back to normal when it changes
				//				if stockedAmount < item.Amount {
				// We're low on this item, increase toward normal.
				//						changed = true
				//					item.Amount--
				//					if item.Amount == 0 {
				//						shop.Inventory.Remove(item)
				//					}
				//				}
				if shop.Stock.Count(item.ID) > item.Amount {
					item.Amount++
				} else {
					item.Amount--
				}
				//				if stockedAmount > item.Amount {
				// We are overstocked with this item, decrease toward normal.
				//					item.Amount++
				//						changed = true
				//				}

				//				}
				return item.Amount == 0
			})
			if changed {
				shop.Players.RangePlayers(func(player *Player) bool {
					player.SendPacket(ShopOpen(shop))
					return false
				})
			}
			//			if changed {
			//				shop.Players.Range(func(player *Player) {
			//					if player.HasState(StateShopping) && player.CurrentShop() == shop {
			//						player.SendPacket(ShopOpen(shop))
			//					}
			//				})
			//			}
		}
		return false
	})
	return shop
}

func (s *ShopItems) Add(item *Item) {
	if !s.Contains(item.ID) {
		s.Lock()
		defer s.Unlock()
		s.set = append(s.set, item)
	} else {
		s.Get(item.ID).Amount += 1
		//		s.RLock()
		//		defer s.RUnlock()
		//		for _, shopItem := range s.set {
		//			if shopItem.ID == item.ID {
		//				for i := 0; i < item.Amount; i++ {
		//					shopItem.Amount += 1
		//				}
		//			}
		//		}
	}
}

func (s *ShopItems) Size() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.set)
}

func (s *ShopItems) Range(fn func(*Item) bool) {
	s.Lock()
	defer s.Unlock()
	p := 0
	for _, i := range s.set {
		if !fn(i) {
			s.set[p] = i
			p++
		}

	}
	s.set = s.set[:p]
}

func (s *ShopItems) Remove(removingItem *Item) {
	s.Range(func(item *Item) bool {
		if item.ID == removingItem.ID {
			item.Amount--
			return item.Amount == 0
		}
		return false
	})
}

func (s *ShopItems) RemoveID(id, amount int, remove bool) {
	s.Range(func(item *Item) bool {
		if item.ID == id {
			item.Amount -= amount
			return remove && item.Amount <= 0
		}

		return false
	})
}

// Returns: the shop item in this collection  with the given ID, or if it can't find one, creates a new one with
// Amount=0, adds it to the collection for later usage, and returns it
func (s *ShopItems) Get(id int) *Item {
	s.RLock()
	defer s.RUnlock()
	for _, item := range s.set {
		if item.ID == id {
			return item
		}
	}
	return &Item{ID: id}
}

// Ensures safe access when requesting whether this collection contains a specific item by ID.
//
// Returns: true if this shop items collection has any items with the provided ID, otherwise returns false.
func (s *ShopItems) Contains(id int) bool {
	return s.Get(id) != nil
}

// Ensures safe access when requesting the current count of a specific item by ID in this shops inventory.
//
// Returns: the inventorys amount of the specified item id, or 0 is no items matched this ID.
func (s *ShopItems) Count(id int) int {
	item := s.Get(id)
	if item == nil {
		return 0
	}
	return item.Amount
}

//Clone will make a new ShopItems collection populated with clones of the values in the existing collection set.
// Very useful to easily define an initial inventory and stock when defining a shop.
func (s *ShopItems) Clone() *ShopItems {
	clone := &ShopItems{}
	clone.Lock()
	s.Range(func(item *Item) bool {
		clone.set = append(clone.set, &Item{ID: item.ID, Amount: item.Amount})
		return false
	})
	clone.Unlock()
	return clone
}

//Clone makes a clone of the receiver shop and returns it.
func (s *Shop) Clone() *Shop {
	return &Shop{BuysUnstocked: s.BuysUnstocked, BasePurchasePercent: s.BasePurchasePercent, BaseSalePercent: s.BaseSalePercent, Stock: s.Stock.Clone(), Inventory: s.Inventory.Clone()}
}

//DeltaPercentMod calculates the percentage to scale the item's price up or down from its respective base percentage.
// The formula simply subtracts the shop's base stocked amount of item from item's amount, and multiplies the difference
// by ShopNormalDeltaRate.
func (s *Shop) DeltaPercentMod(item *Item) int {
	return item.DeltaAmount(s.Stock.Get(item.ID)) * ShopNormalDeltaRate
}

func (s *Shop) DeltaPercentModID(id int) int {
	return s.Inventory.Get(id).DeltaAmount(s.Stock.Get(id)) * ShopNormalDeltaRate
	//return (s.Stock.Count(id) - s.Inventory.Count(id)) * ShopNormalDeltaRate
}

func (s *Shop) Remove(id int, amount int) bool {
	if s.Inventory.Count(id) < amount {
		return false
	}
	s.Inventory.RemoveID(id, amount, !s.Stock.Contains(id))
	return true
}
