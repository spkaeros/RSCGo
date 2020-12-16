package world

import (
	"fmt"

	"go.uber.org/atomic"

	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/entity"
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
	if boundary {
		// TODO: have Object be interface implemented by two 
		// nearly identical, yet distinct, data structures:
		// scene objects, and boundary objects...
	}
	return &Object{ID: id, Direction: byte(direction), Boundary: boundary,
		Entity: Entity{
			Location: NewLocation(x, y),
			Index:    int(ObjectCounter.Swap(ObjectCounter.Load() + 1)),
		},
	}
}

//Name checks if an object definition exists for this object, and if so returns the name associated with it.
func (o *Object) Name() string {
	if !o.Defined() {
		return "nil"
	}
	if o.Boundary {
		return definitions.BoundaryObjects[o.ID].Name
	}
	return definitions.ScenaryObjects[o.ID].Name
}

//Name checks if an object definition exists for this object, and if so returns the name associated with it.
func (o *Object) Command1() string {
	return o.Command(0)
}

func (o *Object) Command2() string {
	return o.Command(1)
}

func (o *Object) Command(click int) string {
	if !o.Defined() {
		return "nil"
	}
	if o.Boundary {
		return definitions.BoundaryObjects[o.ID].Commands[click]
	}
	return definitions.ScenaryObjects[o.ID].Commands[click]
}

//ClipType returns a unique identifier representing what kind of collisions with other entities
// to account for when it is used in-game.
//
// Doors (incl. gates) cause directional blocking only, while open doors do not block, similar to how
// ferns and signs and portraits and etc; these type of objects cause no blocking of any directions.
// doors are types 2 for shut and 3 for open
//
// The only other type is a solid object, e.g a chest or a stall or a table, which causes the
// entire tile(s) this object stands on to be blocked regardless of the origin of the other entity!
// solid objects are type 1
// 
// TODO: type 0 do anything or just a nil-like value?
func (o *Object) ClipType() int {
	if !o.Defined() {
		return 0
	}
	if o.Boundary {
		if definitions.BoundaryObjects[o.ID].Solid {
			return 2
		}
		return 3
	}
	return definitions.ScenaryObjects[o.ID].SolidityType
}

func (o *Object) Defined() bool {
	var set interface{Size() int} = definitions.BoundaryObjects
	if !o.Boundary {
		set = definitions.ScenaryObjects
	}
	if o.ID < 0 || o.ID >= set.Size() {
		return false
	}
	return true
}

//Width The width measured in game tiles that this object takes up in the game world.
func (o *Object) Width() int {
	if !o.Defined() {
		// no such object, so we take up 0x0 tiles
		return 0
	}
	if o.Boundary {
		// no large ass door boundarys exist, we take up 1x1 tiles
		return 1
	}
	return definitions.ScenaryObjects[o.ID].Width
}

//Height The height measured in game tiles that this object takes up in the game world.
func (o *Object) Height() int {
	if !o.Defined() {
		// no such object, so we take up 0x0 tiles
		return 0
	}
	if o.Boundary {
		// no large ass door boundarys exist, we take up 1x1 tiles
		return 1
	}
	return definitions.ScenaryObjects[o.ID].Height
}

func (o *Object) Boundaries() [2]entity.Location {
	dir := o.Direction
	minX, maxX, minY, maxY := o.X(), o.X(), o.Y(), o.Y()
	if o.Boundary { // are bounds quadirectional or something why did I do this?
		if dir == 0 { // only we expand one way here?
			maxX++
		}
		if dir == 1 { // and here also I guess?
			maxY++
		}
		if dir == 2 {
			minX++
			maxY++
		}
		if dir == 3 {
			maxX++
			maxY++
		}
		
		// if dir == 2 || dir == 3 { // diags possibly?  Expand in each direction the actionable positions
			// minX -= 1
			// minY -= 1
			// maxX += 1
			// maxY += 1
		// }
		return [2]entity.Location{NewLocation(minX, minY), NewLocation(maxX, maxY)}
	}
	// scenary is octadirectional for sure
	width, height := o.Width(), o.Height()
	// so this I think covers 1,2,3 and 5,6,7??  We only check for NORTH or SOUTH directly here?
	if dir != 0 && dir != 4 {
		width, height = height, width
	}
	maxX, maxY = width + o.X() - 1, height + o.Y() - 1

	// open or closed door...
	if o.ClipType() >= 2 {
		if dir == byte(North) /* 0 */ {
			width++
			minX--
		}
		if dir == byte(West) /* 2 */ {
			height++
			//minX--
		}
		if dir == byte(East) /* 6 */ {
			minY--
			height++
		}
		if dir == byte(South) /* 4 */ {
			width++
		}
		maxX = width + o.X() - 1
		maxY = height + o.Y() - 1
	}
	return [2]entity.Location{NewLocation(minX, minY), NewLocation(maxX, maxY)}
}
