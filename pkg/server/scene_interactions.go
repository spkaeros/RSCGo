package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/world"
)

type actionHandler func(p *world.Player, args ...interface{})

var objectHandlers = make(map[interface{}]actionHandler)
var object2Handlers = make(map[interface{}]actionHandler)

func init() {
	doors := make(map[int]int)
	doors[59] = 60
	doors[57] = 58
	doors[63] = 64
	for k, v := range doors {
		// Add value->key to handle close as well as open.
		doors[v] = k
	}
	objectHandlers["open"] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		object, ok := args[0].(*world.Object)
		if !ok {
			LogWarning.Println("Handler for this argument type not found.")
			return
		}
		if newID, ok := doors[object.ID]; ok {
			world.ReplaceObject(object, newID)
		}
	}
	object2Handlers["close"] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		object, ok := args[0].(*world.Object)
		if !ok {
			LogWarning.Println("Handler for this argument type not found.")
			return
		}
		if newID, ok := doors[object.ID]; ok {
			world.ReplaceObject(object, newID)
		}
	}
	PacketHandlers["objectaction"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			LogInfo.Println("Object not found.")
			return
		}
		c.player.RunDistancedAction(object.Location, func() {
			c.player.ResetPath()
			if c.player.State != world.Idle || !world.GetRegion(x, y).Objects.Contains(object) || !c.player.WithinRange(object.Location, 1) {
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
		object := world.GetObject(x, y)
		if object == nil {
			LogInfo.Println("Object not found.")
			return
		}
		c.player.RunDistancedAction(object.Location, func() {
			c.player.ResetPath()
			if c.player.State != world.Idle || !world.GetRegion(x, y).Objects.Contains(object) || !c.player.WithinRange(object.Location, 1) {
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
