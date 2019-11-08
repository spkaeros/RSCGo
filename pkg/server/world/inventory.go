package world

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"fmt"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"go.uber.org/atomic"
	"sync"
	"time"
)

//Item Represents a single item in the game.
type Item struct {
	ID     int
	Amount int
	Index  int
	Worn bool
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

	newItem := &Item{id, qty, curSize, false}
	i.Lock.Lock()
	i.List = append(i.List, newItem)
	i.Lock.Unlock()
	return curSize
}

//Remove Removes item at index from this inventory.
func (i *Inventory) Remove(index int) bool {
	size := i.Size()
	if index >= size {
		log.Suspicious.Printf("Attempted removing item out of inventory bounds.  index:%d,size:%d,capacity:%d\n", index, size, i.Capacity)
		return false
	}
	i.Lock.Lock()
	if index < size-1 {
		copy(i.List[index:], i.List[index+1:])
	}
	i.List[size-1] = nil
	i.List = i.List[:size-1]
	i.Lock.Unlock()
	return true
}

//RemoveByID Removes amt items from this inventory by ID, returns the items index if successful, otherwise returns -1
func (i *Inventory) RemoveByID(id, amt int) int {
	for idx := 0; idx < i.Size(); idx++ {
		item := i.Get(idx)
		if item.ID == id {
			if item.Amount > amt {
				item.Amount -= amt
			} else {
				i.Remove(idx)
			}
			return idx
		}
	}
	return -1
}

//Get Returns a reference to the item at index if it exists, otherwise returns nil.
func (i *Inventory) Get(index int) *Item {
	i.Lock.RLock()
	defer i.Lock.RUnlock()
	if index >= len(i.List) {
		return nil
	}

	return i.List[index]
}

//Get Returns a reference to the item at index if it exists, otherwise returns nil.
func (i *Inventory) GetByID(ID int) *Item {
	i.Lock.RLock()
	defer i.Lock.RUnlock()

	for _, item := range i.List {
		if item.ID == ID {
			return item
		}
	}

	return nil
}

//RemoveAll Removes all of the items in offer from this inventory, returns count of items removed.
func (i *Inventory) RemoveAll(offer *Inventory) int {
	offer.Lock.RLock()
	count := 0
	for _, item := range offer.List {
		if i.RemoveByID(item.ID, item.Amount) != -1 {
			count++
		}
	}
	offer.Lock.RUnlock()
	return count
}

//Clear Clears all items out of the inventory.
func (i *Inventory) Clear() {
	i.Lock.Lock()
	i.List = i.List[:0]
	i.Lock.Unlock()
}


//Equals Returns true if o1 is an object reference with identical characteristics to o.
func (i *Item) Equals(i1 objects.Object) bool {
	if i1, ok := i1.(*Item); ok {
		// We can ignore index, right?
		return i1.ID == i.ID && i.Index == i1.Index && i.Amount == i1.Amount
	}

	return false
}

func (i *Item) TypeName() string {
	return "world.Item"
}

func (i *Item) Copy() objects.Object {
	return &Item{i.ID, i.Amount, i.Index, i.Worn}
}

func (i *Item) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	return nil, objects.ErrInvalidOperator
}

func (i *Item) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", i.ID, i.Amount, i.Index)
}

func (i *Item) IndexGet(index objects.Object) (objects.Object, error) {
	switch index := index.(type) {
	case *objects.String:
		switch index.Value {
		case "id":
			return &objects.Int{Value: int64(i.ID)}, nil
		case "amount":
			return &objects.Int{Value: int64(i.Amount)}, nil
		case "index":
			return &objects.Int{Value: int64(i.Index)}, nil
		}
	}
	return nil, objects.ErrInvalidIndexType
}

func (i *Item) IsFalsy() bool {
	return i.ID < 0 || i.Index < 0
}