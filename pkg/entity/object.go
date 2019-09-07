package entity

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	location  *Location
	index     int
	Command   string
}

func (o *Object) Index() int {
	return o.index
}

func (o *Object) SetIndex(idx int) {
	o.index = idx
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
func NewObject(id, direction, x, y int, boundary bool, command string) *Object {
	return &Object{id, direction, boundary, &Location{x, y}, -1, command}
}
