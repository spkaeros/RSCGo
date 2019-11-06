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

var itemAffectedTypes = make(map[int][]int)

func init() {
	itemAffectedTypes[32] = []int{32, 33}
	itemAffectedTypes[33] = []int{32, 33}
	itemAffectedTypes[64] = []int{64, 322}
	itemAffectedTypes[512] = []int{512, 640, 644}
	itemAffectedTypes[8] = []int{8, 24, 8216}
	itemAffectedTypes[1024] = []int{1024}
	itemAffectedTypes[128] = []int{128, 640, 644}
	itemAffectedTypes[644] = []int{128, 512, 640, 644}
	itemAffectedTypes[640] = []int{128, 512, 640, 644}
	itemAffectedTypes[2048] = []int{2048}
	itemAffectedTypes[16] = []int{16, 24, 8216}
	itemAffectedTypes[256] = []int{256, 322}
	itemAffectedTypes[322] = []int{64, 256, 322}
	itemAffectedTypes[24] = []int{8, 16, 24, 8216}
	itemAffectedTypes[8216] = []int{8, 16, 24, 8216}
	PacketHandlers["invwield"] = func(c clients.Client, p *packetbuilders.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", c, index)
			return
		}
		if item := c.Player().Items.Get(index); item != nil {
			if e := db.GetEquipmentDefinition(item.ID); e != nil {
				item.Worn = true
				c.Player().TransAttrs.SetVar("self", false)
				for _, otherItem := range c.Player().Items.List {
					if otherE := db.GetEquipmentDefinition(otherItem.ID); otherE != nil {
						for _, i := range itemAffectedTypes[e.Type] {
							if otherItem != item && otherItem.Worn && i == otherE.Type {
								c.Player().SetAimPoints(c.Player().AimPoints() - otherE.Aim)
								c.Player().SetPowerPoints(c.Player().PowerPoints() - otherE.Power)
								c.Player().SetArmourPoints(c.Player().ArmourPoints() - otherE.Armour)
								c.Player().SetMagicPoints(c.Player().MagicPoints() - otherE.Magic)
								c.Player().SetPrayerPoints(c.Player().PrayerPoints() - otherE.Prayer)
								c.Player().SetRangedPoints(c.Player().RangedPoints() - otherE.Ranged)
								otherItem.Worn = false
								value := 0
								if otherE.Position == 2 {
									value = c.Player().Appearance.Legs
								}
								c.Player().Equips[otherE.Position] = value
							}
						}
					}
				}
				c.Player().SetAimPoints(c.Player().AimPoints() + e.Aim)
				c.Player().SetPowerPoints(c.Player().PowerPoints() + e.Power)
				c.Player().SetArmourPoints(c.Player().ArmourPoints() + e.Armour)
				c.Player().SetMagicPoints(c.Player().MagicPoints() + e.Magic)
				c.Player().SetPrayerPoints(c.Player().PrayerPoints() + e.Prayer)
				c.Player().SetRangedPoints(c.Player().RangedPoints() + e.Ranged)
				c.Player().Equips[e.Position] = e.Sprite
				c.Player().AppearanceTicket++
				c.SendPacket(packetbuilders.EquipmentStats(c.Player()))
				c.SendPacket(packetbuilders.InventoryItems(c.Player()))
			}
		}
	}
	PacketHandlers["removeitem"] = func(c clients.Client, p *packetbuilders.Packet) {
		index := p.ReadShort()
		if index < 0 || index >= 30 {
			log.Suspicious.Printf("Player[%v] tried to wield an item with invalid index: %d\n", c, index)
			return
		}
		// TODO: Wielding
		if item := c.Player().Items.Get(index); item != nil {
			if !item.Worn {
				return
			}
			if e := db.GetEquipmentDefinition(item.ID); e != nil {
				c.Player().TransAttrs.SetVar("self", false)
				item.Worn = false
				c.Player().SetAimPoints(c.Player().AimPoints() - e.Aim)
				c.Player().SetPowerPoints(c.Player().PowerPoints() - e.Power)
				c.Player().SetArmourPoints(c.Player().ArmourPoints() - e.Armour)
				c.Player().SetMagicPoints(c.Player().MagicPoints() - e.Magic)
				c.Player().SetPrayerPoints(c.Player().PrayerPoints() - e.Prayer)
				c.Player().SetRangedPoints(c.Player().RangedPoints() - e.Ranged)
				value := 0
				if e.Position == 2 {
					value = c.Player().Appearance.Legs
				}
				c.Player().Equips[e.Position] = value
				c.Player().AppearanceTicket++
				c.SendPacket(packetbuilders.EquipmentStats(c.Player()))
				c.SendPacket(packetbuilders.InventoryItems(c.Player()))
			}
		}
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