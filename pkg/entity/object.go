package entity

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	location  *Location
	Index     int
}

func (o *Object) X() int {
	return o.location.X
}

func (o *Object) Y() int {
	return o.location.Y
}

func (o *Object) Location() *Location {
	return o.location
}

func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{id, direction, boundary, &Location{x, y}, -1}
}
