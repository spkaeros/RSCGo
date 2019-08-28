package entity

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	location  *Location
	Index     int
}

//X Returns the objects X coordinate.
func (o *Object) X() int {
	return o.location.X
}

//Y Returns the objects Y coordinate.
func (o *Object) Y() int {
	return o.location.Y
}

//Location Returns the objects location in the game world.
func (o *Object) Location() *Location {
	return o.location
}

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{id, direction, boundary, &Location{x, y}, -1}
}
