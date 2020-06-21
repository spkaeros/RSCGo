package world

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/tasks"
	"github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/log"
)

//DefaultDrop returns the default item ID all mobs should drop on death
const DefaultDrop = 20

type ItemBubble struct {
	Owner *Player
	Item  int
}

//Item Represents a single item in the game.
type Item struct {
	ID     int
	Amount int
	Worn   bool
	Index  int
}

//Name returns the receivers name
func (i *Item) Name() string {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return "nil"
	}
	return definitions.Items[i.ID].Name
}

//Price type alias for item prices
type Price int

//Price returns the receivers base price
func (i *Item) Price() Price {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return -1
	}
	return Price(definitions.Items[i.ID].BasePrice)
}

//DeltaAmount returns the difference between the amount of o and the amount of the receiver
func (i *Item) DeltaAmount(o *Item) int {
	return o.Amount - i.Amount
}

//ScalePrice returns the receivers base price, scaled by percent%.
func (i *Item) ScalePrice(percent int) int {
	return int(i.Price().Scale(percent))
}

//Scale Calculates and returns the value for the requested percentage of the receiver price.
//
// In other words, for sleeping bag with basePrice=30
//	player.Inventory.GetByID(1263).PriceScaled(100)
// would be 30 and is the same as calling Price().
// Any percent that is higher than 100 will scale the price up.  E.g:
//	player.Inventory.GetByID(1263).PriceScaled(130)
// would be 39; since 130%(??) of the base price is the value we want, to reach that from the base price(100%(30)), we
// can add 30%(9) to the base price(100%(30)), which gives us 130%(39), the value we want.
//
// This is the same way RSC general stores priced items they sold.  Additionally, though,
//
// Any percent that is lower than 100 will scale the price down.  E.g:
//	player.Inventory.GetByID(1263).PriceScaled(40)
// would be 12. Since 40%(??) of the base price is our target, to reach that from 100%(30), we subtract 60%(18) from it.
// This is the same percentage used for RSC general stores initial sale prices.
//, we'd. is how we mimic canonical RSClassic general store pricing.
//
// Upper bound for percent intended to basically not exist; in practice it's limited by the data type of the argument.
// Lower bound for percent is 10, anything lower will be treated as if it were 10%.
func (p Price) Scale(percent int) Price {
	return Price(int(math.Max(10, float64(percent)))*int(p)) / 100
}

//Command Returns the item command, or nil if none
func (i *Item) Command() string {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return "nil"
	}
	return definitions.Items[i.ID].Command
}

//WieldPos Returns the item equip slot, or -1 if none
func (i *Item) WieldPos() int {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return -1
	}
	def := definitions.Equip(i.ID)
	if def == nil {
		return -1
	}
	return def.Position
}

//Stackable Returns true if the item is stackable, false otherwise.
func (i *Item) Stackable() bool {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return false
	}
	return definitions.Items[i.ID].Stackable
}

//GroundItem Represents a single ground item within the game.
type GroundItem struct {
	ID, Amount int
	Owner string
	*entity.AttributeList
	Entity
}

//ItemIndexer Ensures unique indexes for ground items.
//  TODO: Proper indexing
var ItemIndexer = atomic.NewUint32(0)

//Visibility This is a special state attribute to indicate who the receiver item is visible to.
// Value 0 means the item has expired and is no longer visible to anybody.
//
// Value 1 means the item is visible to only the Owner of it and game administrators(rank=2), e.g if you kill someone or something,
// this will be the value for the first minute or two after it is created, and then will change to value 2.
// NOTE: If this is the current value, but the belongsTo attribute is not set e.g nobody owns this item, it will update
// itself to value 2 prior to when it normally would.
//
// Value 2 means the item is visible to all players.  This is the value when e.g the game starts and makes the worlds
// default item spawns, or an NPC kills a player and they drop their items and/or bones...This is the state that most
// transient ground items will likely spend the most time in before unsetting the visibility attribute(same as value=0)
// and thus disappearing.
func (i *GroundItem) Visibility() int {
	return i.VarInt("visibility", 0)
}

//SpawnedTime Returns: the time this item was spawned into the game world.
func (i *GroundItem) SpawnedTime() time.Time {
	return i.VarTime("spawnTime")
}

//NewPersistentGroundItem Returns a new ground item that respawns at a set rate after pickup.
func NewPersistentGroundItem(id, amount, x, y, respawn int) *GroundItem {
	item := &GroundItem{ID: id, Amount: amount,
		AttributeList: entity.NewAttributeList(),
		Entity: Entity{
			Location: NewLocation(x, y),
			Index:    int(ItemIndexer.Swap(ItemIndexer.Load() + 1)),
		},
	}
	item.SetVar("visibility", 2)
	item.SetVar("respawnTime", respawn)
	item.SetVar("persistent", true)
	return item
}

//NewGroundItem Creates a new ground item in the game world and returns a reference to it.
func NewGroundItem(id, amount, x, y int) *GroundItem {
	item := &GroundItem{ID: id, Amount: amount,
		AttributeList: entity.NewAttributeList(),
		Entity: Entity{
			Location: NewLocation(x, y),
			Index:    int(ItemIndexer.Swap(ItemIndexer.Load() + 1)),
		},
	}
	item.SetVar("visibility", 1)
	tasks.TickList.Add(func() bool {
		item.Inc("ticker", 1)
		curTick := item.VarInt("ticker", 0)
		// Visiblity is scoped to item owner but I guess it doesn't have an owner.
		// Oh well, we'll just let everyone see it early.
		if item.Visibility() == 1 && len(item.Owner) == 0 {
			item.SetVar("visibility", 2)
		}
		// This keeps track of how many times we've ticked for ~71 sec since we started.
		stage := curTick / 110
		// ~71 sec.
		if curTick%110 == 0 {
			// item only seen by owner and administrators

			if item.Visibility() == 1 {
				if stage == 1 {
					// 25% chance to stay visibility=1 until 2nd pass at ~142s...
					if Chance(25) {
						return false
					}

					item.SetVar("visibility", 2)
					return false
				}
			}
			// Time for everyone to see it!
		}
		if stage >= 3 {
			item.Remove()
		}
		return item.Visibility() == 0 || stage >= 3
	})

	item.SetVar("spawnedTime", time.Now())
	return item
}

//NewGroundItemFor Creates a new ground item with an Owner in the game world and returns a reference to it.
func NewGroundItemFor(owner uint64, id, amount, x, y int) *GroundItem {
	item := NewGroundItem(id, amount, x, y)
	item.Owner = strutil.Base37.Decode(owner)
	return item
}

//Name returns the receivers name
func (i *GroundItem) Name() string {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return "nil"
	}
	return definitions.Items[i.ID].Name
}

//Price returns the receivers base price
func (i *GroundItem) Price() Price {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return -1
	}
	return Price(definitions.Items[i.ID].BasePrice)
}

//DeltaAmount returns the difference between the amount of o and the amount of the receiver
func (i *GroundItem) DeltaAmount(o *Item) int {
	return o.Amount - i.Amount
}

//ScalePrice returns the receivers base price, scaled by percent%.
func (i *GroundItem) ScalePrice(percent int) int {
	return int(i.Price().Scale(percent))
}

//Command Returns the command for this item, or nil if none.
func (i *GroundItem) Command() string {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return "nil"
	}
	return definitions.Items[i.ID].Command
}

//WieldPos Returns the equip slot for this item, or -1 if none.
func (i *GroundItem) WieldPos() int {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return -1
	}
	def := definitions.Equip(i.ID)
	if def == nil {
		return -1
	}
	return def.Position
}

//Stackable returns true if the items stackable, otherwise returns false.
func (i *GroundItem) Stackable() bool {
	if i.ID >= len(definitions.Items) || i.ID < 0 {
		return false
	}
	return definitions.Items[i.ID].Stackable
}

//Remove removes the ground item from the world.
func (i *GroundItem) Remove() {
	i.UnsetVar("visibility")
	RemoveItem(i)
	if i.VarBool("persistent", false) {
		go func() {
			time.Sleep(time.Second * time.Duration(i.VarInt("respawnTime", 10)))
			i.SetVar("visibility", 2)
			AddItem(i)
		}()
	}
}

//VisibleTo Returns true if the ground item is visible to this player, otherwise returns false.
func (i *GroundItem) VisibleTo(p *Player) bool {
	if i.Visibility() == 0 {
		// removing from world
		return false
	}

	if p.Rank() == 2 {
		// admins see everything >.<
		return true
	}

	if i.Visibility() == 1 {
		// Owner of item is only one we currently want seeing this
		return i.Owner == p.Username()
	}

	return i.Visibility() == 2
}

//Inventory Represents an inventory of items in the game.
type Inventory struct {
	List            []*Item
	Owner           *Player
	Capacity        int
	stackEverything bool
	Lock            sync.RWMutex
}

type itemSorter []*Item

func (s itemSorter) Len() int {
	return len(s)
}

func (s itemSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s itemSorter) Less(i, j int) bool {
	if !s[j].Stackable() && !s[i].Stackable() {
		return s[i].Price() > s[j].Price()
	}
	return !s[i].Stackable() && s[j].Stackable()
}

//Clone returns a copy of this inventory in a new inventory.
func (i *Inventory) Clone() *Inventory {
	var newList []*Item
	i.Range(func(item *Item) bool {
		newList = append(newList, item)
		return true
	})
	return &Inventory{List: newList, Capacity: i.Capacity, stackEverything: i.stackEverything}
}

//DeathDrops returns a list of items to drop upon dying.  It decides what to keep using a few simple rules:
// If the item is stackable, it gets dropped no matter what, even if it is the only item and keep is 3.
// If the item isn't stackable, the inventory is first sorted by descending BasePrice, and the first `keep` items are
// sliced off of the top of this sorted list which leaves us with the 30-keep least valuable items.
func (i *Inventory) DeathDrops(keep int) []*GroundItem {
	// clone so we don't modify the players inventory during the sorting process
	var pile []*GroundItem
	if keep <= 0 {
		i.Lock.RLock()
		for _, item := range i.List {
			pile = append(pile, NewGroundItem(item.ID, item.Amount, i.Owner.X(), i.Owner.Y()))
		}
		i.Lock.RUnlock()
	} else {
		deathItems := i.Clone()
		sort.Sort(itemSorter(deathItems.List))
		if len(deathItems.List) < keep {
			keep = len(deathItems.List)
		}
		for idx := keep; idx > 0; idx-- {
			if deathItems.List[idx-1].Stackable() {
				keep--
			}
		}
		for _, item := range deathItems.List[keep:] {
			pile = append(pile, NewGroundItem(item.ID, item.Amount, i.Owner.X(), i.Owner.Y()))
		}
	}
	return pile
}

//Range Calls fn for each item in the inventory list.
func (i *Inventory) Range(fn func(*Item) bool) int {
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	for idx, item := range i.List {
		if !fn(item) {
			return idx
		}
	}
	return -1
}

//RangeRev Calls fn for each item in the inventory list, in reverse.
func (i *Inventory) RangeRev(fn func(*Item) bool) int {
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	for idx := len(i.List) - 1; idx >= 0; idx-- {
		if !fn(i.List[idx]) {
			return idx
		}
	}
	return -1
}

//Equipped Returns true if the item with the given ID exists and is wielded in this inventory
func (i *Inventory) Equipped(id int) bool {
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	for _, item := range i.List {
		if item.ID == id && item.Worn {
			return true
		}
	}
	return false
}

//Size Returns the number of items currently in this inventory.
func (i *Inventory) Size() int {
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	return len(i.List)
}

//CanHold returns true if this inventory can hold the specified amount of the item with the specified ID
func (i *Inventory) CanHold(id, amount int) bool {
	var slotsReq int
	if definitions.Items[id].Stackable || i.stackEverything {
		if i.GetByID(id) == nil {
			slotsReq += 1 + (amount / math.MaxInt32)
		} else {
			for i.GetByID(id).Amount+amount > math.MaxInt32 {
				slotsReq++
				amount -= math.MaxInt32
			}
		}
	} else {
		slotsReq++
	}
	return i.Size()+slotsReq-1 < i.Capacity
}

//Add Puts an item into the inventory with the specified id and quantity, and returns its index.
func (i *Inventory) Add(id int, qty int) int {
	if qty == 0 {
		return -1
	}
	if !i.CanHold(id, qty) {
		AddItem(NewGroundItemFor(i.Owner.UsernameHash(), id, qty, i.Owner.X(), i.Owner.Y()))
		i.Owner.Message("Your inventory is full, the " + definitions.Items[id].Name + " drops to the ground!")
		return -1
	}
	if item := i.GetByID(id); (i.stackEverything || definitions.Items[id].Stackable) && item != nil {
		if item.Amount < 0 {
			log.Suspicious.Println(errors.NewArgsError("*Inventory.Add(id,amt) Resulting item amount less than zero: " + strconv.FormatUint(uint64(item.Amount+qty), 10)))
		}
		if item.Amount+qty > math.MaxInt32 {
			item.Amount = math.MaxInt32
			return i.GetIndex(id)
		}
		item.Amount += qty
		return i.GetIndex(id)
	}

	if qty < 0 {
		return -1
	}

	newItem := &Item{ID: id, Amount: qty}
	i.Lock.Lock()
	i.List = append(i.List, newItem)
	i.Lock.Unlock()
	return i.Size() - 1
}

//Remove Removes item at index from this inventory.
func (i *Inventory) Remove(index int) bool {
	item := i.Get(index)
	if item == nil {
		log.Cheatf("Attempted removing non-existent. item:%v\n", index)
		return false
	}
	if i.Owner != nil && i.Owner.Connected() {
		i.Owner.DequipItem(item)
	}
	i.Lock.Lock()
	defer i.Lock.Unlock()
	size := len(i.List)
	if index >= size {
		log.Cheatf("Attempted removing item out of inventory bounds.  index:%d,size:%d,capacity:%d\n", index, size, i.Capacity)
		return false
	}
	if index >= size-1 {
		i.List = i.List[:index]
		return true
	}
	i.List = append(i.List[:index], i.List[index+1:]...)
	return true
}

//RemoveByID Removes amt items from this inventory by ID, returns the items index if successful, otherwise returns -1
func (i *Inventory) RemoveByID(id, amt int) int {
	if i.CountID(id) < amt {
		return -1
	}
	index := i.GetIndex(id)
	if i.stackEverything || definitions.Items[id].Stackable {
		if i.Get(index).Amount == amt {
			i.Remove(index)
		} else {
			i.Get(index).Amount -= amt
		}
	} else {
		for j := 0; j < amt; j++ {
			i.Remove(i.GetIndex(id))
		}
	}
	if i.Owner != nil && i.Owner.Connected() {
		i.Owner.SendInventory()
	}
	return index
}

//Get Returns a reference to the item at index if it exists, otherwise returns nil.
func (i *Inventory) Get(index int) *Item {
	if index >= i.Size() || index < 0 {
		return nil
	}
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	return i.List[index]
}

//GetByID Returns a reference to the item at index if it exists, otherwise returns nil.
func (i *Inventory) GetByID(ID int) *Item {
	return i.Get(i.GetIndex(ID))
}

//GetIndex returns the index of the first item with the provided ID.
func (i *Inventory) GetIndex(ID int) int {
	return i.RangeRev(func(item *Item) bool {
		if item.ID == ID {
			return false
		}
		return true
	})
}

//CountID Returns the total amount of all the items with this ID in this inventory.
func (i *Inventory) CountID(id int) int {
	count := 0
	i.Range(func(item *Item) bool {
		if item.ID == id {
			count += item.Amount
		}
		return true
	})
	if count == 0 {
		return -1
	}
	return count
}

//RemoveAll Removes all of the items in offer from this inventory, returns count of items removed.
func (i *Inventory) RemoveAll(offer *Inventory) int {
	count := 0
	offer.Range(func(item *Item) bool {
		if i.RemoveByID(item.ID, item.Amount) > -1 {
			count++
		}
		return true
	})
	return count
}

//Clear Clears all items out of the inventory.
func (i *Inventory) Clear() {
	i.Lock.Lock()
	defer i.Lock.Unlock()
	i.List = i.List[:0]
}

func (i *Item) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", i.Index, i.ID, i.Amount)
}
