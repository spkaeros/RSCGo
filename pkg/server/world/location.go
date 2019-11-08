package world

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/rand"
	"go.uber.org/atomic"
)

const (
	//North Represents north.
	North int = iota
	//NorthWest Represents north-west.
	NorthWest
	//West Represents west.
	West
	//SouthWest Represents south-west.
	SouthWest
	//South represents south.
	South
	//SouthEast represents south-east
	SouthEast
	//East Represents east.
	East
	//NorthEast Represents north-east.
	NorthEast
	//LeftFighting Represents fighting stance on the left hand side
	LeftFighting
	//RightFighting Represents fighting stance on the right hand side
	RightFighting
	//MaxX Width of the game
	MaxX = 944
	//MaxY Height of the game
	MaxY = 3776
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
	X *atomic.Uint32
	Y *atomic.Uint32
}

//DeathSpot The spot where mobs go to die.
var DeathSpot = NewLocation(0, 0)

//NewLocation Returns a reference to a new instance of the Location data structure.
func NewLocation(x, y int) Location {
	return Location{X: atomic.NewUint32(uint32(x)), Y: atomic.NewUint32(uint32(y))}
}

func (l Location) incX() {
	l.X.Store(l.X.Load() + 1)
}

func (l Location) incY() {
	l.Y.Store(l.Y.Load() + 1)
}

func (l Location) decX() {
	l.X.Store(l.X.Load() - 1)
}

func (l Location) decY() {
	l.Y.Store(l.Y.Load() - 1)
}

//NewRandomLocation Returns a new random location within the specified bounds.  bounds[0] should be lowest corner, and
// bounds[1] should be the highest corner.
func NewRandomLocation(bounds [2]Location) Location {
	return NewLocation(rand.Int31N(int(bounds[0].X.Load()), int(bounds[1].X.Load())), rand.Int31N(int(bounds[0].Y.Load()), int(bounds[1].Y.Load())))
}

//String Returns a string representation of the location
func (l *Location) String() string {
	return fmt.Sprintf("[%d,%d]", l.X.Load(), l.Y.Load())
}

//IsValid Returns true if the tile at x,y is within world boundaries, false otherwise.
func (l Location) IsValid() bool {
	return l.X.Load() <= MaxX && l.Y.Load() <= MaxY
}

//Equals Returns true if this location points to the same location as o
func (l *Location) Equals(o interface{}) bool {
	if o, ok := o.(*Location); ok {
		return l.X.Load() == o.X.Load() && l.Y.Load() == o.Y.Load()
	}
	if o, ok := o.(Location); ok {
		return l.X.Load() == o.X.Load() && l.Y.Load() == o.Y.Load()
	}
	return false
}

//DeltaX Returns the difference between this locations X coord and the other locations X coord
func (l *Location) DeltaX(other Location) (deltaX uint32) {
	ourX := l.X.Load()
	theirX := other.X.Load()
	if ourX > theirX {
		deltaX = ourX - theirX
	} else if theirX > ourX {
		deltaX = theirX - ourX
	}
	return
}

//DeltaY Returns the difference between this locations Y coord and the other locations Y coord
func (l *Location) DeltaY(other Location) (deltaY uint32) {
	ourY := l.Y.Load()
	theirY := other.Y.Load()
	if ourY > theirY {
		deltaY = ourY - theirY
	} else if theirY > ourY {
		deltaY = theirY - ourY
	}
	return
}

//LongestDelta Returns the largest difference in coordinates between receiver and other
func (l *Location) LongestDelta(other Location) uint32 {
	deltaX, deltaY := l.DeltaX(other), l.DeltaY(other)
	if deltaX > deltaY {
		return deltaX
	}
	return deltaY
}

//WithinRange Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) WithinRange(other Location, radius int) bool {
	return int(l.LongestDelta(other)) <= radius
}

//Plane Calculates and returns the plane that this location is on.
func (l *Location) Plane() int {
	return int(l.Y.Load()+100) / 944 // / 1000
}

//Above Returns the location directly above this one, if any.  Otherwise, if we are on the top floor, returns itself.
func (l *Location) Above() Location {
	return Location{X: l.X, Y: atomic.NewUint32(l.PlaneY(true))}
}

//Below Returns the location directly below this one, if any.  Otherwise, if we are on the bottom floor, returns itself.
func (l *Location) Below() Location {
	return Location{X: l.X, Y: atomic.NewUint32(l.PlaneY(false))}
}

//PlaneY Updates the location's Y coordinate, going up by one plane if up is true, else going down by one plane.  Valid planes: ground=0, 2nd story=1, 3rd story=2, basement=3
func (l *Location) PlaneY(up bool) uint32 {
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
	return uint32(newPlane*944) + (l.Y.Load() % 944)
}

//NextTileToward Returns the next tile toward the final destination of this pathway from currentLocation
func (l Location) NextTileToward(other Location) Location {
	destX, destY := other.X.Load(), other.Y.Load()
	currentX, currentY := l.X.Load(), l.Y.Load()
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
