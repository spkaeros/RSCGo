package world

import (
	"sync"
	"github.com/spkaeros/rscgo/pkg/game/entity"
)

//Entity A stationary scene entity within the game world.
type Entity struct {
	entity.Location
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
	return e.Near(location, 0)
}

//entityList Represents a entityList of scene entities.
type entityList struct {
	set  []entity.Entity
	sync.RWMutex
}

//NearbyPlayers creates a slice of *Player and populates it with players within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyPlayers(e entity.MobileEntity) []entity.MobileEntity {
	l.RLock()
	defer l.RUnlock()
	var players []entity.MobileEntity
	for _, v := range l.set {
		if p, ok := v.(*Player); ok && p.ServerIndex() != e.ServerIndex() && p.Near(e, e.SessionCache().VarInt("viewRadius", 16)) {
			players = append(players, p)
		}
	}
	return players
}

//NearbyNpcs creates a slice of *NPC and populates it with npcs within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyNpcs(e entity.MobileEntity) []entity.MobileEntity {
	l.RLock()
	defer l.RUnlock()
	var npcs []entity.MobileEntity
	for _, v := range l.set {
		if n, ok := v.(*NPC); ok && n.Near(e, e.SessionCache().VarInt("viewRadius", 16)) {
			npcs = append(npcs, n)
		}
	}
	
	return npcs
}

//NearbyObjects creates a slice of *Object and populates it with objects within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyObjects(e entity.MobileEntity) []entity.Entity {
	l.RLock()
	defer l.RUnlock()
	var objects []entity.Entity
	for _, o1 := range l.set {
		if o, ok := o1.(*Object); ok && e.Near(o, e.SessionCache().VarInt("viewRadius", 16)>>1) {
			objects = append(objects, o)
		}
	}
	return objects
}

//NearbyObjects creates a slice of *Object and populates it with objects within p's view area that are in this region,
// and returns it.
func (l *entityList) Range(fn func(entity.Entity)) {
	l.RLock()
	defer l.RUnlock()
	for _, o1 := range l.set {
		fn(o1)
	}
}

//NearbyItems creates a slice of *GroundItem and populates it with ground items within p's view area that are in this region,
// and returns it.
func (l *entityList) NearbyItems(e entity.MobileEntity) []entity.Entity {
	l.RLock()
	defer l.RUnlock()
	var items []entity.Entity
	for _, i := range l.set {
		if i, ok := i.(*GroundItem); ok && i.Near(e, e.SessionCache().VarInt("viewRadius", 16)>>1) {
			if p := AsPlayer(e); p != nil && !i.VisibleTo(p) {
				continue
			}
			items = append(items, i)
		}
	}
	
	return items
}

//Add puts e entity into the collection set.
func (l *entityList) Add(e entity.Entity) {
	l.Lock()
	defer l.Unlock()
	l.set = append(l.set, e)
}

//Contains checks the collection set for e and returns true if it finds it.
func (l *entityList) Contains(e entity.Entity) bool {
	l.RLock()
	defer l.RUnlock()
	for _, v := range l.set {
		if (v.ServerIndex() == e.ServerIndex() && v.Type() == e.Type()) || e == v {
			// Pointers should be comparable?
			return true
		}
	}

	return false
}

//Remove removes e from the collection set.
func (l *entityList) Remove(e entity.Entity) {
	l.Lock()
	defer l.Unlock()
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
