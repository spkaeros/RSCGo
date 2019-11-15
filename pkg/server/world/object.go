package world

import (
	"fmt"
	"go.uber.org/atomic"
)

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	*Entity
}

func (o *Object) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", o.ID, o.X.Load(), o.Y.Load())
}

var ObjectCounter = atomic.NewUint32(0)

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{ID: id, Direction: direction, Boundary: boundary,
		Entity: &Entity{
			Location: Location{X: atomic.NewUint32(uint32(x)), Y: atomic.NewUint32(uint32(y))},
			Index:    int(ObjectCounter.Swap(ObjectCounter.Load() + 1)),
		},
	}
}
