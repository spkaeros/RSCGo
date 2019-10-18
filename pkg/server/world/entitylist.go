package world

import (
	"log"
	"os"
	"sync"
)

//Entity A stationary scene entity within the game world.
type Entity struct {
	Location
	Index int
}

//AtLocation Returns true if the entity is at the specified location, otherwise returns false
func (e *Entity) AtLocation(location *Location) bool {
	return e.AtCoords(location.X, location.Y)
}

//AtCoords Returns true if the entity is at the specified coordinates, otherwise returns false
func (e *Entity) AtCoords(x, y int) bool {
	return e.X == x && e.Y == y
}

//LogWarning Log interface for warnings.
var LogWarning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)

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
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(&v.Location) <= 15 {
			players = append(players, v)
		}
	}
	return players
}

//RemovingPlayers Might remove
func (l *List) RemovingPlayers(p *Player) []*Player {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var players []*Player
	for _, v := range l.List {
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(&v.Location) > 15 {
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
		if o1, ok := o1.(*Object); ok && p.LongestDelta(&o1.Location) <= 20 {
			objects = append(objects, o1)
		}
	}
	return objects
}

//RemovingObjects Might remove
func (l *List) RemovingObjects(p *Player) []*Object {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var objects []*Object
	for _, o1 := range l.List {
		if o1, ok := o1.(*Object); ok && p.LongestDelta(&o1.Location) > 20 {
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
			elems[i] = elems[last]
			l.List = elems[:last]
			return
		}
	}
}
