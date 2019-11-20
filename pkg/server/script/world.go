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
		"getPlayerCount": clients.Size,
		"getPlayer": clients.FromIndex,
		"getPlayerByName": clients.FromUserHash,
		"replaceObject": world.ReplaceObject,
		"getObjectAt": world.GetObject,
		"getNpc": world.GetNpc,
		"newPathway": world.NewPathwayToCoords,
		"newLocation": world.NewLocation,
		"checkCollisions": world.IsTileBlocking,
		"objectDefs": world.Objects,
		"objects": world.Npcs,
		"boundaryDefs": world.Boundarys,
		"npcDefs": world.NpcDefs,
		"npcs": world.Npcs,
		"itemDefs": world.ItemDefs,
		"commands": CommandHandlers,
		"addCommand": func(name string, fn func(args []string)) {
			CommandHandlers[name] = func(c clients.Client, args []string) {
				fn(args)
			}
		},
		"broadcast": func(fn func(interface{})) {
			clients.Range(func (c clients.Client) {
				fn(c)
			})
		},
		"announce": func(msg string) {
			clients.Range(func (c clients.Client) {
				c.SendPacket(packetbuilders.ServerMessage("@que@" + msg))
			})
		},
	}, map[string]interface{}{
		"client": reflect.TypeOf(clients.Client(nil)),
		"player": reflect.TypeOf(&world.Player{}),
		"object": reflect.TypeOf(&world.Object{}),
		"item": reflect.TypeOf(&world.Item{}),
		"groundItem": reflect.TypeOf(&world.GroundItem{}),
		"npc": reflect.TypeOf(&world.NPC{}),
		"location": reflect.TypeOf(world.Location{}),
	})
	if err != nil {
		log.Warning.Println("Error initializing VM parameters:", err)
		return nil
	}
	env, err = env.AddPackage("log", map[string]interface{}{
		"debug": log.Info.Println,
		"debugf": log.Info.Printf,
		"warn":  log.Warning.Println,
		"warnf": log.Warning.Printf,
		"err":   log.Error.Println,
		"errf": log.Error.Printf,
		"cheat": log.Suspicious.Println,
		"cheatf": log.Suspicious.Printf,
		"cmd":   log.Commands.Println,
		"cmdf": log.Commands.Printf,
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

