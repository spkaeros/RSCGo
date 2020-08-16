package world

import (
	"fmt"
	"math"
	
	"go.uber.org/atomic"
	
	"github.com/spkaeros/rscgo/pkg/rand"
)

type Direction = int

const (
	North Direction = iota
	NorthWest
	West
	SouthWest
	South
	SouthEast
	East
	NorthEast
	// TODO: Check is right
	LeftFighting
	RightFighting
)

//OrderedDirections This is an array containing all of the directions a mob can walk in, ordered by path finder precedent.
// West, East, North, South, SouthWest, SouthEast, NorthWest, NorthEast
var OrderedDirections = [...]Direction{2, 6, 0, 4, 3, 5, 1, 7}

type Plane = int

const (
	//PlaneGround Represents the value for the ground-level plane
	PlaneGround Plane = iota
	//PlaneSecond Represents the value for the second-story plane
	PlaneSecond
	//PlaneThird Represents the value for the third-story plane
	PlaneThird
	//PlaneBasement Represents the value for the basement plane
	PlaneBasement
)

//tileNode self-referential linked node for use within grid-based path finder algorithms
type tileNode struct {
	//parent the tile node that this node comes from in the path
	parent              *tileNode
	//loc the location that this node points at within the game
	loc                 Location
	//hCost,gCost,nCost the various cost values, containing the heuristically calculated travel costs for this node
	hCost, gCost, nCost float64
	//index the priority queue slot of this node
	index               int
	//open represents whether this node has been opened or not
	open,closed       bool
}

//gCostFrom calculates the travel cost of traversing to this node from the specified neighbor node
func (n *tileNode) gCostFrom(neighbor *tileNode) float64 {
	stepPrice := 1.0
	if n.loc.DeltaX(neighbor.loc)+n.loc.DeltaY(neighbor.loc) > 1 {
		stepPrice = math.Sqrt2
	}
	return n.gCost + stepPrice
}

//tileQueue represents a priority queue of tiles used for path finding
// designed to prioritize the tiles with the least cost
type tileQueue []*tileNode

//Len returns the length of this priority queue
func (q tileQueue) Len() int {
	return len(q)
}

//Less returns true if the node at index `i`'s total cost is less than the node at index `j`'s total cost.
func (q tileQueue) Less(i, j int) bool {
	return q[i].nCost < q[j].nCost
}

//Swap swaps the nodes at i and j with each other
func (q tileQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

//Push pushes the pointer x onto the priority queue, and stores its index
func (q *tileQueue) Push(x interface{}) {
	n := len(*q)
	node := x.(*tileNode)
	node.index = n
	*q = append(*q, node)
}

//Pop removes the top-most prioritized node from the queue, then returns it.
func (q *tileQueue) Pop() interface{} {
	old := *q
	n := len(old)
	if n == 0 {
		return nil
	}
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*q = old[0 : n-1]
	return node
}

//Location A tile in the game world.
type Location struct {
	x *atomic.Uint32
	y *atomic.Uint32
}

//Clone returns a new Location that points to the same coordinates as the receiver.
func (l *Location) Clone() *Location {
	return NewLocation(l.X(), l.Y())
}

//X atomically gets the X coordinate for this Location
func (l Location) X() int {
	if l.x == nil {
		return -1
	}
	
	return int(l.x.Load())
}

//Y atomically gets the Y coordinate for this Location
func (l Location) Y() int {
	if l.y == nil {
		return -1
	}
	
	return int(l.y.Load())
}

//SetX atomically stores a new X coordinate into this Location
func (l Location) SetX(x int) {
	l.x.Store(uint32(x))
}

//SetY atomically stores a new Y coordinate into this Location
func (l Location) SetY(y int) {
	l.y.Store(uint32(y))
}

func (l Location) checkWildH() bool {
	return l.X() < 336 && l.X() > 47
}

func (l Location) wildernessDepth() float64 {
	return float64(427-l.Y()) / 6.0
}
//Wilderness calculates and returns the wilderness level of this Location.
func (l Location) Wilderness() int {
	// check X boundaries
	if !l.checkWildH() {
		return 0
	}
	// get precise wild lvl, ensure it's not 0 or lower
	depth := l.wildernessDepth()
	if depth <= 0 {
		return 0
	}
	// truncate wild lvl to int precision; add 1
	return int(depth)+1

	// 2203-1776=427 is the length of our wilderness zone, the first level starts at y=426, every level is 6 tiles
	// return (2203-(l.Y()+1776))/6 + 1
	// return ((426 - l.Y()) / 6) + 1
}

var (
	//DeathSpot The spot where NPCs go to be dead.
	DeathPoint = NewLocation(0, 0)
	//SpawnPoint The default spawn point, where new players start and dead players respawn.
	SpawnPoint = Lumbridge.Clone()
	//Lumbridge Lumbridge teleport point
	Lumbridge = NewLocation(122, 647)
	//Varrock Varrock teleport point
	Varrock = NewLocation(122, 647)
	//Edgeville Edgeville teleport point
	Edgeville = NewLocation(220, 445)
)

//NewLocation Returns a reference to a new instance of the Location data structure.
func NewLocation(x, y int) *Location {
	return &Location{x: atomic.NewUint32(uint32(x)), y: atomic.NewUint32(uint32(y))}
}

func (l Location) Point() Location {
	return *l.Clone()
}

func (l Location) DirectionTo(destX, destY int) Direction {
	sprites := [3][3]Direction{
		{SouthWest, West, NorthWest},
		{South, -1, North},
		{SouthEast, East, NorthEast},
	}
	xIndex, yIndex := l.X()-destX+1, l.Y()-destY+1
	if xIndex >= 3 || yIndex >= 3 || yIndex < 0 || xIndex < 0 {
		xIndex, yIndex = 1, 2 // North
	}
	return sprites[xIndex][yIndex]
}

//NewRandomLocation Returns a new random location within the specified bounds.  bounds[0] should be lowest corner, and
// bounds[1] should be the highest corner.
func NewRandomLocation(bounds [2]Location) Location {
	return *NewLocation(rand.Rng.Intn(bounds[1].X()-bounds[0].X())+bounds[0].X(), rand.Rng.Intn(bounds[1].Y()-bounds[0].Y())+bounds[0].Y())
}

//String Returns a string representation of the location
func (l *Location) String() string {
	return fmt.Sprintf("[%d,%d]", l.X(), l.Y())
}

func (l Location) Within(minX, maxX, minY, maxY int) bool {
	return l.WithinArea([2]Location { *NewLocation(minX, minY), *NewLocation(maxX, maxY) })
}

//IsValid Returns true if the tile at x,y is within world boundaries, false otherwise.
func (l Location) IsValid() bool {
	return l.WithinArea([2]Location { *NewLocation(0, 0), *NewLocation(MaxX, MaxY)})
}

func (l *Location) NextStep(d Location) Location {
	next := l.Step(l.DirectionToward(d))
	if !l.Reachable(next) {
//	if !l.Reachable(d) {
			if l.X() < d.X() {
				if next = l.Step(West); l.Reachable(next) {
					return next
				}
			}
			if l.X() > d.X() {
				if next = l.Step(East); l.Reachable(next) {
					return next
				}
			}
			if l.Y() < d.Y() {
				if next = l.Step(South); l.Reachable(next) {
					return next
				}
				next = l.Step(South)
			}
			if l.Y() > d.Y() {
				if next = l.Step(North); l.Reachable(next) {
					return next
				}
				next = l.Step(North)
			}
	}
	return next
}

func (l *Location) Step(dir Direction) Location {
	loc := l.Clone()
	if dir == 0 || dir == 1 || dir == 7 {
		loc.y.Dec()
	} else if dir == 4 || dir == 5 || dir == 3 {
		loc.y.Inc()
	}
	if dir == 1 || dir == 2 || dir == 3 {
		loc.x.Inc()
	} else if dir == 5 || dir == 6 || dir == 7 {
		loc.x.Dec()
	}
	return *loc
}

//Equals Returns true if this location points to the same location as o
func (l *Location) Equals(o interface{}) bool {
	if e, ok := o.(Entity); ok && e != nil {
		return e.Location() == l
	}
	return false
}

func (l *Location) Delta(other Location) (delta int) {
	return l.LongestDelta(other)
}

//DeltaX Returns the difference between this locations x coord and the other locations x coord
func (l *Location) DeltaX(other Location) (deltaX int) {
	deltaX = int(math.Abs(float64(other.X()) - float64(l.X())))
	// if ourX > theirX {
		// deltaX = ourX - theirX
	// } else if theirX > ourX {
		// deltaX = theirX - ourX
	// }
	return
}

//DeltaY Returns the difference between this locations y coord and the other locations y coord
func (l *Location) DeltaY(other Location) (deltaY int) {
	deltaY = int(math.Abs(float64(other.Y()) - float64(l.Y())))
	// if ourY > theirY {
		// deltaY = ourY - theirY
	// } else if theirY > ourY {
		// deltaY = theirY - ourY
	// }
	return
}

//LongestDelta Returns the largest difference in coordinates between receiver and other
func (l *Location) LongestDelta(other Location) int {
	if x, y := l.DeltaX(other), l.DeltaY(other); x > y {
		return x
	} else {
		return y
	}
}

//LongestDeltaCoords returns the number of tiles the coordinates provided
func (l *Location) LongestDeltaCoords(x, y int) int {
	return l.LongestDelta(*NewLocation(x, y))
}

func (l Location) EuclideanDistance(other Location) float64 {
	return math.Sqrt(math.Pow(float64(l.DeltaX(other)), 2) + math.Pow(float64(l.DeltaY(other)), 2))
}

//WithinRange Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) WithinRange(other Location, radius int) bool {
	return l.Near(other, radius)
}

//WithinRange Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) WithinRadius(other Entity, radius int) bool {
	return l.Near(*other.Location(), radius)
}

//EntityWithin Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) Near(other Location, radius int) bool {
	return l.LongestDeltaCoords(other.X(), other.Y()) <= radius
}

//Plane Calculates and returns the plane that this location is on.
func (l *Location) Plane() int {
	return int(l.y.Load()) / 944 // / 1000
}

//Above Returns the location directly above this one, if any.  Otherwise, if we are on the top floor, returns itself.
func (l *Location) Above() Location {
	return *NewLocation(l.X(), l.PlaneY(true))
}

//Below Returns the location directly below this one, if any.  Otherwise, if we are on the bottom floor, returns itself.
func (l *Location) Below() Location {
	return *NewLocation(l.X(), l.PlaneY(false))
}

func (l *Location) DirectionToward(end Location) Direction {
	tile := l.NextTileToward(end)
	return l.DirectionTo(tile.X(), tile.Y())
}

//PlaneY Updates the location's y coordinate, going up by one plane if up is true, else going down by one plane.  Valid planes: ground=0, 2nd story=1, 3rd story=2, basement=3
func (l *Location) PlaneY(up bool) int {
	curPlane := l.Plane()
	var newPlane int
	if up {
		switch curPlane {
		case PlaneBasement:
			newPlane = 0
		case PlaneThird:
			newPlane = curPlane
		default:
			newPlane = curPlane + 1
		}
	} else {
		switch curPlane {
		case PlaneGround:
			newPlane = PlaneBasement
		case PlaneBasement:
			newPlane = curPlane
		default:
			newPlane = curPlane - 1
		}
	}
	return newPlane*944 + l.Y()%944
}

//NextTileToward Returns the next tile toward the final destination of this pathway from currentLocation
func (l Location) NextTileToward(dst Location) Location {
	nextStep := l.Clone()
	if delta := l.X() - dst.X(); delta < 0 {
		nextStep.x.Inc()
	} else if delta > 0 {
		nextStep.x.Dec()
	}

	if delta := l.Y() - dst.Y(); delta < 0 {
		nextStep.y.Inc()
	} else if delta > 0 {
		nextStep.y.Dec()
	}
	return *nextStep
}

func (l *Location) CanReach(bounds [2]Location) bool {
	x, y := l.X(), l.Y()

	if x >= bounds[0].X() && x <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y() {
		return true
	}
	if x-1 >= bounds[0].X() && x-1 <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y() &&
		(CollisionData(x-1, y)&ClipWest) == 0 {
		return true
	}
	if x+1 >= bounds[0].X() && x+1 <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y() &&
		(CollisionData(x+1, y)&ClipEast) == 0 {
		return true
	}
	if x >= bounds[0].X() && x <= bounds[1].X() && bounds[0].Y() <= y-1 && bounds[1].Y() >= y-1 &&
		(CollisionData(x, y-1)&ClipSouth) == 0 {
		return true
	}
	if x >= bounds[0].X() && x <= bounds[1].X() && bounds[0].Y() <= y+1 && bounds[1].Y() >= y+1 &&
		(CollisionData(x, y+1)&ClipNorth) == 0 {
		return true
	}
	return false
}

//Hash returns a unique identifier for the tile this location points at.
// This is a perfect hashing function; every tile in the game gets a unique hashcode with it.
func (l Location) Hash() int {
	return (l.X() << 13) | l.Y()
}

func (l Location) WithinArea(area [2]Location) bool {
	return l.X() >= area[0].X() && l.X() <= area[1].X() && l.Y() >= area[0].Y() && l.Y() <= area[1].Y()
}

//ParseDirection Tries to parse the direction indicated in s.  If it can not match any direction, returns the zero-value for direction: north.
func ParseDirection(s string) Direction {
	switch s {
	case "northeast":
		return NorthEast
	case "ne":
		return NorthEast
	case "northwest":
		return NorthWest
	case "nw":
		return NorthWest
	case "east":
		return East
	case "e":
		return East
	case "west":
		return West
	case "w":
		return West
	case "south":
		return South
	case "s":
		return South
	case "southeast":
		return SouthEast
	case "se":
		return SouthEast
	case "southwest":
		return SouthWest
	case "sw":
		return SouthWest
	case "n":
		return North
	case "north":
		return North
	}

	return North
}
