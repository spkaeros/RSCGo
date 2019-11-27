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
	PacketHandlers["invwield"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", player, index)
			return
		}
		if item := player.Items.Get(index); item != nil {
			if item.Worn {
				return
			}
			player.SendPacket(packetbuilders.Sound("click"))
			player.EquipItem(item)
			player.SendPacket(packetbuilders.EquipmentStats(player))
			player.SendPacket(packetbuilders.InventoryItems(player))
		}
	}
	PacketHandlers["removeitem"] = func(player *world.Player, p *packet.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", player, index)
			return
		}
		// TODO: Wielding
		if item := player.Items.Get(index); item != nil {
			if !item.Worn {
				return
			}
			player.SendPacket(packetbuilders.Sound("click"))
			player.DequipItem(item)
			player.SendPacket(packetbuilders.EquipmentStats(player))
			player.SendPacket(packetbuilders.InventoryItems(player))
		}
	}
	PacketHandlers["takeitem"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		id := p.ReadShort()
		p.ReadShort() // Useless, this variable is for what affect we are applying to the ground item, e.g casting, using item with
		if x < 0 || x >= world.MaxX || y < 0 || y >= world.MaxY {
			log.Suspicious.Printf("Player[%v] attempted to pick up an item at an invalid location: [%d,%d]\n", player, x, y)
			return
		}
		if id < 0 || id > 1289 {
			log.Suspicious.Printf("Player[%v] attempted to pick up an item with an invalid ID: %d\n", player, id)
			return
		}

		player.SetDistancedAction(func() bool {
			item := world.GetItem(x, y, id)
			if item == nil || !item.VisibleTo(player) {
				log.Suspicious.Printf("Player[%v] attempted to pick up an item that doesn't exist: %d,%d,%d\n", player, id, x, y)
				return true
			}
			if !player.WithinRange(item.Location, 0) {
				return false
			}
			player.SendPacket(packetbuilders.Sound("takeobject"))
			player.ResetPath()
			if player.Items.Size() >= 30 {
				player.SendPacket(packetbuilders.ServerMessage("You do not have room for that item in your inventory."))
				return true
			}
			item.Remove()
			player.Items.Add(item.ID, item.Amount)
			player.SendPacket(packetbuilders.InventoryItems(player))
			return true
		})
	}
	PacketHandlers["dropitem"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		index := p.ReadShort()
		item := player.Items.Get(index)
		if item != nil {
			player.SetDistancedAction(func() bool {
				if player.FinishedPath() {
					if player.Items.Remove(index, item.Amount) {
						world.AddItem(world.NewGroundItemFor(player.UserBase37, item.ID, item.Amount, player.X(), player.Y()))
						player.SendPacket(packetbuilders.InventoryItems(player))
						player.SendPacket(packetbuilders.Sound("dropobject"))
					}
					return true
				}
				return false
			})
		}
	}
	PacketHandlers["invaction1"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		index := p.ReadShort()
		item := player.Items.Get(index)
		if item != nil {
			player.AddState(world.MSBusy)
			go func() {
				defer func() {
					player.RemoveState(world.MSBusy)
				}()
				for _, fn := range script.InvTriggers {
					ran, err := fn(context.Background(), reflect.ValueOf(player), reflect.ValueOf(item))
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
				player.SendPacket(packetbuilders.DefaultActionMessage)
				//				if !script.Run("invAction", player, "item", item) {
				//					player.SendPacket(packetbuilders.DefaultActionMessage)
				//				}
			}()
		}
	}
}
