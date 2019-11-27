package clients

import (
	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"sync"
)

//Client Represents a client
type Client interface {
	SendPacket(packet *packet.Packet)
	Player() *world.Player
	Message(string)
	Destroy()
	HandleRegister(chan byte)
	HandleLogin(chan byte)
	TradeOpen()
	UpdatePlane()
	Teleport(int, int)
	TeleBubble(int, int)
	UpdateStat(int)
	UpdatePositions()
	ResetUpdateFlags()
}

//Clients Collection containing all of the active client, by index and username hash, guarded by a mutex
var Clients = &struct {
	usernames map[uint64]Client
	indices   map[int]Client
	lock      sync.RWMutex
}{usernames: make(map[uint64]Client), indices: make(map[int]Client)}

//FromUserHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func FromUserHash(hash uint64) (Client, bool) {
	Clients.lock.RLock()
	result, ok := Clients.usernames[hash]
	Clients.lock.RUnlock()
	return result, ok
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func ContainsHash(hash uint64) bool {
	_, ret := FromUserHash(hash)
	return ret
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func FromIndex(index int) (Client, bool) {
	Clients.lock.RLock()
	result, ok := Clients.indices[index]
	Clients.lock.RUnlock()
	return result, ok
}

//Add Puts a client into the map.
func Put(c Client) {
	nextIndex := NextIndex()
	Clients.lock.Lock()
	c.Player().Index = nextIndex
	Clients.usernames[c.Player().UserBase37] = c
	Clients.indices[nextIndex] = c
	Clients.lock.Unlock()
}

//Remove Removes a client from the map.
func Remove(c Client) {
	Clients.lock.Lock()
	delete(Clients.usernames, c.Player().UserBase37)
	delete(Clients.indices, c.Player().Index)
	Clients.lock.Unlock()
}

//Range Calls action for every active client in the collection.
func Range(action func(Client)) {
	Clients.lock.RLock()
	for _, c := range Clients.indices {
		if c != nil && c.Player().TransAttrs.VarBool("connected", false) {
			action(c)
		}
	}
	Clients.lock.RUnlock()
}

//Size Returns the size of the active client collection.
func Size() int {
	Clients.lock.RLock()
	defer Clients.lock.RUnlock()
	return len(Clients.usernames)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func NextIndex() int {
	Clients.lock.RLock()
	defer Clients.lock.RUnlock()
	for i := 0; i < config.MaxPlayers(); i++ {
		if _, ok := Clients.indices[i]; !ok {
			return i
		}
	}
	return -1
}

//BroadcastLogin Broadcasts the login status of player to the whole server.
func BroadcastLogin(player *world.Player, online bool) {
	Range(func(c Client) {
		if c.Player().Friends(player.UserBase37) {
			if !player.FriendBlocked() || player.Friends(c.Player().UserBase37) {
				c.Player().FriendList[player.UserBase37] = online
				c.SendPacket(packetbuilders.FriendUpdate(player.UserBase37, online))
			}
		}
	})
}
