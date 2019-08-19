package entity

//Pathway Represents a path for a mobile entity to traverse across the virtual world.
type Pathway struct {
	StartX, StartY  int
	WaypointsX      []int
	WaypointsY      []int
	CurrentWaypoint int
}

//WaypointXoffset Returns the offset for the X coordinate of the specified waypoint.
func (p *Pathway) WaypointXoffset(w int) int {
	if w >= len(p.WaypointsX) || w == -1 {
		return 0
	}
	return p.WaypointsX[w]
}

//WaypointX Returns the X coordinate of the specified waypoint.
func (p *Pathway) WaypointX(w int) int {
	return p.StartX + p.WaypointXoffset(w)
}

//WaypointYoffset Returns the offset for the Y coordinate of the specified waypoint.
func (p *Pathway) WaypointYoffset(w int) int {
	if w >= len(p.WaypointsY) || w == -1 {
		return 0
	}
	return p.WaypointsY[w]
}

//WaypointY Returns the Y coordinate of the specified waypoint.
func (p *Pathway) WaypointY(w int) int {
	return p.StartY + p.WaypointYoffset(w)
}

//Waypoint Returns the locattion of the specified waypoint
func (p *Pathway) Waypoint(w int) *Location {
	return &Location{p.WaypointX(w), p.WaypointY(w)}
}

//Start Returns the locattion of the start of the path
func (p *Pathway) Start() *Location {
	return &Location{p.StartX, p.StartY}
}

//NextTile Returns the next tile for the mob to move to in the pathway.
func (p *Pathway) NextTile(startX, startY int) *Location {
	destX := p.WaypointX(p.CurrentWaypoint)
	destY := p.WaypointY(p.CurrentWaypoint)
	newLocation := &Location{destX, destY}
	switch {
	case startX > destX:
		newLocation.x = startX - 1
		break
	case startX < destX:
		newLocation.x = startX + 1
		break
	}
	switch {
	case startY > destY:
		newLocation.y = startY - 1
		break
	case startY < destY:
		newLocation.y = startY + 1
		break
	}
	return newLocation
}
