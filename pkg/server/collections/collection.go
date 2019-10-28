package collections

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"sync"
)

//ClientMap A thread-safe concurrent collection type for storing client references.
type ClientMap struct {
	usernames map[uint64]Client
	indices   map[int]Client
	lock      sync.RWMutex
}

//Client Represents a client
type Client interface {
	SendPacket(packet *packetbuilders.Packet)
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
var Clients = &ClientMap{usernames: make(map[uint64]Client), indices: make(map[int]Client)}

//FromUserHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func (m *ClientMap) FromUserHash(hash uint64) (Client, bool) {
	m.lock.RLock()
	result, ok := m.usernames[hash]
	m.lock.RUnlock()
	return result, ok
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func (m *ClientMap) ContainsHash(hash uint64) bool {
	_, ret := m.FromUserHash(hash)
	return ret
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func (m *ClientMap) FromIndex(index int) (Client, bool) {
	m.lock.RLock()
	result, ok := m.indices[index]
	m.lock.RUnlock()
	return result, ok
}

//Put Puts a client into the map.
func (m *ClientMap) Put(c Client) {
	nextIndex := m.NextIndex()
	m.lock.Lock()
	c.Player().Index = nextIndex
	m.usernames[c.Player().UserBase37] = c
	m.indices[nextIndex] = c
	m.lock.Unlock()
}

//Remove Removes a client from the map.
func (m *ClientMap) Remove(c Client) {
	m.lock.Lock()
	delete(m.usernames, c.Player().UserBase37)
	delete(m.indices, c.Player().Index)
	m.lock.Unlock()
}

//Range Calls action for every active client in the collection.
func (m *ClientMap) Range(action func(Client)) {
	m.lock.RLock()
	for _, c := range m.indices {
		if c != nil && c.Player().TransAttrs.VarBool("connected", false) {
			action(c)
		}
	}
	m.lock.RUnlock()
}

//Size Returns the size of the active client collection.
func (m *ClientMap) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.usernames)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func (m *ClientMap) NextIndex() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for i := 0; i < config.MaxPlayers(); i++ {
		if _, ok := m.indices[i]; !ok {
			return i
		}
	}
	return -1
}

//BroadcastLogin Broadcasts the login status of player to the whole server.
func BroadcastLogin(player *world.Player, online bool) {
	Clients.Range(func(c Client) {
		if c.Player().Friends(player.UserBase37) {
			if !player.FriendBlocked() || player.Friends(c.Player().UserBase37) {
				c.Player().FriendList[player.UserBase37] = online
				c.SendPacket(packetbuilders.FriendUpdate(player.UserBase37, online))
			}
		}
	})
}