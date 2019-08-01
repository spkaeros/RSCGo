package entity

type Direction uint8

const (
	MAX_X           = 944
	MAX_Y           = 3776
	North Direction = iota
	NorthWest
	West
	SouthWest
	South
	SouthEast
	East
	NorthEast
)

type Location struct {
	x, y int
}

type Locatable interface {
	Location() Location
}

// Any data structure that represents something in the game world should be able to implement this interface,
type Entity interface {
	Locatable
	// Returns the direction the given Entity is facing in the game world.
	Direction() int
}
