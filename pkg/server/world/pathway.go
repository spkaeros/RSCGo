package world

import (
	"math"
)

//Pathway Represents a path for a mobile entity to traverse across the virtual world.
type Pathway struct {
	StartX, StartY  uint32
	WaypointsX      []int
	WaypointsY      []int
	CurrentWaypoint int
}

//NewPathwayToCoords returns a new Pathway pointing to the specified location.  Will attempt traversal to l via a
// simple algorithm: if curX < destX then increase, if curX > destX then decrease, same for Y, until equal.
// TODO: No clipping is attempted yet, and no path waypoints are generated to avoid obstacles yet.  Gotta do it
func NewPathwayToCoords(destX, destY uint32) *Pathway {
	return NewPathway(destX, destY, []int{}, []int{})
}

//NewPathwayToLocation returns a new Pathway pointing to the specified location.  Will attempt traversal to l via a
// simple algorithm: if curX < destX then increase, if curX > destX then decrease, same for Y, until equal.
// TODO: No clipping is attempted yet, and no path waypoints are generated to avoid obstacles yet.  Gotta do it
func NewPathwayToLocation(l Location) *Pathway {
	return NewPathwayToCoords(l.X.Load(), l.Y.Load())
}

//NewPathway returns a new Pathway with the specified variables.  destX and destY are a straight line, and waypoints define turns from that point.
func NewPathway(destX, destY uint32, waypointsX, waypointsY []int) *Pathway {
	return &Pathway{StartX: destX, StartY: destY, WaypointsX: waypointsX, WaypointsY: waypointsY, CurrentWaypoint: -1}
}

//CountWaypoints Returns the length of the largest waypoint slice within this path.
func (p *Pathway) CountWaypoints() int {
	xCount, yCount := len(p.WaypointsX), len(p.WaypointsY)
	if xCount >= yCount {
		return xCount
	}
	return yCount
}

//WaypointX Returns the X coordinate of the specified waypoint, by taking the waypointX delta at w, and adding it to StartX.
// If w is out of bounds, returns the StartX coordinate, aka the X coord to start turning at.
func (p *Pathway) WaypointX(w int) uint32 {
	offset := func(w int) int {
		if w >= p.CountWaypoints() || w < 0 {
			return 0
		}
		return p.WaypointsX[w]
	}(w)
	return p.StartX + uint32(offset)
}

//WaypointY Returns the Y coordinate of the specified waypoint, by taking the waypointY delta at w, and adding it to StartY.
// If w is out of bounds, returns the StartY coordinate, aka the Y coord to start turning at.
func (p *Pathway) WaypointY(w int) uint32 {
	offset := func(w int) int {
		if w >= p.CountWaypoints() || w < 0 {
			return 0
		}
		return p.WaypointsY[w]
	}(w)
	return p.StartY + uint32(offset)
}

//NextWaypointTile Returns the next destination within our path.  If our current waypoint is out of bounds, it will return
// the same value as StartingTile.
func (p *Pathway) NextWaypointTile() Location {
	return NewLocation(int(p.WaypointX(p.CurrentWaypoint)), int(p.WaypointY(p.CurrentWaypoint)))
}

//StartingTile Returns the location of the start of the path,  This location is actually not our starting location,
// but the first tile that we begin traversing our waypoint deltas from.  Required to walk to this location to start
// traversing waypoints,
func (p *Pathway) StartingTile() Location {
	return NewLocation(int(p.StartX), int(p.StartY))
}

//AddWaypoint Prepends a waypoint to this path.
func (p *Pathway) AddWaypoint(x, y int) *Pathway{
	p.WaypointsX = append([]int{x}, p.WaypointsX...)
	p.WaypointsY = append([]int{y}, p.WaypointsY...)
	return p
}

func MakePath(start, end Location) *Pathway {
	nodes := make(map[int]*Node)
	startNode := &Node{
		cost:   0,
		open:   true,
		parent: nil,
		loc:    start,
	}
	endNode := &Node{
		cost:   0,
		open:   true,
		parent: nil,
		loc:    end,
	}
	if IsTileBlocking(end.CurX(), end.CurY(), 0x40, false) {
		return NewPathwayToLocation(end)
	}
	nodes[start.CurX() << 32 | start.CurY() << 16] = startNode
	nodes[end.CurX() << 32 | end.CurY() << 16] = endNode
	open := []*Node{startNode}
	for len(open) > 0 {
		active := getCheapestNode(&open)
		position := active.loc
		if position.LongestDelta(end) == 0 {
			break
		}
		for i, n := range open {
			if n == active {
				open = append(open[:i], open[i+1:]...)
				break
			}
		}

		active.open = false
		x, y := position.CurX(), position.CurY()
		for nextX := x - 1; nextX <= x + 1; nextX++ {
			for nextY := y - 1; nextY <= y + 1; nextY++ {
				if nextX == x && nextY == y {
					continue
				}

				adj := NewLocation(nextX, nextY)
				sprites := [3][3]int{{SouthWest, West, NorthWest}, {South, -1, North}, {SouthEast, East, NorthEast}}
				xIndex, yIndex := position.CurX()-adj.CurX()+1, position.CurY()-adj.CurY()+1
				bit := 4
				bit2 := 1
				if xIndex < 0 || xIndex >= 3 {
					continue
				}
				if yIndex < 0 || yIndex >= 3 {
					continue
				}
				dir := sprites[xIndex][yIndex]
				switch dir {
				case North:
					bit = 4
					bit2 = 1
				case South:
					bit = 1
					bit2 = 4
				case East:
					bit = 8
					bit2 = 2
				case West:
					bit = 2
					bit2 = 8
				case NorthEast:
					bit = 4 | 8
					bit2 = 1 | 8
				case NorthWest:
					bit = 4 | 2
					bit2 = 1 | 2
				case SouthEast:
					bit = 1 | 8
					bit2 = 4 | 8
				case SouthWest:
					bit = 1 | 2
					bit2 = 4 | 2
				}
				if !IsTileBlocking(position.CurX(), position.CurY(), byte(bit2), true) && !IsTileBlocking(adj.CurX(), adj.CurY(), byte(bit), false) {
					switch dir {
					case NorthEast:
						if IsTileBlocking(position.CurX(), position.CurY()-1, byte(bit), false) {
							continue
						}
						if IsTileBlocking(position.CurX()-1, position.CurY(), byte(bit), false) {
							continue
						}
					case NorthWest:
						if IsTileBlocking(position.CurX(), position.CurY()-1, byte(bit), false) {
							continue
						}
						if IsTileBlocking(position.CurX()+1, position.CurY(), byte(bit), false) {
							continue
						}
					case SouthEast:
						if IsTileBlocking(position.CurX(), position.CurY()+1, byte(bit), false) {
							continue
						}
						if IsTileBlocking(position.CurX()-1, position.CurY(), byte(bit), false) {
							continue
						}
					case SouthWest:
						if IsTileBlocking(position.CurX(), position.CurY()+1, byte(bit), false) {
							continue
						}
						if IsTileBlocking(position.CurX()+1, position.CurY(), byte(bit), false) {
							continue
						}
					}
					node, ok := nodes[adj.CurX() << 32 | adj.CurY() << 16]//&Node{loc: adj, open: true}
					if !ok {
						node = &Node{loc:adj, open:true}
						nodes[adj.CurX() << 32 | adj.CurY() << 16] = node
					}
					compareNodes(active, node, &open, end)
				}
			}
		}
	}

	path := &Pathway{StartX: 0, StartY: 0}

	active := endNode
	if active.parent != nil {
		position := active.loc
		for start.LongestDelta(position) > 0 {
			path.AddWaypoint(position.CurX(), position.CurY())
			active = active.parent
			position = active.loc
		}
	}
	return path
}

func cost(start, end Location) int {
	deltaX, deltaY := start.DeltaX(end), start.DeltaY(end)
	shortL, longL := deltaX, deltaY
	if deltaX > deltaY {
		shortL = deltaY
		longL = deltaX
	}
	return shortL * 14 + ((longL - shortL) * 10)
}

func compareNodes(active, other *Node, open *[]*Node, end Location) {
	gCost := active.gCost + cost(active.loc, other.loc)
	cost := cost(other.loc, end)
	fCost := gCost + cost
	if other.cost > fCost {
		for i, n := range *open {
			if n == other {
				*open = append((*open)[:i], (*open)[i+1:]...)
				break
			}
		}
		other.open = false
	} else if other.open && !inOpen(other, open) {
		other.gCost = gCost
		other.hCost = cost
		other.cost = fCost
		other.parent = active
		*open = append(*open, other)
	}
}

func inOpen(node *Node, open *[]*Node) bool {
	for _, n := range *open {
		if node == n {
			return true
		}
	}
	return false
}

func getCheapestNode(open *[]*Node) *Node {
	var node *Node
	min := math.MaxInt32
	for _, n := range *open {
		if !n.open {
			continue
		}
		if n.cost < min {
			min = n.cost
			node = n
		}
	}
	return node
}

type Node struct {
	hCost, gCost int
	cost         int
	open         bool
	parent       *Node
	loc          Location
}

//NextTileToward Returns the next tile toward the final destination of this pathway from currentLocation
func (p *Pathway) NextTileFrom(currentLocation Location) Location {
	dest := p.NextWaypointTile()
	destX, destY := dest.X.Load(), dest.Y.Load()
	currentX, currentY := currentLocation.X.Load(), currentLocation.Y.Load()
	destination := NewLocation(int(currentX), int(currentY))
	switch {
	case currentX > destX:
		destination.decX()
	case currentX < destX:
		destination.incX()
	}
	switch {
	case currentY > destY:
		destination.decY()
	case currentY < destY:
		destination.incY()
	}
	return destination
}
