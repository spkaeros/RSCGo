package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/script"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["objectaction"] = func(c clients.Client, p *packetbuilders.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Object not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant object at %d,%d\n", c, x, y)
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().WithinRange(object.Location, 1) {
				objectAction(c, object, false)
				return true
			}
			return false
		})
	}
	PacketHandlers["objectaction2"] = func(c clients.Client, p *packetbuilders.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Object not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant object at %d,%d\n", c, x, y)
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().WithinRange(object.Location, 1) {
				objectAction(c, object, true)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction2"] = func(c clients.Client, p *packetbuilders.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", c, x, y)
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().WithinRange(object.Location, 1) {
				boundaryAction(c, object, true)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction"] = func(c clients.Client, p *packetbuilders.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", c, x, y)
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().WithinRange(object.Location, 1) {
				boundaryAction(c, object, false)
				return true
			}
			return false
		})
	}
}

func objectAction(c clients.Client, object *world.Object, rightClick bool) {
	if c.Player().State != world.MSIdle || world.GetObject(int(object.X.Load()), int(object.Y.Load())) != object || !c.Player().WithinRange(object.Location, 1) {
		// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
		return
	}
	c.Player().State = world.MSBusy
	defer func() {
		c.Player().State = world.MSIdle
	}()
	go func() {
		for _, s := range script.ObjectTriggers {
			script.SetScriptVariable(s, "player", c)
			script.SetScriptVariable(s, "object", object)
			if rightClick {
				script.SetScriptVariable(s, "cmd", db.Objects[object.ID].Commands[1])
			} else {
				script.SetScriptVariable(s, "cmd", db.Objects[object.ID].Commands[0])
			}
			if script.RunScript(s) {
				return
			}
		}
		c.SendPacket(packetbuilders.DefaultActionMessage)
	}()
}

func boundaryAction(c clients.Client, object *world.Object, rightClick bool) {
	//	c.Player().ResetPath()
	if c.Player().State != world.MSIdle || world.GetObject(int(object.X.Load()), int(object.Y.Load())) != object || !c.Player().WithinRange(object.Location, 1) {
		// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
		return
	}
	c.Player().State = world.MSBusy
	defer func() {
		c.Player().State = world.MSIdle
	}()
	go func() {
		for _, s := range script.BoundaryTriggers {
			script.SetScriptVariable(s, "player", c)
			script.SetScriptVariable(s, "object", object)
			if rightClick {
				script.SetScriptVariable(s, "cmd", db.Boundarys[object.ID].Commands[1])
			} else {
				script.SetScriptVariable(s, "cmd", db.Boundarys[object.ID].Commands[0])
			}
			if script.RunScript(s) {
				return
			}
		}
		c.SendPacket(packetbuilders.DefaultActionMessage)
	}()
}
