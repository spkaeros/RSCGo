/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package script

import (
	"os"
	"reflect"
	"time"

	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(*world.Player, []string))

func init() {
	env.Packages["world"] = map[string]reflect.Value{
		"getPlayer":       reflect.ValueOf(world.Players.FromIndex),
		"getPlayerByName": reflect.ValueOf(world.Players.FromUserHash),
		"players":         reflect.ValueOf(world.Players),
		"replaceObject":   reflect.ValueOf(world.ReplaceObject),
		"addObject":       reflect.ValueOf(world.AddObject),
		"removeObject":    reflect.ValueOf(world.RemoveObject),
		"addNpc":          reflect.ValueOf(world.AddNpc),
		"removeNpc":       reflect.ValueOf(world.RemoveNpc),
		"addItem":         reflect.ValueOf(world.AddItem),
		"removeItem":      reflect.ValueOf(world.RemoveItem),
		"getObjectAt":     reflect.ValueOf(world.GetObject),
		"getNpc":          reflect.ValueOf(world.GetNpc),
		"checkCollisions": reflect.ValueOf(world.IsTileBlocking),
		"tileData":        reflect.ValueOf(world.CollisionData),
		"kickPlayer": reflect.ValueOf(func(client *world.Player) {
			client.SendPacket(world.Logout)
			client.Destroy()
		}),
		"updateStarted": reflect.ValueOf(func() bool {
			return !world.UpdateTime.IsZero()
		}),
		"announce": reflect.ValueOf(func(msg string) {
			world.Players.Range(func(player *world.Player) {
				player.Message("@que@" + msg)
			})
		}),
		"walkTo": reflect.ValueOf(func(target *world.Player, x, y int) {
			target.WalkTo(world.NewLocation(x, y))
		}),
		"systemUpdate": reflect.ValueOf(func(t int) {
			world.UpdateTime = time.Now().Add(time.Second * time.Duration(t))
			go func() {
				time.Sleep(time.Second * time.Duration(t))
				world.Players.Range(func(player *world.Player) {
					player.SendPacket(world.Logout)
					player.Destroy()
				})
				time.Sleep(300 * time.Millisecond)
				os.Exit(200)
			}()
			world.Players.Range(func(player *world.Player) {
				player.SendUpdateTimer()
			})
		}),
		"teleport": reflect.ValueOf(func(target *world.Player, x, y int, bubble bool) {
			if bubble {
				target.SendPacket(world.TeleBubble(0, 0))
				for _, nearbyPlayer := range target.NearbyPlayers() {
					nearbyPlayer.SendPacket(world.TeleBubble(target.X()-nearbyPlayer.X(), target.Y()-nearbyPlayer.Y()))
				}
			}
			plane := target.Plane()
			target.Teleport(x, y)
			if target.Plane() != plane {
				target.SendPacket(world.PlaneInfo(target))
			}
		}),
		"newShop":        reflect.ValueOf(world.NewShop),
		"newGeneralShop": reflect.ValueOf(world.NewGeneralShop),
		"getShop":        reflect.ValueOf(world.Shops.Get),
		"hasShop":        reflect.ValueOf(world.Shops.Contains),
	}
	env.PackageTypes["world"] = map[string]reflect.Type{
		"players":    reflect.TypeOf(world.Players),
		"player":     reflect.TypeOf(&world.Player{}),
		"object":     reflect.TypeOf(&world.Object{}),
		"item":       reflect.TypeOf(&world.Item{}),
		"groundItem": reflect.TypeOf(&world.GroundItem{}),
		"npc":        reflect.TypeOf(&world.NPC{}),
		"location":   reflect.TypeOf(world.Location{}),
	}
	env.Packages["ids"] = map[string]reflect.Value{
		"COOKEDMEAT":               reflect.ValueOf(132),
		"BURNTMEAT":                reflect.ValueOf(134),
		"FLIER":                    reflect.ValueOf(201),
		"LEATHER_GLOVES":           reflect.ValueOf(16),
		"BOOTS":                    reflect.ValueOf(17),
		"SEAWEED":                  reflect.ValueOf(622),
		"OYSTER":                   reflect.ValueOf(793),
		"CASKET":                   reflect.ValueOf(549),
		"RAW_RAT_MEAT":             reflect.ValueOf(503),
		"RAW_SHRIMP":               reflect.ValueOf(349),
		"RAW_ANCHOVIES":            reflect.ValueOf(351),
		"RAW_TROUT":                reflect.ValueOf(358),
		"RAW_SALMON":               reflect.ValueOf(356),
		"RAW_PIKE":                 reflect.ValueOf(363),
		"RAW_SARDINE":              reflect.ValueOf(354),
		"RAW_HERRING":              reflect.ValueOf(361),
		"RAW_BASS":                 reflect.ValueOf(550),
		"RAW_MACKEREL":             reflect.ValueOf(552),
		"RAW_COD":                  reflect.ValueOf(554),
		"RAW_LOBSTER":              reflect.ValueOf(372),
		"RAW_SWORDFISH":            reflect.ValueOf(369),
		"RAW_TUNA":                 reflect.ValueOf(366),
		"HOLY_SYMBOL_OF_SARADOMIN": reflect.ValueOf(385),
		"RAW_SHARK":                reflect.ValueOf(545),
		"WOODEN_SHIELD":            reflect.ValueOf(4),
		"BRONZE_LSWORD":            reflect.ValueOf(70),
		"NET":                      reflect.ValueOf(376),
		"BIG_NET":                  reflect.ValueOf(548),
		"LOBSTER_POT":              reflect.ValueOf(375),
		"FISHING_ROD":              reflect.ValueOf(377),
		"FLYFISHING_ROD":           reflect.ValueOf(378),
		"OILY_FISHING_ROD":         reflect.ValueOf(589),
		"RAW_LAVA_EEL":             reflect.ValueOf(591),
		"HARPOON":                  reflect.ValueOf(379),
		"FISHING_BAIT":             reflect.ValueOf(380),
		"FEATHER":                  reflect.ValueOf(381),
		"BRONZE_PICKAXE":           reflect.ValueOf(156),
		"IRON_PICKAXE":             reflect.ValueOf(1258),
		"STEEL_PICKAXE":            reflect.ValueOf(1259),
		"MITHRIL_PICKAXE":          reflect.ValueOf(1260),
		"ADAM_PICKAXE":             reflect.ValueOf(1261),
		"RUNE_PICKAXE":             reflect.ValueOf(1262),
		"TIN_ORE":                  reflect.ValueOf(202),
		"SLEEPING_BAG":             reflect.ValueOf(1263),
		"NEEDLE":                   reflect.ValueOf(39),
		"THREAD":                   reflect.ValueOf(43),
		"FIRE_RUNE":                reflect.ValueOf(31),
		"WATER_RUNE":               reflect.ValueOf(32),
		"AIR_RUNE":                 reflect.ValueOf(33),
		"EARTH_RUNE":               reflect.ValueOf(34),
		"MIND_RUNE":                reflect.ValueOf(35),
		"BODY_RUNE":                reflect.ValueOf(36),
		"LIFE_RUNE":                reflect.ValueOf(37),
		"DEATH_RUNE":               reflect.ValueOf(38),
		"NATURE_RUNE":              reflect.ValueOf(40),
		"CHAOS_RUNE":               reflect.ValueOf(41),
		"LAW_RUNE":                 reflect.ValueOf(42),
		"COSMIC_RUNE":              reflect.ValueOf(46),
		"BLOOD_RUNE":               reflect.ValueOf(619),
		"AIR_STAFF":                reflect.ValueOf(101),
		"WATER_STAFF":              reflect.ValueOf(102),
		"EARTH_STAFF":              reflect.ValueOf(103),
		"FIRE_STAFF":               reflect.ValueOf(197),
		"FIRE_BATTLESTAFF":         reflect.ValueOf(615),
		"WATER_BATTLESTAFF":        reflect.ValueOf(616),
		"AIR_BATTLESTAFF":          reflect.ValueOf(617),
		"EARTH_BATTLESTAFF":        reflect.ValueOf(618),
		"E_FIRE_BATTLESTAFF":       reflect.ValueOf(682),
		"E_WATER_BATTLESTAFF":      reflect.ValueOf(683),
		"E_AIR_BATTLESTAFF":        reflect.ValueOf(684),
		"E_EARTH_BATTLESTAFF":      reflect.ValueOf(685),
		"BONES":                    reflect.ValueOf(20),
		"BAT_BONES":                reflect.ValueOf(604),
		"DRAGON_BONES":             reflect.ValueOf(614),
		"RUNE_2H":                  reflect.ValueOf(81),
		"RUNE_CHAIN":               reflect.ValueOf(400),
		"RUNE_PLATEBODY":           reflect.ValueOf(401),
		"RUNE_PLATETOP":            reflect.ValueOf(407),
		"BLACK_DAGGER":             reflect.ValueOf(423),
		"GOLD_RING":                reflect.ValueOf(283),
		"SILK":                     reflect.ValueOf(200),
		"IRON_ORE_CERTIFICATE":     reflect.ValueOf(517),
		"CHOCOLATE_SLICE":          reflect.ValueOf(336),
		"CHOCOLATE_BAR":            reflect.ValueOf(337),
		"SPINACH_ROLL":             reflect.ValueOf(179),
		"DRAGON_SWORD":             reflect.ValueOf(593),
		"DRAGON_AXE":               reflect.ValueOf(594),
		"CHARGED_DSTONE_AMMY":      reflect.ValueOf(597),
		"DRAGON_HELMET":            reflect.ValueOf(795),
		"DRAGON_SHIELD":            reflect.ValueOf(1278),
		"EASTER_EGG":               reflect.ValueOf(677),
		"CHRISTMAS_CRACKER":        reflect.ValueOf(575),
		"PARTYHAT_RED":             reflect.ValueOf(576),
		"PARTYHAT_YELLOW":          reflect.ValueOf(577),
		"PARTYHAT_BLUE":            reflect.ValueOf(578),
		"PARTYHAT_GREEN":           reflect.ValueOf(579),
		"PARTYHAT_PINK":            reflect.ValueOf(580),
		"PARTYHAT_WHITE":           reflect.ValueOf(581),
		"GREEN_MASK":               reflect.ValueOf(828),
		"RED_MASK":                 reflect.ValueOf(831),
		"BLUE_MASK":                reflect.ValueOf(832),
		"SANTA_HAT":                reflect.ValueOf(971),
		"PRESENT":                  reflect.ValueOf(980),
		"GNOME_BALL":               reflect.ValueOf(981),
		"BLURITE_ORE":              reflect.ValueOf(266),
		"CLAY":                     reflect.ValueOf(149),
		"COPPER_ORE":               reflect.ValueOf(150),
		"IRON_ORE":                 reflect.ValueOf(151),
		"GOLD":                     reflect.ValueOf(152),
		"SILVER":                   reflect.ValueOf(383),
		"GOLD2":                    reflect.ValueOf(690),
		"MITHRIL_ORE":              reflect.ValueOf(153),
		"ADAM_ORE":                 reflect.ValueOf(154),
		"RUNITE_ORE":               reflect.ValueOf(409),
		"COAL":                     reflect.ValueOf(155),
	}
	env.Packages["bind"] = map[string]reflect.Value{
		"onLogin": reflect.ValueOf(func(fn func(player *world.Player)) {
			LoginTriggers = append(LoginTriggers, fn)
		}),
		"invOnBoundary": reflect.ValueOf(func(fn func(player *world.Player, boundary *world.Object, item *world.Item) bool) {
			InvOnBoundaryTriggers = append(InvOnBoundaryTriggers, fn)
		}),
		"invOnPlayer": reflect.ValueOf(func(pred func(*world.Item) bool, fn func(player *world.Player, target *world.Player, item *world.Item)) {
			InvOnPlayerTriggers = append(InvOnPlayerTriggers, ItemOnPlayerTrigger{pred, fn})
		}),
		"invOnObject": reflect.ValueOf(func(fn func(player *world.Player, boundary *world.Object, item *world.Item) bool) {
			InvOnObjectTriggers = append(InvOnObjectTriggers, fn)
		}),
		"object": reflect.ValueOf(func(pred func(*world.Object, int) bool, fn func(player *world.Player, object *world.Object, click int)) {
			ObjectTriggers = append(ObjectTriggers, ObjectTrigger{pred, fn})
		}),
		"item": reflect.ValueOf(func(check func(item *world.Item) bool, fn func(player *world.Player, item *world.Item)) {
			ItemTriggers = append(ItemTriggers, ItemTrigger{check, fn})
		}),
		"boundary": reflect.ValueOf(func(pred func(*world.Object, int) bool, fn func(player *world.Player, object *world.Object, click int)) {
			BoundaryTriggers = append(BoundaryTriggers, ObjectTrigger{pred, fn})
		}),
		"npc": reflect.ValueOf(func(predicate func(npc *world.NPC) bool, fn func(player *world.Player, npc *world.NPC)) {
			NpcTriggers = append(NpcTriggers, NpcTrigger{predicate, fn})
		}),
		"npcAttack": reflect.ValueOf(func(pred NpcActionPredicate, fn func(player *world.Player, npc *world.NPC)) {
			NpcAtkTriggers = append(NpcAtkTriggers, world.NpcBlockingTrigger{pred, fn})
		}),
		"npcKilled": reflect.ValueOf(func(pred NpcActionPredicate, fn func(player *world.Player, npc *world.NPC)) {
			world.NpcDeathTriggers = append(world.NpcDeathTriggers, world.NpcBlockingTrigger{pred, fn})
		}),
		"command": reflect.ValueOf(func(name string, fn func(p *world.Player, args []string)) {
			CommandHandlers[name] = fn
		}),
	}
	env.Packages["log"] = map[string]reflect.Value{
		"debug":  reflect.ValueOf(log.Info.Println),
		"debugf": reflect.ValueOf(log.Info.Printf),
		"warn":   reflect.ValueOf(log.Warning.Println),
		"warnf":  reflect.ValueOf(log.Warning.Printf),
		"err":    reflect.ValueOf(log.Error.Println),
		"errf":   reflect.ValueOf(log.Error.Printf),
		"cheat":  reflect.ValueOf(log.Suspicious.Println),
		"cheatf": reflect.ValueOf(log.Suspicious.Printf),
		"cmd":    reflect.ValueOf(log.Commands.Println),
		"cmdf":   reflect.ValueOf(log.Commands.Printf),
	}
}

func WorldModule() *env.Env {
	e := env.NewEnv()
	e.Define("sleep", time.Sleep)
	e.Define("runAfter", time.AfterFunc)
	e.Define("after", time.After)
	e.Define("tMinute", time.Second*60)
	e.Define("tHour", time.Second*60*60)
	e.Define("tSecond", time.Second)
	e.Define("tMillis", time.Millisecond)
	e.Define("ChatDelay", time.Millisecond*1800)
	e.Define("tNanos", time.Nanosecond)
	e.Define("ATTACK", world.StatAttack)
	e.Define("DEFENSE", world.StatDefense)
	e.Define("STRENGTH", world.StatStrength)
	e.Define("HITPOINTS", world.StatHits)
	e.Define("RANGED", world.StatRanged)
	e.Define("PRAYER", world.StatPrayer)
	e.Define("MAGIC", world.StatMagic)
	e.Define("COOKING", world.StatCooking)
	e.Define("WOODCUTTING", world.StatWoodcutting)
	e.Define("FLETCHING", world.StatFletching)
	e.Define("FISHING", world.StatFishing)
	e.Define("FIREMAKING", world.StatFiremaking)
	e.Define("CRAFTING", world.StatCrafting)
	e.Define("SMITHING", world.StatSmithing)
	e.Define("MINING", world.StatMining)
	e.Define("HERBLAW", world.StatHerblaw)
	e.Define("AGILITY", world.StatAgility)
	e.Define("THIEVING", world.StatThieving)
	e.Define("PRAYER_RAPID_RESTORE", 6)
	e.Define("PRAYER_RAPID_HEAL", 7)
	e.Define("ZeroTime", time.Time{})
	e.Define("itemDefs", world.ItemDefs)
	e.Define("objectDefs", world.ObjectDefs)
	e.Define("boundaryDefs", world.BoundaryDefs)
	e.Define("npcDefs", world.NpcDefs)
	e.Define("lvlToExp", world.LevelToExperience)
	e.Define("expToLvl", world.ExperienceToLevel)
	e.Define("withinWorld", world.WithinWorld)
	e.Define("skillIndex", world.SkillIndex)
	e.Define("skillName", world.SkillName)
	e.Define("newNpc", world.NewNpc)
	e.Define("newObject", world.NewObject)
	e.Define("base37", strutil.Base37.Encode)
	e.Define("rand", rand.Int31N)
	e.Define("North", world.North)
	e.Define("NorthEast", world.NorthEast)
	e.Define("NorthWest", world.NorthWest)
	e.Define("South", world.South)
	e.Define("SouthEast", world.SouthEast)
	e.Define("SouthWest", world.SouthWest)
	e.Define("East", world.East)
	e.Define("West", world.West)
	e.Define("parseDirection", world.ParseDirection)
	e.Define("contains", func(s []int64, elem int64) bool {
		for _, v := range s {
			if v == elem {
				return true
			}
		}
		return false
	})
	e.Define("gatheringSuccess", func(req, cur int) bool {
		if cur < req {
			return false
		}
		return float64(rand.Int31N(1, 128)) <= (float64(cur)+40)-(float64(req)*1.5)
	})
	e.Define("roll", world.Chance)
	e.Define("boundedRoll", world.BoundedChance)
	e.Define("weightedChance", world.WeightedChoice)

	e.Define("npcPredicate", func(ids ...interface{}) func(*world.NPC) bool {
		return func(npc *world.NPC) bool {
			for _, id := range ids {
				if cmd, ok := id.(string); ok {
					if npc.Name() == cmd {
						return true
					}
				} else if id, ok := id.(int64); ok {
					if npc.ID == int(id) {
						return true
					}
				}
			}
			return false
		}
	})
	e.Define("npcBlockingPredicate", func(ids ...interface{}) func(*world.Player, *world.NPC) bool {
		return func(player *world.Player, npc *world.NPC) bool {
			for _, id := range ids {
				if cmd, ok := id.(string); ok {
					if npc.Name() == cmd {
						return true
					}
				} else if id, ok := id.(int64); ok {
					if npc.ID == int(id) {
						return true
					}
				}
			}
			return false
		}
	})
	e.Define("itemPredicate", func(ids ...interface{}) func(*world.Item) bool {
		return func(item *world.Item) bool {
			for _, id := range ids {
				if cmd, ok := id.(string); ok {
					if item.Command() == cmd {
						return true
					}
				} else if id, ok := id.(int64); ok {
					if item.ID == int(id) {
						return true
					}
				}
			}
			return false
		}
	})
	e.Define("objectPredicate", func(ids ...interface{}) func(*world.Object, int) bool {
		return func(object *world.Object, click int) bool {
			for _, id := range ids {
				if cmd, ok := id.(string); ok {
					if world.ObjectDefs[object.ID].Commands[click] == cmd {
						return true
					}
				} else if id, ok := id.(int64); ok {
					if object.ID == int(id) {
						return true
					}
				}
			}
			return false
		}
	})
	e = core.Import(e)
	return e
}
