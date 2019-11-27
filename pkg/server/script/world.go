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
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"reflect"
	"time"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(clients.Client, []string))

func WorldModule() *vm.Env {
	env, err := vm.NewEnv().AddPackage("world", map[string]interface{}{
		"getPlayerCount":  clients.Size,
		"getPlayers":      clients.Clients,
		"getPlayer":       clients.FromIndex,
		"getPlayerByName": clients.FromUserHash,
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
		"newLocation":     world.NewLocation,
		"checkCollisions": world.IsTileBlocking,
		"objectDefs":      world.Objects,
		"objects":         world.Npcs,
		"boundaryDefs":    world.Boundarys,
		"npcDefs":         world.NpcDefs,
		"npcs":            world.Npcs,
		"itemDefs":        world.ItemDefs,
		"commands":        CommandHandlers,
		"addCommand": func(name string, fn func(p *world.Player, args []string)) {
			CommandHandlers[name] = func(c clients.Client, args []string) {
				fn(c.Player(), args)
			}
		},
		"broadcast": func(fn func(interface{})) {
			clients.Range(func(c clients.Client) {
				fn(c)
			})
		},
		"announce": func(msg string) {
			clients.Range(func(c clients.Client) {
				c.SendPacket(packetbuilders.ServerMessage("@que@" + msg))
			})
		},
		"parseDirection": world.ParseDirection,
		"North":          world.North,
		"South":          world.South,
		"East":           world.East,
		"West":           world.West,
		"NorthWest":      world.NorthWest,
		"NorthEast":      world.NorthEast,
		"SouthWest":      world.SouthWest,
		"SouthEast":      world.SouthEast,
		"ATTACK":         world.StatAttack,
		"DEFENSE":        world.StatDefense,
		"STRENGTH":       world.StatStrength,
		"HITPOINTS":      world.StatHits,
		"RANGED":         world.StatRanged,
		"PRAYER":         world.StatPrayer,
		"MAGIC":          world.StatMagic,
		"COOKING":        world.StatCooking,
		"WOODCUTTING":    world.StatWoodcutting,
		"FLETCHING":      world.StatFletching,
		"FISHING":        world.StatFishing,
		"FIREMAKING":     world.StatFiremaking,
		"CRAFTING":       world.StatCrafting,
		"SMITHING":       world.StatSmithing,
		"MINIG":          world.StatMining,
		"HERBLAW":        world.StatHerblaw,
		"AGILITY":        world.StatAgility,
		"THIEVING":       world.StatThieving,
		"IDLE":           world.MSIdle,
		"BUSY":           world.MSBusy,
		"MENUCHOOSING":   world.MSMenuChoosing,
		"teleport": func(player *world.Player, x, y int) {
			player.Teleport(x, y)
		},
		"openOptionMenu": func(player *world.Player, options ...string) int {
			player.SendPacket(packetbuilders.OptionMenuOpen(options...))
			player.State = world.MSMenuChoosing
			select {
			case reply := <-player.OptionMenuC:
				if reply < 0 || int(reply) > len(options) {
					return -1
				}

				for _, player2 := range player.NearbyPlayers() {
					player2.SendPacket(packetbuilders.PlayerMessage(player, options[reply]))
				}
				player.SendPacket(packetbuilders.PlayerMessage(player, options[reply]))
				time.Sleep(time.Millisecond * 1800)
				return int(reply)
			case <-time.After(time.Second * 10):
				return -1
			}
		},
		"closeOptionMenu": func(player *world.Player, questions ...string) {
			player.SendPacket(packetbuilders.OptionMenuClose)
		},
		"rand": rand.Int31N,
		"npcChat": func(sender *world.NPC, target *world.Player, msgs ...string) {
			for _, msg := range msgs {
				for _, player := range target.NearbyPlayers() {
					player.SendPacket(packetbuilders.NpcMessage(sender, msg, target))
				}
				target.SendPacket(packetbuilders.NpcMessage(sender, msg, target))
				//sender.ChatTarget = target.Index
				//sender.ChatMessage = msg
				time.Sleep(time.Millisecond * 1800)
			}
		},
		"playerChat": func(sender *world.Player, msg ...string) {
			for _, s := range msg {
				for _, player := range sender.NearbyPlayers() {
					player.SendPacket(packetbuilders.PlayerMessage(sender, s))
				}
				sender.SendPacket(packetbuilders.PlayerMessage(sender, s))
				time.Sleep(time.Millisecond * 1800)
			}
		},
		"playerDamage": func(target *world.Player, damage int) {
			for _, player := range target.NearbyPlayers() {
				player.SendPacket(packetbuilders.PlayerDamage(target, damage))
			}
			target.SendPacket(packetbuilders.PlayerDamage(target, damage))
		},
		"sendSound": func(target *world.Player, sound string) {
			target.SendPacket(packetbuilders.Sound(sound))
		},
		"sendStats": func(target *world.Player) {
			target.SendPacket(packetbuilders.PlayerStats(target))
		},
		"sendStat": func(target *world.Player, idx int) {
			target.SendPacket(packetbuilders.PlayerStat(target, idx))
		},
		"sendInventory": func(target *world.Player) {
			target.SendPacket(packetbuilders.InventoryItems(target))
		},
		"sendDeath": func(target *world.Player) {
			target.SendPacket(packetbuilders.Death)
		},
		"sendPlane": func(target *world.Player) {
			target.SendPacket(packetbuilders.PlaneInfo(target))
		},
		"sendMessage": func(target *world.Player, msg string) {
			target.SendPacket(packetbuilders.ServerMessage(msg))
		},
		"getSkillIndex": func(name string) int {
			return world.ParseSkill(name)
		},
		"expToLvl": world.ExperienceToLevel,
		"lvlToExp": world.LevelToExperience,
	}, map[string]interface{}{
		"clientMap":  reflect.TypeOf(clients.Clients),
		"client":     reflect.TypeOf(clients.Client(nil)),
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
		"BigInformationBox":     packetbuilders.BigInformationBox,
		"BoundaryLocations":     packetbuilders.BoundaryLocations,
		"CannotLogout":          packetbuilders.CannotLogout,
		"ChangeAppearance":      packetbuilders.ChangeAppearance,
		"ClientSettings":        packetbuilders.ClientSettings,
		"Death":                 packetbuilders.Death,
		"DefaultActionMessage":  packetbuilders.DefaultActionMessage,
		"EquipmentStats":        packetbuilders.EquipmentStats,
		"Fatigue":               packetbuilders.Fatigue,
		"FightMode":             packetbuilders.FightMode,
		"FriendList":            packetbuilders.FriendList,
		"FriendUpdate":          packetbuilders.FriendUpdate,
		"IgnoreList":            packetbuilders.IgnoreList,
		"InventoryItems":        packetbuilders.InventoryItems,
		"ItemLocations":         packetbuilders.ItemLocations,
		"LoginBox":              packetbuilders.LoginBox,
		"LoginResponse":         packetbuilders.LoginResponse,
		"Logout":                packetbuilders.Logout,
		"NPCPositions":          packetbuilders.NPCPositions,
		"NpcDamage":             packetbuilders.NpcDamage,
		"ObjectLocations":       packetbuilders.ObjectLocations,
		"PlaneInfo":             packetbuilders.PlaneInfo,
		"PlayerAppearances":     packetbuilders.PlayerAppearances,
		"PlayerChat":            packetbuilders.PlayerChat,
		"PlayerDamage":          packetbuilders.PlayerDamage,
		"PlayerPositions":       packetbuilders.PlayerPositions,
		"PlayerStat":            packetbuilders.PlayerStat,
		"PlayerStats":           packetbuilders.PlayerStats,
		"PrivacySettings":       packetbuilders.PrivacySettings,
		"PrivateMessage":        packetbuilders.PrivateMessage,
		"ResponsePong":          packetbuilders.ResponsePong,
		"ServerInfo":            packetbuilders.ServerInfo,
		"ServerMessage":         packetbuilders.ServerMessage,
		"TeleBubble":            packetbuilders.TeleBubble,
		"TradeAccept":           packetbuilders.TradeAccept,
		"TradeClose":            packetbuilders.TradeClose,
		"TradeConfirmationOpen": packetbuilders.TradeConfirmationOpen,
		"TradeOpen":             packetbuilders.TradeOpen,
		"TradeTargetAccept":     packetbuilders.TradeTargetAccept,
		"TradeUpdate":           packetbuilders.TradeUpdate,
		"WelcomeMessage":        packetbuilders.WelcomeMessage,
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
	err = env.DefineGlobal("tNanos", time.Nanosecond)
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	core.Import(env)
	packages.DefineImport(env)
	return env
}
