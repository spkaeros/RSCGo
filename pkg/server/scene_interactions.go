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
		p.DistancedAction = func() {
			if p.WithinRange(*o.Location(), 1) {
				newDoor := entity.NewObject(58, o.Direction, o.X(), o.Y(), false)
				region := entity.GetRegionFromLocation(*o.Location())
				region.Objects.RemoveObject(o)
				region.Objects.AddObject(newDoor)
				p.ResetDistancedAction()
			}
		}
	}
	object2Handlers[58] = func(p *entity.Player, o *entity.Object) {
		p.DistancedAction = func() {
			if p.WithinRange(*o.Location(), 1) {
				newDoor := entity.NewObject(57, o.Direction, o.X(), o.Y(), false)
				region := entity.GetRegionFromLocation(*o.Location())
				region.Objects.RemoveObject(o)
				region.Objects.AddObject(newDoor)
				p.ResetDistancedAction()
			}
		}
	}
	PacketHandlers["objectaction"] = func(c *Client, p *packets.Packet) {
		targetObject := c.player.LocalObjects.GetObject(p.ReadShort(), p.ReadShort())
		if targetObject != nil {
			if handler, ok := objectHandlers[targetObject.ID]; ok {
				handler(c.player, targetObject)
			} else if handler, ok := objectCommandHandlers[ObjectDefinitions[targetObject.ID].Commands[0]]; ok {
				handler(c.player, targetObject)
			}
		}
	}
	PacketHandlers["objectaction2"] = func(c *Client, p *packets.Packet) {
		targetObject := c.player.LocalObjects.GetObject(p.ReadShort(), p.ReadShort())
		if targetObject != nil {
			if handler, ok := object2Handlers[targetObject.ID]; ok {
				handler(c.player, targetObject)
			} else if handler, ok := objectCommand2Handlers[ObjectDefinitions[targetObject.ID].Commands[1]]; ok {
				handler(c.player, targetObject)
			}
		}
	}
}
