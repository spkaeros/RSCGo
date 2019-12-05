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
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/packages"
	"github.com/mattn/anko/vm"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"math"
	"os"
	"reflect"
	"time"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(*world.Player, []string))

func WorldModule() *vm.Env {
	env, err := vm.NewEnv().AddPackage("world", map[string]interface{}{
		"getPlayerCount":  players.Size,
		"getPlayers":      players.Players,
		"getPlayer":       players.FromIndex,
		"getPlayerByName": players.FromUserHash,
		"replaceObject":   world.ReplaceObject,
		"addObject":       world.AddObject,
		"removeObject":    world.RemoveObject,
		"addNpc":          world.AddNpc,
		"removeNpc":       world.RemoveNpc,
		"addItem":         world.AddItem,
		"removeItem":      world.RemoveItem,
		"getObjectAt":     world.GetObject,
		"newObject":       world.NewObject,
		"getNpc":          world.GetNpc,
		"newPathway":      world.NewPathwayToCoords,
		"newNpc":          world.NewNpc,
		"newLocation":     world.NewLocation,
		"checkCollisions": world.IsTileBlocking,
		"tileData":        world.CollisionData,
		"objectDefs":      world.Objects,
		"objects":         world.Npcs,
		"boundaryDefs":    world.BoundaryDefs,
		"npcDefs":         world.NpcDefs,
		"npcs":            world.Npcs,
		"itemDefs":        world.ItemDefs,
		"commands":        CommandHandlers,
		"kick": func(client *world.Player) {
			client.SendPacket(world.Logout)
			client.Destroy()
		},
		"updateStarted": func() bool {
			return !world.UpdateTime.IsZero()
		},
		"broadcast": func(fn func(interface{})) {
			players.Range(func(player *world.Player) {
				fn(player)
			})
		},
		"announce": func(msg string) {
			players.Range(func(player *world.Player) {
				player.Message("@que@" + msg)
			})
		},
		"parseDirection":     world.ParseDirection,
		"North":              world.North,
		"South":              world.South,
		"East":               world.East,
		"West":               world.West,
		"NorthWest":          world.NorthWest,
		"NorthEast":          world.NorthEast,
		"SouthWest":          world.SouthWest,
		"SouthEast":          world.SouthEast,
		"ATTACK":             world.StatAttack,
		"DEFENSE":            world.StatDefense,
		"STRENGTH":           world.StatStrength,
		"HITPOINTS":          world.StatHits,
		"RANGED":             world.StatRanged,
		"PRAYER":             world.StatPrayer,
		"MAGIC":              world.StatMagic,
		"COOKING":            world.StatCooking,
		"WOODCUTTING":        world.StatWoodcutting,
		"FLETCHING":          world.StatFletching,
		"FISHING":            world.StatFishing,
		"FIREMAKING":         world.StatFiremaking,
		"CRAFTING":           world.StatCrafting,
		"SMITHING":           world.StatSmithing,
		"MINIG":              world.StatMining,
		"HERBLAW":            world.StatHerblaw,
		"AGILITY":            world.StatAgility,
		"THIEVING":           world.StatThieving,
		"IDLE":               world.MSIdle,
		"BUSY":               world.MSBusy,
		"MENUCHOOSING":       world.MSOptionMenu,
		"CHATTING":           world.MSChatting,
		"BANKING":            world.MSBanking,
		"TRADING":            world.MSTrading,
		"DUELING":            world.MSDueling,
		"FIGHTING":           world.MSFighting,
		"BATCHING":           world.MSBatching,
		"SLEEPING":           world.MSSleeping,
		"CHANGINGAPPEARANCE": world.MSChangingAppearance,
		"rand":               rand.Int31N,
		"walkTo": func(target *world.Player, x, y int) {
			target.WalkTo(world.NewLocation(x, y))
		},
		"gatheringSuccess": func(req, cur int) bool {
			roll := float64(rand.Int31N(1, 128))
			if cur < req {
				return false
			}
			threshold := math.Max(float64(1), float64(cur) * 40 - float64(req) * 1.5)
			if 127 < threshold {
				threshold = 127
			}
			return roll <= threshold
		},
		"systemUpdate": func(t int) {
			world.UpdateTime = time.Now().Add(time.Second * time.Duration(t))
			go func() {
				time.Sleep(time.Second * time.Duration(t))
				players.Range(func(player *world.Player) {
					player.SendPacket(world.Logout)
					player.Destroy()
				})
				time.Sleep(300 * time.Millisecond)
				os.Exit(200)
			}()
			players.Range(func(player *world.Player) {
				player.SendUpdateTimer()
			})
		},
		"base37": strutil.Base37.Encode,
		"teleport": func(target *world.Player, x, y int, bubble bool) {
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
		},
		"getSkillIndex": func(name string) int {
			return world.SkillIndex(name)
		},
		"expToLvl":    world.ExperienceToLevel,
		"lvlToExp":    world.LevelToExperience,
		"withinWorld": world.WithinWorld,
	}, map[string]interface{}{
		"clientMap":  reflect.TypeOf(players.Players),
		"player":     reflect.TypeOf(&world.Player{}),
		"object":     reflect.TypeOf(&world.Object{}),
		"item":       reflect.TypeOf(&world.Item{}),
		"groundItem": reflect.TypeOf(&world.GroundItem{}),
		"npc":        reflect.TypeOf(&world.NPC{}),
		"location":   reflect.TypeOf(world.Location{}),
	})
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	env, err = env.AddPackage("ids", map[string]interface{}{
		"COOKEDMEAT":    132,
		"BURNTMEAT":     134,
		"RAW_RAT_MEAT":  503,
		"RAW_SHRIMP":  349,
		"WOODEN_SHIELD": 4,
		"BRONZE_LSWORD": 70,
		"NET": 376,
		"BRONZE_PICKAXE": 156,
		"SLEEPING_BAG": 1263,
		"NEEDLE": 39,
		"THREAD": 43,
		"FIRE_RUNE": 31,
		"WATER_RUNE": 32,
		"AIR_RUNE": 33,
		"EARTH_RUNE": 34,
		"MIND_RUNE": 35,
		"BODY_RUNE": 36,
		"LIFE_RUNE": 37,
		"DEATH_RUNE": 38,
		"NATURE_RUNE": 40,
		"CHAOS_RUNE": 41,
		"LAW_RUNE": 42,
		"COSMIC_RUNE": 46,
		"BLOOD_RUNE": 619,
		"AIR_STAFF": 101,
		"WATER_STAFF": 102,
		"EARTH_STAFF": 103,
		"FIRE_STAFF": 197,
		"FIRE_BATTLESTAFF": 615,
		"WATER_BATTLESTAFF": 616,
		"AIR_BATTLESTAFF": 617,
		"EARTH_BATTLESTAFF": 618,
		"E_FIRE_BATTLESTAFF": 682,
		"E_WATER_BATTLESTAFF": 683,
		"E_AIR_BATTLESTAFF": 684,
		"E_EARTH_BATTLESTAFF": 685,
		"BONES": 20,
		"BAT_BONES": 604,
		"DRAGON_BONES": 614,
		"RUNE_2H": 81,
		"RUNE_CHAIN": 400,
		"RUNE_PLATEBODY": 401,
		"RUNE_PLATETOP": 407,
		"DRAGON_SWORD": 593,
		"DRAGON_AXE": 594,
		"CHARGED_DSTONE_AMMY": 597,
		"DRAGON_HELMET": 795,
		"DRAGON_SHIELD": 1278,
		"EASTER_EGG": 677,
		"CHRISTMAS_CRACKER": 575,
		"PARTYHAT_RED": 576,
		"PARTYHAT_YELLOW": 577,
		"PARTYHAT_BLUE": 578,
		"PARTYHAT_GREEN": 579,
		"PARTYHAT_PINK": 580,
		"PARTYHAT_WHITE": 581,
		"GREEN_MASK": 828,
		"RED_MASK": 831,
		"BLUE_MASK": 832,
		"SANTA_HAT": 971,
		"PRESENT": 980,
		"GNOME_BALL": 981,
	}, nil)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	env, err = env.AddPackage("packets", map[string]interface{}{
		"BigInformationBox":     world.BigInformationBox,
		"BoundaryLocations":     world.BoundaryLocations,
		"CannotLogout":          world.CannotLogout,
		"OpenChangeAppearance":  world.OpenChangeAppearance,
		"ClientSettings":        world.ClientSettings,
		"Death":                 world.Death,
		"DefaultActionMessage":  world.DefaultActionMessage,
		"EquipmentStats":        world.EquipmentStats,
		"Fatigue":               world.Fatigue,
		"FightMode":             world.FightMode,
		"FriendList":            world.FriendList,
		"FriendUpdate":          world.FriendUpdate,
		"IgnoreList":            world.IgnoreList,
		"InventoryItems":        world.InventoryItems,
		"ItemLocations":         world.ItemLocations,
		"LoginBox":              world.LoginBox,
		"LoginResponse":         world.LoginResponse,
		"Logout":                world.Logout,
		"NPCPositions":          world.NPCPositions,
		"NpcDamage":             world.NpcDamage,
		"ObjectLocations":       world.ObjectLocations,
		"PlaneInfo":             world.PlaneInfo,
		"PlayerAppearances":     world.PlayerAppearances,
		"PlayerChat":            world.PlayerChat,
		"PlayerDamage":          world.PlayerDamage,
		"PlayerPositions":       world.PlayerPositions,
		"PlayerStat":            world.PlayerStat,
		"PlayerStats":           world.PlayerStats,
		"PrivacySettings":       world.PrivacySettings,
		"PrivateMessage":        world.PrivateMessage,
		"ResponsePong":          world.ResponsePong,
		"ServerMessage":         world.ServerMessage,
		"TeleBubble":            world.TeleBubble,
		"TradeAccept":           world.TradeAccept,
		"TradeClose":            world.TradeClose,
		"TradeConfirmationOpen": world.TradeConfirmationOpen,
		"TradeOpen":             world.TradeOpen,
		"TradeTargetAccept":     world.TradeTargetAccept,
		"TradeUpdate":           world.TradeUpdate,
		"WelcomeMessage":        world.WelcomeMessage,
	}, nil)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	env, err = env.AddPackage("bind", map[string]interface{}{
		"onLogin": func(fn func(player *world.Player)) {
			LoginTriggers = append(LoginTriggers, fn)
		},
		"invOnBoundary": func(fn func(player *world.Player, boundary *world.Object, item *world.Item) bool) {
			InvOnBoundaryTriggers = append(InvOnBoundaryTriggers, fn)
		},
		"invOnObject": func(fn func(player *world.Player, boundary *world.Object, item *world.Item) bool) {
			InvOnObjectTriggers = append(InvOnObjectTriggers, fn)
		},
		"npc": func(ident interface{}, fn func(player *world.Player, npc *world.NPC)) {
			if id, ok := ident.(int64); ok {
				NpcTriggers[id] = fn
			}
			if ids, ok := ident.([]interface{}); ok {
				for _, id := range ids {
					NpcTriggers[id.(int64)] = fn
				}
			}
			if name, ok := ident.(string); ok {
				NpcTriggers[name] = fn
			}
			//			NpcTriggers[id] = fn
		},
		"npcAttack": func(ident interface{}, fn func(player *world.Player, npc *world.NPC) bool) {
			log.Info.Printf("%T", ident)
			if id, ok := ident.(int64); ok {
				NpcAtkTriggers[id] = fn
			}
			if ids, ok := ident.([]interface{}); ok {
				for _, id := range ids {
					NpcAtkTriggers[id.(int64)] = fn
				}
			}
			if name, ok := ident.(string); ok {
				NpcAtkTriggers[name] = fn
			}
			//			NpcAtkTriggers[id] = fn
		},
		"npcKilled": func(ident interface{}, fn func(player *world.Player, npc *world.NPC)) {
			log.Info.Printf("%T", ident)
			if id, ok := ident.(int64); ok {
				NpcDeathTriggers[id] = fn
			}
			if ids, ok := ident.([]interface{}); ok {
				for _, id := range ids {
					NpcDeathTriggers[id.(int64)] = fn
				}
			}
			if name, ok := ident.(string); ok {
				NpcDeathTriggers[name] = fn
			}
			//			NpcDeathTriggers[id] = fn
		},
		"command": func(name string, fn func(p *world.Player, args []string)) {
			CommandHandlers[name] = fn
		},
	}, nil)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	env, err = env.AddPackage("log", map[string]interface{}{
		"debug":  log.Info.Println,
		"debugf": log.Info.Printf,
		"warn":   log.Warning.Println,
		"warnf":  log.Warning.Printf,
		"err":    log.Error.Println,
		"errf":   log.Error.Printf,
		"cheat":  log.Suspicious.Println,
		"cheatf": log.Suspicious.Printf,
		"cmd":    log.Commands.Println,
		"cmdf":   log.Commands.Printf,
	}, nil)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("sleep", time.Sleep)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("runAfter", time.AfterFunc)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("after", time.After)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("tMinute", time.Second*60)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("tHour", time.Second*60*60)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("tSecond", time.Second)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("tMillis", time.Millisecond)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("ChatDelay", time.Millisecond*1800)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	err = env.DefineGlobal("tNanos", time.Nanosecond)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	core.Import(env)
	packages.DefineImport(env)
	return env
}
