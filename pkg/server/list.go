package server

import (
	"fmt"
)

const MaxClients = 2048

var ActiveClients = &ClientList{clients: make([]*Client, MaxClients)}

//ClientList Data structure representing a single instance of a Client list.
type ClientList struct {
	clients []*Client
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
	for i, c := range list.clients {
		c.Unregister()
		list.clients[i] = nil
	}
}

//Add Add a Client to the `list`.  If list is full, log it as a warning.
func (list *ClientList) Add(c *Client) {
	idx := list.NextIndex()
	if idx != -1 {
		list.clients[idx] = c
		c.index = idx
	} else {
		fmt.Println("WARNING: Client list appears to be full.  Could not insert new Client to Client list.")
	}
}

//Get(int) *Client Returns the Client at the specific index in the receiver ClientList.
func (list *ClientList) Get(idx int) *Client {
	c := list.clients[idx]
	if c == nil {
		fmt.Println("WARNING: Tried to Get Client that does not exist from receiver ClientList.")
		return nil
	}
	return c
}

//Remove Remove a Client from the specified `list`, by the index of the Client.
func (list *ClientList) Remove(index int) {
	list.clients[index] = nil
}
