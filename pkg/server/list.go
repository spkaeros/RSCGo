package server

import (
	"fmt"
)

const MaxClients = 2048

var ActiveClients = &ClientList{make(map[int] *client)}

//ClientList Data structure representing a single instance of a client list.
type ClientList struct {
	clients  map[int] *client
}

//NextIndex Returns the lowest available index in the specified receiver list.
func (list *ClientList) NextIndex() int {
	for i := 0; i < MaxClients; i++ {
		if _, ok := list.clients[i]; !ok {
			return i
		}
	}

	return -1
}

//Clear Clears the receiver ClientList instance, and unregisters all of the clients safely.
func (list *ClientList) Clear() {
	for i, c := range list.clients {
		c.unregister()
		delete(list.clients, i)
	}
}

//Add Add a client to the `list`.  If list is full, log it as a warning.
func (list *ClientList) Add(c *client) {
	idx := list.NextIndex()
	if idx != -1 {
		list.clients[idx] = c
		c.index = idx
	} else {
		fmt.Println("WARNING: client list appears to be full.  Could not insert new client to client list.")
	}
}

//Get(int) *client Returns the client at the specific index in the receiver ClientList.
func (list *ClientList) Get(idx int) *client {
	client, ok := list.clients[idx]
	if !ok {
		fmt.Println("WARNING: Tried to Get client that does not exist from receiver ClientList.")
		return nil
	}
	return client
}

//Remove Remove a client from the specified `list`, by the index of the client.
func (list *ClientList) Remove(index int) {
	delete(list.clients, index)
}