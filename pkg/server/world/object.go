package world

import (
	"fmt"
	"go.uber.org/atomic"
)


//ObjectDefinition This represents a single definition for a single object in the game.
type ObjectDefinition struct {
	ID            int
	Name          string
	Commands      []string
	Description   string
	Type          int
	Width, Height int
	Length        int
}

//Objects This holds the defining characteristics for all of the game's scene objects, ordered by ID.
var Objects []ObjectDefinition

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	*Entity
}

func (o *Object) String() string {
	return fmt.Sprintf("[%v, (%v, %v)]", o.ID, o.X(), o.Y())
}

var ObjectCounter = atomic.NewUint32(0)

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{ID: id, Direction: direction, Boundary: boundary,
		Entity: &Entity{
			Location: NewLocation(x, y),
			Index:    int(ObjectCounter.Swap(ObjectCounter.Load() + 1)),
		},
	}
}

func (o *Object) Boundaries() [2]Location {
	dir := o.Direction
	minX := o.X()
	minY := o.Y()
	maxX := minX
	maxY := minY
	if !o.Boundary {
		width := Objects[o.ID].Width
		height := Objects[o.ID].Height
		if dir != 0 && dir != 4 {
			width = Objects[o.ID].Height
			height = Objects[o.ID].Width
		}

		if Objects[o.ID].Type == 2 || Objects[o.ID].Type == 3 {
			if dir == 0 {
				width++
				minX--
			}
			if dir == 2 {
				height++
			}
			if dir == 6 {
				minY--
				height++
			}
			if dir == 4 {
				width++
			}
		}
		maxX = width + o.X() - 1
		maxY = height + o.Y() - 1
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
