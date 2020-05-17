package world

import (
	"fmt"

	"go.uber.org/atomic"

	"github.com/spkaeros/rscgo/pkg/definitions"
)

//Object Represents a game object in the world.
type Object struct {
	Entity
	ID        int
	Direction byte
	Boundary  bool
}

func (o *Object) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", o.ID, o.X(), o.Y())
}

var ObjectCounter = atomic.NewUint32(0)

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{ID: id, Direction: byte(direction), Boundary: boundary,
		Entity: Entity{
			Location: NewLocation(x, y),
			Index:    int(ObjectCounter.Swap(ObjectCounter.Load() + 1)),
		},
	}
}

//Name checks if an object definition exists for this object, and if so returns the name associated with it.
func (o *Object) Name() string {
	if !o.Boundary {
		if o.ID < 0 || o.ID > 1188 {
			return "nil"
		}
		return definitions.ScenaryObjects[o.ID].Name
	}
	if o.ID < 0 || o.ID >= len(definitions.BoundaryObjects) {
		return "nil"
	}
	return definitions.BoundaryObjects[o.ID].Name
}

//Name checks if an object definition exists for this object, and if so returns the name associated with it.
func (o *Object) Command1() string {
	if !o.Boundary {
		if o.ID < 0 || o.ID > 1188 {
			return "nil"
		}
		return definitions.ScenaryObjects[o.ID].Commands[0]
	}
	if o.ID < 0 || o.ID >= len(definitions.BoundaryObjects) {
		return "nil"
	}
	return definitions.BoundaryObjects[o.ID].Commands[0]
}

func (o *Object) Command2() string {
	if !o.Boundary {
		if o.ID < 0 || o.ID > 1188 {
			return "nil"
		}
		return definitions.ScenaryObjects[o.ID].Commands[1]
	}
	if o.ID < 0 || o.ID >= len(definitions.BoundaryObjects) {
		return "nil"
	}
	return definitions.BoundaryObjects[o.ID].Commands[1]
}

func (o *Object) Width() int {
	if !o.Boundary {
		if o.ID < 0 || o.ID > 1188 {
			return 1
		}
		return definitions.ScenaryObjects[o.ID].Width
	}
	return 1
}

func (o *Object) Height() int {
	if !o.Boundary {
		if o.ID < 0 || o.ID > 1188 {
			return 1
		}
		return definitions.ScenaryObjects[o.ID].Height
	}
	return 1
}

func (o *Object) Boundaries() [2]Location {
	dir := o.Direction
	minX := o.X()
	minY := o.Y()
	maxX := minX
	maxY := minY
	if !o.Boundary {
		width := o.Width()
		height := o.Height()
		if dir != 0 && dir != 4 {
			width = o.Height()
			height = o.Width()
		}
		maxX = width + o.X() - 1
		maxY = height + o.Y() - 1

		if definitions.ScenaryObjects[o.ID].CollisionType == 2 || definitions.ScenaryObjects[o.ID].CollisionType == 3 {
			if dir == 0 {
				width++
				minX--
			}
			if dir == 2 {
				height++
				//minX--
			}
			if dir == 6 {
				minY--
				height++
			}
			if dir == 4 {
				width++
			}
			maxX = width + o.X() - 1
			maxY = height + o.Y() - 1
		}
	} else {
		if dir == 0 {
			minY--
		}
		if dir == 1 {
			minX--
		}
		if dir == 2 || dir == 3 {
			minX--
			minY--
			maxX++
			maxY++
		}
	}
	return [2]Location{NewLocation(minX, minY), NewLocation(maxX, maxY)}
}
