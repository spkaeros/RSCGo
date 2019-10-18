package world

import (
	"fmt"
	"sync"
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
	X, Y int
	lock sync.RWMutex
}

//DeathSpot The spot where mobs go to die.
var DeathSpot = NewLocation(0, 0)

//NewLocation Returns a reference to a new instance of the Location data structure.
func NewLocation(x, y int) *Location {
	return &Location{X: x, Y: y}
}

//String Returns a string representation of the location
func (l *Location) String() string {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return fmt.Sprintf("[%d,%d]", l.X, l.Y)
}

//Equals Returns true if this location points to the same location as o
func (l *Location) Equals(o interface{}) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if o, ok := o.(*Location); ok {
		return l.X == o.X && l.Y == o.Y
	}
	if o, ok := o.(Location); ok {
		return l.X == o.X && l.Y == o.Y
	}
	return false
}

//DeltaX Returns the difference between this locations X coord and the other locations X coord
func (l *Location) DeltaX(other *Location) (deltaX int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if l.X > other.X {
		deltaX = l.X - other.X
	} else if other.X > l.X {
		deltaX = other.X - l.X
	}
	return
}

//DeltaY Returns the difference between this locations Y coord and the other locations Y coord
func (l *Location) DeltaY(other *Location) (deltaY int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if l.Y > other.Y {
		deltaY = l.Y - other.Y
	} else if other.Y > l.Y {
		deltaY = other.Y - l.Y
	}
	return
}

//LongestDelta Returns the largest difference in coordinates between receiver and other
func (l *Location) LongestDelta(other *Location) int {
	deltaX, deltaY := l.DeltaX(other), l.DeltaY(other)
	if deltaX > deltaY {
		return deltaX
	}
	return deltaY
}

//WithinRange Returns true if the other location is within radius tiles of the receiver location, otherwise false.
func (l *Location) WithinRange(other *Location, radius int) bool {
	return l.LongestDelta(other) <= radius
}

//Plane Calculates and returns the plane that this location is on.
func (l *Location) Plane() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return (l.Y + 100) / 944 // / 1000
}

//Above Returns the location directly above this one, if any.  Otherwise, if we are on the top floor, returns itself.
func (l *Location) Above() Location {
	return Location{X: l.X, Y: l.PlaneY(true)}
}

//Below Returns the location directly below this one, if any.  Otherwise, if we are on the bottom floor, returns itself.
func (l *Location) Below() Location {
	return Location{X: l.X, Y: l.PlaneY(false)}
}

//PlaneY Updates the location's Y coordinate, going up by one plane if up is true, else going down by one plane.  Valid planes: ground=0, 2nd story=1, 3rd story=2, basement=3
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
	l.lock.RLock()
	defer l.lock.RUnlock()
	return (newPlane * 944) + (l.Y % 944)
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
