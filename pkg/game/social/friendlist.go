/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package social

import (
	"sync"

	"github.com/spkaeros/rscgo/pkg/strutil"
)

type friendSet map[uint64]bool

//FriendsList A map with base37 integer representations of player usernames mapped to their login status.
// The responsibility of managing the login statuses falls upon the caller; it will not attempt to grab data from
// anywhere else to automatically update statuses of the players it lists, it is just concurrency-safe data-set
// with an intuitive API designed to promote external mutability of its contents.
type FriendsList struct {
	sync.RWMutex
	friendSet
	Owner uint64
}

// TODO: Should I remove the owner username hash from this?  Is it really beneficial to use like this
// TODO: or is it considered harmful?
//New returns a new friends list reference.
func New() *FriendsList {
	return &FriendsList{friendSet: make(friendSet)}
}

func (f *FriendsList) Contains(name string) bool {
	f.RLock()
	defer f.RUnlock()
	_, ok := f.friendSet[strutil.Base37.Encode(name)]
	return ok
}

//Status Returns the online status for a given player username.
func (f *FriendsList) Status(name string) bool {
	f.RLock()
	defer f.RUnlock()
	status, ok := f.friendSet[strutil.Base37.Encode(name)]
	if ok {
		return status
	}
	return false
}

//Status Returns the online status for a given player username hash.
func (f *FriendsList) StatusHash(name uint64) bool {
	f.RLock()
	defer f.RUnlock()
	status, ok := f.friendSet[name]
	if ok {
		return status
	}
	return false
}

func (f *FriendsList) ContainsHash(hash uint64) bool {
	f.RLock()
	defer f.RUnlock()
	_, ok := f.friendSet[hash]
	return ok
}

//Add Deprecated: See Set(string,bool)
func (f *FriendsList) Add(name string) {
	f.Lock()
	defer f.Unlock()
	f.friendSet[strutil.Base37.Encode(name)] = false
}

//ToggleStatus will flip the boolean value mapped to name then return true, or if no such entry exists, does nothing and returns false.
func (f *FriendsList) ToggleStatus(name string) bool {
	f.Lock()
	defer f.Unlock()
	toggleHash := strutil.Base37.Encode(name)
	for hash, status := range f.friendSet {
		if hash == toggleHash {
			f.friendSet[hash] = !status
			return true
		}
	}
	return false
}

//Set maps an online status (true/false) to a player by their username hash.
func (f *FriendsList) Set(name string, val bool) {
	f.Lock()
	defer f.Unlock()
	f.friendSet[strutil.Base37.Encode(name)] = val
}

func (f *FriendsList) Remove(name string) {
	f.Lock()
	defer f.Unlock()
	hash := strutil.Base37.Encode(name)
	delete(f.friendSet, hash)
	//	if p, ok := Players.FromUserHash(hash); ok && p.FriendList.contains(f.Owner) {
	//		p.SendPacket(FriendUpdate(f.Owner, false))
	//	}
}

//NameSet returns all the players names in this collection.
func (f *FriendsList) NameSet() (nameList []string) {
	f.RLock()
	defer f.RUnlock()
	for name := range f.friendSet {
		nameList = append(nameList, strutil.Base37.Decode(name))
	}
	return
}

//EntrySet returns a copy of this collection.
func (f *FriendsList) EntrySet() (nameList map[string]bool) {
	f.RLock()
	defer f.RUnlock()
	nameList = make(map[string]bool)
	for name, status := range f.friendSet {
		nameList[strutil.Base37.Decode(name)] = status
	}
	return
}

//ForEach runs fn(key,val) for each entry in the collection
func (f *FriendsList) ForEach(fn func(string, bool) bool) {
	f.RLock()
	defer f.RUnlock()
	for name, status := range f.friendSet {
		if fn(strutil.Base37.Decode(name), status) {
			break
		}
	}
}

//Size returns the length of the entry set.
func (f *FriendsList) Size() int {
	f.RLock()
	defer f.RUnlock()
	return len(f.friendSet)
}
