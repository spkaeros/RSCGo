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
	"os"
	"runtime/pprof"
	"reflect"
	"time"
	"strings"
	"fmt"
	"strconv"
	"sync"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"github.com/mattn/anko/parser"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/tasks"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/strutil"

	// Defines various package-related scripting utilities
	_ "github.com/mattn/anko/packages"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(*Player, []string))

const serverPrefix = "@que@@whi@[@cya@SERVER@whi@]: "

func init() {
	env.Packages["state"] = map[string]reflect.Value{
		"UsingItem":			  reflect.ValueOf(MSItemAction),
		"InDuel":				  reflect.ValueOf(StateDueling),
		"InTrade":				  reflect.ValueOf(StateTrading),
		"DoingThing":			  reflect.ValueOf(StateAction),
		"ChatMenu":				  reflect.ValueOf(StateChatChoosing),
		"ChangingLooks":	  reflect.ValueOf(StateChangingLooks),
		"Chatting":				  reflect.ValueOf(StateChatting),
		"OptionMenu":			  reflect.ValueOf(MSItemAction),
		"Banking":				  reflect.ValueOf(StateBanking),
		"Shopping":				  reflect.ValueOf(StateShopping),
		"Batching":				  reflect.ValueOf(MSBatching),
	}
	env.Packages["world"] = map[string]reflect.Value{
		"getPlayer":              reflect.ValueOf(Players.FindIndex),
		"OrderedDirections":              reflect.ValueOf(OrderedDirections),
		"getPlayerByName":        reflect.ValueOf(Players.FindHash),
		"getNpcNear":             reflect.ValueOf(NpcNearest),
		"getGridNpc":             reflect.ValueOf(NpcVisibleFrom),
		"players":                reflect.ValueOf(Players),
		"getEquipmentDefinition": reflect.ValueOf(definitions.Equip),
		"replaceObject":          reflect.ValueOf(ReplaceObject),
		"addObject":              reflect.ValueOf(AddObject),
		"removeObject":           reflect.ValueOf(RemoveObject),
		"addNpc":                 reflect.ValueOf(AddNpc),
		"removeNpc":              reflect.ValueOf(RemoveNpc),
		"getItem":                reflect.ValueOf(GetItem),
		"itemActions":            reflect.ValueOf(&ItemTriggers),
		"unhandledMessage":       reflect.ValueOf(DefaultActionMessage),
		"addItem":                reflect.ValueOf(AddItem),
		"removeItem":             reflect.ValueOf(RemoveItem),
		"maxX":                   reflect.ValueOf(MaxX),
		"maxY":                   reflect.ValueOf(MaxY),
		"getObjectAt":            reflect.ValueOf(GetObject),
		"getNpc":                 reflect.ValueOf(GetNpc),
		"attackNpcCalls":         reflect.ValueOf(NpcAtkTriggers),
		"checkCollisions":        reflect.ValueOf(IsTileBlocking),
		"tileData":               reflect.ValueOf(CollisionData),
		"kickPlayer": reflect.ValueOf(func(client *Player) {
			client.Unregister()
		}),
		"updateStarted": reflect.ValueOf(func() bool {
			return !UpdateTime.IsZero()
		}),
		"announce": reflect.ValueOf(func(msg string) {
			Players.Range(func(player *Player) {
				player.Message("@que@" + msg)
			})
		}),
		"walkTo": reflect.ValueOf(func(target *Player, x, y int) {
			target.WalkTo(NewLocation(x, y))
		}),
		"systemUpdate": reflect.ValueOf(func(t int) {
			start := time.Now()
			UpdateTime = time.Now().Add(time.Second * time.Duration(t))
			tasks.Schedule(1, func() bool {
				if time.Since(start) >= time.Duration(t) * time.Second {
					Players.Range(func(player *Player) {
						player.Unregister()
					})
					time.Sleep(2 * time.Second)
					os.Exit(200)
					return true
				}
				if CurrentTick() % 10 == 0 {
					Players.Range(func(player *Player) {
						player.SendUpdateTimer()
					})
				}
				return false
			})
		}),
		"teleport": reflect.ValueOf(func(target *Player, x, y int, bubble bool) {
			if bubble {
				target.SendPacket(TeleBubble(0, 0))
				target.LocalPlayers.RangePlayers(func(p1 *Player) bool {
					p1.SendPacket(TeleBubble(target.X()-p1.X(), target.Y()-p1.Y()))
					return false
				})
				// nearList := &MobList{mobSet: target.NearbyPlayers()}
				// nearList.RangePlayers(func(p1 *Player) bool {
					// p1.SendPacket(TeleBubble(target.X()-p1.X(), target.Y()-p1.Y()))
					// return false
				// })
			}
			plane := target.Plane()
			target.ResetPath()
			target.Teleport(x, y)
			if target.Plane() != plane {
				target.SendPacket(PlaneInfo(target))
			}
		}),
		"newGroundItemFor": reflect.ValueOf(NewGroundItemFor),
		"newGroundItem":    reflect.ValueOf(NewGroundItem),
		"curTick":		    reflect.ValueOf(tasks.Ticks.Load()),
		"newShop":          reflect.ValueOf(NewShop),
		"newItem":          reflect.ValueOf(NewItem),
		"newLocation":      reflect.ValueOf(NewLocation),
		"newGeneralShop":   reflect.ValueOf(NewGeneralShop),
		"getShop":          reflect.ValueOf(Shops.Get),
		"hasShop":          reflect.ValueOf(Shops.Contains),
	}
	env.Packages["net"] = map[string]reflect.Value {
		"barePacket": reflect.ValueOf(net.NewEmptyPacket),
		"logout": reflect.ValueOf(Logout),
		"cannotLogout": reflect.ValueOf(CannotLogout),
		"bankUpdateItem": reflect.ValueOf(BankUpdateItem),
		"shopOpen": reflect.ValueOf(ShopOpen),
	}
	env.PackageTypes["world"] = map[string]reflect.Type{
		"players":    reflect.TypeOf(Players),
		"player":     reflect.TypeOf(&Player{}),
		"object":     reflect.TypeOf(&Object{}),
		"item":       reflect.TypeOf(&Item{}),
		"groundItem": reflect.TypeOf(&GroundItem{}),
		"npc":        reflect.TypeOf(&NPC{}),
		"location":   reflect.TypeOf(Location{}),
	}
	env.Packages["packets"] = map[string]reflect.Value{
		"ping": reflect.ValueOf(67),
		"tradeRequest": reflect.ValueOf(142),
		"tradeDecline": reflect.ValueOf(230),
		"tradeAccept": reflect.ValueOf(55),
		"tradeAccept2": reflect.ValueOf(104),
		"tradeUpdate": reflect.ValueOf(46),
		"follow": reflect.ValueOf(165),
		"closeBank": reflect.ValueOf(212),
		"depositBank": reflect.ValueOf(23),
		"withdrawBank": reflect.ValueOf(22),
		"menuAnswer": reflect.ValueOf(116),
		"equip": reflect.ValueOf(169),
		"npcChat": reflect.ValueOf(153),
		"npcAction": reflect.ValueOf(202),
		"sceneAction": reflect.ValueOf(136),
		"sceneAction2": reflect.ValueOf(79),
		"boundaryAction": reflect.ValueOf(127),
		"boundaryAction2": reflect.ValueOf(14),
		"invOnScene": reflect.ValueOf(115),
		"unequip": reflect.ValueOf(170),
		"dropItem": reflect.ValueOf(246),
		"recoverAccount": reflect.ValueOf(220),
		"closeStream": reflect.ValueOf(31),
		"logout": reflect.ValueOf(102),
		"pickupItem": reflect.ValueOf(247),
		"itemAction": reflect.ValueOf(90),
		"spellOnNpc": reflect.ValueOf(50),
		"spellOnGroundItem": reflect.ValueOf(249),
		"spellOnInvItem": reflect.ValueOf(4),
		"spellOnPlayer": reflect.ValueOf(229),
		"spellOnSelf": reflect.ValueOf(137),
		"ticketRequests": reflect.ValueOf(163),
		"appearance": reflect.ValueOf(235),
		"report": reflect.ValueOf(206),
		"attackNpc": reflect.ValueOf(190),
		"attackPlayer": reflect.ValueOf(171),
		"fightMode": reflect.ValueOf(29),
		"duelRequest": reflect.ValueOf(103),
		"duelDecline": reflect.ValueOf(197),
		"duelAccept": reflect.ValueOf(176),
		"duelAccept2": reflect.ValueOf(77),
		"duelUpdate": reflect.ValueOf(33),
		"duelSettings": reflect.ValueOf(8),
		"settings": reflect.ValueOf(111),
		"privacySettings": reflect.ValueOf(64),
		"shopSell": reflect.ValueOf(221),
		"shopBuy": reflect.ValueOf(236),
		"shopClose": reflect.ValueOf(166),
		"command": reflect.ValueOf(38),
		"chat": reflect.ValueOf(216),
		"cancelRecoverys": reflect.ValueOf(196),
		"changePassword": reflect.ValueOf(25),
		"changeRecoverys": reflect.ValueOf(203),
		"recoverys": reflect.ValueOf(208),
		"prayerOn": reflect.ValueOf(60),
		"prayerOff": reflect.ValueOf(254),
		"walkRequest": reflect.ValueOf(187),
		"walkAction": reflect.ValueOf(16),
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
		"BANANA":                   reflect.ValueOf(249),
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
		"DSTONE_AMULET_C":          reflect.ValueOf(597),
		"DSTONE_AMULET":            reflect.ValueOf(522),
		"DSTONE_AMULET_U":          reflect.ValueOf(522),
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
		"login": reflect.ValueOf(func(fn func(player *Player)) {
			LoginTriggers = append(LoginTriggers, fn)
		}),
		"invOnBoundary": reflect.ValueOf(func(fn func(player *Player, boundary *Object, item *Item) bool) {
			InvOnBoundaryTriggers = append(InvOnBoundaryTriggers, fn)
		}),
		"invOnPlayer": reflect.ValueOf(func(pred func(*Item) bool, fn func(player *Player, target *Player, item *Item)) {
			InvOnPlayerTriggers = append(InvOnPlayerTriggers, ItemOnPlayerTrigger{pred, fn})
		}),
		"invOnObject": reflect.ValueOf(func(fn func(player *Player, boundary *Object, item *Item) bool) {
			InvOnObjectTriggers = append(InvOnObjectTriggers, fn)
		}),
		"object": reflect.ValueOf(func(pred func(*Object, int) bool, fn func(player *Player, object *Object, click int)) {
			ObjectTriggers = append(ObjectTriggers, ObjectTrigger{pred, fn})
		}),
		"item": reflect.ValueOf(func(check func(item *Item) bool, fn func(player *Player, item *Item)) {
			ItemTriggers = append(ItemTriggers, ItemTrigger{check, fn})
		}),
		"boundary": reflect.ValueOf(func(pred func(*Object, int) bool, fn func(player *Player, object *Object, click int)) {
			BoundaryTriggers = append(BoundaryTriggers, ObjectTrigger{pred, fn})
		}),
		"npc": reflect.ValueOf(func(predicate func(npc *NPC) bool, fn func(player *Player, npc *NPC)) {
			NpcTalkList = append(NpcTalkList, NpcTrigger{predicate, fn})
		}),
		"spell": reflect.ValueOf(func(ident interface{}, fn func(player *Player, spell interface{})) {
			switch ident.(type) {
			case int64:
				SpellTriggers[int(ident.(int64))] = fn
			case int:
				SpellTriggers[int(ident.(int))] = fn
			default:
				log.Debugf("%v, %T", ident, ident)
			}
		}),
		"commands": reflect.ValueOf(CommandHandlers),
		"chatNpcs": reflect.ValueOf(&NpcTalkList),
		"sceneActions": reflect.ValueOf(&ObjectTriggers),
		"invSceneActions": reflect.ValueOf(&InvOnObjectTriggers),
		"boundaryActions": reflect.ValueOf(&BoundaryTriggers),
		"spells": reflect.ValueOf(SpellTriggers),
		"packet": reflect.ValueOf(func(ident interface{}, fn func(player *Player, packet interface{})) {
			switch ident.(type) {
			case int64:
				PacketTriggers[byte(ident.(int64))] = fn
			case int:
				PacketTriggers[byte(ident.(int))] = fn
			default:
				log.Debugf("%v, %T", ident, ident)
			}
		}),
		"npcAttack": reflect.ValueOf(func(pred NpcActionPredicate, fn func(player *Player, npc *NPC)) {
			NpcAtkTriggers = append(NpcAtkTriggers, NpcBlockingTrigger{pred, fn})
		}),
		"npcKilled": reflect.ValueOf(func(pred NpcActionPredicate, fn func(player *Player, npc *NPC)) {
			NpcDeathTriggers = append(NpcDeathTriggers, NpcBlockingTrigger{pred, fn})
		}),
		"command": reflect.ValueOf(func(name string, fn func(p *Player, args []string)) {
			CommandHandlers[name] = fn
		}),
	}
	env.Packages["log"] = map[string]reflect.Value{
		"print":  reflect.ValueOf(fmt.Println),
		"printf": reflect.ValueOf(fmt.Printf),

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

func ScriptEnv() *env.Env {
	e := env.NewEnv()
	parser.EnableErrorVerbose()
	e.Define("sleep", time.Sleep)
	e.Define("runAfter", time.AfterFunc)
	e.Define("after", time.After)
	e.Define("schedule", tasks.TickList.Schedule)
	e.Define("onTick", tasks.TickList.Add)
	e.Define("newChatMessage", NewChatMessage)
	e.Define("newProjectile", NewProjectile)
	e.Define("eventsPlayer", playerEvents)
	e.Define("eventsNpc", npcEvents)
	e.Define("Minute", TicksMinute)
	e.Define("Hour", TicksHour)
	e.Define("Second", time.Second)
	e.Define("Millisecond", time.Millisecond)
	e.Define("tNanos", time.Nanosecond)
	e.Define("ChatDelay", TickMillis*3)
	e.Define("encryptMsg", strutil.Encipher)
	e.Define("decryptMsg", strutil.Decipher)
	e.Define("ATTACK", entity.StatAttack)
	e.Define("DEFENSE", entity.StatDefense)
	e.Define("STRENGTH", entity.StatStrength)
	e.Define("HITPOINTS", entity.StatHits)
	e.Define("RANGED", entity.StatRanged)
	e.Define("PRAYER", entity.StatPrayer)
	e.Define("MAGIC", entity.StatMagic)
	e.Define("COOKING", entity.StatCooking)
	e.Define("WOODCUTTING", entity.StatWoodcutting)
	e.Define("FLETCHING", entity.StatFletching)
	e.Define("FISHING", entity.StatFishing)
	e.Define("FIREMAKING", entity.StatFiremaking)
	e.Define("CRAFTING", entity.StatCrafting)
	e.Define("SMITHING", entity.StatSmithing)
	e.Define("MINING", entity.StatMining)
	e.Define("HERBLAW", entity.StatHerblaw)
	e.Define("AGILITY", entity.StatAgility)
	e.Define("THIEVING", entity.StatThieving)
	e.Define("PRAYER_THICK_SKIN", 0)
	e.Define("PRAYER_BURST_OF_STRENGTH", 1)
	e.Define("PRAYER_CLARITY_OF_THOUGHT", 2)
	e.Define("PRAYER_ROCK_SKIN", 3)
	e.Define("PRAYER_SUPERHUMAN_STRENGTH", 4)
	e.Define("PRAYER_IMPROVED_REFLEXES", 5)
	e.Define("PRAYER_RAPID_RESTORE", 6)
	e.Define("PRAYER_RAPID_HEAL", 7)
	e.Define("PRAYER_PROTECT_ITEM", 8)
	e.Define("PRAYER_STEEL_SKIN", 9)
	e.Define("PRAYER_ULTIMATE_STRENGTH", 10)
	e.Define("PRAYER_INCREDIBLE_REFLEXES", 11)
	e.Define("PRAYER_PARALYZE_MONSTER", 12)
	e.Define("PRAYER_PROTECT_FROM_MISSILES", 13)
	e.Define("ZeroTime", time.Time{})
	e.Define("itemDefs", definitions.Items)
	e.Define("objectDefs", definitions.ScenaryObjects)
	e.Define("objectDef", definitions.Scenary)
	e.Define("boundaryDefs", definitions.BoundaryObjects)
	e.Define("npcDefs", definitions.Npcs)
	e.Define("lvlToExp", entity.LevelToExperience)
	e.Define("expToLvl", entity.ExperienceToLevel)
	e.Define("withinWorld", WithinWorld)
	e.Define("skillIndex", entity.SkillIndex)
	e.Define("skillName", entity.SkillName)
	e.Define("newNpc", NewNpc)
	e.Define("newObject", NewObject)
	e.Define("newPath", NewPathway)
	e.Define("base37", strutil.Base37.Encode)
	e.Define("fromBase37", strutil.Base37.Decode)
	e.Define("appraiseItem", func(shop *Shop, id int) int {
		realPrice := int(Price(definitions.Item(id).BasePrice).Scale(shop.BasePurchasePercent + shop.DeltaPercentModID(id)))
		return realPrice
	})

	e.Define("rand", func(low, high int) int {
		return int(rand.Rng.Float64()*float64(high-low+1))+low
	})
	e.Define("randExcl", func(low, high int) int {
		return int(rand.Rng.Float64()*float64(high-low))+low
	})
	e.Define("randIncl", func(low, high int) int {
		return int(rand.Rng.Float64()*float64(high-low+1))+low
	})
	e.Define("random", rand.Rng.Uint8)
	e.Define("randomWord", rand.Rng.Uint32)
	e.Define("randomLong", rand.Rng.Uint64)
	e.Define("randomFloat", rand.Rng.Float64)
	e.Define("randomBytes", rand.Rng.NextBytes)
	
	e.Define("NORTH", North)
	e.Define("NORTHEAST", NorthEast)
	e.Define("NORTHWEST", NorthWest)
	e.Define("SOUTH", South)
	e.Define("SOUTHEAST", SouthEast)
	e.Define("SOUTHWEST", SouthWest)
	e.Define("EAST", East)
	e.Define("WEST", West)
	e.Define("parseDirection", ParseDirection)
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
		return rand.Rng.Float64()*127.0+1.0 <= float64(cur)+40.0-float64(req)*1.5
	})
	e.Define("roll", Chance)
	e.Define("parseArgs", strutil.ParseArgs)
	e.Define("boundedRoll", BoundedChance)
	e.Define("weightedChance", WeightedChoice)
	e.Define("statRoll", Statistical)
	e.Define("CurTick", CurrentTick)
	e.Define("npcPredicate", func(ids ...interface{}) func(*NPC) bool {
		return func(npc *NPC) bool {
			for _, id := range ids {
				switch id.(type) {
				case string:
					if npc.Name() == id.(string) {
						return true
					}
				case int64:
					if npc.ID == int(id.(int64)) {
						return true
					}
				case int:
					if npc.ID == id.(int) {
						return true
					}
				}
			}
			return false
		}
	})
	e.Define("npcBlockingPredicate", func(ids ...interface{}) func(*Player, *NPC) bool {
		return func(player *Player, npc *NPC) bool {
			for _, id := range ids {
				switch id.(type) {
				case string:
					if npc.Name() == id.(string) {
						return true
					}
				case int64:
					if npc.ID == int(id.(int64)) {
						return true
					}
				case int:
					if npc.ID == id.(int) {
						return true
					}
				}
			}
			return false
		}
	})

	e.Define("fuzzyItem", func(input string) (itemList []map[string]interface{}) {
		for id, item := range definitions.Items {
			if fuzzy.MatchNormalized(strings.ToLower(input), strings.ToLower(item.Name)) {
				 rank := fuzzy.LevenshteinDistance(input, item.Name)
				itemList = append(itemList, map[string]interface{}{"name": item.Name, "id": id, "rank": rank})
				for idx := len(itemList)-1; idx > 0; idx-- {
					if itemList[idx]["rank"].(int) <= itemList[idx-1]["rank"].(int) {
						itemList[idx-1], itemList[idx] = itemList[idx], itemList[idx-1]
					}
				}
			}
		}

		
		return itemList
	})

	e.Define("itemPredicate", func(ids ...interface{}) func(*Item) bool {
		return func(item *Item) bool {
			for _, id := range ids {
				switch id.(type) {
				case string:
					if item.Command() == id.(string) || item.Name() == id.(string) {
						return true
					}
				case int64:
					if item.ID == int(id.(int64)) {
						return true
					}
				case int:
					if item.ID == id.(int) {
						return true
					}
				default:
					break
				}
			}
			return false
		}
	})
	e.Define("objectPredicate", func(ids ...interface{}) func(*Object, int) bool {
		return func(object *Object, click int) bool {
			for _, id := range ids {
				switch id.(type) {
				case string:
					if definitions.ScenaryObjects[object.ID].Commands[click] == id.(string) ||
							definitions.ScenaryObjects[object.ID].Name == id.(string) {
						return true
					}
				case int64:
					if object.ID == int(id.(int64)) {
						return true
					}
				case int:
					if object.ID == id.(int) {
						return true
					}
				default:
					break
				}
			}
			return false
		}
	})
	e.Define("toPlayer", AsPlayer)
	e.Define("toNpc", AsNpc)
	e = core.Import(e)
	return e
}
func init() {
	CommandHandlers["shutdown"] = func(player *Player, args []string) {
		wait := sync.WaitGroup{}
		Players.Range(func(p1 *Player) {
			wait.Add(1)
			go func() {
				defer wait.Done()
				p1.Message(serverPrefix + "Shutting down.")
				time.Sleep(TickMillis)
				p1.WriteNow(*Logout)
				p1.Destroy()
			}()
		})
		wait.Wait()
		os.Exit(1)
	}
	CommandHandlers["memdump"] = func(player *Player, args []string) {
		file, err := os.Create("rscgo.mprof")
		if err != nil {
			log.Warning.Println("Could not open file to dump memory profile:", err)
			player.Message(serverPrefix + "Error encountered opening profile output file.")
			return
		}
		err = pprof.WriteHeapProfile(file)
		if err != nil {
			log.Warning.Println("Could not write heap profile to file::", err)
			player.Message(serverPrefix + "Error encountered writing profile output file.")
			return
		}
		err = file.Close()
		if err != nil {
			log.Warning.Println("Could not close heap file::", err)
			player.Message(serverPrefix + "Error encountered closing profile output file.")
			return
		}
		log.Command(player.Username() + " dumped memory profile of the game to rscgo.mprof")
		player.Message(serverPrefix + "Dumped memory profile.")
	}
	CommandHandlers["cpudump"] = func(player *Player, args []string) {
		if len(args) < 1 {
			player.Message(serverPrefix + "Invalid args.  Usage: /pprof <start|stop>")
			return
		}
		switch args[0] {
		case "start":
			file, err := os.Create("rscgo.pprof")
			if err != nil {
				log.Warn("Could not open file to dump CPU profile:", err)
				player.Message(serverPrefix + "Error encountered opening profile output file.")
				return
			}
			err = pprof.StartCPUProfile(file)
			if err != nil {
				log.Warning.Println("Could not start CPU profile:", err)
				player.Message(serverPrefix + "Error encountered starting CPU profile.")
				return
			}
			log.Command(player.Username() + " began profiling CPU time.")
			player.Message(serverPrefix + "CPU profiling started.")
		case "stop":
			pprof.StopCPUProfile()
			log.Command(player.Username() + " has finished profiling CPU time, output should be in rscgo.pprof")
			player.Message(serverPrefix + "CPU profiling finished.")
		default:
			player.Message(serverPrefix + "Invalid args.  Usage: /pprof <start|stop>")
		}
	}
	CommandHandlers["run"] = func(player *Player, args []string) {
		line := strings.Join(args, " ")
		env := ScriptEnv()
		env.Define("p", player)
		env.Define("target", player.TargetMob)
		env.Define("npc", func() *NPC {
			return AsNpc(player.TargetMob())
		})
		env.Define("n", func() *NPC {
			return AsNpc(player.TargetMob())
		})
		env.Define("p1", func() *Player {
			return AsPlayer(player.TargetMob())
		})
		env.Define("player", player)
		env.Define("dataService", DefaultPlayerService)
		ret, err := vm.Execute(env, nil, "bind = import(\"bind\")\nworld = import(\"world\")\nlog = import(\"log\")\nids = import(\"ids\")\npackets = import(\"packets\")\nnet = import(\"net\")\nstate = import(\"state\")\n\n"+line)
		if err != nil {
			player.Message(serverPrefix + "Error: " + err.Error())
			log.Debug("Anko Error: " + err.Error())
			return
		}
		switch ret.(type) {
		case string:
			player.Message(serverPrefix + "string(" + ret.(string) + ")")
		case int64:
			player.Message(serverPrefix + "int(" + strconv.Itoa(int(ret.(int64))) + ")")
		case int:
			player.Message(serverPrefix + "int(" + strconv.Itoa(ret.(int)) + ")")
		case bool:
			if ret.(bool) {
				player.Message(serverPrefix + "bool(TRUE)")
			} else {
				player.Message(serverPrefix + "bool(FALSE)")
			}
		default:
			player.Message(serverPrefix + fmt.Sprintf("%v", ret))
		}
		log.Debugf("%v\n", ret)
	}
	CommandHandlers["reload"] = func(player *Player, args []string) {
		Clear()
		RunScripts()
		player.Message(serverPrefix + "Reloaded runtime content scripts from ./scripts/")
		log.Debugf("Triggers[\n\t%d item actions,\n\t%d scenary actions,\n\t%d boundary actions,\n\t%d npc actions,\n\t%d item->boundary actions,\n\t%d item->scenary actions,\n\t%d attacking NPC actions,\n\t%d killing NPC actions\n];\n", len(ItemTriggers), len(ObjectTriggers), len(BoundaryTriggers), len(NpcTalkList), len(InvOnBoundaryTriggers), len(InvOnObjectTriggers), len(NpcAtkTriggers), len(NpcDeathTriggers))
	}
}
