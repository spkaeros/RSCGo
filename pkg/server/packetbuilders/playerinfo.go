package packetbuilders

import (
	"github.com/spkaeros/rscgo/pkg/server/world"
)

//ChangeAppearance The appearance changing window.
var ChangeAppearance = NewOutgoingPacket(59)

//InventoryItems Builds a packet containing the players inventory items.
func InventoryItems(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(53)
	p.AddByte(uint8(player.Items.Size()))
	player.Items.Range(func(item *world.Item) bool {
		if item.Worn {
			p.AddShort(uint16(item.ID + 0x8000))
		} else {
			p.AddShort(uint16(item.ID))
		}
		if world.ItemDefs[item.ID].Stackable {
			p.AddInt2(uint32(item.Amount))
		}
		return true
	})
	return
}

//FightMode Builds a packet with the players fight mode information in it.
func FightMode(player *world.Player) (p *Packet) {
	// TODO: 204
	p = NewOutgoingPacket(132)
	p.AddByte(byte(player.FightMode()))
	return p
}

//Fatigue Builds a packet with the players fatigue percentage in it.
func Fatigue(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(114)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the server to match RSC
	p.AddShort(uint16(player.Fatigue() / 100))
	return p
}

//ClientSettings Builds a packet containing the players client settings, e.g camera mode, mouse mode, sound fx...
func ClientSettings(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(240)
	// TODO: Right IDs?
	if player.GetClientSetting(0) {
		p.AddByte(1)
	} else {
		p.AddByte(0)
	}
	if player.GetClientSetting(2) {
		p.AddByte(1)
	} else {
		p.AddByte(0)
	}
	if player.GetClientSetting(3) {
		p.AddByte(1)
	} else {
		p.AddByte(0)
	}

	//	p.AddByte(0) // Camera auto/manual?
	//	p.AddByte(0) // Mouse buttons 1 or 2?
	//	p.AddByte(1) // Sound effects on/off?
	return
}

//PlayerStats Builds a packet containing all the player's stat information and returns it.
func PlayerStats(player *world.Player) *Packet {
	p := NewOutgoingPacket(156)
	for i := 0; i < 18; i++ {
		p.AddByte(uint8(player.Skillset.Current[i]))
	}

	for i := 0; i < 18; i++ {
		p.AddByte(uint8(player.Skillset.Maximum[i]))
	}

	for i := 0; i < 18; i++ {
		p.AddInt(uint32(player.Skillset.Experience[i]))
	}
	return p
}

//PlayerStat Builds a packet containing player's stat information for skill at idx and returns it.
func PlayerStat(player *world.Player, idx int) *Packet {
	p := NewOutgoingPacket(159)
	p.AddByte(byte(idx))
	p.AddInt(uint32(player.Skillset.Experience[idx]))
	return p
}

//EquipmentStats Builds a packet with the players equipment statistics in it.
func EquipmentStats(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(153)
	p.AddByte(uint8(player.ArmourPoints()))
	p.AddByte(uint8(player.AimPoints()))
	p.AddByte(uint8(player.PowerPoints()))
	p.AddByte(uint8(player.MagicPoints()))
	p.AddByte(uint8(player.PrayerPoints()))
	p.AddByte(uint8(player.RangedPoints()))
	return
}
