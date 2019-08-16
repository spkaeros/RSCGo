package server

const MaxClients = 2048

var ActiveClients = &ClientList{clients: make([]interface{}, MaxClients)}

//ClientList Data structure representing a single instance of a Client list.
type ClientList struct {
	clients []interface{}
}

//NextIndex Returns the lowest available index in the specified receiver list.
func (list *ClientList) NextIndex() int {
	for i, c := range list.clients {
		if c == nil {
			return i
		}
	}

	return -1
}

//Clear Clears the receiver ClientList instance, and unregisters all of the clients safely.
func (list *ClientList) Clear() {
	for i := range list.clients {
		list.clients[i] = nil
	}
}

//Add Add a Client to the `list`.  If list is full, log it as a warning.
func (list *ClientList) Add(c interface{}) {
	idx := list.NextIndex()
	if idx != -1 {
		list.clients[idx] = c
		c.(*Client).index = idx
	} else {
		LogWarning.Printf("WARNING: Client list appears to be full.  Could not insert new Client to Client list.")
	}
}

//Get Returns the Client at the specific index in the receiver ClientList.
func (list *ClientList) Get(idx int) *Client {
	c := list.clients[idx]
	if c == nil {
		LogWarning.Printf("WARNING: Tried to Get Client that does not exist from receiver ClientList.")
		return nil
	}
	return c.(*Client)
}

//Remove Remove a Client from the specified `list`, by the index of the Client.
func (list *ClientList) Remove(index int) {
	c := list.clients[index]
	if c != nil {
		list.clients[index] = nil
		LogWarning.Printf("Removed client: %v\n", c)
	} else {
		LogWarning.Printf("WARNING: Tried removing nil client at index %d\n", index)
	}
}
