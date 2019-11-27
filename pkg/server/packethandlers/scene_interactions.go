package packethandlers

import (
	"context"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"reflect"
)

func init() {
	PacketHandlers["objectaction"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
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
		c.SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if (c.NextTo(bounds[1]) || c.NextTo(bounds[0])) && c.X() >= bounds[0].X() && c.Y() >= bounds[0].Y() && c.X() <= bounds[1].X() && c.Y() <= bounds[1].Y() {
					c.ResetPath()
					objectAction(c, object)
					return true
				}
				return false
			}
			if c.NextTo(object.Location) && (c.WithinRange(bounds[0], 1) || c.WithinRange(bounds[1], 1)) {
				c.ResetPath()
				objectAction(c, object)
				return true
			}
			return false

		})
	}
	PacketHandlers["objectaction2"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
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
		c.SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if (c.NextTo(bounds[1]) || c.NextTo(bounds[0])) && c.X() >= bounds[0].X() && c.Y() >= bounds[0].Y() && c.X() <= bounds[1].X() && c.Y() <= bounds[1].Y() {
					c.ResetPath()
					objectAction(c, object)
					return true
				}
				return false
			}
			if (c.NextTo(bounds[1]) || c.NextTo(bounds[0])) && (c.WithinRange(bounds[0], 1) || c.WithinRange(bounds[1], 1)) {
				c.ResetPath()
				objectAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction2"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
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
		c.SetDistancedAction(func() bool {
			if (c.NextTo(bounds[1]) || c.NextTo(bounds[0])) && c.X() >= bounds[0].X() && c.Y() >= bounds[0].Y() && c.X() <= bounds[1].X() && c.Y() <= bounds[1].Y() {
				c.ResetPath()
				boundaryAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
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
		c.SetDistancedAction(func() bool {
			if (c.NextTo(bounds[1]) || c.NextTo(bounds[0])) && c.X() >= bounds[0].X() && c.Y() >= bounds[0].Y() && c.X() <= bounds[1].X() && c.Y() <= bounds[1].Y() {
				c.ResetPath()
				boundaryAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["talktonpc"] = func(c *world.Player, p *packet.Packet) {
		idx := p.ReadShort()
		npc := world.GetNpc(idx)
		if npc == nil {
			return
		}
		if c.Busy() {
			return
		}
		c.SetDistancedAction(func() bool {
			if c.NextTo(npc.Location) && c.WithinRange(npc.Location, 1) && !npc.Busy() {
				if c.Location.Equals(npc.Location) {
					if c.NextTo(world.NewLocation(c.X(), c.Y()-1)) {
						c.SetCoords(c.X(), c.Y()-1, false)
					} else if c.NextTo(world.NewLocation(c.X(), c.Y()+1)) {
						c.SetCoords(c.X(), c.Y()+1, false)
					} else if c.NextTo(world.NewLocation(c.X()-1, c.Y())) {
						c.SetCoords(c.X()-1, c.Y(), false)
					} else if c.NextTo(world.NewLocation(c.X()+1, c.Y())) {
						c.SetCoords(c.X()+1, c.Y(), false)
					} else {
						c.SetPath(world.MakePath(c.Location, npc.Location))
						return false
					}
				}
				c.ResetPath()
				npc.ResetPath()
				c.SetDirection(c.DirectionTo(npc.X(), npc.Y()))
				npc.SetDirection(npc.DirectionTo(c.X(), c.Y()))
				c.AddState(world.MSChatting)
				npc.AddState(world.MSChatting)
				go func() {
					defer func() {
						c.RemoveState(world.MSChatting)
						npc.RemoveState(world.MSChatting)
					}()
					for _, fn := range script.NpcTriggers {
						ran, err := fn(context.Background(), reflect.ValueOf(c), reflect.ValueOf(npc))
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
				c.SetPath(world.MakePath(c.Location, npc.Location))
			}
			return false
		})
	}
}

func objectAction(c *world.Player, object *world.Object) {
	if c.Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	c.AddState(world.MSBusy)

	go func() {
		defer func() {
			c.RemoveState(world.MSBusy)
		}()
		for _, fn := range script.ObjectTriggers {
			ran, err := fn(context.Background(), reflect.ValueOf(c), reflect.ValueOf(object))
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

func boundaryAction(c *world.Player, object *world.Object) {
	if c.Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	c.AddState(world.MSBusy)
	go func() {
		defer func() {
			c.RemoveState(world.MSBusy)
		}()
		for _, fn := range script.BoundaryTriggers {
			ran, err := fn(context.Background(), reflect.ValueOf(c), reflect.ValueOf(object))
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
