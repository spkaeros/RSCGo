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
	v := list.Values[idx]
	if v == nil {
		logWarning.Printf("Tried to get value that does not exist from list at index %d\n", idx)
		return nil
	}
	return v
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
