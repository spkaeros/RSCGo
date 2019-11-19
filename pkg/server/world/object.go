package world

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/log"
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
	return fmt.Sprintf("[%v, (%v, %v)]", o.ID, o.X.Load(), o.Y.Load())
}

var ObjectCounter = atomic.NewUint32(0)

func (o *Object) WithinRange(location Location, r int) bool {
	var width = int(0)
	var height = int(0)

	cantWalk := false
	if !o.Boundary {
		def := Objects[o.ID]
		width = int(def.Height)-1
		height = int(def.Width)-1
		if o.Direction == 0 || o.Direction == 4 {
			height = int(def.Height)-1
			width = int(def.Width)-1
		}
		if def.Type != 1 {
			cantWalk = true
		}
	}
	width++
	height++
	deltaX := int(o.DeltaX(location))
	deltaY := int(o.DeltaY(location))
	if deltaX < 0 {
		deltaX = 0
	}
	if deltaY < 0 {
		deltaY = 0
	}
	if cantWalk {
		if (o.Direction == 0 || o.Direction == 4) && location.Y.Load() >= o.Y.Load() {
			deltaY++
		}
		if (o.Direction == 0 || o.Direction == 4) && location.X.Load() >= o.X.Load() {
			deltaX++
		}
		if (o.Direction == 2 || o.Direction == 6) && location.Y.Load() <= o.Y.Load() {
			deltaY++
		}
		if (o.Direction == 2 || o.Direction == 6) && location.X.Load() <= o.X.Load() {
			deltaX++
		}
	}

	log.Info.Println(deltaX, deltaY, width, height)
	return deltaX <= r+width && deltaY <= r+height
}

//NewObject Returns a reference to a new instance of a game object.
func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{ID: id, Direction: direction, Boundary: boundary,
		Entity: &Entity{
			Location: Location{X: atomic.NewUint32(uint32(x)), Y: atomic.NewUint32(uint32(y))},
			Index:    int(ObjectCounter.Swap(ObjectCounter.Load() + 1)),
		},
	}
}
