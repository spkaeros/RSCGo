package world

//Pathway Represents a path for a mobile entity to traverse across the virtual world.
type Pathway struct {
	StartX, StartY  int
	WaypointsX      []int
	WaypointsY      []int
	CurrentWaypoint int
}

//NewPathwayToCoords returns a new Pathway pointing to the specified location.  Will attempt traversal to l via a
// simple algorithm: if curX < destX then increase, if curX > destX then decrease, same for y, until equal.
// TODO: No clipping is attempted yet, and no path waypoints are generated to avoid obstacles yet.  Gotta do it
func NewPathwayToCoords(destX, destY int) *Pathway {
	return NewPathway(destX, destY, []int{}, []int{})
}

//NewPathwayToLocation returns a new Pathway pointing to the specified location.  Will attempt traversal to l via a
// simple algorithm: if curX < destX then increase, if curX > destX then decrease, same for y, until equal.
// TODO: No clipping is attempted yet, and no path waypoints are generated to avoid obstacles yet.  Gotta do it
func NewPathwayToLocation(l Location) *Pathway {
	return NewPathwayToCoords(l.X(), l.Y())
}

//NewPathway returns a new Pathway with the specified variables.  destX and destY are a straight line, and waypoints define turns from that point.
func NewPathway(destX, destY int, waypointsX, waypointsY []int) *Pathway {
	return &Pathway{StartX: destX, StartY: destY, WaypointsX: waypointsX, WaypointsY: waypointsY, CurrentWaypoint: -1}
}

func MakePath(start, end Location) *Pathway {
	return NewPathfinder(start, end).MakePath()
}

//CountWaypoints Returns the length of the largest waypoint slice within this path.
func (p *Pathway) CountWaypoints() int {
	xCount, yCount := len(p.WaypointsX), len(p.WaypointsY)
	if xCount >= yCount {
		return xCount
	}
	return yCount
}

//WaypointX Returns the x coordinate of the specified waypoint, by taking the waypointX delta at w, and adding it to StartX.
// If w is out of bounds, returns the StartX coordinate, aka the x coord to start turning at.
func (p *Pathway) WaypointX(w int) int {
	offset := func(w int) int {
		if w >= p.CountWaypoints() || w < 0 {
			return 0
		}
		return p.WaypointsX[w]
	}(w)
	return p.StartX + offset
}

//WaypointY Returns the y coordinate of the specified waypoint, by taking the waypointY delta at w, and adding it to StartY.
// If w is out of bounds, returns the StartY coordinate, aka the y coord to start turning at.
func (p *Pathway) WaypointY(w int) int {
	offset := func(w int) int {
		if w >= p.CountWaypoints() || w < 0 {
			return 0
		}
		return p.WaypointsY[w]
	}(w)
	return p.StartY + offset
}

//NextWaypointTile Returns the next destination within our path.  If our current waypoint is out of bounds, it will return
// the same value as StartingTile.
func (p *Pathway) NextWaypointTile() Location {
	return NewLocation(p.WaypointX(p.CurrentWaypoint), p.WaypointY(p.CurrentWaypoint))
}

//StartingTile Returns the location of the start of the path,  This location is actually not our starting location,
// but the first tile that we begin traversing our waypoint deltas from.  Required to walk to this location to start
// traversing waypoints,
func (p *Pathway) StartingTile() Location {
	return NewLocation(p.StartX, p.StartY)
}

//AddWaypoint Prepends a waypoint to this path.
func (p *Pathway) AddWaypoint(x, y int) *Pathway {
	p.WaypointsX = append([]int{x}, p.WaypointsX...)
	p.WaypointsY = append([]int{y}, p.WaypointsY...)
	return p
}

//NextTileToward Returns the next tile toward the final destination of this pathway from currentLocation
func (p *Pathway) NextTileFrom(currentLocation Location) Location {
	dest := p.NextWaypointTile()
	destX, destY := dest.X(), dest.Y()
	currentX, currentY := currentLocation.X(), currentLocation.Y()
	destination := NewLocation(currentX, currentY)
	switch {
	case currentX > destX:
		destination.x.Dec()
	case currentX < destX:
		destination.x.Inc()
	}
	switch {
	case currentY > destY:
		destination.y.Dec()
	case currentY < destY:
		destination.y.Inc()
	}
	return destination
}
