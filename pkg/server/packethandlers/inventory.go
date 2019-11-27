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
	PacketHandlers["invwield"] = func(c *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", c, index)
			return
		}
		if item := c.Items.Get(index); item != nil {
			if item.Worn {
				return
			}
			c.SendPacket(packetbuilders.Sound("click"))
			c.EquipItem(item)
			c.SendPacket(packetbuilders.EquipmentStats(c))
			c.SendPacket(packetbuilders.InventoryItems(c))
		}
	}
	PacketHandlers["removeitem"] = func(c *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", c, index)
			return
		}
		// TODO: Wielding
		if item := c.Items.Get(index); item != nil {
			if !item.Worn {
				return
			}
			c.SendPacket(packetbuilders.Sound("click"))
			c.DequipItem(item)
			c.SendPacket(packetbuilders.EquipmentStats(c))
			c.SendPacket(packetbuilders.InventoryItems(c))
		}
	}
	PacketHandlers["takeitem"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		id := p.ReadShort()
		p.ReadShort() // Useless, this variable is for what affect we are applying to the ground item, e.g casting, using item with
		if x < 0 || x >= world.MaxX || y < 0 || y >= world.MaxY {
			log.Suspicious.Printf("Player[%v] attempted to pick up an item at an invalid location: [%d,%d]\n", c, x, y)
			return
		}
		if id < 0 || id > 1289 {
			log.Suspicious.Printf("Player[%v] attempted to pick up an item with an invalid ID: %d\n", c, id)
			return
		}

		c.SetDistancedAction(func() bool {
			item := world.GetItem(x, y, id)
			if item == nil || !item.VisibleTo(c) {
				log.Suspicious.Printf("Player[%v] attempted to pick up an item that doesn't exist: %d,%d,%d\n", c, id, x, y)
				return true
			}
			if !c.WithinRange(item.Location, 0) {
				return false
			}
			c.SendPacket(packetbuilders.Sound("takeobject"))
			c.ResetPath()
			if c.Items.Size() >= 30 {
				c.SendPacket(packetbuilders.ServerMessage("You do not have room for that item in your inventory."))
				return true
			}
			item.Remove()
			c.Items.Add(item.ID, item.Amount)
			c.SendPacket(packetbuilders.InventoryItems(c))
			return true
		})
	}
	PacketHandlers["dropitem"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
			return
		}
		index := p.ReadShort()
		item := c.Items.Get(index)
		if item != nil {
			c.SetDistancedAction(func() bool {
				if c.FinishedPath() {
					if c.Items.Remove(index, item.Amount) {
						world.AddItem(world.NewGroundItemFor(c.UserBase37, item.ID, item.Amount, c.X(), c.Y()))
						c.SendPacket(packetbuilders.InventoryItems(c))
						c.SendPacket(packetbuilders.Sound("dropobject"))
					}
					return true
				}
				return false
			})
		}
	}
	PacketHandlers["invaction1"] = func(c *world.Player, p *packet.Packet) {
		if c.Busy() {
			return
		}
		index := p.ReadShort()
		item := c.Items.Get(index)
		if item != nil {
			c.AddState(world.MSBusy)
			go func() {
				defer func() {
					c.RemoveState(world.MSBusy)
				}()
				for _, fn := range script.InvTriggers {
					ran, err := fn(context.Background(), reflect.ValueOf(c), reflect.ValueOf(item))
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
				//				if !script.Run("invAction", c, "item", item) {
				//					c.SendPacket(packetbuilders.DefaultActionMessage)
				//				}
			}()
		}
	}
}
