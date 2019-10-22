package world

import "go.uber.org/atomic"

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	Entity
}

//Equals Returns true if o1 is an object reference with identical characteristics to o.
func (o *Object) Equals(o1 interface{}) bool {
	if o1, ok := o1.(*Object); ok {
		// We can ignore index, right?
		return o1.ID == o.ID && o1.X == o.X && o1.Y == o.Y && o1.Direction == o.Direction && o1.Boundary == o.Boundary
	}

	return false
}

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{id, direction, boundary, Entity{Location{X: atomic.NewUint32(uint32(x)), Y: atomic.NewUint32(uint32(y))}, -1}}
}
