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
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//FriendList Builds a packet with the players friend entityList information in it.
func FriendList(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(71)
	p.AddUint8(byte(player.FriendList.Size()))
	for s := range player.FriendList.EntrySet() {
		hash := strutil.Base37.Encode(s)
		p.AddUint64(hash)
		
		p1, ok := Players.FindHash(hash)
		if p1 != nil && ok && (p1.FriendList.Contains(player.Username()) || !p1.FriendBlocked()) {
			p.AddUint8(0xFF)
		} else {
			p.AddUint8(0)
		}
	}
	return p
}

//PrivateMessage Builds a packet with a private message from hash with content msg.
func PrivateMessage(hash uint64, msg string) (p *net.Packet) {
	p = net.NewEmptyPacket(120)
	p.AddUint64(hash)
	p.AddUint32(rand.Rng.Uint32()) // unique Message ID to prevent duplicate messages somehow arriving or something idk
	// for _, c := range strutil.ChatFilter.Pack(msg) {
	for _, c := range []byte(strutil.ChatFilter.Format(msg)) {
		p.AddUint8(c)
	}
	return p
}

//func CreateProjectile(owner *Player, target entity.MobileEntity, projectileID int) (p *net.Packet) {
//	p := net.NewEmptyPacket(234)
//	return p
//}

//IgnoreList Builds a packet with the players ignore entityList information in it.
func IgnoreList(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(109)
	p.AddUint8(byte(len(player.IgnoreList)))
	for _, hash := range player.IgnoreList {
		p.AddUint64(hash)
	}
	return p
}

//FriendUpdate Builds a packet with an online status update for the player with the specified hash
func FriendUpdate(hash uint64, online bool) (p *net.Packet) {
	p = net.NewEmptyPacket(149)
	p.AddUint64(hash)
	if online {
		p.AddUint8(0xFF)
	} else {
		p.AddUint8(0)
	}
	return
}

func NpcEvents(player *Player) (p *net.Packet) {
	updateSize := 0
	
	p = net.NewEmptyPacket(104)
	p.AddUint16(uint16(updateSize))
	list, ok := player.Var("npcSplatQ")
	if ok {
		list, ok := list.([]HitSplat)
		if !ok {
			return nil
		}
		newList := make([]HitSplat, 0, len(list))
		for _, splat := range list {
			if splat.Owner.IsNpc() {
				p.AddUint16(uint16(splat.Owner.ServerIndex()))
				p.AddUint8(2)
				p.AddUint8(uint8(splat.Damage))
				p.AddUint8(uint8(splat.Owner.Skills().Current(entity.StatHits)))
				p.AddUint8(uint8(splat.Owner.Skills().Maximum(entity.StatHits)))
				updateSize++
			} else {
				newList = append(newList, splat)
			}
		}
		player.SetVar("npcSplatQ", newList)
	}
	list, ok = player.Var("npcChatQ")
	if ok {
		list, ok := list.([]ChatMessage)
		if !ok {
			return nil
		}
		newList := make([]ChatMessage, 0, len(list))
		for _, msg := range list {
			if msg.Owner.IsNpc() {
				p.AddUint16(uint16(msg.Owner.ServerIndex()))
				p.AddUint8(1)
				p.AddUint16(uint16(msg.Target.ServerIndex()))
				if len(msg.string) > 255 {
					msg.string = msg.string[:255]
				}
				message := strutil.ChatFilter.Format(msg.string)
				// messageRaw := strutil.ChatFilter.Pack(message)
				p.AddUint8(uint8(len(message)))
				for _, c := range message {
					p.AddUint8(byte(c))
				}
				updateSize++
			} else {
				newList = append(newList, msg)
			}
		}
		player.SetVar("npcChatQ", newList)
	}
	p.SetUint16At(1, uint16(updateSize))
	if updateSize <= 0 {
		return nil
	}
	return
}

//ShopClose A net to tell the client to close any open shop interface.
var ShopClose = net.NewEmptyPacket(137)

//ShopOpen Builds a packet to open a shop interface with the data about this shop.
func ShopOpen(shop *Shop) (p *net.Packet) {
	p = net.NewEmptyPacket(101)
	p.AddUint8(uint8(shop.Inventory.Size()))
	p.AddBoolean(shop.BuysUnstocked)
	p.AddUint8(uint8(shop.BasePurchasePercent))
	p.AddUint8(uint8(shop.BaseSalePercent))

	shop.Inventory.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint16(uint16(item.Amount))
		p.AddUint8(uint8(shop.DeltaPercentMod(item)))
		return false
	})
	return p
}

func SleepWord(player *Player) (p *net.Packet) {
	// TODO: Figure this out
	return net.NewEmptyPacket(117)
}

func SleepFatigue(player *Player) (p *net.Packet) {
	return net.NewEmptyPacket(244).AddUint16(uint16(player.VarInt("sleepFatigue", 0)))
}

var SleepClose = net.NewEmptyPacket(84)

var SleepWrong = net.NewEmptyPacket(194)

func NpcMessage(sender *NPC, message string, target *Player) (p *net.Packet) {
	target.QueueNpcChat(sender, target, message)
/*	p = net.NewEmptyPacket(104)
	p.AddUint16(1)
	p.AddUint16(uint16(sender.Index))
	p.AddUint8(1)
	p.AddUint16(uint16(target.Index))
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	// messageRaw := strutil.ChatFilter.Pack(message)
	messageRaw := message
	p.AddUint8(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddUint8(byte(c))
	}
	return*/
	// return net.NewEmptyPacket(6)
	return nil
}

//PrivacySettings Builds a packet containing the players privacy settings for display in the settings menu.
func PrivacySettings(player *Player) (p *net.Packet) {
	return net.NewEmptyPacket(51).AddBoolean(player.ChatBlocked()).AddBoolean(player.FriendBlocked()).AddBoolean(player.TradeBlocked()).AddBoolean(player.DuelBlocked())
}

func OptionMenuOpen(questions ...string) (p *net.Packet) {
	p = net.NewEmptyPacket(245)
	p.AddUint8(uint8(len(questions)))
	for _, question := range questions {
		p.AddUint8(uint8(len(question)))
		p.AddBytes([]byte(question))
	}
	return p
}

var OptionMenuClose = net.NewEmptyPacket(252)

//NPCPositions Builds a packet containing view area NPC position and sprite information
func NPCPositions(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(79)
	changed := 0
	p.AddBitmask(player.LocalNPCs.Size(), 8)
	var removing = NewMobList()
	player.LocalNPCs.RangeNpcs(func(n *NPC) bool {
		changed++
		n.RLock()
		mask := n.SyncMask
		if !player.WithinRange(player.Location, player.VarInt("viewRadius", 16)) || mask&SyncRemoved == SyncRemoved || n.Location.Equals(DeathPoint) || n.VarBool("removed", false) {
			p.AddBitmask(1, 1)
			p.AddBitmask(1, 1)
			p.AddBitmask(3, 2)
			removing.Add(n)
		} else if mask&SyncMoved == SyncMoved {
			p.AddBitmask(1, 1)
			p.AddBitmask(0, 1)
			p.AddBitmask(n.Direction(), 3)
		} else if mask&SyncSprite == SyncSprite {
			p.AddBitmask(1, 1)
			p.AddBitmask(1, 1)
			p.AddBitmask(n.Direction(), 4)
		} else {
			p.AddBitmask(0, 1)
			changed--
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
			if player.VarInt("viewRadius", 16) > 1 {
				player.Dec("viewRadius", 1)
			}
			break
		}
		if player.VarInt("viewRadius", 16) < 16 {
			player.Inc("viewRadius", 1)
		}
		newCount++
		player.LocalNPCs.Add(n)
		p.AddBitmask(n.Index, 12)
		// bitwise trick avoids branching to do a manual addition, and maintains binary compatibility with the original protocol
		p.AddSignedBits(n.X()-player.X(), 5)
		p.AddSignedBits(n.Y()-player.Y(), 5)
		p.AddBitmask(n.Direction(), 4)
		p.AddBitmask(n.ID, 10)
		changed++
	}
	if changed <= 0 {
		return nil
	}
	return
}

func PrayerStatus(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(206)
	for i := 0; i < len(player.Mob.Prayers); i++ {
		p.AddBoolean(player.PrayerActivated(i))
	}
	return p
}

//PlayerPositions Builds a packet containing view area player position and sprite information, including ones own information, and returns it.
// If no players need to be updated, returns nil.
func PlayerPositions(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(191)
	// Note: x coords can be held in 10 bits and y can be held in 12 bits
	//  Presumably, Jagex used 11 and 13 to evenly fill 3 bytes of data?
	p.AddBitmask(player.X(), 11)
	p.AddBitmask(player.Y(), 13)
	p.AddBitmask(player.Direction(), 4)
	p.AddBitmask(player.LocalPlayers.Size(), 8)
	changed := 0
	if player.SyncMask&SyncNeedsPosition != 0 {
		changed++
		player.PostTickables.Add(func() bool {
			player.ResetRegionRemoved()
			player.ResetRegionMoved()
			player.ResetSpriteUpdated()
			return true
		})
	}
	var removing []*Player
	player.LocalPlayers.RangePlayers(func(p1 *Player) bool {
		p1.RLock()
		defer      p1.RUnlock()
		changed++
		mask := p1.SyncMask
		if mask&(SyncRemoved|SyncMoved|SyncSprite) == 0 {
			p.AddBitmask(0, 1)
			changed--
			return false
		}
		p.AddBitmask(1, 1)
		if mask&SyncMoved == SyncMoved {
			p.AddBitmask(0, 1)
			p.AddBitmask(p1.Direction(), 3)
			p1.PostTickables.Add(func(p1 *Player) bool {
				p1.ResetRegionMoved()
				p1.ResetSpriteUpdated()
				return true
			})
			return false
		}
		if p1.LongestDelta(player.Location) >= player.VarInt("viewRadius", 16) || mask&(SyncRemoved|SyncSprite) != 0 {
			p.AddBitmask(1, 1)
			if mask&SyncSprite != 0 {
				p.AddBitmask(p1.Direction(), 4)
			} else {
				p.AddBitmask(3, 2)
				removing = append(removing, p1)
			}
			p1.PostTickables.Add(func(p1 *Player) bool {
				if mask&SyncSprite != 0 {
					p1.ResetSpriteUpdated()
					return true
				}
				p1.ResetRegionRemoved()
				p1.ResetRegionMoved()
				return true
			})
		}
		return false
	})
	for _, p1 := range removing {
		player.LocalPlayers.Remove(p1)
	}
	newPlayerCount := 0
	player.NewPlayers().RangePlayers(func(p1 *Player) bool {
		if player.LocalPlayers.Size() >= 255 {
			// We can only support so many players.  This might even be too much
			return false
		}
		if newPlayerCount >= 25 {
			// Shrink view area when too many new players in one tick
			if player.VarInt("viewRadius", 16) > 1 {
				player.Dec("viewRadius", 1)
			}
			return false
		} else if player.VarInt("viewRadius", 16) < 16 {
			// Grow view area back out after it had been shrunk
			player.Inc("viewRadius", 1)
		}
		newPlayerCount++
		player.LocalPlayers.Add(p1)
		p.AddBitmask(p1.Index, 11)
		// bitwise trick avoids branching to do a manual addition, and maintains binary compatibility with the original protocol
		p.AddSignedBits(p1.X()-player.X(), 5)
		p.AddSignedBits(p1.Y()-player.Y(), 5)
		p.AddBitmask(p1.Direction(), 4)
		if ticket, ok := player.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket() || p1.SyncMask&(SyncRemoved|SyncAppearance)!=0 {
			p.AddBitmask(1, 1)
			player.AppearanceReq = append(player.AppearanceReq, p1)
		} else {
			p.AddBitmask(0, 1)
		}
		return false
	})
//	if changed+newPlayerCount <= 0 {
//		return nil
//	}
	return
}

//PlayerAppearances Builds a packet with the view-area player appearance profiles in it.
func PlayerAppearances(ourPlayer *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(234)
	p.AddUint16(0)
	updateSize := 0
	list, ok := ourPlayer.Var("bubbleQ")
	if ok {
		for _, bubble := range list.([]ItemBubble) {
			p.AddUint16(uint16(bubble.Owner.ServerIndex())) // Index
			p.AddUint8(0) // Update Type
			p.AddUint16(uint16(bubble.Item)) // Item ID
			updateSize++
		}
		ourPlayer.UnsetVar("bubbleQ")
	}
	list, ok = ourPlayer.Var("publicChatQ")
	if ok {
		for _, msg := range list.([]ChatMessage) {
			p.AddUint16(uint16(msg.Owner.ServerIndex())) // Index
			p.AddUint8(1) // Update Type
			// TODO: Is this better or is end of message indicator better
			p.AddUint8(uint8(len(msg.string))) // Count of UTF-8 characters in message
			p.AddBytes([]byte(msg.string)) // UTF-8 encoded message
			updateSize++
		}
		ourPlayer.UnsetVar("publicChatQ")
	}
	list, ok = ourPlayer.Var("hitsplatQ")
	if ok {
		for _, splat := range list.([]HitSplat) {
			p.AddUint16(uint16(splat.Owner.ServerIndex())) // Index
			p.AddUint8(2) // Update Type
			p.AddUint8(uint8(splat.Damage)) // How much damage was done
			p.AddUint8(uint8(splat.Owner.Skills().Current(entity.StatHits))) // Current hitpoints level, for healthbar percentage 
			p.AddUint8(uint8(splat.Owner.Skills().Maximum(entity.StatHits))) // Maximum hitpoints level, for healthbar percentage
			updateSize++
		}
		ourPlayer.UnsetVar("hitsplatQ")
	}
	list, ok = ourPlayer.Var("projectileQ")
	if ok {
		for _, shot := range list.([]Projectile) {
			p.AddUint16(uint16(shot.Owner.ServerIndex())) // Index
			updateType := 3
			if shot.Target.IsPlayer() {
				updateType = 4
			}
			p.AddUint8(uint8(updateType)) // Update Type

			p.AddUint16(uint16(shot.Kind)) // Projectile Type
			p.AddUint16(uint16(shot.Target.ServerIndex())) // Projectile target index
			updateSize++
		}
		ourPlayer.UnsetVar("projectileQ")
	}
	list, ok = ourPlayer.Var("questChatQ")
	if ok {
		for _, msg := range list.([]ChatMessage) {
			p.AddUint16(uint16(msg.Owner.ServerIndex())) // Index
			p.AddUint8(6) // Update Type
			// Format chat messages to match the rules of Jagex chat format
			// Examples: First letters capitalized for every sentence, color-codes are properly identified, etc.
			msg.string = strutil.ChatFilter.Format(msg.string)
			// Too long messages are truncated to 255 bytes
			if len(msg.string) > 0xFF {
				msg.string = msg.string[:0xFF]
			}
			// Deprecated below call; Go defaults string encoding to UTF-8 and I updated the clients to use UTF-8 as well
			// messageRaw := strutil.ChatFilter.Encode(message)
			p.AddUint8(uint8(len(msg.string)))
			p.AddBytes([]byte(msg.string))
			updateSize++
		}
		ourPlayer.UnsetVar("questChatQ")
	}


	var appearanceList []*Player
	if ourPlayer.SyncMask&(SyncRemoved|SyncAppearance) != 0 {
		ourPlayer.PostTickables.Add(func() bool {
			ourPlayer.ResetAppearanceChanged()
			return true
		})
		appearanceList = append(appearanceList, ourPlayer)
	}

	appearanceList = append(appearanceList, ourPlayer.AppearanceReq...)
	ourPlayer.AppearanceReq = ourPlayer.AppearanceReq[:0]
	ourPlayer.LocalPlayers.RangePlayers(func(p1 *Player) bool {
		if ticket, ok := ourPlayer.KnownAppearances[p1.ServerIndex()]; !ok || ticket != p1.AppearanceTicket() {//||
			// p1.SyncMask&(SyncRemoved|SyncAppearance) != 0 {
			appearanceList = append(appearanceList, p1)
		}
		return false
	})
	for _, player := range appearanceList {
		p.AddUint16(uint16(player.Index)) // index
		p.AddUint8(5) // update type
		// This ticket is to track changes to the players around us
		// Everytime this ticket changes, we must send this block out regionally,
		// containing data that identifies all of the owning players characteristics
		p.AddUint16(uint16(player.AppearanceTicket())) // appearance uuid
		p.AddUint64(player.UsernameHash()) // base37 encoded username
		ourPlayer.KnownAppearances[player.Index] = player.AppearanceTicket()
		sprites := player.Equips()
		p.AddUint8(uint8(len(sprites))) // length of equipped item sprites  If length less than 12 any ones after length will get set to 0
		for i := 0; i < len(sprites); i++ {
			p.AddUint8(uint8(sprites[i]))
		}

		// The below colors will set the human character animation colors used for this player,
		// it will not apply to any equipment on top of said human character
		// They are simple array indexes corresponding to arrays built in the client
		p.AddUint8(uint8(player.Appearance.HeadColor))
		p.AddUint8(uint8(player.Appearance.BodyColor))
		p.AddUint8(uint8(player.Appearance.LegsColor))
		p.AddUint8(uint8(player.Appearance.SkinColor))

		// Combat level is the publically shown level of this player; it gives a general
		// idea of how good this player is in combat, it's calculated from the levels of the 
		// first 6 skill types
		p.AddUint8(uint8(player.Skills().CombatLevel()))
		p.AddBoolean(player.Skulled())
		updateSize++
	}
	if updateSize <= 0 {
		return nil
	}
	p.SetUint16At(1, uint16(updateSize))
	return
}

//ClearDistantChunks iterates through a players transient `distantChunks` attribute and sends them to the client to signal
// a removal of all stationary entities within an 8x8 chunk of tiles surrounding the cached location.
func ClearDistantChunks(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(211)
	ichunks, ok := player.Var("distantChunks")
	if !ok {
		return nil
	}
	chunks := ichunks.([]Location)
	if len(chunks) <= 0 {
		return nil
	}
	for _, chunk := range chunks {
		p.AddUint16(uint16(chunk.X() - player.X()))
		p.AddUint16(uint16(chunk.Y() - player.Y()))
	}
	player.UnsetVar("distantChunks")
	return
}

//ObjectLocations Builds a packet with the view-area object positions in it, relative to the player.
// If no new objects are available and no existing local objects are removed from area, returns nil.
func ObjectLocations(player *Player) (p *net.Packet) {
	changed := 0
	p = net.NewEmptyPacket(48)
	var removing = []*Object{}
	for _, o := range player.LocalObjects.set {
		if o, ok := o.(*Object); ok {
			if o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, player.VarInt("viewRadius", 16)+5) || GetObject(o.X(), o.Y()) != o {
				if !player.WithinRange(o.Location, 144) {
					// suddenly this local entity is now miles away which isn't very local
					if chunks, ok := player.Var("distantChunks"); ok {
						player.SetVar("distantChunks", append(chunks.([]Location), o.Location.Clone()))
					} else {
						player.SetVar("distantChunks", []Location{o.Location.Clone()})
					}
				} else {
					p.AddUint16(60000)
					p.AddUint8(byte(o.X() - player.X()))
					p.AddUint8(byte(o.Y() - player.Y()))
					changed++
				}
				removing = append(removing, o)
			}
		}
	}
	for _, o := range removing {
		player.LocalObjects.Remove(o)
	}
	for _, o := range player.NewObjects() {
		if o.Boundary {
			continue
		}
		p.AddUint16(uint16(o.ID))
		p.AddUint8(byte(o.X() - player.X()))
		p.AddUint8(byte(o.Y() - player.Y()))
		player.LocalObjects.Add(o)
		changed++
	}
	if changed == 0 {
		return nil
	}
	return
}

//BoundaryLocations Builds a packet with the view-area boundary positions in it, relative to the player.
// If no new objects are available and no existing local boundarys are removed from area, returns nil.
func BoundaryLocations(player *Player) (p *net.Packet) {
	changed := 0
	p = net.NewEmptyPacket(91)
	var removing = []*Object{}
	for _, o := range player.LocalObjects.set {
		if o, ok := o.(*Object); ok {
			if !o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, player.VarInt("viewRadius", 16)+5) || GetObject(o.X(), o.Y()) != o {
				if !player.WithinRange(o.Location, 144) {
					if chunks, ok := player.Var("distantChunks"); ok {
						player.SetVar("distantChunks", append(chunks.([]Location), o.Location.Clone()))
					} else {
						player.SetVar("distantChunks", []Location{o.Location.Clone()})
					}
				} else {
					// network protocol does not support actual removal of previously existing boundary objects
					// so instead, we replace with an invisible boundary that does not block.
					// This is seen in canonical game most notably when slicing a spider web with a weapon
					p.AddUint16(16)
					p.AddUint8(uint8(o.X() - player.X()))
					p.AddUint8(uint8(o.Y() - player.Y()))
					p.AddUint8(o.Direction)
					changed++
				}
				removing = append(removing, o)
			}
		}
	}
	for _, o := range removing {
		player.LocalObjects.Remove(o)
	}
	for _, o := range player.NewObjects() {
		if !o.Boundary {
			continue
		}
		p.AddUint16(uint16(o.ID))
		p.AddUint8(byte(o.X() - player.X()))
		p.AddUint8(byte(o.Y() - player.Y()))
		p.AddUint8(o.Direction)
		player.LocalObjects.Add(o)
		changed++
	}
	if changed == 0 {
		return nil
	}
	return
}

//ItemLocations Builds a packet with the view-area item positions in it, relative to the player.
// If no new items are available and no existing items are removed from area, returns nil.
func ItemLocations(player *Player) (p *net.Packet) {
	changed := 0
	p = net.NewEmptyPacket(99)
	var removing = []*GroundItem{}
	for _, i := range player.LocalItems.set {
		if i, ok := i.(*GroundItem); ok {
			x, y := i.X(), i.Y()
			if !player.WithinRange(i.Location, player.VarInt("viewRadius", 16)) {
				if !player.WithinRange(i.Location, 144) {
					if chunks, ok := player.Var("distantChunks"); ok {
						player.SetVar("distantChunks", append(chunks.([]Location), i.Location.Clone()))
					} else {
						player.SetVar("distantChunks", []Location{i.Location.Clone()})
					}
				} else {
					// If first byte is 0xFF, all ground items at this location get cleared
					p.AddUint8(255)
					p.AddUint8(byte(x - player.X()))
					p.AddUint8(byte(y - player.Y()))
					changed++
				}
				removing = append(removing, i)
			} else if !i.VisibleTo(player) || !Region(x, y).Items.Contains(i) {
				p.AddUint16(uint16(i.ID | 0x8000)) // turn remove by ID bit on
				p.AddUint8(byte(x - player.X()))
				p.AddUint8(byte(y - player.Y()))
				removing = append(removing, i)
				changed++
			}
		}
	}
	for _, i := range removing {
		player.LocalItems.Remove(i)
	}
	for _, i := range player.NewItems() {
		p.AddUint16(uint16(i.ID))
		p.AddUint8(byte(i.X() - player.X()))
		p.AddUint8(byte(i.Y() - player.Y()))
		player.LocalItems.Add(i)
		changed++
	}
	if changed == 0 {
		return nil
	}
	return
}

//OpenChangeAppearance The appearance changing window.
var OpenChangeAppearance = net.NewEmptyPacket(59)

//InventoryItems Builds a packet containing the players inventory items.
func InventoryItems(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(53)
	p.AddUint8(uint8(player.Inventory.Size()))
	player.Inventory.Range(func(item *Item) bool {
		if item.Worn {
			// turn equipped bit on
			p.AddUint16(uint16(item.ID | 0x8000))
		} else {
			p.AddUint16(uint16(item.ID))
		}
		if definitions.Items[item.ID].Stackable {
			p.AddSmart08_32(item.Amount)
		}
		return true
	})
	return
}

//FightMode Builds a packet with the players fight mode information in it.
func FightMode(player *Player) (p *net.Packet) {
	// TODO: add to 204
	p = net.NewEmptyPacket(132)
	p.AddUint8(byte(player.FightMode()))
	return p
}

//Fatigue Builds a packet with the players fatigue percentage in it.
func Fatigue(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(114)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the game to match RSC
	p.AddUint16(uint16(player.Fatigue() / 100))
	return p
}

//ClientSettings Builds a packet containing the players client settings, e.g camera mode, mouse mode, sound fx...
func ClientSettings(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(240)
	// TODO: Right IDs?
	p.AddBoolean(player.GetClientSetting(0))
	p.AddBoolean(player.GetClientSetting(2))
	p.AddBoolean(player.GetClientSetting(3))

	//	p.AddUint8(0) // Camera auto/manual?
	//	p.AddUint8(0) // Mouse buttons 1 or 2?
	//	p.AddUint8(1) // Sound effects on/off?
	return
}

//PlayerStats Builds a packet containing all the player's stat information and returns it.
func PlayerStats(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(156)
	for i := 0; i < 18; i++ {
		p.AddUint8(uint8(player.Skills().Current(i)))
	}

	for i := 0; i < 18; i++ {
		p.AddUint8(uint8(player.Skills().Maximum(i)))
	}

	for i := 0; i < 18; i++ {
		p.AddUint32(uint32(player.Skills().Experience(i)))
	}
	return p
}

//PlayerStat Builds a packet containing player's stat information for skill at idx and returns it.
func PlayerExperience(player *Player, idx int) (p *net.Packet) {
	p = net.NewEmptyPacket(33)
	p.AddUint8(byte(idx))
	p.AddUint32(uint32(player.Skills().Experience(idx)) )
	return p
}

func PlayerCombatPoints(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(242)
	p.AddUint32(uint32(player.Attributes.VarInt("combatPoints", 0)))
	return p
}

//PlayerStat Builds a packet containing player's stat information for skill at idx and returns it.
func PlayerStat(player *Player, idx int) (p *net.Packet) {
	p = net.NewEmptyPacket(159)
	p.AddUint8(byte(idx))
	p.AddUint8(byte(player.Skills().Current(idx)))
	p.AddUint8(byte(player.Skills().Maximum(idx)))
	p.AddUint32(uint32(player.Skills().Experience(idx)))
	return p
}

//EquipmentStats Builds a packet with the players equipment statistics in it.
func EquipmentStats(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(153)
	p.AddUint8(uint8(player.ArmourPoints()))
	p.AddUint8(uint8(player.AimPoints()))
	p.AddUint8(uint8(player.PowerPoints()))
	p.AddUint8(uint8(player.MagicPoints()))
	p.AddUint8(uint8(player.PrayerPoints()))
	p.AddUint8(uint8(player.RangedPoints()))
	p.AddUint8(uint8(player.VarInt("questPoints", 0xFF)))
	return
}

var BankClose = net.NewEmptyPacket(203)

func BankOpen(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(42)
	p.AddUint8(uint8(player.bank.Size()))
	p.AddUint8(uint8(player.bank.Capacity))
	player.bank.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddSmart08_32(item.Amount)
		return false
	})
	return p
}

func BankUpdateItem(index, id, amount int) (p *net.Packet) {
	p = net.NewEmptyPacket(249)
	p.AddUint8(uint8(index))
	p.AddUint16(uint16(id))
	p.AddSmart08_32(amount)
	return p
}

//DuelOpen Builds a packet to open a duel negotiation window
func DuelOpen(targetIndex int) (p *net.Packet) {
	return net.NewEmptyPacket(176).AddUint16(uint16(targetIndex))
}

//DuelUpdate Builds a packet to update a duel offer
func DuelUpdate(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(6)
	p.AddUint8(uint8(player.DuelOffer.Size()))
	player.DuelOffer.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint32(uint32(item.Amount))
		return true
	})
	return
}

//DuelTargetAccept Builds a packet to change duel targets accepted status
func DuelTargetAccept(accepted bool) (p *net.Packet) {
	return net.NewEmptyPacket(253).AddBoolean(accepted)
}

//DuelOptions Builds a packet to update duel fight options
func DuelOptions(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(30)
	p.AddBoolean(!player.VarBool("duelCanRetreat", true))
	p.AddBoolean(!player.VarBool("duelCanMagic", true))
	p.AddBoolean(!player.VarBool("duelCanPrayer", true))
	p.AddBoolean(!player.VarBool("duelCanEquip", true))
	return p
}

//DuelConfirmationOpen Builds a packet to open the duel confirmation page
func DuelConfirmationOpen(player, other *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(172)

	p.AddUint64(other.UsernameHash())
	
	p.AddUint8(uint8(other.DuelOffer.Size()))
	other.DuelOffer.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint32(uint32(item.Amount))
		return true
	})

	p.AddUint8(uint8(player.DuelOffer.Size()))
	player.DuelOffer.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint32(uint32(item.Amount))
		return true
	})

	p.AddBoolean(!player.VarBool("duelCanRetreat", true))
	p.AddBoolean(!player.VarBool("duelCanMagic", true))
	p.AddBoolean(!player.VarBool("duelCanPrayer", true))
	p.AddBoolean(!player.VarBool("duelCanEquip", true))
	return
}

var DuelClose = net.NewEmptyPacket(225)

//TradeClose Closes a trade window
var TradeClose = net.NewEmptyPacket(128)

//TradeOpen Builds a packet to open a trade window
func TradeOpen(targetIndex int) (p *net.Packet) {
	return net.NewEmptyPacket(92).AddUint16(uint16(targetIndex))
}

//TradeUpdate Builds a packet to update a trade offer
func TradeUpdate(player *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(97)
	p.AddUint8(uint8(player.TradeOffer.Size()))
	player.TradeOffer.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint32(uint32(item.Amount))
		return true
	})
	return
}

//TradeTargetAccept Builds a packet to change trade targets accepted status
func TradeTargetAccept(accepted bool) (p *net.Packet) {
	return net.NewEmptyPacket(162).AddBoolean(accepted)
}

//TradeAccept Builds a packet to change trade targets accepted status
func TradeAccept(accepted bool) (p *net.Packet) {
	return net.NewEmptyPacket(15).AddBoolean(accepted)
}

//TradeConfirmationOpen Builds a packet to open the trade confirmation page
func TradeConfirmationOpen(player, other *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(20)

	p.AddUint64(other.UsernameHash())
	p.AddUint8(uint8(other.TradeOffer.Size()))
	other.TradeOffer.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint32(uint32(item.Amount))
		return true
	})

	p.AddUint8(uint8(player.TradeOffer.Size()))
	player.TradeOffer.Range(func(item *Item) bool {
		p.AddUint16(uint16(item.ID))
		p.AddUint32(uint32(item.Amount))
		return true
	})

	return p
}

//Logout Resets client to login welcome screen
var Logout = net.NewEmptyPacket(4)

//WelcomeMessage Welcome to the game on login
var WelcomeMessage = ServerMessage("Welcome to RuneScape")

//Death The 'Oh dear...You are dead' fade-to-black graphic effect when you die.
var Death = net.NewEmptyPacket(83)

//ResponsePong Response to a RSC protocol ping net
var ResponsePong = net.NewEmptyPacket(9)

//CannotLogout Message that you can not logout right now.
var CannotLogout = net.NewEmptyPacket(183)

//DefaultActionMessage This is a message to inform the player that the action they were trying to perform didn't do anything.
var DefaultActionMessage = ServerMessage("Nothing interesting happens.")

//ServerMessage Builds a packet containing a game message to display in the chat box.
func ServerMessage(msg string) (p *net.Packet) {
	p = net.NewEmptyPacket(131)
	p.AddBytes([]byte(msg))
	return
}

//TeleBubble Builds a packet to draw a teleport bubble at the specified offsets.
func TeleBubble(offsetX, offsetY int) (p *net.Packet) {
	p = net.NewEmptyPacket(36)
	p.AddUint8(0) // type, 0 is mobs, 1 is stationary entities, e.g telegrab
	p.AddUint8(uint8(offsetX))
	p.AddUint8(uint8(offsetY))
	return
}

//SystemUpdate A packet with the time until servers next system update, measured in server ticks (640ms intervals)
func SystemUpdate(t int64) (p *net.Packet) {
	p = net.NewEmptyPacket(52)
	p.AddUint16(uint16(t / 640))
	return p
}

func Sound(name string) (p *net.Packet) {
	return net.NewEmptyPacket(204).AddBytes([]byte(name))
}

//LoginBox Builds a packet to create a welcome box on the client with the inactiveDays since login, and lastIP connected from.
func LoginBox(inactiveDays int, lastIP string) (p *net.Packet) {
	p = net.NewEmptyPacket(182)
	p.AddUint32(uint32(strutil.IPToInteger(lastIP))) // IP
	p.AddUint16(uint16(inactiveDays))                // Last logged in
	// TODO: Recoverys
	p.AddUint8(201) // recovery questions set days, 200 = unset, 201 = set
	// TODO: Message center
	p.AddUint16(0) // Unread messages, number minus one, 0 does not render anything
	return p
}

//BigInformationBox Builds a packet to trigger the opening of a large black text window with msg as its contents
func BigInformationBox(msg string) (p *net.Packet) {
	return net.NewEmptyPacket(222).AddBytes([]byte(msg))
}

//InformationBox Builds a packet to trigger the opening of a small black text window with msg as its contents
func InformationBox(msg string) (p *net.Packet) {
	return net.NewEmptyPacket(89).AddBytes([]byte(msg))
}

//HandshakeResponse Builds a bare net with the login response code.
func HandshakeResponse(v int) (p *net.Packet) {
	return &net.Packet{FrameBuffer: []byte{ byte(v) }}
}

//PlaneInfo Builds a packet to update information about the client environment, e.g height, player index...
func PlaneInfo(player *Player) (p *net.Packet) {
	playerInfo := net.NewEmptyPacket(25)
	playerInfo.AddUint16(uint16(player.Index))
	playerInfo.AddUint16(2304) // alleged width, tiles per sector also...
	playerInfo.AddUint16(1776) // alleged height

	playerInfo.AddUint16(uint16(player.Plane())) // plane

	playerInfo.AddUint16(944) // REAL plane height
	return playerInfo
}
