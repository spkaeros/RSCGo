package packethandlers

import (
	"context"
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"reflect"
)

func init() {
	PacketHandlers["objectaction"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Object not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant object at %d,%d\n", c, x, y)
			return
		}
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if (c.Player().NextTo(bounds[1]) || c.Player().NextTo(bounds[0])) && c.Player().X() >= bounds[0].X() && c.Player().Y() >= bounds[0].Y() && c.Player().X() <= bounds[1].X() && c.Player().Y() <= bounds[1].Y() {
					c.Player().ResetPath()
					objectAction(c, object)
					return true
				}
				return false
			}
			if c.Player().NextTo(object.Location) && (c.Player().WithinRange(bounds[0], 1) || c.Player().WithinRange(bounds[1], 1)) {
				c.Player().ResetPath()
				objectAction(c, object)
				return true
			}
			return false

		})
	}
	PacketHandlers["objectaction2"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Object not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant object at %d,%d\n", c, x, y)
			return
		}
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if (c.Player().NextTo(bounds[1]) || c.Player().NextTo(bounds[0])) && c.Player().X() >= bounds[0].X() && c.Player().Y() >= bounds[0].Y() && c.Player().X() <= bounds[1].X() && c.Player().Y() <= bounds[1].Y() {
					c.Player().ResetPath()
					objectAction(c, object)
					return true
				}
				return false
			}
			if (c.Player().NextTo(bounds[1]) || c.Player().NextTo(bounds[0])) && (c.Player().WithinRange(bounds[0], 1) || c.Player().WithinRange(bounds[1], 1)) {
				c.Player().ResetPath()
				objectAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction2"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", c, x, y)
			return
		}
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if (c.Player().NextTo(bounds[1]) || c.Player().NextTo(bounds[0])) && c.Player().X() >= bounds[0].X() && c.Player().Y() >= bounds[0].Y() && c.Player().X() <= bounds[1].X() && c.Player().Y() <= bounds[1].Y() {
				c.Player().ResetPath()
				boundaryAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction"] = func(c clients.Client, p *packet.Packet) {
		if c.Player().Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", c, x, y)
			return
		}
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if (c.Player().NextTo(bounds[1]) || c.Player().NextTo(bounds[0])) && c.Player().X() >= bounds[0].X() && c.Player().Y() >= bounds[0].Y() && c.Player().X() <= bounds[1].X() && c.Player().Y() <= bounds[1].Y() {
				c.Player().ResetPath()
				boundaryAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["talktonpc"] = func(c clients.Client, p *packet.Packet) {
		idx := p.ReadShort()
		npc := world.GetNpc(idx)
		if npc == nil {
			return
		}
		if c.Player().Busy() {
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().NextTo(npc.Location) && c.Player().WithinRange(npc.Location, 1) && !npc.Busy() {
				if c.Player().Location.Equals(npc.Location) {
					if c.Player().NextTo(world.NewLocation(c.Player().X(), c.Player().Y()-1)) {
						c.Player().SetCoords(c.Player().X(), c.Player().Y()-1, false)
					} else if c.Player().NextTo(world.NewLocation(c.Player().X(), c.Player().Y()+1)) {
						c.Player().SetCoords(c.Player().X(), c.Player().Y()+1, false)
					} else if c.Player().NextTo(world.NewLocation(c.Player().X()-1, c.Player().Y())) {
						c.Player().SetCoords(c.Player().X()-1, c.Player().Y(), false)
					} else if c.Player().NextTo(world.NewLocation(c.Player().X()+1, c.Player().Y())) {
						c.Player().SetCoords(c.Player().X()+1, c.Player().Y(), false)
					} else {
						c.Player().SetPath(world.MakePath(c.Player().Location, npc.Location))
						return false
					}
				}
				c.Player().ResetPath()
				npc.ResetPath()
				c.Player().SetDirection(c.Player().DirectionTo(npc.X(), npc.Y()))
				npc.SetDirection(npc.DirectionTo(c.Player().X(), c.Player().Y()))
				c.Player().AddState(world.MSChatting)
				npc.AddState(world.MSChatting)
				go func() {
					defer func() {
						c.Player().RemoveState(world.MSChatting)
						npc.RemoveState(world.MSChatting)
					}()
					for _, fn := range script.NpcTriggers {
						ran, err := fn(context.Background(), reflect.ValueOf(c.Player()), reflect.ValueOf(npc))
						if !ran.IsValid() {
							continue
						}
						if !err.IsNil() {
							log.Info.Println(err)
							continue
						}
						if ran.Bool() {
							return
						}
					}
					c.SendPacket(packetbuilders.ServerMessage("The " + world.NpcDefs[npc.ID].Name + " does not appear interested in talking"))
				}()
				return true
			} else {
				c.Player().SetPath(world.MakePath(c.Player().Location, npc.Location))
			}
			return false
		})
	}
}

func objectAction(c clients.Client, object *world.Object) {
	if c.Player().Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	c.Player().AddState(world.MSBusy)

	go func() {
		defer func() {
			c.Player().RemoveState(world.MSBusy)
		}()
		for _, fn := range script.ObjectTriggers {
			ran, err := fn(context.Background(), reflect.ValueOf(c.Player()), reflect.ValueOf(object))
			if !ran.IsValid() {
				continue
			}
			if !err.IsNil() {
				log.Info.Println(err)
				continue
			}
			if ran.Bool() {
				return
			}
		}
		c.SendPacket(packetbuilders.DefaultActionMessage)
		//		if !script.Run("objectAction", c, "object", object) {
		//			c.SendPacket(packetbuilders.DefaultActionMessage)
		//		}
	}()
}

func boundaryAction(c clients.Client, object *world.Object) {
	if c.Player().Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	c.Player().AddState(world.MSBusy)
	go func() {
		defer func() {
			c.Player().RemoveState(world.MSBusy)
		}()
		for _, fn := range script.BoundaryTriggers {
			ran, err := fn(context.Background(), reflect.ValueOf(c.Player()), reflect.ValueOf(object))
			if !ran.IsValid() {
				continue
			}
			if !err.IsNil() {
				log.Info.Println(err)
				continue
			}
			if ran.Bool() {
				return
			}
		}
		c.SendPacket(packetbuilders.DefaultActionMessage)
		//		if !script.Run("boundaryAction", c, "object", object) {
		//			c.SendPacket(packetbuilders.DefaultActionMessage)
		//		}
	}()
}
