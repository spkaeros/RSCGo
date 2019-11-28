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
	"os"
	"reflect"
	"time"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(*world.Player, []string))

var UpdateTime time.Time

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
		"tileData":        world.ClipData,
		"objectDefs":      world.Objects,
		"objects":         world.Npcs,
		"boundaryDefs":    world.Boundarys,
		"npcDefs":         world.NpcDefs,
		"npcs":            world.Npcs,
		"itemDefs":        world.ItemDefs,
		"commands":        CommandHandlers,
		"kick": func(client *world.Player) {
			client.SendPacket(world.Logout)
			client.Destroy()
		},
		"addCommand": func(name string, fn func(p *world.Player, args []string)) {
			CommandHandlers[name] = func(player *world.Player, args []string) {
				fn(player, args)
			}
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
		"rand": rand.Int31N,
		"walkTo": func(target *world.Player, x, y int) {
			target.WalkTo(world.NewLocation(x, y))
		},
		"systemUpdate": func(t int) {
			UpdateTime = time.Now().Add(time.Second * time.Duration(t))
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
				player.SendPacket(world.SystemUpdate(t))
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
			return world.ParseSkill(name)
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
	env, err = env.AddPackage("packets", map[string]interface{}{
		"BigInformationBox":     world.BigInformationBox,
		"BoundaryLocations":     world.BoundaryLocations,
		"CannotLogout":          world.CannotLogout,
		"OpenChangeAppearance":      world.OpenChangeAppearance,
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
