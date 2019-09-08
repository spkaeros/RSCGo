package entity

type Entity struct {
	Location
	Index int
}

//EntityList Represents a list of Entity scene entities.
type EntityList struct {
	List []interface{}
}

func (l *EntityList) AddEntity(e *Entity) {
	l.List = append(l.List, e)
}

func (l *EntityList) RemoveEntity(e Entity) {
	entitys := l.List
	for i, v := range l.List {
		if v, ok := v.(Entity); ok && v.Index == e.Index {
			last := len(entitys) - 1
			entitys[i] = entitys[last]
			l.List = entitys[:last]
			return
		}
	}
}

//AddPlayer Add a player to the region.
func (l *EntityList) AddPlayer(p *Player) {
	l.List = append(l.List, p)
}

func (l *EntityList) NearbyPlayers(p *Player) []*Player {
	var players []*Player
	for _, v := range l.List {
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(v.Location) <= 15 {
			players = append(players, v)
		}
	}
	return players
}

func (l *EntityList) RemovingPlayers(p *Player) []*Player {
	var players []*Player
	for _, v := range l.List {
		if v, ok := v.(*Player); ok && v.Index != p.Index && p.LongestDelta(v.Location) > 15 {
			players = append(players, p)
		}
	}
	return players
}

func (l *EntityList) NearbyObjects(p *Player) []*Object {
	var objects []*Object
	for _, o1 := range l.List {
		if o1, ok := o1.(*Object); ok && o1.Index() != p.Index && p.LongestDelta(*o1.location) <= 20 {
			objects = append(objects, o1)
		}
	}
	return objects
}

func (l *EntityList) RemovingObjects(p *Player) []*Object {
	var objects []*Object
	for _, o1 := range l.List {
		if o1, ok := o1.(*Object); ok && p.LongestDelta(*o1.location) > 20 {
			objects = append(objects, o1)
		}
	}
	return objects
}

//ContainsPlayer Returns true if the receiver list contains the player specified, false otherwise.
func (l *EntityList) ContainsPlayer(p *Player) bool {
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
func (l *EntityList) RemovePlayer(p *Player) {
	players := l.List
	for i, v := range players {
		if v, ok := v.(*Player); ok {
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
func (l *EntityList) AddObject(p *Object) {
	l.List = append(l.List, p)
}

//RemoveObject Remove an object from the list.
func (l *EntityList) RemoveObject(p *Object) {
	objects := l.List
	for i, v := range objects {
		if v, ok := v.(*Object); ok {
			if v.Location().LongestDelta(*p.location) == 0 {
				last := len(objects) - 1
				objects[i] = objects[last]
				l.List = objects[:last]
				return
			}
		}
	}
}

//GetObject Looks for an object at given coordinates.  If found, returns a reference to it.  Otherwise, returns nil
func (l *EntityList) GetObject(x, y int) *Object {
	for _, v := range l.List {
		if v, ok := v.(*Object); ok {
			if v.Location().X == x && v.Location().Y == y {
				return v
			}
		}
	}
	return nil
}

//ContainsObject Returns true if the receiver list contains the player specified, false otherwise.
func (l *EntityList) ContainsObject(o *Object) bool {
	for _, v := range l.List {
		if v, ok := v.(*Object); ok {
			if v.Location().LongestDelta(*o.location) == 0 && v.ID == o.ID && v.Direction == o.Direction {
				return true
			}
		}
	}
	return false
}
