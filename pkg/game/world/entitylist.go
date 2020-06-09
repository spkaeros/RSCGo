package world

import (
	"sync"
	"github.com/spkaeros/rscgo/pkg/game/entity"
)

//Entity A stationary scene entity within the game world.
type Entity struct {
	Location
	Index int
}

func (e *Entity) Point() entity.Location {
	return entity.Location(e)
}

func (e *Entity) ServerIndex() int {
	return e.Index
}

//AtLocation Returns true if the entity is at the specified location, otherwise returns false
func (e *Entity) AtLocation(location Location) bool {
	return e.Location.Equals(location)
}

//entityList Represents a entityList of scene entities.
type entityList struct {
	set  []interface{}
	lock sync.RWMutex
}

//NearbyPlayers creates a slice of *Player and populates it with players within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyPlayers(p *Player) []*Player {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var players []*Player
	for _, v := range l.set {
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(v.Location) < p.VarInt("viewRadius", 16) {
			players = append(players, v)
		}
	}
	return players
}

//NearbyNpcs creates a slice of *NPC and populates it with npcs within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyNpcs(p *Player) []*NPC {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var npcs []*NPC
	for _, v := range l.set {
		if v, ok := v.(*NPC); ok && p.LongestDelta(v.Location) < p.VarInt("viewRadius", 16) {
			npcs = append(npcs, v)
		}
	}
	return npcs
}

//NearbyObjects creates a slice of *Object and populates it with objects within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyObjects(p *Player) []*Object {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var objects []*Object
	for _, o1 := range l.set {
		if o1, ok := o1.(*Object); ok && p.LongestDelta(o1.Location) < p.VarInt("viewRadius", 16)+5 {
			objects = append(objects, o1)
		}
	}
	return objects
}

//NearbyItems creates a slice of *GroundItem and populates it with ground items within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyItems(p *Player) []*GroundItem {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var items []*GroundItem
	for _, i := range l.set {
		if i, ok := i.(*GroundItem); ok && i.VisibleTo(p) && p.WithinRange(i.Location, p.VarInt("viewRadius", 16)) {
			items = append(items, i)
		}
	}
	return items
}

//Add puts e entity into the collection set.
func (l *entityList) Add(e interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.set = append(l.set, e)
}

//Contains checks the collection set for e and returns true if it finds it.
func (l *entityList) Contains(e interface{}) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for _, v := range l.set {
		if v == e {
			// Pointers should be comparable?
			return true
		}
	}

	return false
}

//Remove removes e from the collection set.
func (l *entityList) Remove(e interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	elems := l.set
	for i, v := range elems {
		if v == e {
			last := len(elems) - 1
			if i < last {
				copy(elems[i:], elems[i+1:])
			}
			elems[last] = nil
			l.set = elems[:last]
			return
		}
	}
}
