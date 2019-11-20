package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
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
			if object.WithinRange(c.Player().Location, 0) {
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
			if object.WithinRange(c.Player().Location, 0) {
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
			if object.WithinRange(c.Player().Location, 0) {
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
			if object.WithinRange(c.Player().Location, 0) {
				boundaryAction(c, object, false)
				return true
			}
			return false
		})
	}
}

func objectAction(c clients.Client, object *world.Object, rightClick bool) {
	if c.Player().State != world.MSIdle || world.GetObject(object.CurX(), object.CurY()) != object || !object.WithinRange(c.Player().Location, 0) {
		// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
		return
	}
	c.Player().State = world.MSBusy
	defer func() {
		c.Player().State = world.MSIdle
	}()

	c.Player().ResetPath()
	go func() {
		env := script.WorldModule()
		err := env.Define("client", c)
		if err != nil {
			log.Info.Println("Error initializing scripting environment:", err)
			return
		}
		err = env.Define("player", c.Player())
		if err != nil {
			log.Info.Println("Error initializing scripting environment:", err)
			return
		}
		err = env.Define("object", object)
		if err != nil {
			log.Info.Println("Error initializing scripting environment:", err)
			return
		}
		if rightClick {
			err = env.Define("cmd", world.Objects[object.ID].Commands[1])
			if err != nil {
				log.Info.Println("Error initializing scripting environment:", err)
				return
			}
		} else {
			err = env.Define("cmd", world.Objects[object.ID].Commands[0])
			if err != nil {
				log.Info.Println("Error initializing scripting environment:", err)
				return
			}
		}
		for _, s := range script.Scripts {
			scriptTriggered, err := env.Execute(s + `
objectAction()`)
			if err != nil {
//				log.Info.Println(err)
				continue
			}
			if scriptTriggered.(bool) {
				return
			}
		}
		c.SendPacket(packetbuilders.DefaultActionMessage)
	}()
}

func boundaryAction(c clients.Client, object *world.Object, rightClick bool) {
	//	c.Player().ResetPath()
	if c.Player().State != world.MSIdle || world.GetObject(object.CurX(), object.CurY()) != object || !object.WithinRange(c.Player().Location, 0) {
		// If somehow we became busy, the object changed before arriving, or somehow this action fired without actually arriving at the object, we do nothing.
		return
	}
	c.Player().State = world.MSBusy
	defer func() {
		c.Player().State = world.MSIdle
	}()
	c.Player().ResetPath()
	go func() {
		env := script.WorldModule()
		err := env.Define("client", c)
		if err != nil {
			log.Info.Println("Error initializing scripting environment:", err)
			return
		}
		err = env.Define("player", c.Player())
		if err != nil {
			log.Info.Println("Error initializing scripting environment:", err)
			return
		}
		err = env.Define("object", object)
		if err != nil {
			log.Info.Println("Error initializing scripting environment:", err)
			return
		}
		if rightClick {
			err = env.Define("cmd", world.Boundarys[object.ID].Commands[1])
			if err != nil {
				log.Info.Println("Error initializing scripting environment:", err)
				return
			}
		} else {
			err = env.Define("cmd", world.Boundarys[object.ID].Commands[0])
			if err != nil {
				log.Info.Println("Error initializing scripting environment:", err)
				return
			}
		}
		for _, s := range script.Scripts {
			scriptTriggered, err := env.Execute(s + `
boundaryAction()`)
			if err != nil {
//				log.Info.Println(err)
				continue
			}
			if scriptTriggered.(bool) {
				return
			}
		}
		c.SendPacket(packetbuilders.DefaultActionMessage)
	}()
}
