package world

import (
	"fmt"
	"math"

	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"go.uber.org/atomic"
)
//OrderedDirections This is an array containing all of the directions a mob can walk in, ordered by path finder precedent.
// West, East, North, South, SouthWest, SouthEast, NorthWest, NorthEast
var OrderedDirections = [...]int{2, 6, 0, 4, 3, 5, 1, 7}

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


const (
	//PlaneGround Represents the value for the ground-level plane
	PlaneGround int = iota
	//PlaneSecond Represents the value for the second-story plane
	PlaneSecond
	//PlaneThird Represents the value for the third-story plane
	PlaneThird
	//PlaneBasement Represents the value for the basement plane
	PlaneBasement
)

//Location A tile in the game world.
type Location struct {
	x *atomic.Uint32
	y *atomic.Uint32
}

func (l Location) Clone() Location {
	return NewLocation(l.X(), l.Y())
}

func (l Location) X() int {
	if l.x == nil {
		return -1
	}
	
	return int(l.x.Load())
}

func (l Location) Y() int {
	if l.y == nil {
		return -1
	}
	
	return int(l.y.Load())
}

func (l Location) SetX(x int) {
	l.x.Store(uint32(x))
}

func (l Location) SetY(y int) {
	l.y.Store(uint32(y))
}

func (l Location) Wilderness() int {
	/* max wilderness X */
	if l.X() > 344 {
		return 0
	}
	return (2203-(l.Y()+1776))/6 + 1
}

var (
	//DeathSpot The spot where NPCs go to be dead.
	DeathPoint = NewLocation(0, 0)
	//SpawnPoint The default spawn point, where new players start and dead players respawn.
	SpawnPoint = Lumbridge
	//Lumbridge Lumbridge teleport point
	Lumbridge = NewLocation(122, 647)
	//Varrock Varrock teleport point
	Varrock = NewLocation(122, 647)
	//Edgeville Edgeville teleport point
	Edgeville = NewLocation(220, 445)
)

//NewLocation Returns a reference to a new instance of the Location data structure.
func NewLocation(x, y int) Location {
	return Location{x: atomic.NewUint32(uint32(x)), y: atomic.NewUint32(uint32(y))}
}

func (l Location) Point() entity.Location {
	return l.Clone()
}

func (l Location) DirectionTo(destX, destY int) int {
	sprites := [3][3]int{{SouthWest, West, NorthWest}, {South, -1, North}, {SouthEast, East, NorthEast}}
	xIndex, yIndex := l.X()-destX+1, l.Y()-destY+1
	if xIndex >= 3 || yIndex >= 3 || yIndex < 0 || xIndex < 0 {
		xIndex, yIndex = 1, 2 // North
	}
	return sprites[xIndex][yIndex]
}

//NewRandomLocation Returns a new random location within the specified bounds.  bounds[0] should be lowest corner, and
// bounds[1] should be the highest corner.
func NewRandomLocation(bounds [2]Location) Location {
	return NewLocation(rand.Rng.Intn(bounds[1].X()-bounds[0].X())+bounds[0].X(), rand.Rng.Intn(bounds[1].Y()-bounds[0].Y())+bounds[0].Y())
}

//String Returns a string representation of the location
func (l *Location) String() string {
	return fmt.Sprintf("[%d,%d]", l.X(), l.Y())
}

func (l Location) Within(minX, maxX, minY, maxY int) bool {
	return l.WithinArea([2]Location { NewLocation(minX, minY), NewLocation(maxX, maxY) })
}

//IsValid Returns true if the tile at x,y is within world boundaries, false otherwise.
func (l Location) IsValid() bool {
	return l.WithinArea([2]Location { NewLocation(0, 0), NewLocation(MaxX, MaxY)})
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

func (l *Location) Step(dir int) Location {
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
	return loc
}

//Equals Returns true if this location points to the same location as o
func (l *Location) Equals(o interface{}) bool {
	switch o.(type) {
	case Location:
		return l.LongestDelta(o.(Location)) == 0
	case *Location:
		return l.LongestDelta(*o.(*Location)) == 0
	case *Player:
		return l.LongestDelta(o.(*Player).Mob.Entity.Location) == 0
	case Player:
		return l.LongestDelta(o.(Player).Mob.Entity.Location) == 0
	case *NPC:
		return l.LongestDelta(o.(*NPC).Mob.Entity.Location) == 0
	case NPC:
		return l.LongestDelta(o.(NPC).Mob.Entity.Location) == 0
	case *Object:
		return l.LongestDelta(o.(*Object).Entity.Location) == 0
	case Object:
		return l.LongestDelta(o.(Object).Entity.Location) == 0
	case *GroundItem:
		return l.LongestDelta(o.(*GroundItem).Entity.Location) == 0
	case GroundItem:
		return l.LongestDelta(o.(GroundItem).Entity.Location) == 0
	case *Mob:
		return l.LongestDelta(o.(*Mob).Entity.Location) == 0
	case Mob:
		return l.LongestDelta(o.(Mob).Entity.Location) == 0
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
	return l.LongestDelta(NewLocation(x, y))
}

func (l Location) EuclideanDistance(other Location) float64 {
	return math.Sqrt(math.Pow(float64(l.DeltaX(other)), 2) + math.Pow(float64(l.DeltaY(other)), 2))
}

//WithinRange Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) WithinRange(other entity.Location, radius int) bool {
	return l.Near(other, radius)
}

//EntityWithin Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) Near(other entity.Location, radius int) bool {
	return l.LongestDeltaCoords(other.X(), other.Y()) <= radius
}

//Plane Calculates and returns the plane that this location is on.
func (l *Location) Plane() int {
	return int(l.y.Load()+100) / 944 // / 1000
}

//Above Returns the location directly above this one, if any.  Otherwise, if we are on the top floor, returns itself.
func (l *Location) Above() Location {
	return NewLocation(l.X(), l.PlaneY(true))
}

//Below Returns the location directly below this one, if any.  Otherwise, if we are on the bottom floor, returns itself.
func (l *Location) Below() Location {
	return NewLocation(l.X(), l.PlaneY(false))
}

func (l *Location) DirectionToward(end Location) int {
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
	return nextStep
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

func (l Location) WithinArea(area [2]Location) bool {
	return l.X() >= area[0].X() && l.X() <= area[1].X() && l.Y() >= area[0].Y() && l.Y() <= area[1].Y()
}

//ParseDirection Tries to parse the direction indicated in s.  If it can not match any direction, returns the zero-value for direction: north.
func ParseDirection(s string) int {
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
