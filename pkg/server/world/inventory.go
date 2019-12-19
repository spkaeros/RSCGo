package world

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"go.uber.org/atomic"
	"sort"
	"sync"
	"time"
)

//ItemDefinition This represents a single definition for a single item in the game.
type ItemDefinition struct {
	ID          int
	Name        string
	Description string
	Command     string
	BasePrice   int
	Stackable   bool
	Quest       bool
	Members     bool
}

//ItemDefs This holds the defining characteristics for all of the game's items, ordered by ID.
var ItemDefs []ItemDefinition

type EquipmentDefinition struct {
	ID       int
	Sprite   int
	Type     int
	Armour   int
	Magic    int
	Prayer   int
	Ranged   int
	Aim      int
	Power    int
	Position int
	Female   bool
}

var EquipmentDefs []EquipmentDefinition

func GetEquipmentDefinition(id int) *EquipmentDefinition {
	for _, e := range EquipmentDefs {
		if e.ID == id {
			return &e
		}
	}

	return nil
}

//Item Represents a single item in the game.
type Item struct {
	ID     int
	Amount int
	Worn   bool
}

func (i *Item) Name() string {
	if i.ID >= len(ItemDefs) || i.ID < 0 {
		return "nil"
	}
	return ItemDefs[i.ID].Name
}

func (i *Item) Price() int {
	if i.ID >= len(ItemDefs) || i.ID < 0 {
		return -1
	}
	return ItemDefs[i.ID].BasePrice
}

func (i *Item) Command() string {
	if i.ID >= len(ItemDefs) || i.ID < 0 {
		return "nil"
	}
	return ItemDefs[i.ID].Command
}

func (i *Item) WieldPos() int {
	if i.ID >= len(ItemDefs) || i.ID < 0 {
		return -1
	}
	def := GetEquipmentDefinition(i.ID)
	if def == nil {
		return -1
	}
	return def.Position
}

func (i *Item) Stackable() bool {
	if i.ID >= len(ItemDefs) || i.ID < 0 {
		return false
	}
	return ItemDefs[i.ID].Stackable
}

//GroundItem Represents a single ground item within the game.
type GroundItem struct {
	owner     uint64
	removed   bool
	spawnTime time.Time
	lock      sync.RWMutex
	Item
	*Entity
}

var itemIndexer = atomic.NewUint32(0)

//NewGroundItem Creates a new ground item in the game world and returns a reference to it.
func NewGroundItem(id, amount, x, y int) *GroundItem {
	gi := &GroundItem{owner: strutil.MaxBase37 + 5000, spawnTime: time.Now(), removed: false,
		Item: Item{
			ID:     id,
			Amount: amount,
		}, Entity: &Entity{
			Location: NewLocation(x, y),
			Index:    int(itemIndexer.Swap(itemIndexer.Load() + 1)),
		},
	}
	go func() {
		time.Sleep(time.Minute * 3)
		gi.Remove()
	}()
	return gi
}

//NewGroundItemFor Creates a new ground item with an owner in the game world and returns a reference to it.
func NewGroundItemFor(owner uint64, id, amount, x, y int) *GroundItem {
	gi := &GroundItem{owner: owner, spawnTime: time.Now(), removed: false,
		Item: Item{
			ID:     id,
			Amount: amount,
		}, Entity: &Entity{
			Location: NewLocation(x, y),
			Index:    int(itemIndexer.Swap(itemIndexer.Load() + 1)),
		},
	}
	go func() {
		time.Sleep(time.Minute * 3)
		gi.Remove()
	}()
	return gi
}

//Remove removes the ground item from the world.
func (i *GroundItem) Remove() {
	i.removed = true
	RemoveItem(i)
}

//VisibleTo Returns true if the ground item is visible to this player, otherwise returns false.
func (i *GroundItem) VisibleTo(p *Player) bool {
	if i.removed {
		return false
	}
	if i.owner > strutil.MaxBase37 || p.UsernameHash() == i.owner || p.Rank() == 2 {
		return true
	}
	return time.Since(i.spawnTime) > time.Minute
}

//Inventory Represents an inventory of items in the game.
type Inventory struct {
	List            []*Item
	Owner           *Player
	Capacity        int
	stackEverything bool
	Lock            sync.RWMutex
}

//Shop represents an NPC owned shop in the game
type Shop struct {
	General bool
	SellMultiplier int
	BuyMultiplier int
	InitialStock Inventory
	*Inventory
}

//NewShop returns a reference to a newly created shop with the specified parameters.
func NewShop(general bool, sellMul, buyMul int, items map[int]int) *Shop {
	capacity := len(items)
	if general {
		capacity = 40
	}
	initStock := &Inventory{Capacity: capacity, stackEverything: true}
	for id, amount := range items {
		initStock.Add(id, amount)
	}
	return &Shop{general, sellMul, buyMul, *initStock.Clone(), initStock}
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
func (i *Inventory) DeathDrops(keep int) *Inventory {
	// clone so we don't modify the players inventory during the sorting process
	if keep <= 0 {
		return i.Clone()
	}
	deathItems := i.Clone()
	sort.Sort(itemSorter(deathItems.List))
	if len(deathItems.List) < keep {
		keep = len(deathItems.List)
	}
	for idx := keep; idx > 0; idx-- {
		if deathItems.List[idx-1].Stackable()  {
			keep--
		}
	}
	return &Inventory{List: deathItems.List[keep:], Capacity: 30}
}

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

//Add Puts an item into the inventory with the specified id and quantity, and returns its index.
func (i *Inventory) Add(id int, qty int) int {
	curSize := i.Size()
	if item := i.GetByID(id); (i.stackEverything || ItemDefs[id].Stackable) && item != nil {
		item.Amount += qty
		return i.GetIndex(id)
	}
	if curSize >= i.Capacity {
		return -1
	}

	newItem := &Item{id, qty, false}
	i.Lock.Lock()
	i.List = append(i.List, newItem)
	i.Lock.Unlock()
	return curSize
}

//Remove Removes item at index from this inventory.
func (i *Inventory) Remove(index int) bool {
	item := i.Get(index)
	if item == nil {
		log.Suspicious.Printf("Attempted removing non-existant. item:%v\n", index)
		return false
	}
	if i.Owner != nil && i.Owner.Connected() {
		i.Owner.DequipItem(item)
	}
	i.Lock.Lock()
	defer i.Lock.Unlock()
	size := len(i.List)
	if index >= size {
		log.Suspicious.Printf("Attempted removing item out of inventory bounds.  index:%d,size:%d,capacity:%d\n", index, size, i.Capacity)
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
	if i.stackEverything || ItemDefs[id].Stackable {
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

//Get Returns a reference to the item at index if it exists, otherwise returns nil.
func (i *Inventory) GetByID(ID int) *Item {
	return i.Get(i.GetIndex(ID))
}

//Get returns the index of the first item with the provided ID.
func (i *Inventory) GetIndex(ID int) int {
	return i.Range(func(item *Item) bool {
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
	return fmt.Sprintf("[%v, (%v, %v)]", i.ID, i.Amount)
}
