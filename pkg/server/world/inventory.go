package world

//Item Represents a single item in the game.
type Item struct {
	ID     int
	Amount int
	Index  int
}

//Inventory Represents an inventory of items in the game.
type Inventory struct {
	List     []*Item
	Capacity int
}

//Put Puts an item into the inventory with the specified id and quantity, and returns its index.
func (i *Inventory) Put(id int, qty int) int {
	if len(i.List) >= i.Capacity {
		return -1
	}

	newItem := &Item{id, qty, len(i.List)}
	i.List = append(i.List, newItem)
	return newItem.Index
}

func (i *Inventory) Remove(index int) bool {
	curSize := len(i.List)
	if curSize < index {
		return false
	}

	i.List = append(i.List[:index], i.List[index+1:]...)
	return true
}
