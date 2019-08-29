package list

import (
	"log"
	"os"
)

//logWarning Log interface for warnings.
var logWarning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)

//List Data structure representing a single instance of a list.
type List struct {
	Values []interface{}
	offset int
}

//New Create a new instance of a List and return a reference to it.
func New(cap int) *List {
	return &List{Values: make([]interface{}, cap)}
}

//nextIndex Returns the lowest available index in the receiver list.
func (list *List) nextIndex() int {
	for i, v := range list.Values {
		if v == nil {
			return i
		}
	}

	return -1
}

//Clear Clears the receiver List
func (list *List) Clear() {
	for i := range list.Values {
		list.Values[i] = nil
	}
}

//Add Add a value to list and return the index we stored it at.  If list is full, log it as a warning.
func (list *List) Add(v interface{}) int {
	if idx := list.nextIndex(); idx != -1 {
		list.Values[idx] = v
		return idx
	}
	logWarning.Println("List appears to be full.  Could not insert new value to list.")
	return -1
}

//Get Returns the value at the specified index in list.
func (list *List) Get(idx int) interface{} {
	if idx >= len(list.Values) || idx < 0 {
		return nil
	}
	return list.Values[idx]
}

//Remove Remove a value from the specified `list`, by its index.  Returns true if removed the value at index,
// otherwise returns false.
func (list *List) Remove(index int) bool {
	if v := list.Values[index]; v != nil {
		list.Values[index] = nil
		return true
	}

	logWarning.Printf("Tried removing value that doesn't exist at index %d\n", index)
	return false
}

//Size Returns the number of non-nil elements in the list's backing slice.
func (list *List) Size() (total int) {
	for _, v := range list.Values {
		if v != nil {
			total++
		}
	}
	return
}

//Next Returns the next value in the list.
func (list *List) Next() interface{} {
	next := list.Get(list.offset)
	list.offset++
	return next
}

//Previous Returns the previous value in the list.
func (list *List) Previous() interface{} {
	list.offset--
	return list.Get(list.offset)
}

//HasNext Returns true if there is another value after the current one in the list.
func (list *List) HasNext() bool {
	return list.offset != list.Size()+1
}

//ResetIterator Returns the previous value in the list.
func (list *List) ResetIterator() {
	list.offset = 0
}
