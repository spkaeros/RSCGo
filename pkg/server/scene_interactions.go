package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

type actionHandler func(p *world.Player, args ...interface{})
type actionsMap map[interface{}]actionHandler

var objectHandlers = make(actionsMap)
var object2Handlers = make(actionsMap)

var boundaryHandlers = make(actionsMap)
var boundary2Handlers = make(actionsMap)

func init() {
	//TODO: This whole entire file is messy and could use tidying.
	// Actually, to that end, I will be implementing a scripting language of some sort, so I'll leave it for now.
	oDoors := make(map[int]int)
	oDoors[59] = 60
	oDoors[57] = 58
	oDoors[63] = 64
	for k, v := range oDoors {
		// Add value->key to handle close as well as open.
		oDoors[v] = k
	}
	bDoors := make(map[int]int)
	bDoors[2] = 1
	for k, v := range bDoors {
		// Add value->key to handle close as well as open.
		bDoors[v] = k
	}
	objectHandlers[19] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		if p.Skillset.Current[5] < p.Skillset.Maximum[5] {
			c, _ := Clients.FromIndex(p.Index)
			p.Skillset.Current[5] = p.Skillset.Maximum[5]
			c.UpdateStat(5)
			c.Message("You recharge your prayer points at the altar.")
		}
	}
	objectHandlers["climb-up"] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		if nextLocation := p.Above(); !nextLocation.Equals(p.Location) {
			c, _ := Clients.FromIndex(p.Index)
			p.SetLocation(nextLocation)
			c.UpdatePlane()
		}
	}
	objectHandlers["climb-down"] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		if nextLocation := p.Below(); !nextLocation.Equals(p.Location) {
			c, _ := Clients.FromIndex(p.Index)
			p.SetLocation(nextLocation)
			c.UpdatePlane()
		}
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
		if newID, ok := oDoors[object.ID]; ok {
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
		if newID, ok := oDoors[object.ID]; ok {
			world.ReplaceObject(object, newID)
		}
	}
	boundaryHandlers["open"] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		object, ok := args[0].(*world.Object)
		if !ok {
			LogWarning.Println("Handler for this argument type not found.")
			return
		}
		if object.ID == 109 {
			// Quest hut by wilderness in between edgeville and varrock
			dest := world.Location{161, 465}
			if p.Y >= dest.Y {
				dest.Y--
			}
			go p.EnterDoor(object, dest)
		}
		if newID, ok := bDoors[object.ID]; ok {
			world.ReplaceObject(object, newID)
		}
	}
	boundary2Handlers["close"] = func(p *world.Player, args ...interface{}) {
		if len(args) <= 0 {
			LogWarning.Println("Must provide at least 1 argument to action handlers.")
			return
		}

		object, ok := args[0].(*world.Object)
		if !ok {
			LogWarning.Println("Handler for this argument type not found.")
			return
		}
		if newID, ok := bDoors[object.ID]; ok {
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
			objectAction(c, object, false)
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
			objectAction(c, object, true)
		})
	}
	PacketHandlers["boundaryaction2"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			LogInfo.Println("Boundary not found.")
			return
		}
		c.player.RunDistancedAction(object.Location, func() {
			boundaryAction(c, object, true)
		})
	}
	PacketHandlers["boundaryaction"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			LogInfo.Println("Boundary not found.")
			return
		}
		c.player.RunDistancedAction(object.Location, func() {
			boundaryAction(c, object, false)
		})
	}
}

func objectAction(c *Client, object *world.Object, rightClick bool) {
	c.player.ResetPath()
	if c.player.State != world.MSIdle || world.GetObject(object.X, object.Y) != object || !c.player.WithinRange(object.Location, 1) {
		// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
		return
	}
	handlers := objectHandlers
	command := ObjectDefinitions[object.ID].Commands[0]
	if rightClick {
		handlers = object2Handlers
		command = ObjectDefinitions[object.ID].Commands[1]
	}
	if handler, ok := handlers[object.ID]; ok {
		// If there is a handler for this specific ID, call it, and that's all we have to do.
		handler(c.player, object)
		return
	}
	if handler, ok := handlers[command]; ok {
		// Otherwise, check for handlers associated by commands.
		handler(c.player, object)
		return
	}
	// Give up, concluding there isn't a handler for this object action
	c.outgoingPackets <- packets.DefaultActionMessage
}

func boundaryAction(c *Client, object *world.Object, rightClick bool) {
	c.player.ResetPath()
	if c.player.State != world.MSIdle || world.GetObject(object.X, object.Y) != object || !c.player.WithinRange(object.Location, 1) {
		// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
		return
	}
	handlers := boundaryHandlers
	command := BoundaryDefinitions[object.ID].Commands[0]
	if rightClick {
		handlers = boundary2Handlers
		command = BoundaryDefinitions[object.ID].Commands[1]
	}
	if handler, ok := handlers[object.ID]; ok {
		// If there is a handler for this specific ID, call it, and that's all we have to do.
		handler(c.player, object)
		return
	}
	if handler, ok := handlers[command]; ok {
		// Otherwise, check for handlers associated by commands.
		handler(c.player, object)
		return
	}
	// Give up, concluding there isn't a handler for this object action
	c.outgoingPackets <- packets.DefaultActionMessage
}
