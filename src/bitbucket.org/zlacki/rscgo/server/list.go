package server

import (
	"fmt"
)

var activeClients = &clientList{}

//clientList Data structure representing a single instance of a client list.
type clientList struct {
	clients [2048]*client
}

//findLowestAvailableIndex Returns the lowest available index in the specified receiver list.
func (list *clientList) findLowestAvailableIndex() int {
	for i, c := range list.clients {
		if c == nil {
			return i
		}
	}

	return -1
}

//clear Clears the receiver clientList instance, and unregisters all of the clients safely.
func (list *clientList) clear() {
	for i, c := range list.clients {
		c.unregister()
		list.clients[i] = nil
	}
}

//add add a client to the `list`.  If list is full, log it as a warning.
func (list *clientList) add(c *client) {
	idx := list.findLowestAvailableIndex()
	if idx != -1 {
		list.clients[idx] = c
		c.index = idx
	} else {
		fmt.Println("WARNING: client list appears to be full.  Could not insert new client to client list.")
	}
}

//get(int) *client Returns the client at the specific index in the receiver clientList.
func (list *clientList) get(idx int) *client {
	client := list.clients[idx]
	if client == nil {
		fmt.Println("WARNING: Tried to get client that does not exist from receiver clientList.")
		return nil
	}
	return client
}

//remove remove a client from the specified `list`, by the index of the client.
func (list *clientList) remove(index int) {
	list.clients[index] = nil
}
