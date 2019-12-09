package players

import (
	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"sync"
)

//Players Collection containing all of the active client, by index and username hash, guarded by a mutex
var Players = &struct {
	usernames map[uint64]*world.Player
	indices   map[int]*world.Player
	lock      sync.RWMutex
}{usernames: make(map[uint64]*world.Player), indices: make(map[int]*world.Player)}

//FromUserHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func FromUserHash(hash uint64) (*world.Player, bool) {
	Players.lock.RLock()
	result, ok := Players.usernames[hash]
	Players.lock.RUnlock()
	return result, ok
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func ContainsHash(hash uint64) bool {
	_, ret := FromUserHash(hash)
	return ret
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func FromIndex(index int) (*world.Player, bool) {
	Players.lock.RLock()
	result, ok := Players.indices[index]
	Players.lock.RUnlock()
	return result, ok
}

//Add Puts a client into the map.
func Put(player *world.Player) {
	nextIndex := NextIndex()
	Players.lock.Lock()
	player.Index = nextIndex
	Players.usernames[player.UsernameHash()] = player
	Players.indices[nextIndex] = player
	Players.lock.Unlock()
}

//Remove Removes a client from the map.
func Remove(player *world.Player) {
	Players.lock.Lock()
	delete(Players.usernames, player.UsernameHash())
	delete(Players.indices, player.Index)
	Players.lock.Unlock()
}

//Range Calls action for every active client in the collection.
func Range(action func(*world.Player)) {
	Players.lock.RLock()
	for _, c := range Players.indices {
		if c != nil && c.Connected() {
			action(c)
		}
	}
	Players.lock.RUnlock()
}

//Size Returns the size of the active client collection.
func Size() int {
	Players.lock.RLock()
	defer Players.lock.RUnlock()
	return len(Players.usernames)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func NextIndex() int {
	Players.lock.RLock()
	defer Players.lock.RUnlock()
	for i := 0; i < config.MaxPlayers(); i++ {
		if _, ok := Players.indices[i]; !ok {
			return i
		}
	}
	return -1
}

//BroadcastLogin Broadcasts the login status of player to the whole server.
func BroadcastLogin(player *world.Player, online bool) {
	Range(func(rangedPlayer *world.Player) {
		if player.Friends(rangedPlayer.UsernameHash()) {
			if !rangedPlayer.FriendBlocked() || rangedPlayer.Friends(rangedPlayer.UsernameHash()) {
				player.FriendList[rangedPlayer.UsernameHash()] = online
			}
		}
		if rangedPlayer.Friends(player.UsernameHash()) {
			if !player.FriendBlocked() || player.Friends(rangedPlayer.UsernameHash()) {
				rangedPlayer.FriendList[player.UsernameHash()] = online
				rangedPlayer.SendPacket(world.FriendUpdate(player.UsernameHash(), online))
			}
		}
	})
}
