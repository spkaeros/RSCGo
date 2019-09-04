package entity

//Pathway Represents a path for a mobile entity to traverse across the virtual world.
type Pathway struct {
	StartX, StartY  int
	WaypointsX      []int
	WaypointsY      []int
	CurrentWaypoint int
}

//NewPathway returns a new Pathway pointing to the specified coordinates.  Must be a straight line from starting tile.
func NewPathway(destX, destY int) *Pathway {
	return &Pathway{StartX: destX, StartY: destY, CurrentWaypoint: -1}
}

//NewPathwayFromLocation returns a new Pathway pointing to the specified location.  Must be a straight line from starting location.
func NewPathwayFromLocation(l Location) *Pathway {
	return NewPathway(l.X, l.Y)
}

//NewPathwayComplete returns a new Pathway with the specified variables.  destX and destY are a straight line, and waypoints define turns from that point.
func NewPathwayComplete(destX, destY int, waypointsX, waypointsY []int) *Pathway {
	return &Pathway{destX, destY, waypointsX, waypointsY, -1}
}

//waypointXoffset Returns the offset for the X coordinate of the specified waypoint.
func (p *Pathway) waypointXoffset(w int) int {
	if w >= len(p.WaypointsX) || w == -1 {
		return 0
	}
	return p.WaypointsX[w]
}

//waypointX Returns the X coordinate of the specified waypoint.
func (p *Pathway) waypointX(w int) int {
	return p.StartX + p.waypointXoffset(w)
}

//waypointYoffset Returns the offset for the Y coordinate of the specified waypoint.
func (p *Pathway) waypointYoffset(w int) int {
	if w >= len(p.WaypointsY) || w == -1 {
		return 0
	}
	return p.WaypointsY[w]
}

//waypointY Returns the Y coordinate of the specified waypoint.
func (p *Pathway) waypointY(w int) int {
	return p.StartY + p.waypointYoffset(w)
}

//Waypoint Returns the locattion of the specified waypoint
func (p *Pathway) Waypoint(w int) *Location {
	return &Location{p.waypointX(w), p.waypointY(w)}
}

//Start Returns the location of the start of the path
func (p *Pathway) Start() *Location {
	return &Location{p.StartX, p.StartY}
}

//NextTile Returns the next tile for the mob to move to in the pathway.
func (p *Pathway) NextTile(startX, startY int) Location {
	destX := p.waypointX(p.CurrentWaypoint)
	destY := p.waypointY(p.CurrentWaypoint)
	newLocation := &Location{destX, destY}
	switch {
	case startX > destX:
		newLocation.X = startX - 1
		break
	case startX < destX:
		newLocation.X = startX + 1
		break
	}
	switch {
	case startY > destY:
		newLocation.Y = startY - 1
		break
	case startY < destY:
		newLocation.Y = startY + 1
		break
	}
	return *newLocation
}
