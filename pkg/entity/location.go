package entity

import "fmt"

//Direction Direction within gameworld.
type Direction uint8

const (
	//North Represents north.
	North Direction = iota
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

//Location A tile in the game world.
type Location struct {
	X, Y int
}

func NewLocation(x, y int) *Location {
	return &Location{x, y}
}

//String Returns a string representation of the location
func (l *Location) String() string {
	return fmt.Sprintf("[%d,%d]", l.X, l.Y)
}

//LongestDelta Returns the largest difference in coordinates between receiver and other
func (l *Location) LongestDelta(other *Location) int {
	deltaX := l.X - other.X
	if deltaX < 0 {
		deltaX = -deltaX
	}
	deltaY := l.Y - other.Y
	if deltaY < 0 {
		deltaY = -deltaY
	}
	if deltaX > deltaY {
		return deltaX
	}
	return deltaY
}
