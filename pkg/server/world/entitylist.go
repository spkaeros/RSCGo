package world

import (
	"sync"
)

//Entity A stationary scene entity within the game world.
type Entity struct {
	Location
	Index int
}

//AtLocation Returns true if the entity is at the specified location, otherwise returns false
func (e *Entity) AtLocation(location Location) bool {
	return e.AtCoords(location.X.Load(), location.Y.Load())
}

//AtCoords Returns true if the entity is at the specified coordinates, otherwise returns false
func (e *Entity) AtCoords(x, y uint32) bool {
	return e.X.Load() == x && e.Y.Load() == y
}

//List Represents a list of scene entities.
type List struct {
	List []interface{}
	lock sync.RWMutex
}

//NearbyPlayers Might remove
func (l *List) NearbyPlayers(p *Player) []*Player {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var players []*Player
	for _, v := range l.List {
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(v.Location) <= 15 {
			players = append(players, v)
		}
	}
	return players
}

//NearbyNPCs Might remove
func (l *List) NearbyNPCs(p *Player) []*NPC {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var npcs []*NPC
	for _, v := range l.List {
		if v, ok := v.(*NPC); ok && p.LongestDelta(v.Location) <= 15 {
			npcs = append(npcs, v)
		}
	}
	return npcs
}

//RemovingPlayers Might remove
func (l *List) RemovingPlayers(p *Player) []*Player {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var players []*Player
	for _, v := range l.List {
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(v.Location) > 15 {
			players = append(players, p)
		}
	}
	return players
}

//NearbyObjects Might remove
func (l *List) NearbyObjects(p *Player) []*Object {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var objects []*Object
	for _, o1 := range l.List {
		if o1, ok := o1.(*Object); ok && p.LongestDelta(o1.Location) <= 20 {
			objects = append(objects, o1)
		}
	}
	return objects
}

//NearbyItems Might remove
func (l *List) NearbyItems(p *Player) []*GroundItem {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var items []*GroundItem
	for _, i := range l.List {
		if i, ok := i.(*GroundItem); ok && i.VisibleTo(p) && p.WithinRange(i.Location, 21) {
			items = append(items, i)
		}
	}
	return items
}

//RemovingObjects Might remove
func (l *List) RemovingObjects(p *Player) []*Object {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var objects []*Object
	for _, o1 := range l.List {
		if o1, ok := o1.(*Object); ok && p.LongestDelta(o1.Location) > 20 {
			objects = append(objects, o1)
		}
	}
	return objects
}

//Add Add an entity to the list.
func (l *List) Add(e interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.List = append(l.List, e)
}

//Contains Returns true if e is an element of l, otherwise returns false.
func (l *List) Contains(e interface{}) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for _, v := range l.List {
		if v == e {
			// Pointers should be comparable?
			return true
		}
	}

	return false
}

//Remove Removes Entity e from List l.
func (l *List) Remove(e interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	elems := l.List
	for i, v := range elems {
		if v == e {
			last := len(elems) - 1
			if i < last {
				copy(elems[i:], elems[i+1:])
			}
			elems[last] = nil
			l.List = elems[:last]
			return
		}
	}
}
