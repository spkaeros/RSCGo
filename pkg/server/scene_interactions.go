package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

type actionHandler func(p *entity.Player, args ...interface{})

var objectHandlers = make(map[interface{}]actionHandler)
var object2Handlers = make(map[interface{}]actionHandler)

//replaceObject Replaces object with a new game object having all of the same attributes, except its ID will be newID.
func replaceObject(object *entity.Object, newID int) {
	region := entity.GetRegionFromLocation(object.Location)
	region.Objects.Remove(object)
	region.Objects.Add(entity.NewObject(newID, object.Direction, object.X, object.Y, object.Boundary))
}

func init() {
	doors := make(map[int]int)
	doors[59] = 60
	doors[60] = 59
	doors[57] = 58
	doors[58] = 57
	doors[63] = 64
	doors[64] = 63
	objectHandlers["open"] = func(p *entity.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		object, ok := args[0].(*entity.Object)
		if !ok {
			LogWarning.Println("Handler for this argument type not found.")
			return
		}
		if newID, ok := doors[object.ID]; ok {
			replaceObject(object, newID)
		}
	}
	object2Handlers["close"] = func(p *entity.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		object, ok := args[0].(*entity.Object)
		if !ok {
			LogWarning.Println("Handler for this argument type not found.")
			return
		}
		if newID, ok := doors[object.ID]; ok {
			replaceObject(object, newID)
		}
	}
	PacketHandlers["objectaction"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := entity.GetObject(x, y)
		if object == nil {
			LogInfo.Println("Object not found.")
			return
		}
		c.player.RunDistancedAction(object.Location, func() {
			c.player.ResetPath()
			if c.player.State != entity.Idle || !entity.GetRegion(x, y).Objects.Contains(object) || !c.player.WithinRange(object.Location, 1) {
				// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
				return
			}
			if handler, ok := objectHandlers[object.ID]; ok {
				// If there is a handler for this specific ID, call it, and that's all we have to do.
				handler(c.player, object)
				return
			}
			if handler, ok := objectHandlers[ObjectDefinitions[object.ID].Commands[0]]; ok {
				// Otherwise, check for handlers associated by commands.
				handler(c.player, object)
				return
			}
			// Give up, concluding there isn't a handler for this object action
			c.outgoingPackets <- packets.DefaultActionMessage
		})
	}
	PacketHandlers["objectaction2"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := entity.GetObject(x, y)
		if object == nil {
			LogInfo.Println("Object not found.")
			return
		}
		c.player.RunDistancedAction(object.Location, func() {
			c.player.ResetPath()
			if c.player.State != entity.Idle || !entity.GetRegion(x, y).Objects.Contains(object) || !c.player.WithinRange(object.Location, 1) {
				// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
				return
			}
			if handler, ok := object2Handlers[object.ID]; ok {
				// If there is a handler for this specific ID, call it, and that's all we have to do.
				handler(c.player, object)
				return
			}
			if handler, ok := object2Handlers[ObjectDefinitions[object.ID].Commands[1]]; ok {
				// Otherwise, check for handlers associated by commands.
				handler(c.player, object)
				return
			}
			// Give up, concluding there isn't a handler for this object action
			c.outgoingPackets <- packets.DefaultActionMessage
		})
	}
}
