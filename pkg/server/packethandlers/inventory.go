package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/script"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"strings"
)

func init() {
	PacketHandlers["invwield"] = func(c clients.Client, p *packetbuilders.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", c, index)
			return
		}
		// TODO: Wielding
	}
	PacketHandlers["takeitem"] = func(c clients.Client, p *packetbuilders.Packet) {
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

		c.Player().SetDistancedAction(func() bool {
			item := world.GetItem(x, y, id)
			if item == nil || !item.VisibleTo(c.Player()) {
				log.Suspicious.Printf("Player[%v] attempted to pick up an item that doesn't exist: %d,%d,%d\n", c, id, x, y)
				return true
			}
			if !c.Player().WithinRange(item.Location, 0) {
				return false
			}
			if c.Player().Items.Size() >= 30 {
				c.Message("You do not have room for that item in your inventory.")
				return true
			}
			item.Remove()
			c.Player().Items.Put(item.ID, item.Amount)
			c.SendPacket(packetbuilders.InventoryItems(c.Player()))
			return true
		})
	}
	PacketHandlers["dropitem"] = func(c clients.Client, p *packetbuilders.Packet) {
		index := p.ReadShort()
		item := c.Player().Items.Get(index)
		if item != nil {
			if c.Player().Items.Remove(index) {
				world.AddItem(world.NewGroundItemFrom(c.Player().UserBase37, item.ID, item.Amount, int(c.Player().X.Load()), int(c.Player().Y.Load())))
				c.SendPacket(packetbuilders.InventoryItems(c.Player()))
			}
		}
	}
	PacketHandlers["invaction1"] = func(c clients.Client, p *packetbuilders.Packet) {
		index := p.ReadShort()
		item := c.Player().Items.Get(index)
		if item != nil {
			for _, s := range script.ItemTriggers {
				script.SetScriptVariable(s, "player", c)
				script.SetScriptVariable(s, "item", item)
				script.SetScriptVariable(s, "cmd", strings.ToLower(db.Items[item.ID].Command))
				if script.RunScript(s) {
					return
				}
			}
			c.SendPacket(packetbuilders.DefaultActionMessage)
		}
	}
}