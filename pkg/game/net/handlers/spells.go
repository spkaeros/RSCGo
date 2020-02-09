/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package handlers

import (
	"math"
	"strconv"
	"time"

	//	stdrand "math/rand"

	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/rand"
)

type spell struct {
	level int
	kind  int
	runes map[int]int
}

var spellDefs = []spell{
	{
		level: 1,
		kind:  2,
		runes: map[int]int{
			33: 1,
			35: 1,
		},
	},
	{
		level: 3,
		kind:  2,
		runes: map[int]int{
			32: 3,
			34: 2,
			36: 1,
		},
	},
	{
		level: 5,
		kind:  2,
		runes: map[int]int{
			32: 1,
			33: 1,
			35: 1,
		},
	},
	{
		level: 7,
		kind:  3,
		runes: map[int]int{
			32: 1,
			46: 1,
		},
	},
	{
		level: 9,
		kind:  2,
		runes: map[int]int{
			34: 2,
			33: 1,
			35: 1,
		},
	},
	{
		level: 11,
		kind:  2,
		runes: map[int]int{
			32: 3,
			34: 2,
			36: 1,
		},
	},
	{
		level: 13,
		kind:  2,
		runes: map[int]int{
			31: 3,
			33: 2,
			35: 1,
		},
	},
	{
		level: 15,
		kind:  0,
		runes: map[int]int{
			34: 2,
			32: 2,
			40: 1,
		},
	},
	{
		level: 17,
		kind:  2,
		runes: map[int]int{
			33: 2,
			41: 1,
		},
	},
	{
		level: 19,
		kind:  2,
		runes: map[int]int{
			32: 2,
			34: 3,
			36: 1,
		},
	},
	{
		level: 21,
		kind:  3,
		runes: map[int]int{
			31: 3,
			40: 1,
		},
	},
	{
		level: 23,
		kind:  2,
		runes: map[int]int{
			32: 2,
			33: 2,
			41: 1,
		},
	},
	{
		level: 25,
		kind:  0,
		runes: map[int]int{
			31: 1,
			33: 3,
			42: 1,
		},
	},
	{
		level: 27,
		kind:  3,
		runes: map[int]int{
			33: 3,
			46: 1,
		},
	},
	{
		level: 29,
		kind:  2,
		runes: map[int]int{
			34: 3,
			33: 2,
			41: 1,
		},
	},
	{
		level: 31,
		kind:  0,
		runes: map[int]int{
			34: 1,
			33: 3,
			42: 1,
		},
	},
	{
		level: 33,
		kind:  3,
		runes: map[int]int{
			33: 1,
			42: 1,
		},
	},
	{
		level: 35,
		kind:  2,
		runes: map[int]int{
			31: 4,
			33: 3,
			41: 1,
		},
	},
	{
		level: 37,
		kind:  0,
		runes: map[int]int{
			32: 1,
			33: 3,
			42: 1,
		},
	},
	{
		level: 39,
		kind:  2,
		runes: map[int]int{
			34: 2,
			33: 2,
			41: 1,
		},
	},
	{
		level: 41,
		kind:  2,
		runes: map[int]int{
			33: 3,
			38: 1,
		},
	},
	{
		level: 43,
		kind:  3,
		runes: map[int]int{
			31: 4,
			40: 1,
		},
	},
	{
		level: 45,
		kind:  0,
		runes: map[int]int{
			33: 5,
			42: 1,
		},
	},
	{
		level: 47,
		kind:  2,
		runes: map[int]int{
			32: 3,
			33: 3,
			38: 1,
		},
	},
	{
		level: 49,
		kind:  3,
		runes: map[int]int{
			31: 5,
			46: 1,
		},
	},
	{
		level: 50,
		kind:  2,
		runes: map[int]int{
			31: 5,
			38: 1,
		},
	},
	{
		level: 51,
		kind:  0,
		runes: map[int]int{
			32: 2,
			42: 2,
		},
	},
	{
		level: 53,
		kind:  2,
		runes: map[int]int{
			34: 4,
			33: 3,
			38: 1,
		},
	},
	{
		level: 55,
		kind:  3,
		runes: map[int]int{
			31: 5,
			40: 1,
		},
	},
	{
		level: 56,
		kind:  5,
		runes: map[int]int{
			32:  30,
			46:  3,
			611: 1,
		},
	},
	{
		level: 57,
		kind:  3,
		runes: map[int]int{
			34: 10,
			46: 1,
		},
	},
	{
		level: 58,
		kind:  0,
		runes: map[int]int{
			34: 2,
			42: 2,
		},
	},
	{
		level: 59,
		kind:  2,
		runes: map[int]int{
			31: 5,
			33: 4,
			38: 1,
		},
	},
	{
		level: 60,
		kind:  2,
		runes: map[int]int{
			31:  1,
			33:  4,
			619: 2,
		},
	},
	{
		level: 60,
		kind:  2,
		runes: map[int]int{
			31:  2,
			33:  4,
			619: 2,
		},
	},
	{
		level: 60,
		kind:  2,
		runes: map[int]int{
			31:  4,
			33:  1,
			619: 2,
		},
	},
	{
		level: 60,
		kind:  5,
		runes: map[int]int{
			34:  30,
			46:  3,
			611: 1,
		},
	},
	{
		level: 62,
		kind:  2,
		runes: map[int]int{
			33:  5,
			619: 1,
		},
	},
	{
		level: 63,
		kind:  5,
		runes: map[int]int{
			31:  30,
			46:  3,
			611: 1,
		},
	},
	{
		level: 65,
		kind:  2,
		runes: map[int]int{
			32:  7,
			33:  5,
			619: 1,
		},
	},
	{
		level: 66,
		kind:  5,
		runes: map[int]int{
			33:  30,
			46:  3,
			611: 1,
		},
	},
	{
		level: 66,
		kind:  2,
		runes: map[int]int{
			34:  5,
			32:  5,
			825: 1,
		},
	},
	{
		level: 68,
		kind:  3,
		runes: map[int]int{
			32: 15,
			34: 15,
			46: 1,
		},
	},
	{
		level: 70,
		kind:  2,
		runes: map[int]int{
			34:  7,
			33:  5,
			619: 1,
		},
	},
	{
		level: 73,
		kind:  2,
		runes: map[int]int{
			34:  8,
			32:  8,
			825: 1,
		},
	},
	{
		level: 75,
		kind:  2,
		runes: map[int]int{
			31:  7,
			33:  5,
			619: 1,
		},
	},
	{
		level: 80,
		kind:  2,
		runes: map[int]int{
			34:  12,
			32:  12,
			825: 1,
		},
	},
	{
		level: 80,
		kind:  0,
		runes: map[int]int{
			31:  3,
			33:  3,
			619: 3,
		},
	},
}

var dmgs = map[int]int{ // reqLvl mapped to maxDmg
	1:  1,
	5:  2,
	9:  3,
	13: 4,
	17: 5,
	23: 6,
	29: 7,
	35: 8,
	41: 9,
	47: 10,
	53: 11,
	59: 12,
	62: 13,
	65: 14,
	70: 15,
	75: 16,
}

func init() {
	AddHandler("spellnpc", func(player *world.Player, p *net.Packet) {
		targetIndex := int(p.ReadShort())
		target := world.GetNpc(targetIndex)
		if target == nil {
			return
		}
		spellIndex := int(p.ReadShort())
		log.Info.Println("cast on npc:", targetIndex, target.ID, spellIndex)
		handleSpells(player, spellIndex, target)
	})
	AddHandler("spellplayer", func(player *world.Player, p *net.Packet) {
		targetIndex := int(p.ReadShort())
		target, ok := world.Players.FromIndex(targetIndex)
		if !ok {
			return
		}
		spellIndex := int(p.ReadShort())
		log.Info.Println("cast on player:", targetIndex, target.String(), spellIndex)
		handleSpells(player, spellIndex, target)
	})
	AddHandler("spellself", func(player *world.Player, p *net.Packet) {
		idx := int(p.ReadShort())

		log.Info.Println("Cast on self:", idx)
		handleSpells(player, idx, nil)
	})
	AddHandler("spellinvitem", func(player *world.Player, p *net.Packet) {
		itemIndex := int(p.ReadShort())
		spellIndex := int(p.ReadShort())
		log.Info.Println("Cast on invitem:", spellIndex, "on", itemIndex)
	})
	AddHandler("spellgrounditem", func(player *world.Player, p *net.Packet) {
		itemX := int(p.ReadShort())
		itemY := int(p.ReadShort())
		itemID := int(p.ReadShort())
		spellIndex := int(p.ReadShort())
		log.Info.Println(itemX, itemY, itemID, "cast on grounditem:", spellIndex, "on", strconv.Itoa(itemID), "at", strconv.Itoa(itemX)+","+strconv.Itoa(itemY))
		handleSpells(player, spellIndex, nil)
	})
}

func handleSpells(player *world.Player, idx int, target world.MobileEntity) {
	if idx < 0 || idx >= len(spellDefs) {
		return
	}
	s := spellDefs[idx]
	checkAndRemoveRunes := func() bool {
		for id, amt := range s.runes {
			if player.Inventory.CountID(id) < amt {
				log.Suspicious.Println(player, "casted spell with not enough runes")
				player.Message("You don't have all the reagents you need for this spell")
				return false
			}
		}
		for id, amt := range s.runes {
			player.Inventory.RemoveByID(id, amt)
		}
		return true
	}
	finalize := func() {
		player.TransAttrs.SetVar("lastSpell", time.Now())
		player.PlaySound("spellok")
		player.Message("Cast spell successfully")
		player.ResetPath()
	}
	if player.Skills().Current(world.StatMagic) < s.level {
		player.Message("Your magic ability is not high enough for this spell.")
		player.ResetPath()
		return
	}
	log.Info.Println(s)
	if lvDelta := player.Skills().Current(world.StatMagic) - s.level; lvDelta < 0 || (lvDelta < 10-int(math.Min(math.Max((float64(player.MagicPoints())-5)/5, 0), 5)) &&
		rand.Int31N(0, (lvDelta+2)*2) == 0) {
		player.PlaySound("spellfail")
		player.Message("The spell fails! You may try again in 20 seconds")
		player.ResetPath()
		return
	}

	switch s.kind {
	case 0: // Self
		switch idx {
		case 7:
			count := player.Inventory.CountID(20)
			if count <= 0 {
				player.Message("You aren't holding any bones!")
				return
			}
			if !checkAndRemoveRunes() {
				return
			}
			player.Inventory.RemoveByID(20, count)
			for i := 0; i < count; i++ {
				player.Inventory.Add(249, 1)
			}
			player.SendInventory()
			finalize()
			return
		case 12: // Varrock Teleport
			if !checkAndRemoveRunes() {
				return
			}
			teleport(player, 120, 504)
			finalize()
		case 15: // Lumbridge Teleport
			if !checkAndRemoveRunes() {
				return
			}
			teleport(player, 120, 648)
			finalize()
		case 18: // Falador Teleport
			if !checkAndRemoveRunes() {
				return
			}
			teleport(player, 312, 552)
			finalize()
		case 22: // Camelot Teleport
			if !checkAndRemoveRunes() {
				return
			}
			teleport(player, 465, 456)
			finalize()
		case 26: // Ardougne Teleport
			//			player.Teleport(588, 621, true)
			player.Message("You don't know how to cast this spell yet")
			player.Message("You need to do the plague city quest")
			return
		case 31: // Watchtower Teleport
			//			player.Teleport(493, 3525, true)
			player.Message("You cannot cast this spell")
			player.Message("You need to finish the watchtower quest first")
			return
		case 47:
			fallthrough
		default:
			player.Message("@ora@Not yet implemented.")
			return
		}
	case 2: // combat spell
		if target == nil {
			return
		}
		// if it is in our damage defs, it's an offensive spell without any special fx
		if val, ok := dmgs[s.level]; ok {
			player.SetDistancedAction(func() bool {
				if player.Busy() || target == nil {
					return true
				}

				if player.WithinRange(world.NewLocation(target.X(), target.Y()), 4) {
					steps := 0
					xOff := player.X()
					yOff := player.Y()
					for steps < 10 {
						if xOff == target.X() && yOff == target.Y() {
							break
						}
						if yOff > target.Y() {
							yOff--
							if world.IsTileBlocking(xOff, yOff, world.ClipSouth, false) {
								return false
							}
							steps++
						} else if yOff < target.Y() {
							yOff++
							if world.IsTileBlocking(xOff, yOff, world.ClipNorth, false) {
								return false
							}
							steps++
						}
						if xOff > target.X() {
							xOff--
							if world.IsTileBlocking(xOff, yOff, world.ClipWest, false) {
								return false
							}
							steps++
						} else if xOff < target.X() {
							xOff++
							if world.IsTileBlocking(xOff, yOff, world.ClipEast, false) {
								return false
							}
							steps++
						}
					}
					// reaching here means made it to target within 4 steps without hitting a barrier
					player.ResetPath()
					if !checkAndRemoveRunes() {
						return true
					}
					finalize()

					dmg := float64(val)
					probs := map[int]float64{}
					rat := 45.0 + float64(player.MagicPoints())
					peak := (dmg / 100.0) * rat
					dip := peak / 3.0

					curProb := 100.0 * dmg
					for i := 0.0; i <= dmg; i++ {
						probs[int(i)] = curProb
						if i < dip || i > peak {
							curProb -= (dmg * 100) / 3
						} else {
							curProb += (dmg * 100) / 3
						}
					}
					hit := int(math.Min(float64(target.Skills().Current(world.StatHits)), float64(world.WeightedChoice(probs))))
					target.Skills().DecreaseCur(world.StatHits, hit)
					if target.Skills().Current(world.StatHits) <= 0 {
						target.Killed(player)
						return true
					}
					projectile := world.CreateProjectile(player,target,1)
					for _, v := range player.NearbyPlayers() {
						
						v.SendPacket(projectile)
					}
					player.SendPacket(projectile)
					target.Damage(hit)
					return true
				}
				player.WalkTo(world.NewLocation(target.X(), target.Y()))
				return false
			})
		} else {
			// TODO: Handle these spells
			// Curses (e.g any stat-lowering spells), godspells, and crumble undead are all that fit this description I think?
		}
	case 3: // enchant spell
		switch idx {
		//		case 16: // telegrab
		default:
			player.Message("@ora@Not yet implemented.")
			return
		}
	case 5: // charge orb
		switch idx {
		default:
			player.Message("@ora@Not yet implemented.")
			return
		}
	default: // charge orb
		return
	}
}

func teleport(target *world.Player, x, y int) {
	target.SendPacket(world.TeleBubble(0, 0))
	for _, nearbyPlayer := range target.NearbyPlayers() {
		nearbyPlayer.SendPacket(world.TeleBubble(target.X()-nearbyPlayer.X(), target.Y()-nearbyPlayer.Y()))
	}
	plane := target.Plane()
	target.Teleport(x, y)
	if target.Plane() != plane {
		target.SendPacket(world.PlaneInfo(target))
	}
}
