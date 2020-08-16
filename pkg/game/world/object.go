package world

import (
	"fmt"

	"go.uber.org/atomic"

	"github.com/spkaeros/rscgo/pkg/definitions"
)

//Object Represents a game object in the world.
type Object struct {
	Index int
	ID        int
	direction Direction
	Boundary  bool
	position *Location
}


func (o *Object) Direction() Direction {
	return o.direction
}

func (o *Object) X() int {
	return o.position.X()
}

func (o *Object) Y() int {
	return o.position.Y()
}

func (o *Object) ServerIndex() int {
	return o.Index
}

func (o Object) Location() *Location {
	return o.position
}


func (o *Object) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", o.ID, o.Location().X(), o.Location().Y())
}

var ObjectCounter = atomic.NewUint32(0)

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id int, direction Direction, x, y int, boundary bool) *Object {
	return &Object{ID: id, direction: Direction(direction), Boundary: boundary,
		position: NewLocation(x, y),
		Index:    int(ObjectCounter.Swap(ObjectCounter.Load() + 1)),
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

func (o *Object) IsObject() bool {
	return true
}

func (o *Object) IsNpc() bool {
	return false
}
func (o *Object) IsPlayer() bool {
	return false
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
	dir := o.Direction()
	minX := o.Location().X()
	minY := o.Location().Y()
	maxX := minX
	maxY := minY
	if !o.Boundary {
		width := o.Width()
		height := o.Height()
		if dir != 0 && dir != 4 {
			width = o.Height()
			height = o.Width()
		}
		maxX = width + o.Location().X() - 1
		maxY = height + o.Location().Y() - 1

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
			maxX = width + o.Location().X() - 1
			maxY = height + o.Location().Y() - 1
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
	return [2]Location{*NewLocation(minX, minY), *NewLocation(maxX, maxY)}
}
