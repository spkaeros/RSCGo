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
	"strconv"

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
	events, ok := player.VarChecked(npcEvents).([]interface{})
	p = net.NewEmptyPacket(104)
	eventCount := uint16(len(events))
	if eventCount <= 0 {
		return nil
	}
	p.AddUint16(eventCount)
	if ok && eventCount > 0 {
		var newList []interface{}
		for _, e := range events {
			if splat, ok := e.(*HitSplat); ok {
				if AsNpc(splat.Owner) != nil {
					p.AddUint16(uint16(splat.Owner.ServerIndex()))
					p.AddUint8(2)
					p.AddUint8(uint8(splat.Damage))
					p.AddUint8(uint8(splat.Owner.Skills().Current(entity.StatHits)))
					p.AddUint8(uint8(splat.Owner.Skills().Maximum(entity.StatHits)))
				} else {
					newList = append(newList, splat)
				}
			}
			if msg, ok := e.(ChatMessage); ok {
				if msg.Owner.IsNpc() {
					p.AddUint16(uint16(msg.Owner.ServerIndex()))
					p.AddUint8(1)
					p.AddUint16(uint16(msg.Target.ServerIndex()))
					message := strutil.ChatFilter.Format(msg.string)
					size := len(message)
					if size > 0 && size < 128 {
						p.AddUint8(uint8(size))
					} else if size >= 0 && size < 0x8000 {
						p.AddUint16(uint16(size))
					}
					for _, c := range message {
						p.AddUint8(byte(c))
					}
					p.AddBytes([]byte(message))
				} else {
					newList = append(newList, msg)
				}
			}
		}
		player.SetVar(npcEvents, newList)
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
	// var local []entity.MobileEntity
	var local = player.LocalNPCs.mobSet[:0]
	// var local = make(mobSet, 0, len(player.NearbyNpcs()))
	player.LocalNPCs.RangeNpcs(func(n *NPC) bool {
		local = append(local, n)
		changed++
		n.RLock()
		mask := n.SyncMask
		if !player.Near(n, player.ViewRadius()) || mask&SyncRemoved == SyncRemoved || n.Point().Equals(DeathPoint) || n.VarBool("removed", false) {
			// p.AddBitmask(1, 1)
			// p.AddBitmask(1, 1)
			// p.AddBitmask(3, 2)
			p.AddBitmask(0xF, 4)
			local = local[:len(local)-1]
		} else if mask&SyncMoved == SyncMoved {
			// p.AddBitmask(1, 1)
			// p.AddBitmask(0, 1)
			p.AddBitmask(2, 2)
			p.AddBitmask(n.Direction(), 3)
		} else if mask&SyncSprite == SyncSprite {
			// p.AddBitmask(1, 1)
			// p.AddBitmask(1, 1)
			p.AddBitmask(3, 2)
			p.AddBitmask(n.Direction(), 4)
		} else {
			p.AddBitmask(0, 1)
			changed--
		}
		n.RUnlock()
		return false
	})
	player.LocalNPCs.Set(local)

	newCount := 0
	player.NewNPCs().RangeNpcs(func(n *NPC) bool {
		if player.LocalNPCs.Size() >= 255 {
			return true
		}
		if newCount >= 25 {
			if player.ViewRadius() > 1 {
				player.Dec("viewRadius", 1)
			}
			return true
		}
		if player.ViewRadius() < 16 {
			player.Inc("viewRadius", 1)
		}
		newCount++
		p.AddBitmask(n.ServerIndex(), 12)
		// bitwise trick avoids branching to do a manual addition, and maintains binary compatibility with the original protocol
		p.AddBitmask(n.X() - player.X() & 0x1F, 5)
		p.AddBitmask(n.Y() - player.Y() & 0x1F, 5)
		p.AddBitmask(n.Direction(), 4)
		p.AddBitmask(n.ID, 10)
		changed++
		player.LocalNPCs.Add(n)
		return false
	})
	// if newCount + player.LocalNPCs.Size() <= 0 {
		// return nil
	// }
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
	// if player.SyncMask&SyncNeedsPosition != 0 {
		// changed++
	// }
	// var removing []*Player
	var local mobSet
	player.LocalPlayers.RangePlayers(func(p1 *Player) bool {
		local = append(local, p1)
		changed++
		p1.RLock()
		mask := p1.SyncMask
		if !player.Near(p1, player.ViewRadius()) || mask&SyncRemoved == SyncRemoved {
			// flips on the next 4 bits
			p.AddBitmask(0xF, 4)
			local = local[:len(local)-1]
		} else if mask&SyncMoved == SyncMoved {
			// Sets next 2 bits as 1, 0 consecutively
			p.AddBitmask(2, 2)
			p.AddBitmask(p1.Direction(), 3)
		} else if mask&SyncSprite == SyncSprite {
			p.AddBitmask(3, 2)
			p.AddBitmask(p1.Direction(), 4)
		} else {
			p.AddBitmask(0, 1)
			changed--
		}
		p1.RUnlock()
		return false
	})
	player.LocalPlayers.Set(local)
	newPlayerCount := 0
	player.NewPlayers().RangePlayers(func(p1 *Player) bool {
		if player.LocalPlayers.Size() >= 255 {
			// We can only support so many players.  This might even be too much
			return false
		}
		if newPlayerCount >= 25 {
			// Shrink view area when too many new players in one tick
			if player.ViewRadius() > 1 {
				player.Dec("viewRadius", 1)
			}
			return true
		} else if player.ViewRadius() < 16 {
			// Grow view area back out after it had been shrunk
			player.Inc("viewRadius", 1)
		}
		newPlayerCount++
		player.LocalPlayers.Add(p1)
		p.AddBitmask(p1.ServerIndex(), 11)
		// bitwise trick avoids branching to do a manual addition, and maintains binary compatibility with the original protocol
		p.AddSignedBits(p1.X()-player.X(), 5)
		p.AddSignedBits(p1.Y()-player.Y(), 5)
		p.AddBitmask(p1.Direction(), 4)
		// if ticket, ok := player.KnownAppearances[p1.ServerIndex()]; !ok || ticket != p1.AppearanceTicket() || p1.SyncMask&(SyncRemoved|SyncAppearance) != 0 {
		if ticket, hasPlayerTicket := player.Var(strconv.Itoa(p1.ServerIndex()) + "=" + strconv.Itoa(p1.AppearanceTicket())); !hasPlayerTicket || ticket != p1.AppearanceTicket() || p1.SyncMask&(SyncRemoved|SyncAppearance) != 0 {
			player.enqueue(playerEvents, map[string]int {"index": int(p1.ServerIndex()), "ticket": int(p1.AppearanceTicket())})
			p.AddBitmask(0, 1)
			return false
		} else {
			p.AddBitmask(1, 1)
		}
		return false
	})
	return
}

//PlayerAppearances Builds a packet with the view-area player appearance profiles in it.
func PlayerAppearances(ourPlayer *Player) (p *net.Packet) {
	p = net.NewEmptyPacket(234)
	// updateSize := 0
	list, ok := ourPlayer.VarChecked(playerEvents).([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	p.AddUint16(uint16(len(list)))
	if ok {
		for _, e := range list {
			if bubble, ok := e.(ItemBubble); ok {
				p.AddUint16(uint16(bubble.Owner.ServerIndex())) // Index
				p.AddUint8(0)                                   // Update Type
				p.AddUint16(uint16(bubble.Item))                // Item ID
			}
			if msg, ok := e.(ChatMessage); ok {
				if msg.Target == nil {
					p.AddUint16(uint16(msg.Owner.ServerIndex())) // Index
					p.AddUint8(1)                                // Update Type
					// TODO: Is this better or is end of message indicator better
					size := uint8(len(msg.string))
					if size > 84 {
						size = 84
						msg.string = msg.string[:size]
					}
					p.AddUint8(size)               // Count of UTF-8 characters in message
					p.AddBytes([]byte(msg.string)) // UTF-8 encoded message
				} else {
					p.AddUint16(uint16(msg.Owner.ServerIndex())) // Index
					p.AddUint8(6)                                // Update Type
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
				}
			}
			if splat, ok := e.(*HitSplat); ok {
				p.AddUint16(uint16(splat.Owner.ServerIndex()))                   // Index
				p.AddUint8(2)                                                    // Update Type
				p.AddUint8(uint8(splat.Damage))                                  // How much damage was done
				p.AddUint8(uint8(splat.Owner.Skills().Current(entity.StatHits))) // Current hitpoints level, for healthbar percentage
				p.AddUint8(uint8(splat.Owner.Skills().Maximum(entity.StatHits))) // Maximum hitpoints level, for healthbar percentage
			}
			if shot, ok := e.(Projectile); ok {
				p.AddUint16(uint16(shot.Owner.ServerIndex())) // Index
				if shot.Target.IsNpc() {
					p.AddUint8(3)
				} else if shot.Target.IsPlayer() {
					p.AddUint8(4)
				} else {
					p.Rewind(2)
					continue
				}
				p.AddUint16(uint16(shot.Kind))                 // Projectile Type, this is large bit-length for such small data
				p.AddUint16(uint16(shot.Target.ServerIndex())) // Projectile target index
			}
			if ticket, ok := e.(map[string]int); ok && ticket != nil {
				idx, ok := ticket["index"]
				if !ok {
					// log.Debug("bad ticket: no index; required for any tickets")
					continue
				}
				ticketID, ok := ticket["ticket"]
				if !ok {
					// log.Debug("bad ticket: no ticket ID; required for any tickets")
					continue
				}

				p1, ok := Players.FindIndex(int(idx))
				if p1.AppearanceTicket() != int(ticketID) || !ok {
					// log.Debug("bad ticket ID; player at idx does not match for it")
					continue
				}
				ourPlayer.SetVar("ticket" + strconv.Itoa(p1.ServerIndex()), ticketID)
				p.AddUint16(uint16(p1.ServerIndex())) // index
				p.AddUint8(5)                             // update type
				// This ticket is to track changes to the players around us
				// Everytime this ticket changes, we must send this block out regionally,
				// containing data that identifies all of the owning players characteristics
				p.AddUint16(uint16(ticketID)) // appearance uuid
				p.AddUint64(p1.UsernameHash())             // base37 encoded username
				// ourPlayer.KnownAppearances[player.ServerIndex()] = player.AppearanceTicket()
				sprites := p1.Equips()
				p.AddUint8(uint8(len(sprites))) // length of equipped item sprites  If length less than 12 any ones after length will get set to 0
				for i := 0; i < len(sprites); i++ {
					p.AddUint8(uint8(sprites[i]))
				}

				// The below colors will set the human character animation colors used for this player,
				// it will not apply to any equipment on top of said human character
				// They are simple array indexes corresponding to arrays built in the client
				p.AddUint8(uint8(p1.Appearance.HeadColor))
				p.AddUint8(uint8(p1.Appearance.BodyColor))
				p.AddUint8(uint8(p1.Appearance.LegsColor))
				p.AddUint8(uint8(p1.Appearance.SkinColor))

				// Combat level is the publically shown level of this player; it gives a general
				// idea of how good this player is in combat, it's calculated from the levels of the
				// first 6 skill types
				p.AddUint8(uint8(p1.Skills().CombatLevel()))
				p.AddBoolean(p1.Skulled())
			}
		}
		ourPlayer.UnsetVar(playerEvents)
	}
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
	chunks := ichunks.([]entity.Location)
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
	var local []entity.Entity
	player.LocalObjects.Range(func(e entity.Entity) {
		if o, ok := e.(*Object); ok {
			local = append(local, o)
			if o.Boundary {
				return
			}
			if !player.Near(o, player.ViewRadius()<<1) || GetObject(o.X(), o.Y()) != o {
				if !player.Near(o, player.ViewRadius()*3) {
					// suddenly this local entity is now miles away which isn't very local
					if chunks, ok := player.Var("distantChunks"); ok {
						player.SetVar("distantChunks", append(chunks.([]entity.Location), o.Location.Clone()))
					} else {
						player.SetVar("distantChunks", []entity.Location{o.Location.Clone()})
					}
				} else {
					p.AddUint16(60000)
					p.AddUint8(byte(o.X() - player.X()))
					p.AddUint8(byte(o.Y() - player.Y()))
					changed++
				}
				local = local[:len(local)-1]
			}
		}
		return
	})
	player.LocalObjects.Lock()
	player.LocalObjects.set = local
	player.LocalObjects.Unlock()
	newo := 0
	for _, o := range player.NewObjects() {
		if o.Boundary {
			continue
		}
		if newo >= 300 {
			return
		}
		newo++
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
	var local []entity.Entity
	for _, o := range player.LocalObjects.set {
		if o, ok := o.(*Object); ok {
			if !o.Boundary {
				local = append(local, o)
				continue
			}
			if !player.Near(o, player.ViewRadius()*3) {
				if !player.Near(o, player.ViewRadius()*3) {
					if !player.Near(o, player.ViewRadius()*9) {
						if chunks, ok := player.Var("distantChunks"); ok {
							player.SetVar("distantChunks", append(chunks.([]entity.Location), o.Location.Clone()))
						} else {
							player.SetVar("distantChunks", []entity.Location{o.Location.Clone()})
						}
					} else {
						p.AddUint8(0xFF)
						p.AddUint8(uint8(o.X()))
						p.AddUint8(uint8(o.Y()))
						changed++
					}
				}
			} else if o1 := GetObject(o.X(), o.Y()); o1 != o {
				// network protocol does not support actual removal of previously existing boundary objects
				// so instead, we replace with an invisible boundary that does not block.
				// This is seen in canonical game most notably when slicing a spider web with a weapon
				changed++
				if o1 == nil {
					p.AddUint16(0x10)
					p.AddUint8(uint8(player.TheirDeltaX(o)))
					p.AddUint8(uint8(player.TheirDeltaY(o)))
					p.AddUint8(o.Direction)
					continue
				}
				p.AddUint16(uint16(o1.ID))
				p.AddUint8(uint8(player.TheirDeltaX(o1)))
				p.AddUint8(uint8(player.TheirDeltaY(o1)))
				p.AddUint8(o.Direction)
			} else {
				local = append(local, o)
			}
		}
	}

	player.LocalObjects.Lock()
	player.LocalObjects.set = local
	player.LocalObjects.Unlock()
	
	for _, o := range player.NewObjects() {
		if !o.Boundary {
			continue
		}
		p.AddUint16(uint16(o.ID))
		p.AddUint8(byte(player.TheirDeltaX(o)))
		p.AddUint8(byte(player.TheirDeltaY(o)))
		p.AddUint8(o.Direction)
		player.LocalObjects.Add(o)
		changed++
	}
	if changed <= 0 {
		return nil
	}
	return
}

//ItemLocations Builds a packet with the view-area item positions in it, relative to the player.
// If no new items are available and no existing items are removed from area, returns nil.
func ItemLocations(player *Player) (p *net.Packet) {
	changed := 0
	p = net.NewEmptyPacket(99)
	var local []entity.Entity
	for _, i := range player.LocalItems.set {
		if i, ok := i.(*GroundItem); ok {
			x, y := i.X(), i.Y()
			if !player.Near(i, player.ViewRadius()*3) {
				if !player.Near(i, player.ViewRadius()*9) {
					if chunks, ok := player.Var("distantChunks"); ok {
						player.SetVar("distantChunks", append(chunks.([]entity.Location), i.Location.Clone()))
					} else {
						player.SetVar("distantChunks", []entity.Location{i.Location.Clone()})
					}
				} else {
					// If first byte is 0xFF, all ground items at this location get cleared
					p.AddUint8(255)
					p.AddUint8(byte(player.TheirDeltaX(i)))
					p.AddUint8(byte(player.TheirDeltaY(i)))
					changed++
				}
			} else if !i.VisibleTo(player) || !Region(x, y).Items.Contains(i) {
				p.AddUint16(uint16(i.ID | 0x8000)) // turn remove by ID bit on
				p.AddUint8(byte(player.TheirDeltaX(i)))
				p.AddUint8(byte(player.TheirDeltaY(i)))
				changed++
			} else {
				local = append(local, i)
			}
		}
	}
	player.LocalItems.Lock()
	player.LocalItems.set = local
	player.LocalItems.Unlock()
	for _, i := range player.NewItems() {
		p.AddUint16(uint16(i.ID))
		p.AddUint8(byte(player.TheirDeltaX(i)))
		p.AddUint8(byte(player.TheirDeltaY(i)))
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
			p.AddUint16(uint16(item.ID) | 0x8000)
		} else {
			p.AddUint16(uint16(item.ID))
		}
		if definitions.Items[item.ID].Stackable {
			p.AddSmart0832(item.Amount)
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
	p.AddUint32(uint32(player.Skills().Experience(idx)))
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
		p.AddSmart0832(item.Amount)
		return false
	})
	return p
}

func BankUpdateItem(index, id, amount int) (p *net.Packet) {
	p = net.NewEmptyPacket(249)
	p.AddUint8(uint8(index))
	p.AddUint16(uint16(id))
	p.AddSmart0832(amount)
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
	i, err := strconv.Atoi(strutil.IPToInteger(lastIP).String())
	if err != nil {
		p.AddUint32(127<<24 | 1)
	} else {
		p.AddUint32(uint32(i)) // IP
	}
	p.AddUint16(uint16(inactiveDays)) // Last logged in
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
	return &net.Packet{FrameBuffer: []byte{byte(v)}}
}

//PlaneInfo Builds a packet to update information about the client environment, e.g height, player index...
func PlaneInfo(player *Player) (p *net.Packet) {
	playerInfo := net.NewEmptyPacket(25)
	playerInfo.AddUint16(uint16(player.ServerIndex()))
	playerInfo.AddUint16(2304) // alleged width, tiles per sector also...
	playerInfo.AddUint16(1776) // alleged height

	playerInfo.AddUint16(uint16(player.Plane())) // plane

	playerInfo.AddUint16(944) // REAL plane height
	return playerInfo
}
