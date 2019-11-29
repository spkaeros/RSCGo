/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//FriendList Builds a packet with the players friend list information in it.
func FriendList(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(71)
	p.AddByte(byte(len(player.FriendList)))
	for hash, online := range player.FriendList {
		p.AddLong(hash)
		status := 0
		if online {
			status = 0xFF
		}
		p.AddByte(byte(status)) // 255 for online, 0 for offline.
	}
	return p
}

//PrivateMessage Builds a packet with a private message from hash with content msg.
func PrivateMessage(hash uint64, msg string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(120)
	p.AddLong(hash)
	p.AddInt(rand.Uint32()) // unique Message ID to prevent duplicate messages somehow arriving or something idk
	for _, c := range strutil.ChatFilter.Pack(msg) {
		p.AddByte(c)
	}
	return p
}

//IgnoreList Builds a packet with the players ignore list information in it.
func IgnoreList(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(109)
	p.AddByte(byte(len(player.IgnoreList)))
	for _, hash := range player.IgnoreList {
		p.AddLong(hash)
	}
	return p
}

//FriendUpdate Builds a packet with an online status update for the player with the specified hash
func FriendUpdate(hash uint64, online bool) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(149)
	p.AddLong(hash)
	if online {
		p.AddByte(0xFF)
	} else {
		p.AddByte(0)
	}
	return
}

//PlayerChat Builds a packet containing a view-area chat message from the player with the index sender and returns it.
func PlayerChat(sender int, msg string) *packet.Packet {
	p := packet.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender))
	p.AddByte(1)
	p.AddByte(uint8(len(msg)))
	p.AddBytes([]byte(msg))
	return p
}

//PlayerDamage Builds a packet containing a view-area damage display for this player
func PlayerDamage(victim *Player, damage int) *packet.Packet {
	p := packet.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(victim.Index))
	p.AddByte(2)
	p.AddByte(uint8(damage))
	p.AddByte(uint8(victim.Skills().Current(StatHits)))
	p.AddByte(uint8(victim.Skills().Maximum(StatHits)))
	return p
}

//NpcDamage Builds a packet containing a view-area damage display for this NPC
func NpcDamage(victim *NPC, damage int) *packet.Packet {
	p := packet.NewOutgoingPacket(104)
	p.AddShort(1)
	p.AddShort(uint16(victim.Index))
	p.AddByte(2)
	p.AddByte(uint8(damage))
	p.AddByte(uint8(victim.Skills().Current(StatHits)))
	p.AddByte(uint8(victim.Skills().Maximum(StatHits)))
	return p
}

func NpcMessage(sender *NPC, message string, target *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(104)
	p.AddShort(1)
	p.AddShort(uint16(sender.Index))
	p.AddByte(1)
	p.AddShort(uint16(target.Index))
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	messageRaw := strutil.ChatFilter.Pack(message)
	p.AddByte(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddByte(c)
	}
	return
}

func PlayerMessage(sender *Player, message string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender.Index))
	p.AddByte(6)
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	messageRaw := strutil.ChatFilter.Pack(message)
	p.AddByte(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddByte(c)
	}
	return
}

//PrivacySettings Builds a packet containing the players privacy settings for display in the settings menu.
func PrivacySettings(player *Player) *packet.Packet {
	p := packet.NewOutgoingPacket(51)
	p.AddBool(player.ChatBlocked())
	p.AddBool(player.FriendBlocked())
	p.AddBool(player.TradeBlocked())
	p.AddBool(player.DuelBlocked())
	return p
}

func OptionMenuOpen(questions ...string) *packet.Packet {
	p := packet.NewOutgoingPacket(245)
	p.AddByte(uint8(len(questions)))
	for _, question := range questions {
		p.AddByte(uint8(len(question)))
		p.AddBytes([]byte(question))
	}
	return p
}

var OptionMenuClose = packet.NewOutgoingPacket(252)

//NPCPositions Builds a packet containing view area NPC position and sprite information
func NPCPositions(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(79)
	counter := 0
	p.AddBits(len(player.LocalNPCs.List), 8)
	var removing = List{}
	for _, n := range player.LocalNPCs.List {
		if n, ok := n.(*NPC); ok {
			if n.LongestDelta(player.Location) > 15 || n.TransAttrs.HasMasks("sync", SyncRemoved) {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				removing.List = append(removing.List, n)
				counter++
			} else if n.TransAttrs.HasMasks("sync", SyncMoved, SyncChanged) {
				p.AddBits(1, 1)
				if n.TransAttrs.HasMasks("sync", SyncMoved) {
					p.AddBits(0, 1)
					p.AddBits(n.Direction(), 3)
				} else {
					p.AddBits(1, 1)
					p.AddBits(n.Direction(), 4)
				}
				counter++
			} else {
				p.AddBits(0, 1)
			}
		}
	}
	for _, n := range removing.List {
		player.LocalNPCs.Remove(n)
	}
	newCount := 0
	for _, n := range player.NewNPCs() {
		if len(player.LocalNPCs.List) >= 255 || newCount >= 25 {
			break
		}
		newCount++
		player.LocalNPCs.Add(n)
		p.AddBits(n.Index, 12)
		offsetX := n.X() - player.X()
		if offsetX < 0 {
			offsetX += 32
		}
		offsetY := n.Y() - player.Y()
		if offsetY < 0 {
			offsetY += 32
		}
		p.AddBits(offsetX, 5)
		p.AddBits(offsetY, 5)
		p.AddBits(n.Direction(), 4)
		p.AddBits(n.ID, 10)
		counter++
	}
	if counter <= 0 {
		return nil
	}
	return
}

//PlayerPositions Builds a packet containing view area player position and sprite information, including ones own information, and returns it.
// If no players need to be updated, returns nil.
func PlayerPositions(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(191)
	// Note: x coords can be held in 10 bits and y can be held in 12 bits
	//  Presumably, Jagex used 11 and 13 to evenly fill 3 bytes of data?
	p.AddBits(player.X(), 11)
	p.AddBits(player.Y(), 13)
	p.AddBits(player.Direction(), 4)
	p.AddBits(len(player.LocalPlayers.List), 8)
	counter := 0
	if player.TransAttrs.HasMasks("sync", SyncRemoved, SyncMoved, SyncChanged) || !player.Transients().HasMasks("sync", SyncSelf) {
		counter++
	}
	var removing = List{}
	for _, p1 := range player.LocalPlayers.List {
		if p1, ok := p1.(*Player); ok {
			if p1.LongestDelta(player.Location) > 15 || p1.TransAttrs.HasMasks("sync", SyncRemoved) {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				removing.List = append(removing.List, p1)
				player.AppearanceLock.Lock()
				delete(player.KnownAppearances, p1.Index)
				player.AppearanceLock.Unlock()
				counter++
			} else if p1.TransAttrs.HasMasks("sync", SyncMoved, SyncChanged) {
				p.AddBits(1, 1)
				if p1.TransAttrs.HasMasks("sync", SyncMoved) {
					p.AddBits(0, 1)
					p.AddBits(p1.Direction(), 3)
				} else {
					p.AddBits(1, 1)
					p.AddBits(p1.Direction(), 4)
				}
				counter++
			} else {
				p.AddBits(0, 1)
			}
		}
	}
	for _, p1 := range removing.List {
		player.LocalPlayers.Remove(p1)
	}
	newPlayerCount := 0
	for _, p1 := range player.NewPlayers() {
		if len(player.LocalPlayers.List) >= 255 || newPlayerCount >= 25 {
			// No more than 255 players in view at once, no more than 25 new players at once.
			break
		}
		newPlayerCount++
		p.AddBits(p1.Index, 11)
		offsetX := p1.X() - player.X()
		if offsetX < 0 {
			offsetX += 32
		}
		offsetY := p1.Y() - player.Y()
		if offsetY < 0 {
			offsetY += 32
		}
		p.AddBits(offsetX, 5)
		p.AddBits(offsetY, 5)
		p.AddBits(p1.Direction(), 4)
		player.AppearanceLock.RLock()
		if ticket, ok := player.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket {
			p.AddBits(0, 1)
		} else {
			p.AddBits(1, 1)
		}
		player.AppearanceLock.RUnlock()
		player.LocalPlayers.Add(p1)
		counter++
	}
	if counter <= 0 {
		return nil
	}
	return
}

//PlayerAppearances Builds a packet with the view-area player appearance profiles in it.
func PlayerAppearances(ourPlayer *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(234)
	var appearanceList []*Player
	if !ourPlayer.TransAttrs.HasMasks("sync", SyncSelf) {
		appearanceList = append(appearanceList, ourPlayer)
	}
	ourPlayer.AppearanceLock.Lock()
	appearanceList = append(appearanceList, ourPlayer.AppearanceReq...)
	ourPlayer.AppearanceReq = ourPlayer.AppearanceReq[:0]
	ourPlayer.AppearanceLock.Unlock()
	for _, p1 := range ourPlayer.LocalPlayers.List {
		if p1, ok := p1.(*Player); ok {
			ourPlayer.AppearanceLock.RLock()
			if ticket, ok := ourPlayer.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket {
				appearanceList = append(appearanceList, p1)
			}
			ourPlayer.AppearanceLock.RUnlock()
		}
	}
	if len(appearanceList) <= 0 {
		return nil
	}
	p.AddShort(uint16(len(appearanceList))) // Update size
	for _, player := range appearanceList {
		ourPlayer.AppearanceLock.Lock()
		ourPlayer.KnownAppearances[player.Index] = player.AppearanceTicket
		ourPlayer.AppearanceLock.Unlock()
		p.AddShort(uint16(player.Index))
		p.AddByte(5) // player appearances
		p.AddShort(uint16(player.AppearanceTicket))
		p.AddLong(player.UserBase37)
		p.AddByte(12) // length of sprites.  Anything less than 12 will get padded with 0s
		//		p.AddByte(uint8(player.Appearance.Head))
		//		p.AddByte(uint8(player.Appearance.Body))
		//		p.AddByte(uint8(player.Appearance.Legs))
		ourPlayer.AppearanceLock.RLock()
		for i := 0; i < 12; i++ {
			p.AddByte(uint8(player.Equips[i]))
		}
		ourPlayer.AppearanceLock.RUnlock()
		p.AddByte(uint8(player.Appearance.HeadColor))
		p.AddByte(uint8(player.Appearance.BodyColor))
		p.AddByte(uint8(player.Appearance.LegsColor))
		p.AddByte(uint8(player.Appearance.SkinColor))
		p.AddByte(uint8(player.Skills().CombatLevel()))
		p.AddByte(0) // TODO: skulled
	}
	return
}

//ObjectLocations Builds a packet with the view-area object positions in it, relative to the player.
// If no new objects are available and no existing local objects are removed from area, returns nil.
func ObjectLocations(player *Player) (p *packet.Packet) {
	counter := 0
	p = packet.NewOutgoingPacket(48)
	var removing = List{}
	for _, o := range player.LocalObjects.List {
		if o, ok := o.(*Object); ok {
			if o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, 21) || GetObject(o.X(), o.Y()) != o {
				p.AddShort(60000)
				p.AddByte(byte(o.X() - player.X()))
				p.AddByte(byte(o.Y() - player.Y()))
				//				p.AddByte(byte(o.Direction))
				removing.Add(o)
				counter++
			}
		}
	}
	for _, p1 := range removing.List {
		player.LocalObjects.Remove(p1)
	}
	for _, o := range player.NewObjects() {
		if o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		//		p.AddByte(byte(o.Direction))
		player.LocalObjects.Add(o)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

//BoundaryLocations Builds a packet with the view-area boundary positions in it, relative to the player.
// If no new objects are available and no existing local boundarys are removed from area, returns nil.
func BoundaryLocations(player *Player) (p *packet.Packet) {
	counter := 0
	p = packet.NewOutgoingPacket(91)
	var removing = List{}
	for _, o := range player.LocalObjects.List {
		if o, ok := o.(*Object); ok {
			if !o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, 21) {
				//p.AddShort(65535)
				p.AddByte(255)
				p.AddByte(byte(o.X() - player.X()))
				p.AddByte(byte(o.Y() - player.Y()))
				//p.AddByte(byte(o.Direction))
				removing.Add(o)
				counter++
			}
		}
	}
	for _, p1 := range removing.List {
		player.LocalObjects.Remove(p1)
	}
	for _, o := range player.NewObjects() {
		if !o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		p.AddByte(byte(o.Direction))
		player.LocalObjects.Add(o)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

func NpcAppearances(player *Player) *packet.Packet {
	p := packet.NewOutgoingPacket(104)
	toUpdate := 0
	p.AddShort(0)
	for _, npc := range player.LocalNPCs.List {
		if npc, ok := npc.(*NPC); ok {
			if npc.ChatMessage != "" && npc.ChatTarget > -1 {
				message := npc.ChatMessage
				npc.ChatMessage = ""
				toUpdate++
				p.AddShort(uint16(npc.Index))
				p.AddByte(1)
				p.AddShort(uint16(npc.ChatTarget))
				npc.ChatTarget = -1
				if len(message) > 255 {
					message = message[:255]
				}
				message = strutil.ChatFilter.Format(message)
				messageRaw := strutil.ChatFilter.Pack(message)
				p.AddByte(uint8(len(messageRaw)))
				p.AddBytes(messageRaw)
			}
		}
	}
	if toUpdate > 0 {
		p.SetShort(0, uint16(toUpdate))
	} else {
		return nil
	}
	return p
}

//ItemLocations Builds a packet with the view-area item positions in it, relative to the player.
// If no new items are available and no existing items are removed from area, returns nil.
func ItemLocations(player *Player) (p *packet.Packet) {
	counter := 0
	p = packet.NewOutgoingPacket(99)
	var removing = List{}
	for _, i := range player.LocalItems.List {
		if i, ok := i.(*GroundItem); ok {
			x, y := i.X(), i.Y()
			if !player.WithinRange(i.Location, 21) {
				p.AddByte(255)
				p.AddByte(byte(x - player.X()))
				p.AddByte(byte(y - player.Y()))
				removing.Add(i)
				counter++
			} else if !i.VisibleTo(player) || !GetRegion(x, y).Items.Contains(i) {
				p.AddShort(uint16(i.ID + 0x8000)) // + 32768
				p.AddByte(byte(x - player.X()))
				p.AddByte(byte(y - player.Y()))
				removing.Add(i)
				counter++
			}
		}
	}
	for _, p1 := range removing.List {
		player.LocalItems.Remove(p1)
	}
	for _, i := range player.NewItems() {
		p.AddShort(uint16(i.ID))
		p.AddByte(byte(i.X() - player.X()))
		p.AddByte(byte(i.Y() - player.Y()))
		player.LocalItems.Add(i)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

//OpenChangeAppearance The appearance changing window.
var OpenChangeAppearance = packet.NewOutgoingPacket(59)

//InventoryItems Builds a packet containing the players inventory items.
func InventoryItems(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(53)
	p.AddByte(uint8(player.Inventory.Size()))
	player.Inventory.Range(func(item *Item) bool {
		if item.Worn {
			p.AddShort(uint16(item.ID + 0x8000))
		} else {
			p.AddShort(uint16(item.ID))
		}
		if ItemDefs[item.ID].Stackable {
			p.AddInt2(uint32(item.Amount))
		}
		return true
	})
	return
}

//FightMode Builds a packet with the players fight mode information in it.
func FightMode(player *Player) (p *packet.Packet) {
	// TODO: 204
	p = packet.NewOutgoingPacket(132)
	p.AddByte(byte(player.FightMode()))
	return p
}

//Fatigue Builds a packet with the players fatigue percentage in it.
func Fatigue(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(114)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the server to match RSC
	p.AddShort(uint16(player.Fatigue() / 100))
	return p
}

//ClientSettings Builds a packet containing the players client settings, e.g camera mode, mouse mode, sound fx...
func ClientSettings(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(240)
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
func PlayerStats(player *Player) *packet.Packet {
	p := packet.NewOutgoingPacket(156)
	for i := 0; i < 18; i++ {
		p.AddByte(uint8(player.Skills().Current(i)))
	}

	for i := 0; i < 18; i++ {
		p.AddByte(uint8(player.Skills().Maximum(i)))
	}

	for i := 0; i < 18; i++ {
		p.AddInt(uint32(player.Skills().Experience(i) * 4))
	}
	return p
}

//PlayerStat Builds a packet containing player's stat information for skill at idx and returns it.
func PlayerExperience(player *Player, idx int) *packet.Packet {
	p := packet.NewOutgoingPacket(33)
	p.AddByte(byte(idx))
	p.AddInt(uint32(player.Skills().Experience(idx)) * 4)
	return p
}

//PlayerStat Builds a packet containing player's stat information for skill at idx and returns it.
func PlayerStat(player *Player, idx int) *packet.Packet {
	p := packet.NewOutgoingPacket(159)
	p.AddByte(byte(idx))
	p.AddByte(byte(player.Skills().Current(idx)))
	p.AddByte(byte(player.Skills().Maximum(idx)))
	p.AddInt(uint32(player.Skills().Experience(idx)) * 4)
	return p
}

//EquipmentStats Builds a packet with the players equipment statistics in it.
func EquipmentStats(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(153)
	p.AddByte(uint8(player.ArmourPoints()))
	p.AddByte(uint8(player.AimPoints()))
	p.AddByte(uint8(player.PowerPoints()))
	p.AddByte(uint8(player.MagicPoints()))
	p.AddByte(uint8(player.PrayerPoints()))
	p.AddByte(uint8(player.RangedPoints()))
	return
}

var BankClose = packet.NewOutgoingPacket(203)

func BankOpen(player *Player) *packet.Packet {
	p := packet.NewOutgoingPacket(42)
	p.AddByte(uint8(player.Bank.Size()))
	p.AddByte(uint8(player.Bank.Capacity))
	for _, item := range player.Bank.List {
		p.AddShort(uint16(item.ID))
		p.AddInt2(uint32(item.Amount))
	}
	return p
}

func BankUpdateItem(item *Item) *packet.Packet {
	p := packet.NewOutgoingPacket(249)
	p.AddByte(uint8(item.Index))
	p.AddShort(uint16(item.ID))
	p.AddInt2(uint32(item.Amount))
	return p
}

//TradeClose Closes a trade window
var TradeClose = packet.NewOutgoingPacket(128)

//TradeUpdate Builds a packet to update a trade offer
func TradeUpdate(player *Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(97)
	p.AddByte(uint8(player.TradeOffer.Size()))
	player.TradeOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})
	return
}

//TradeOpen Builds a packet to open a trade window
func TradeOpen(player *Player) *packet.Packet {
	return packet.NewOutgoingPacket(92).AddShort(uint16(player.TradeTarget()))
}

//TradeTargetAccept Builds a packet to change trade targets accepted status
func TradeTargetAccept(accepted bool) *packet.Packet {
	if accepted {
		return packet.NewOutgoingPacket(162).AddByte(1)
	}
	return packet.NewOutgoingPacket(162).AddByte(0)
}

//TradeAccept Builds a packet to change trade targets accepted status
func TradeAccept(accepted bool) *packet.Packet {
	if accepted {
		return packet.NewOutgoingPacket(15).AddByte(1)
	}
	return packet.NewOutgoingPacket(15).AddByte(0)
}

//TradeConfirmationOpen Builds a packet to open the trade confirmation page
func TradeConfirmationOpen(player, other *Player) *packet.Packet {
	p := packet.NewOutgoingPacket(20)

	p.AddLong(other.UserBase37)
	p.AddByte(uint8(other.TradeOffer.Size()))
	other.TradeOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})

	p.AddByte(uint8(player.TradeOffer.Size()))
	player.TradeOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})

	return p
}

//Logout Resets client to login welcome screen
var Logout = packet.NewOutgoingPacket(4)

//WelcomeMessage Welcome to the game on login
var WelcomeMessage = ServerMessage("Welcome to RuneScape")

//Death The 'Oh dear...You are dead' fade-to-black graphic effect when you die.
var Death = packet.NewOutgoingPacket(83)

//ResponsePong Response to a RSC protocol ping packet
var ResponsePong = packet.NewOutgoingPacket(9)

//CannotLogout Message that you can not logout right now.
var CannotLogout = packet.NewOutgoingPacket(183)

//DefaultActionMessage This is a message to inform the player that the action they were trying to perform didn't do anything.
var DefaultActionMessage = ServerMessage("Nothing interesting happens.")

//ServerMessage Builds a packet containing a server message to display in the chat box.
func ServerMessage(msg string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(131)
	p.AddBytes([]byte(msg))
	return
}

//TeleBubble Builds a packet to draw a teleport bubble at the specified offsets.
func TeleBubble(offsetX, offsetY int) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(36)
	p.AddByte(0) // type, 0 is mobs, 1 is stationary entities, e.g telegrab
	p.AddByte(uint8(offsetX))
	p.AddByte(uint8(offsetY))
	return
}

func SystemUpdate(t int) *packet.Packet {
	p := packet.NewOutgoingPacket(52)
	p.AddShort(uint16((t * 50) / 32))
	return p
}

func Sound(name string) *packet.Packet {
	return packet.NewOutgoingPacket(204).AddBytes([]byte(name))
}

//LoginBox Builds a packet to create a welcome box on the client with the inactiveDays since login, and lastIP connected from.
func LoginBox(inactiveDays int, lastIP string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(182)
	p.AddInt(uint32(strutil.IPToInteger(lastIP))) // IP
	p.AddShort(uint16(inactiveDays))              // Last logged in
	p.AddByte(0)                                  // recovery questions set days, 200 = unset, 201 = set
	p.AddShort(1)                                 // Unread messages, number minus one, 0 does not render anything
	p.AddBytes([]byte(lastIP))
	return p
}

//BigInformationBox Builds a packet to trigger the opening of a large black text window with msg as its contents
func BigInformationBox(msg string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(222)
	p.AddBytes([]byte(msg))
	return p
}

//LoginResponse Builds a bare packet with the login response code.
func LoginResponse(v int) *packet.Packet {
	return packet.NewBarePacket([]byte{byte(v)})
}

//PlaneInfo Builds a packet to update information about the client environment, e.g height, player index...
func PlaneInfo(player *Player) *packet.Packet {
	playerInfo := packet.NewOutgoingPacket(25)
	playerInfo.AddShort(uint16(player.Index))
	playerInfo.AddShort(2304) // alleged width, tiles per sector also...
	playerInfo.AddShort(1776) // alleged height

	playerInfo.AddShort(uint16(player.Plane())) // plane

	playerInfo.AddShort(944) // REAL plane height
	return playerInfo
}
