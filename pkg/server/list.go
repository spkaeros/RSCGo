package server

//MaxClients The maxium number of active clients supported by this server.
const MaxClients = 2048

//Indexed An interface that any data structures you want to use in a List need to implement.
type Indexed interface {
	Index() int
	SetIndex(i int)
}

//ActiveClients The list of Client references for all of the currently active connected clients.
var ActiveClients = &List{values: make([]Indexed, MaxClients)}

//List Data structure representing a single instance of a Client list.
type List struct {
	values []Indexed
}

//NextIndex Returns the lowest available index in the receiver list.
func (list *List) NextIndex() int {
	for i, v := range list.values {
		if v == nil {
			return i
		}
	}

	return -1
}

//Clear Clears the receiver List
func (list *List) Clear() {
	for i := range list.values {
		list.values[i] = nil
	}
}

//Add Add a Client to the `list`.  If list is full, log it as a warning.
func (list *List) Add(v Indexed) {
	idx := list.NextIndex()
	if idx != -1 {
		list.values[idx] = v
		v.SetIndex(idx)
	} else {
		LogWarning.Println("List appears to be full.  Could not insert new value to list.")
	}
}

//Get Returns the value at the specific index in the receiver List.
func (list *List) Get(idx int) Indexed {
	v := list.values[idx]
	if v == nil {
		LogWarning.Println("Tried to Get value that does not exist from receiver List.")
		return nil
	}
	return v
}

//Remove Remove a value from the specified `list`, by its index.
func (list *List) Remove(index int) {
	v := list.values[index]
	if v != nil {
		list.values[index] = nil
		LogWarning.Printf("Removed: %v\n", v)
	} else {
		LogWarning.Printf("WARNING: Tried removing value that doesn't exist at index %d\n", index)
	}
}

//Size returns the number of non-nil elements in the list's backing slice.
func (list *List) Size() (total int) {
	for _, v := range list.values {
		if v != nil {
			total++
		}
	}
	return
}
