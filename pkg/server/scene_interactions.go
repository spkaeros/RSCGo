package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

var objectCommandHandlers = make(map[string]func(p *entity.Player, o *entity.Object))
var objectCommand2Handlers = make(map[string]func(p *entity.Player, o *entity.Object))
var objectHandlers = make(map[int]func(p *entity.Player, o *entity.Object))
var object2Handlers = make(map[int]func(p *entity.Player, o *entity.Object))

func init() {
	objectHandlers[57] = func(p *entity.Player, o *entity.Object) {
		if p.WithinRange(*o.Location(), 1) {
			region := entity.GetRegionFromLocation(*o.Location())
			region.Objects.RemoveObject(o)
			region.Objects.AddObject(entity.NewObject(58, o.Direction, o.X(), o.Y(), false))
			return
		}
	}
	object2Handlers[58] = func(p *entity.Player, o *entity.Object) {
		if p.WithinRange(*o.Location(), 1) {
			region := entity.GetRegionFromLocation(*o.Location())
			region.Objects.RemoveObject(o)
			region.Objects.AddObject(entity.NewObject(57, o.Direction, o.X(), o.Y(), false))
			return
		}
	}
	PacketHandlers["objectaction"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		if object := entity.GetRegion(x, y).Objects.GetObject(x, y); object != nil {
			c.player.WalkAction = &entity.DistancedAction{Destination: *object.Location(), Arrived: func() {
				target := entity.GetRegionFromLocation(*object.Location()).Objects.GetObject(object.X(), object.Y())
				if c.player.State != entity.Idle || target != object || !c.player.WithinRange(*target.Location(), 1) {
					return
				}
				c.player.ResetPath()
				handler, ok := objectHandlers[target.ID]
				if ok {
					// If there is a handler for this specific ID, call it
					handler(c.player, target)
					return
				}
				// Otherwise, check for handlers associated by commands, and then by names, before giving up.
				handler, ok = objectCommandHandlers[ObjectDefinitions[target.ID].Commands[0]]
				if ok {
					handler(c.player, target)
					return
				}
				c.outgoingPackets <- packets.ServerMessage("Nothing interesting happens.")
			}}
		}
	}
	PacketHandlers["objectaction2"] = func(c *Client, p *packets.Packet) {
		x := p.ReadShort()
		y := p.ReadShort()
		if object := entity.GetRegion(x, y).Objects.GetObject(x, y); object != nil {
			c.player.WalkAction = &entity.DistancedAction{Destination: *object.Location(), Arrived: func() {
				target := entity.GetRegionFromLocation(*object.Location()).Objects.GetObject(object.X(), object.Y())
				if c.player.State != entity.Idle || target != object || !c.player.WithinRange(*target.Location(), 1) {
					return
				}
				c.player.ResetPath()
				handler, ok := object2Handlers[target.ID]
				if ok {
					// If there is a handler for this specific ID, call it
					handler(c.player, target)
					return
				}
				// Otherwise, check for handlers associated by commands, and then by names, before giving up.
				handler, ok = objectCommand2Handlers[ObjectDefinitions[target.ID].Commands[1]]
				if ok {
					handler(c.player, target)
					return
				}
				c.outgoingPackets <- packets.ServerMessage("Nothing interesting happens.")
			}}
		}
	}
}
