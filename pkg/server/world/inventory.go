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
	Index  int
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
	return &GroundItem{owner: strutil.MaxBase37 + 5000, spawnTime: time.Now(), removed: false,
		Item: Item{
			ID:     id,
			Amount: amount,
		}, Entity: &Entity{
			Location: NewLocation(x, y),
			Index:    int(itemIndexer.Swap(itemIndexer.Load() + 1)),
		},
	}
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
	if i.owner > strutil.MaxBase37 || p.UserBase37 == i.owner || p.Rank == 2 {
		return true
	}
	return time.Since(i.spawnTime) > time.Minute
}

//Inventory Represents an inventory of items in the game.
type Inventory struct {
	List            []*Item
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
	if !s[j].Stackable() && !s[i].Stackable(){
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

func (i *Inventory) DeathDrops(keep int) *Inventory {
	// clone so we don't modify the players inventory during the sorting process
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
	return &Inventory{List: deathItems.List[keep:], Capacity:30}
}

func (i *Inventory) Range(fn func(*Item) bool) int {
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	index := 0
	for _, item := range i.List {
		if !fn(item) {
			break
		}
		index++
	}
	return index
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
		return item.Index
	}
	if curSize >= i.Capacity {
		return -1
	}

	newItem := &Item{id, qty, curSize, false}
	i.Lock.Lock()
	i.List = append(i.List, newItem)
	i.Lock.Unlock()
	return curSize
}

//Remove Removes item at index from this inventory.
func (i *Inventory) Remove(index int, amt int) bool {
	size := i.Size()
	if index >= size {
		log.Suspicious.Printf("Attempted removing item out of inventory bounds.  index:%d,size:%d,capacity:%d\n", index, size, i.Capacity)
		return false
	}
	i.Lock.Lock()
	defer i.Lock.Unlock()
	item := i.List[index]
	if item == nil || item.Amount < amt {
		log.Suspicious.Printf("Attempted removing too much of an item.  item:%v, removeAmt:%v\n", item, amt)
		return false
	}
	item.Amount -= amt
	if (i.stackEverything || ItemDefs[item.ID].Stackable) && item.Amount > 0 {
		return true
	}
	if index >= size-1 {
		i.List = i.List[:index]
		return true
	}
	i.List = append(i.List[:index], i.List[index+1:]...)
	for idx := range i.List {
		i.List[idx].Index = idx
	}
	return true
}

//RemoveByID Removes amt items from this inventory by ID, returns the items index if successful, otherwise returns -1
func (i *Inventory) RemoveByID(id, amt int) int {
	if i.CountID(id) < amt {
		return -1
	}
	if i.stackEverything || ItemDefs[id].Stackable {
		item := i.GetByID(id)
		if i.Remove(item.Index, amt) {
			return item.Index
		}
	} else {
		for j := 0; j < amt; j++ {
			item := i.GetByID(id)
			if i.Remove(item.Index, 1) {
				return item.Index
			}
		}
	}
	return -1
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
	idx := i.Range(func(item *Item) bool {
		if item.ID == ID {
			return false
		}
		return true
	})
	return i.Get(idx)
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
	return fmt.Sprintf("[%v, (%v, %v)]", i.ID, i.Amount, i.Index)
}
