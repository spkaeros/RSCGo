package world

import (
	"fmt"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"go.uber.org/atomic"
)

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	Entity
}

//Equals Returns true if o1 is an object reference with identical characteristics to o.
func (o *Object) Equals(o1 objects.Object) bool {
	if o1, ok := o1.(*Object); ok {
		// We can ignore index, right?
		return o1.ID == o.ID && o1.X == o.X && o1.Y == o.Y && o1.Direction == o.Direction && o1.Boundary == o.Boundary
	}

	return false
}

func (o *Object) TypeName() string {
	return "world.Object"
}

func (o *Object) Copy() objects.Object {
	return NewObject(o.ID, o.Direction, int(o.X.Load()), int(o.Y.Load()), o.Boundary)
}

func (o *Object) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	return nil, objects.ErrInvalidOperator
}

func (o *Object) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", o.ID, o.X.Load(), o.Y.Load())
}

func (o *Object) IndexGet(index objects.Object) (objects.Object, error) {
	switch index := index.(type) {
	case *objects.String:
		switch index.Value {
		case "id":
			return &objects.Int{Value: int64(o.ID)}, nil
		case "x":
			return &objects.Int{Value: int64(o.X.Load())}, nil
		case "y":
			return &objects.Int{Value: int64(o.Y.Load())}, nil
		}
	}
	return nil, objects.ErrInvalidIndexType
}

func (o *Object) IsFalsy() bool {
	return o.X.Load() == 0 || o.Y.Load() == 0
}

var ObjectCounter = atomic.NewUint32(0)

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{id, direction, boundary, Entity{Location{X: atomic.NewUint32(uint32(x)), Y: atomic.NewUint32(uint32(y))}, int(ObjectCounter.Swap(ObjectCounter.Load() + 1))}}
}
