package world

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"go.uber.org/atomic"
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

var Equipment []EquipmentDefinition

//ItemDefs This holds the defining characteristics for all of the game's items, ordered by ID.
var ItemDefs []ItemDefinition

//Item Represents a single item in the game.
type Item struct {
	ID     int
	Amount int
	Index  int
	Worn   bool
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

//NewGroundItemFrom Creates a new ground item with an owner in the game world and returns a reference to it.
func NewGroundItemFrom(owner uint64, id, amount, x, y int) *GroundItem {
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

//Remove Removes the ground item from the world.
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
	List     []*Item
	Capacity int
	Lock     sync.RWMutex
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
	if curSize >= i.Capacity {
		return -1
	}
	if item := i.GetByID(id); ItemDefs[id].Stackable && item != nil {
		item.Amount += qty
		return item.Index
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
	if i.List[index] == nil || i.List[index].Amount < amt {
		log.Suspicious.Printf("Attempted removing too much of an item.  item:%v, removeAmt:%v\n", i.List[index], amt)
		return false
	}
	if ItemDefs[i.List[index].ID].Stackable && i.List[index].Amount > amt {
		i.List[index].Amount -= amt
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
	size := i.Size()
	for idx := 0; idx < size; idx++ {
		if item := i.Get(idx); item != nil && item.ID == id && i.Remove(idx, amt) {
			return idx
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

//RemoveAll Removes all of the items in offer from this inventory, returns count of items removed.
func (i *Inventory) RemoveAll(offer *Inventory) int {
	count := 0
	offer.Range(func(item *Item) bool {
		if i.Remove(item.Index, item.Amount) {
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
