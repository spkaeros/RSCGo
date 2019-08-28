package entity

//Locatable An interface for all locatable scene entities to implement
type Locatable interface {
	Location() *Location
}

//LocatableList Represents a list of locatable scene entities.
type LocatableList struct {
	List []Locatable
}

//AddPlayer Add a player to the region.
func (l *LocatableList) AddPlayer(p *Player) {
	l.List = append(l.List, p)
}

//ContainsPlayer Returns true if the receiver list contains the player specified, false otherwise.
func (l *LocatableList) ContainsPlayer(p *Player) bool {
	for _, v := range l.List {
		if v, ok := v.(*Player); ok {
			if v.Index == p.Index {
				return true
			}
		}
	}
	return false
}

//RemovePlayer Remove a player from the region.
func (l *LocatableList) RemovePlayer(p *Player) {
	players := l.List
	for i, v := range players {
		v, ok := v.(*Player)
		if ok {
			if v.Index == p.Index {
				last := len(players) - 1
				players[i] = players[last]
				l.List = players[:last]
				return
			}
		}
	}
}

//AddObject Add an object to the list.
func (l *LocatableList) AddObject(p *Object) {
	l.List = append(l.List, p)
}

//RemoveObject Remove an object from the list.
func (l *LocatableList) RemoveObject(p *Object) {
	objects := l.List
	for i, v := range objects {
		v, ok := v.(*Object)
		if ok {
			if v.Index == p.Index {
				last := len(objects) - 1
				objects[i] = objects[last]
				l.List = objects[:last]
				return
			}
		}
	}
}

//ContainsObject Returns true if the receiver list contains the player specified, false otherwise.
func (l *LocatableList) ContainsObject(o *Object) bool {
	for _, v := range l.List {
		if v, ok := v.(*Object); ok {
			if v.Location().LongestDelta(o.Location()) == 0 {
				return true
			}
		}
	}
	return false
}
