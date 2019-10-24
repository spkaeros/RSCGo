package world

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"sync"
)

//Item Represents a single item in the game.
type Item struct {
	ID     int
	Amount int
	Index  int
}

//GroundItem Represents a single ground item within the game.
type GroundItem struct {
	Item
	Entity
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

//Put Puts an item into the inventory with the specified id and quantity, and returns its index.
func (i *Inventory) Put(id int, qty int) int {
	curSize := i.Size()
	if curSize >= i.Capacity {
		return -1
	}

	newItem := &Item{id, qty, curSize}
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
