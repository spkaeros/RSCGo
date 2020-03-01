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
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//FriendList Builds a net with the players friend entityList information in it.
func FriendList(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(71)
	p.AddByte(byte(player.FriendList.size()))
	player.FriendList.ForEach(func(s string, b bool) bool {
		hash := strutil.Base37.Encode(s)
		p.AddLong(hash)

		p1, ok := Players.FromUserHash(hash)
		if p1 != nil && ok && (p1.FriendList.contains(player.Username()) || !p1.FriendBlocked()) {
			p.AddByte(0xFF)
		} else {
			p.AddByte(0)
		}
		return false
	})
	return p
}

//PrivateMessage Builds a net with a private message from hash with content msg.
func PrivateMessage(hash uint64, msg string) (p *net.Packet) {
	p = net.NewOutgoingPacket(120)
	p.AddLong(hash)
	p.AddInt(rand.Uint32()) // unique Message ID to prevent duplicate messages somehow arriving or something idk
	// for _, c := range strutil.ChatFilter.Pack(msg) {
	for _, c := range msg {
		p.AddByte(byte(c))
	}
	return p
}

func CreateProjectile(owner *Player, target entity.MobileEntity, projectileID int) *net.Packet {
	p := net.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(owner.Index))
	eventID := 3
	if target.IsPlayer() {
		eventID = 4
	}
	p.AddByte(uint8(eventID))

	p.AddShort(uint16(projectileID))
	p.AddShort(uint16(target.ServerIndex()))

	return p
}

//IgnoreList Builds a net with the players ignore entityList information in it.
func IgnoreList(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(109)
	p.AddByte(byte(len(player.IgnoreList)))
	for _, hash := range player.IgnoreList {
		p.AddLong(hash)
	}
	return p
}

//FriendUpdate Builds a net with an online status update for the player with the specified hash
func FriendUpdate(hash uint64, online bool) (p *net.Packet) {
	p = net.NewOutgoingPacket(149)
	p.AddLong(hash)
	if online {
		p.AddByte(0xFF)
	} else {
		p.AddByte(0)
	}
	return
}

//PlayerChat Builds a net containing a view-area chat message from the player with the index sender and returns it.
func PlayerChat(sender int, msg string) *net.Packet {
	p := net.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender))
	p.AddByte(1)
	p.AddByte(uint8(len(msg)))
	p.AddBytes([]byte(msg))
	return p
}

//PlayerDamage Builds a net containing a view-area damage display for this player
func PlayerDamage(victim *Player, damage int) *net.Packet {
	p := net.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(victim.Index))
	p.AddByte(2)
	p.AddByte(uint8(damage))
	p.AddByte(uint8(victim.Skills().Current(entity.StatHits)))
	p.AddByte(uint8(victim.Skills().Maximum(entity.StatHits)))
	return p
}

//PlayerItemBubble Builds a net containing a view-area item action bubble display for this player
func PlayerItemBubble(player *Player, id int) *net.Packet {
	p := net.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(player.Index))
	p.AddByte(0)
	p.AddShort(uint16(id))
	return p
}

//NpcDamage Builds a net containing a view-area damage display for this NPC
func NpcDamage(victim *NPC, damage int) *net.Packet {
	p := net.NewOutgoingPacket(104)
	p.AddShort(1)
	p.AddShort(uint16(victim.Index))
	p.AddByte(2)
	p.AddByte(uint8(damage))
	p.AddByte(uint8(victim.Skills().Current(entity.StatHits)))
	p.AddByte(uint8(victim.Skills().Maximum(entity.StatHits)))
	return p
}

//ShopClose A net to tell the client to close any open shop interface.
var ShopClose = net.NewOutgoingPacket(137)

//ShopOpen Builds a net to open a shop interface with the data about this shop.
func ShopOpen(shop *Shop) *net.Packet {
	p := net.NewOutgoingPacket(101)
	p.AddByte(uint8(shop.Inventory.Size()))
	p.AddBool(shop.BuysUnstocked)
	p.AddByte(uint8(shop.BasePurchasePercent))
	p.AddByte(uint8(shop.BaseSalePercent))

	shop.Inventory.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddShort(uint16(item.Amount))
		p.AddByte(uint8(shop.DeltaPercentMod(item)))
		return false
	})
	return p
}

func SleepWord(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(117)
	// TODO: Figure this out
	return p
}

func SleepFatigue(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(244)
	p.AddShort(uint16(player.TransAttrs.VarInt("sleepFatigue", 0)))
	return p
}

var SleepClose = net.NewOutgoingPacket(84)

var SleepWrong = net.NewOutgoingPacket(194)

func NpcMessage(sender *NPC, message string, target *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(104)
	p.AddShort(1)
	p.AddShort(uint16(sender.Index))
	p.AddByte(1)
	p.AddShort(uint16(target.Index))
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	// messageRaw := strutil.ChatFilter.Pack(message)
	messageRaw := message
	p.AddByte(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddByte(byte(c))
	}
	return
}

func PlayerMessage(sender *Player, message string) (p *net.Packet) {
	p = net.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender.Index))
	p.AddByte(6)
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	// messageRaw := strutil.ChatFilter.Pack(message)
	messageRaw := message
	p.AddByte(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddByte(byte(c))
	}
	return
}

//PrivacySettings Builds a net containing the players privacy settings for display in the settings menu.
func PrivacySettings(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(51)
	p.AddBool(player.ChatBlocked())
	p.AddBool(player.FriendBlocked())
	p.AddBool(player.TradeBlocked())
	p.AddBool(player.DuelBlocked())
	return p
}

func OptionMenuOpen(questions ...string) *net.Packet {
	p := net.NewOutgoingPacket(245)
	p.AddByte(uint8(len(questions)))
	for _, question := range questions {
		p.AddByte(uint8(len(question)))
		p.AddBytes([]byte(question))
	}
	return p
}

var OptionMenuClose = net.NewOutgoingPacket(252)

//NPCPositions Builds a net containing view area NPC position and sprite information
func NPCPositions(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(79)
	counter := 0
	p.AddBits(player.LocalNPCs.Size(), 8)
	var removing = NewMobList()
	player.LocalNPCs.RangeNpcs(func(n *NPC) bool {
		counter++
		n.RLock()
		if !player.WithinRange(player.Location, player.TransAttrs.VarInt("viewRadius", 16)) || n.SyncMask&SyncRemoved == SyncRemoved || n.Location.Equals(DeathPoint) {
			p.AddBits(1, 1)
			p.AddBits(1, 1)
			p.AddBits(3, 2)
			removing.Add(n)
		} else if n.SyncMask&SyncMoved == SyncMoved {
			p.AddBits(1, 1)
			p.AddBits(0, 1)
			p.AddBits(n.Direction(), 3)
		} else if n.SyncMask&SyncSprite == SyncSprite {
			p.AddBits(1, 1)
			p.AddBits(1, 1)
			p.AddBits(n.Direction(), 4)
		} else {
			p.AddBits(0, 1)
			counter--
		}
		n.RUnlock()
		return false
	})

	removing.RangeNpcs(func(n *NPC) bool {
		player.LocalNPCs.Remove(n)
		return false
	})

	newCount := 0
	for _, n := range player.NewNPCs() {
		if player.LocalNPCs.Size() >= 255 {
			break
		}
		if newCount >= 25 {
			if player.TransAttrs.VarInt("viewRadius", 16) > 1 {
				player.TransAttrs.DecVar("viewRadius", 1)
			}
			break
		} else {
			if player.TransAttrs.VarInt("viewRadius", 16) < 16 {
				player.TransAttrs.IncVar("viewRadius", 1)
			}
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
	return
}

func PrayerStatus(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(206)
	for i := 0; i < 14; i++ {
		p.AddBool(player.PrayerActivated(i))
	}
	return p
}

//PlayerPositions Builds a net containing view area player position and sprite information, including ones own information, and returns it.
// If no players need to be updated, returns nil.
func PlayerPositions(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(191)
	// Note: x coords can be held in 10 bits and y can be held in 12 bits
	//  Presumably, Jagex used 11 and 13 to evenly fill 3 bytes of data?
	p.AddBits(player.X(), 11)
	p.AddBits(player.Y(), 13)
	p.AddBits(player.Direction(), 4)
	p.AddBits(player.LocalPlayers.Size(), 8)
	counter := 0
	player.RLock()
	//	if player.SyncMask&SyncNeedsPosition != 0 {
	counter++
	//		player.ResetTickables = append(player.ResetTickables, func() {
	//			player.ResetRegionRemoved()
	//			player.ResetRegionMoved()
	//			player.ResetSpriteUpdated()
	//		})
	//	}
	player.RUnlock()
	var removing = NewMobList()
	player.LocalPlayers.RangePlayers(func(p1 *Player) bool {
		p1.RLock()
		counter++
		if p1.LongestDelta(player.Location) >= player.TransAttrs.VarInt("viewRadius", 16) || p1.SyncMask&SyncRemoved == SyncRemoved {
			p.AddBits(1, 1)
			p.AddBits(1, 1)
			p.AddBits(3, 2)
			removing.Add(p1)
			player.AppearanceLock.Lock()
			delete(player.KnownAppearances, p1.Index)
			player.AppearanceLock.Unlock()
			//				p1.ResetTickables = append(p1.ResetTickables, func() {
			//					p1.ResetRegionRemoved()
			//p1.ResetRegionMoved()
			//p1.ResetSpriteUpdated()
			//				})
		} else if p1.SyncMask&SyncMoved == SyncMoved {
			p.AddBits(1, 1)
			p.AddBits(0, 1)
			p.AddBits(p1.Direction(), 3)
			//				p1.ResetTickables = append(p1.ResetTickables, func() {
			//					p1.ResetRegionMoved()
			//p1.ResetSpriteUpdated()
			//				})
		} else if p1.SyncMask&SyncSprite == SyncSprite {
			p.AddBits(1, 1)
			p.AddBits(1, 1)
			p.AddBits(p1.Direction(), 4)
			//				p1.ResetTickables = append(p1.ResetTickables, func() {
			//					p1.ResetSpriteUpdated()
			//				})
		} else {
			p.AddBits(0, 1)
			counter--
		}
		p1.RUnlock()
		return false
	})
	removing.Range(func(m entity.MobileEntity) bool {
		player.LocalPlayers.Remove(m)
		return false
	})
	newPlayerCount := 0
	player.NewPlayers()
	for _, p1 := range player.NewPlayers() {
		if len(player.LocalPlayers.mobSet) >= 255 {
			break
		}
		if newPlayerCount >= 25 {
			if player.TransAttrs.VarInt("viewRadius", 16) > 1 {
				player.TransAttrs.DecVar("viewRadius", 1)
			}
			break
		} else {
			if player.TransAttrs.VarInt("viewRadius", 16) < 16 {
				player.TransAttrs.IncVar("viewRadius", 1)
			}
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
		if ticket, ok := player.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket() {
			p.AddBits(0, 1)
		} else {
			p.AddBits(1, 1)
		}
		player.AppearanceLock.RUnlock()
		player.LocalPlayers.Add(p1)
		//		p1.ResetTickables = append(p1.ResetTickables, func() {
		//			p1.ResetRegionMoved()
		//			p1.ResetSpriteUpdated()
		//		})
		counter++
	}
	if counter <= 0 {
		return nil
	}
	return
}

//PlayerAppearances Builds a net with the view-area player appearance profiles in it.
func PlayerAppearances(ourPlayer *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(234)
	var appearanceList []*Player
	ourPlayer.RLock()
	if ourPlayer.SyncMask&SyncAppearance == SyncAppearance {
		ourPlayer.ResetTickables = append(ourPlayer.ResetTickables, func() {
			ourPlayer.ResetAppearanceChanged()
		})
		appearanceList = append(appearanceList, ourPlayer)
	}
	ourPlayer.RUnlock()

	ourPlayer.AppearanceLock.Lock()
	appearanceList = append(appearanceList, ourPlayer.AppearanceReq...)
	ourPlayer.AppearanceReq = ourPlayer.AppearanceReq[:0]
	ourPlayer.AppearanceLock.Unlock()
	ourPlayer.LocalPlayers.Range(func(p1 entity.MobileEntity) bool {
		ourPlayer.AppearanceLock.RLock()
		if ticket, ok := ourPlayer.KnownAppearances[p1.ServerIndex()]; !ok || ticket != p1.(*Player).AppearanceTicket() {
			appearanceList = append(appearanceList, p1.(*Player))
		}
		ourPlayer.AppearanceLock.RUnlock()
		return false
	})
	if len(appearanceList) <= 0 {
		return nil
	}
	p.AddShort(uint16(len(appearanceList))) // Update size
	for _, player := range appearanceList {
		ourPlayer.AppearanceLock.Lock()
		ourPlayer.KnownAppearances[player.Index] = player.AppearanceTicket()
		ourPlayer.AppearanceLock.Unlock()
		p.AddShort(uint16(player.Index))
		p.AddByte(5) // player appearances
		p.AddShort(uint16(player.AppearanceTicket()))
		p.AddLong(player.UsernameHash())
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
		p.AddBool(player.Skulled())
	}
	return
}

//ClearDistantChunks iterates through a players transient `distantChunks` attribute and sends them to the client to signal
// a removal of all stationary entities within an 8x8 chunk of tiles surrounding the cached location.
func ClearDistantChunks(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(211)
	chunks, ok := player.TransAttrs.Var("distantChunks")
	cleaned := 0
	if ok {
		for _, chunk := range chunks.([]Location) {
			p.AddShort(uint16(chunk.X() - player.X()))
			p.AddShort(uint16(chunk.Y() - player.Y()))
			cleaned++
		}
	}
	if cleaned == 0 {
		return nil
	}
	player.TransAttrs.UnsetVar("distantChunks")
	return
}

//ObjectLocations Builds a net with the view-area object positions in it, relative to the player.
// If no new objects are available and no existing local objects are removed from area, returns nil.
func ObjectLocations(player *Player) (p *net.Packet) {
	counter := 0
	p = net.NewOutgoingPacket(48)
	var removing = entityList{}
	for _, o := range player.LocalObjects.set {
		if o, ok := o.(*Object); ok {
			if o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, player.TransAttrs.VarInt("viewRadius", 16)+5) || GetObject(o.X(), o.Y()) != o {
				if !player.WithinRange(o.Location, 128) {
					if chunks, ok := player.TransAttrs.Var("distantChunks"); ok {
						player.TransAttrs.SetVar("distantChunks", append(chunks.([]Location), o.Location.Clone()))
					} else {
						player.TransAttrs.SetVar("distantChunks", []Location{o.Location.Clone()})
					}
					removing.Add(o)
					continue
				}
				p.AddShort(60000)
				p.AddByte(byte(o.X() - player.X()))
				p.AddByte(byte(o.Y() - player.Y()))
				removing.Add(o)
				counter++
			}
		}
	}
	for _, p1 := range removing.set {
		player.LocalObjects.Remove(p1)
	}
	for _, o := range player.NewObjects() {
		if o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		player.LocalObjects.Add(o)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

//BoundaryLocations Builds a net with the view-area boundary positions in it, relative to the player.
// If no new objects are available and no existing local boundarys are removed from area, returns nil.
func BoundaryLocations(player *Player) (p *net.Packet) {
	counter := 0
	p = net.NewOutgoingPacket(91)
	var removing = entityList{}
	for _, o := range player.LocalObjects.set {
		if o, ok := o.(*Object); ok {
			if !o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, player.TransAttrs.VarInt("viewRadius", 16)+5) || GetObject(o.X(), o.Y()) != o {
				if !player.WithinRange(o.Location, 128) {
					if chunks, ok := player.TransAttrs.Var("distantChunks"); ok {
						player.TransAttrs.SetVar("distantChunks", append(chunks.([]Location), o.Location.Clone()))
					} else {
						player.TransAttrs.SetVar("distantChunks", []Location{o.Location.Clone()})
					}
					removing.Add(o)
					continue
				}
				p.AddShort(16)
				xOff := o.X() - player.X()
				yOff := o.Y() - player.Y()
				p.AddByte(uint8(xOff))
				p.AddByte(uint8(yOff))
				p.AddByte(o.Direction)
				removing.Add(o)
				counter++
			}
		}
	}
	for _, p1 := range removing.set {
		player.LocalObjects.Remove(p1)
	}
	for _, o := range player.NewObjects() {
		if !o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		p.AddByte(o.Direction)
		player.LocalObjects.Add(o)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

//ItemLocations Builds a net with the view-area item positions in it, relative to the player.
// If no new items are available and no existing items are removed from area, returns nil.
func ItemLocations(player *Player) (p *net.Packet) {
	counter := 0
	p = net.NewOutgoingPacket(99)
	var removing = entityList{}
	for _, i := range player.LocalItems.set {
		if i, ok := i.(*GroundItem); ok {
			x, y := i.X(), i.Y()
			if !player.WithinRange(i.Location, player.TransAttrs.VarInt("viewRadius", 16)) {
				if !player.WithinRange(i.Location, 128) {
					if chunks, ok := player.TransAttrs.Var("distantChunks"); ok {
						player.TransAttrs.SetVar("distantChunks", append(chunks.([]Location), i.Location.Clone()))
					} else {
						player.TransAttrs.SetVar("distantChunks", []Location{i.Location.Clone()})
					}
					removing.Add(i)
					continue
				}
				p.AddByte(255)
				p.AddByte(byte(x - player.X()))
				p.AddByte(byte(y - player.Y()))
				removing.Add(i)
				counter++
			} else if !i.VisibleTo(player) || !getRegion(x, y).Items.Contains(i) {
				p.AddShort(uint16(i.ID + 0x8000)) // + 32768
				p.AddByte(byte(x - player.X()))
				p.AddByte(byte(y - player.Y()))
				removing.Add(i)
				counter++
			}
		}
	}
	for _, p1 := range removing.set {
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
var OpenChangeAppearance = net.NewOutgoingPacket(59)

//InventoryItems Builds a net containing the players inventory items.
func InventoryItems(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(53)
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

//FightMode Builds a net with the players fight mode information in it.
func FightMode(player *Player) (p *net.Packet) {
	// TODO: 204
	p = net.NewOutgoingPacket(132)
	p.AddByte(byte(player.FightMode()))
	return p
}

//Fatigue Builds a net with the players fatigue percentage in it.
func Fatigue(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(114)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the game to match RSC
	p.AddShort(uint16(player.Fatigue() / 100))
	return p
}

//ClientSettings Builds a net containing the players client settings, e.g camera mode, mouse mode, sound fx...
func ClientSettings(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(240)
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

//PlayerStats Builds a net containing all the player's stat information and returns it.
func PlayerStats(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(156)
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

//PlayerStat Builds a net containing player's stat information for skill at idx and returns it.
func PlayerExperience(player *Player, idx int) *net.Packet {
	p := net.NewOutgoingPacket(33)
	p.AddByte(byte(idx))
	p.AddInt(uint32(player.Skills().Experience(idx)) * 4)
	return p
}

//PlayerStat Builds a net containing player's stat information for skill at idx and returns it.
func PlayerStat(player *Player, idx int) *net.Packet {
	p := net.NewOutgoingPacket(159)
	p.AddByte(byte(idx))
	p.AddByte(byte(player.Skills().Current(idx)))
	p.AddByte(byte(player.Skills().Maximum(idx)))
	p.AddInt(uint32(player.Skills().Experience(idx)) * 4)
	return p
}

//EquipmentStats Builds a net with the players equipment statistics in it.
func EquipmentStats(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(153)
	p.AddByte(uint8(player.ArmourPoints()))
	p.AddByte(uint8(player.AimPoints()))
	p.AddByte(uint8(player.PowerPoints()))
	p.AddByte(uint8(player.MagicPoints()))
	p.AddByte(uint8(player.PrayerPoints()))
	p.AddByte(uint8(player.RangedPoints()))
	return
}

var BankClose = net.NewOutgoingPacket(203)

func BankOpen(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(42)
	p.AddByte(uint8(player.Bank().Size()))
	p.AddByte(uint8(player.Bank().Capacity))
	for _, item := range player.Bank().List {
		p.AddShort(uint16(item.ID))
		p.AddInt2(uint32(item.Amount))
	}
	return p
}

func BankUpdateItem(index, id, amount int) *net.Packet {
	p := net.NewOutgoingPacket(249)
	p.AddByte(uint8(index))
	p.AddShort(uint16(id))
	p.AddInt2(uint32(amount))
	return p
}

//DuelOpen Builds a net to open a duel negotiation window
func DuelOpen(targetIndex int) *net.Packet {
	return net.NewOutgoingPacket(176).AddShort(uint16(targetIndex))
}

//DuelUpdate Builds a net to update a duel offer
func DuelUpdate(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(6)
	p.AddByte(uint8(player.DuelOffer.Size()))
	player.DuelOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})
	return
}

//DuelTargetAccept Builds a net to change duel targets accepted status
func DuelTargetAccept(accepted bool) *net.Packet {
	return net.NewOutgoingPacket(253).AddBool(accepted)
}

//DuelOptions Builds a net to update duel fight options
func DuelOptions(player *Player) *net.Packet {
	p := net.NewOutgoingPacket(30)
	p.AddBool(!player.TransAttrs.VarBool("duelCanRetreat", true))
	p.AddBool(!player.TransAttrs.VarBool("duelCanMagic", true))
	p.AddBool(!player.TransAttrs.VarBool("duelCanPrayer", true))
	p.AddBool(!player.TransAttrs.VarBool("duelCanEquip", true))
	return p
}

//DuelConfirmationOpen Builds a net to open the duel confirmation page
func DuelConfirmationOpen(player, other *Player) *net.Packet {
	p := net.NewOutgoingPacket(172)

	p.AddLong(other.UsernameHash())
	p.AddByte(uint8(other.DuelOffer.Size()))
	other.DuelOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})

	p.AddByte(uint8(player.DuelOffer.Size()))
	player.DuelOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})

	p.AddBool(!player.TransAttrs.VarBool("duelCanRetreat", true))
	p.AddBool(!player.TransAttrs.VarBool("duelCanMagic", true))
	p.AddBool(!player.TransAttrs.VarBool("duelCanPrayer", true))
	p.AddBool(!player.TransAttrs.VarBool("duelCanEquip", true))

	return p
}

var DuelClose = net.NewOutgoingPacket(225)

//TradeClose Closes a trade window
var TradeClose = net.NewOutgoingPacket(128)

//TradeOpen Builds a net to open a trade window
func TradeOpen(targetIndex int) *net.Packet {
	return net.NewOutgoingPacket(92).AddShort(uint16(targetIndex))
}

//TradeUpdate Builds a net to update a trade offer
func TradeUpdate(player *Player) (p *net.Packet) {
	p = net.NewOutgoingPacket(97)
	p.AddByte(uint8(player.TradeOffer.Size()))
	player.TradeOffer.Range(func(item *Item) bool {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
		return true
	})
	return
}

//TradeTargetAccept Builds a net to change trade targets accepted status
func TradeTargetAccept(accepted bool) *net.Packet {
	return net.NewOutgoingPacket(162).AddBool(accepted)
}

//TradeAccept Builds a net to change trade targets accepted status
func TradeAccept(accepted bool) *net.Packet {
	return net.NewOutgoingPacket(15).AddBool(accepted)
}

//TradeConfirmationOpen Builds a net to open the trade confirmation page
func TradeConfirmationOpen(player, other *Player) *net.Packet {
	p := net.NewOutgoingPacket(20)

	p.AddLong(other.UsernameHash())
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
var Logout = net.NewOutgoingPacket(4)

//WelcomeMessage Welcome to the game on login
var WelcomeMessage = ServerMessage("Welcome to RuneScape")

//Death The 'Oh dear...You are dead' fade-to-black graphic effect when you die.
var Death = net.NewOutgoingPacket(83)

//ResponsePong Response to a RSC protocol ping net
var ResponsePong = net.NewOutgoingPacket(9)

//CannotLogout Message that you can not logout right now.
var CannotLogout = net.NewOutgoingPacket(183)

//DefaultActionMessage This is a message to inform the player that the action they were trying to perform didn't do anything.
var DefaultActionMessage = ServerMessage("Nothing interesting happens.")

//ServerMessage Builds a net containing a game message to display in the chat box.
func ServerMessage(msg string) (p *net.Packet) {
	p = net.NewOutgoingPacket(131)
	p.AddBytes([]byte(msg))
	return
}

//TeleBubble Builds a net to draw a teleport bubble at the specified offsets.
func TeleBubble(offsetX, offsetY int) (p *net.Packet) {
	p = net.NewOutgoingPacket(36)
	p.AddByte(0) // type, 0 is mobs, 1 is stationary entities, e.g telegrab
	p.AddByte(uint8(offsetX))
	p.AddByte(uint8(offsetY))
	return
}

func SystemUpdate(t int) *net.Packet {
	p := net.NewOutgoingPacket(52)
	p.AddShort(uint16((t * 50) / 32))
	return p
}

func Sound(name string) *net.Packet {
	return net.NewOutgoingPacket(204).AddBytes([]byte(name))
}

//LoginBox Builds a net to create a welcome box on the client with the inactiveDays since login, and lastIP connected from.
func LoginBox(inactiveDays int, lastIP string) (p *net.Packet) {
	p = net.NewOutgoingPacket(182)
	p.AddInt(uint32(strutil.IPToInteger(lastIP))) // IP
	p.AddShort(uint16(inactiveDays))              // Last logged in
	p.AddByte(0)                                  // recovery questions set days, 200 = unset, 201 = set
	p.AddShort(1)                                 // Unread messages, number minus one, 0 does not render anything
	p.AddBytes([]byte(lastIP))
	return p
}

//BigInformationBox Builds a net to trigger the opening of a large black text window with msg as its contents
func BigInformationBox(msg string) (p *net.Packet) {
	p = net.NewOutgoingPacket(222)
	p.AddBytes([]byte(msg))
	return p
}

//BigInformationBox Builds a net to trigger the opening of a small black text window with msg as its contents
func InformationBox(msg string) (p *net.Packet) {
	p = net.NewOutgoingPacket(89)
	p.AddBytes([]byte(msg))
	return p
}

//HandshakeResponse Builds a bare net with the login response code.
func HandshakeResponse(v int) *net.Packet {
	return net.NewBarePacket([]byte{byte(v)})
}

//PlaneInfo Builds a net to update information about the client environment, e.g height, player index...
func PlaneInfo(player *Player) *net.Packet {
	playerInfo := net.NewOutgoingPacket(25)
	playerInfo.AddShort(uint16(player.Index))
	playerInfo.AddShort(2304) // alleged width, tiles per sector also...
	playerInfo.AddShort(1776) // alleged height

	playerInfo.AddShort(uint16(player.Plane())) // plane

	playerInfo.AddShort(944) // REAL plane height
	return playerInfo
}
