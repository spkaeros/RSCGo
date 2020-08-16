package world

import (
	"sync"
)

type EntityType int

const (
	TypePlayer     EntityType = 1 << iota // 1
	TypeNpc                         // 2
	TypeObject                      // 4
	TypeDoor                        // 8
	TypeItem                        // 16
	TypeGroundItem                  // 32
	
	TypeMob    = TypePlayer | TypeNpc
	TypeEntity = TypeObject | TypeDoor | TypeItem | TypeGroundItem
)


//Entity Anything that is a part of the games physical world, for everyone to see.
type Entity interface {
	//MobileEntity
	Location() *Location
	X() int
	Y() int
	ServerIndex() int
	IsObject() bool
	IsGroundItem() bool
	IsPlayer() bool
	IsNpc() bool
	Type() EntityType
	
//`	type *Location, Location, *Player, Player, *NPC, NPC, *Object, Object, *GroundItem, GroundItem, MobileEntity
}


type entityVec []Entity

//entityList Represents a entityList of scene entities.
type EntityList struct {
	entityVec
	sync.RWMutex
}

func MakeEntityList() *EntityList {
	return &EntityList{}
}

func (l *EntityList) RangeAs(action func (e Entity) bool) int {

	for i,v := range l.entityVec {
		if action(v) {
			return i
		}
	}
	return -1
}

func (l *EntityList) Visible(p *Player, rad int) []Entity {
	l.RLock()
	defer l.RUnlock()
	var entities []Entity
	for _, e := range l.entityVec {
		if e != nil && p.Location().WithinRange(*e.Location(),16) {
			entities = append(entities, e)
		}
	}
	
	return entities
}
//
////NearbyPlayers creates a slice of *Player and populates it with players within p's view area that are in this region,
//// and returns it.
////func (l *entityList) NearbyPlayers(p *Player) []*Player {
////	l.lock.RLock()
////	defer l.lock.RUnlock()
////	var players []*Player
////	for _, v := range l.entityVec {
////		if v, ok := v.(*Player); ok && (v.IsPlayer() && v.ServerIndex() != p.ServerIndex()) && p.Location().LongestDelta(*v.Location()) < p.VarInt("viewRadius", rad) {
////			players = append(players, v)
////		}
////	}
////	return players
////}
//
////NearbyNpcs creates a slice of *NPC and populates it with npcs within p's view area that are in this region,
//// and returns it.
//func (l *entityList) NearbyNpcs(p *Player) []*NPC {
//	l.lock.RLock()
//	defer l.lock.RUnlock()
//	var npcs []*NPC
//	for _, v := range l.entityVec {
//		if v, ok := v.(*NPC); ok && p.LongestDelta(v.Location) < p.VarInt("viewRadius", 16) {
//			npcs = append(npcs, v)
//		}
//	}
//	return npcs
//}
//
////NearbyObjects creates a slice of *Object and populates it with objects within p's view area that are in this region,
//// and returns it.
//func (l *entityList) NearbyObjects(p *Player) []*Object {
//	l.lock.RLock()
//	defer l.lock.RUnlock()
//	var objects []*Object
//	for _, o1 := range l.entityVec {
//		if o1, ok := o1.(*Object); ok && p.LongestDelta(o1.Location) < p.VarInt("viewRadius", 16)+5 {
//			objects = append(objects, o1)
//		}
//	}
//	return objects
//}
//
////NearbyItems creates a slice of *GroundItem and populates it with ground items within p's view area that are in this region,
//// and returns it.
//func (l *entityList) NearbyItems(p *Player) []*GroundItem {
//	l.lock.RLock()
//	defer l.lock.RUnlock()
//	var items []*GroundItem
//	for _, i := range l.entityVec {
//		if i, ok := i.(*GroundItem); ok && i.VisibleTo(p) && p.WithinRange(i.Location, p.VarInt("viewRadius", 16)) {
//			items = append(items, i)
//		}
//	}
//	return items
//}
//
////Add puts e entity into the collection set.
//func (l *entityList) Add(e interface{}) {
//	l.lock.Lock()
//	defer l.lock.Unlock()
//	l.entityVec = append(l.entityVec, e)
//}

//Contains checks the collection set for e and returns true if it finds it.
func (l *EntityList) Contains(e Entity) bool {
	l.RLock()
	defer l.RUnlock()
	elems := l.entityVec
	for _, v := range elems {
		if v == e {
			// Pointers should be comparable?
			return true
		}
	}
	
	return false
}

//Remove removes e from the collection set.
func (l *EntityList) Remove(e Entity) {
	l.Lock()
	defer l.Unlock()
	elems := l.entityVec
	for i, v := range elems {
		if v == e {
			last := len(elems) - 1
			if i < last {
				copy(elems[i:], elems[i+1:])
			}
			elems[last] = nil
			l.entityVec = elems[:last]
			return
		}
	}
}

func (set *entityVec) Put(e Entity) {
	*set = append(*set, e)
}
