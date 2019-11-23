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
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if c.Player().CurX() >= bounds[0].CurX() && c.Player().CurY() >= bounds[0].CurY() && c.Player().CurX() <= bounds[1].CurX() && c.Player().CurY() <= bounds[1].CurY() {
					c.Player().ResetPath()
					objectAction(c, object)
					return true
				}
				return false
			}
			if c.Player().WithinRange(bounds[0], 1) || c.Player().WithinRange(bounds[1], 1) {
				c.Player().ResetPath()
				objectAction(c, object)
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
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if c.Player().CurX() >= bounds[0].CurX() && c.Player().CurY() >= bounds[0].CurY() && c.Player().CurX() <= bounds[1].CurX() && c.Player().CurY() <= bounds[1].CurY() {
					c.Player().ResetPath()
					objectAction(c, object)
					return true
				}
				return false
			}
			if c.Player().WithinRange(bounds[0], 1) || c.Player().WithinRange(bounds[1], 1) {
				c.Player().ResetPath()
				objectAction(c, object)
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
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if c.Player().CurX() >= bounds[0].CurX() && c.Player().CurY() >= bounds[0].CurY() && c.Player().CurX() <= bounds[1].CurX() && c.Player().CurY() <= bounds[1].CurY() {
				c.Player().ResetPath()
				boundaryAction(c, object)
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
		bounds := object.Boundaries()
		c.Player().SetDistancedAction(func() bool {
			if c.Player().CurX() >= bounds[0].CurX() && c.Player().CurY() >= bounds[0].CurY() && c.Player().CurX() <= bounds[1].CurX() && c.Player().CurY() <= bounds[1].CurY() {
				c.Player().ResetPath()
				boundaryAction(c, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["talktonpc"] = func(c clients.Client, p *packetbuilders.Packet) {
		idx := p.ReadShort()
		npc := world.GetNpc(idx)
		if npc == nil {
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().WithinRange(npc.Location, 1) {
				c.Player().ResetPath()
				c.SendPacket(packetbuilders.ServerMessage("The " + world.NpcDefs[npc.ID].Name + " does not appear interested in talking"))
				return true
			} else {
				c.Player().SetPath(world.MakePath(c.Player().Location, npc.Location))
			}
			return false
		})
	}
}

func objectAction(c clients.Client, object *world.Object) {
	if c.Player().State != world.MSIdle || world.GetObject(object.CurX(), object.CurY()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	c.Player().State = world.MSBusy

	go func() {
		defer func() {
			c.Player().State = world.MSIdle
		}()
		if !script.Run("objectAction", c, "object", object) {
			c.SendPacket(packetbuilders.DefaultActionMessage)
		}
	}()
}

func boundaryAction(c clients.Client, object *world.Object) {
	if c.Player().State != world.MSIdle || world.GetObject(object.CurX(), object.CurY()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	c.Player().State = world.MSBusy
	go func() {
		defer func() {
			c.Player().State = world.MSIdle
		}()
		if !script.Run("boundaryAction", c, "object", object) {
			c.SendPacket(packetbuilders.DefaultActionMessage)
		}
	}()
}
