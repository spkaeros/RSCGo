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
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/anko/vm"
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	game.AddHandler("command", func(player *world.Player, p *net.Packet) {
		raw := string(p.FrameBuffer[:len(p.FrameBuffer)-1])
		args := strutil.ParseArgs(raw)
		// prevent `::` freezing player
		if len(args) <= 0 {
			return
		}
		handler, ok := world.CommandHandlers[strings.ToLower(args[0])]
		if !ok {
			player.Message("@que@Command not found.  Double check your spelling, and try again.")
			log.Command("%v sent invalid command: ::%v\n", player.Username(), strings.ToLower(args[0]))
			return
		}
		log.Commandf("%v: ::%v\n", player.Username(), raw)
		handler(player, args[1:])
	})
	world.CommandHandlers["shutdown"] = func(player *world.Player, args []string) {
		var wg sync.WaitGroup
		world.Players.Range(func(p1 *world.Player) {
			go func() {
				defer wg.Done()
				p1.Message("Shutting down.")
				wg.Add(1)
				p1.Destroy()
			}()
		})
		wg.Wait()
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
	world.CommandHandlers["memdump"] = func(player *world.Player, args []string) {
		file, err := os.Create("rscgo.mprof")
		if err != nil {
			log.Warning.Println("Could not open file to dump memory profile:", err)
			player.Message("Error encountered opening profile output file.")
			return
		}
		err = pprof.WriteHeapProfile(file)
		if err != nil {
			log.Warning.Println("Could not write heap profile to file::", err)
			player.Message("Error encountered writing profile output file.")
			return
		}
		err = file.Close()
		if err != nil {
			log.Warning.Println("Could not close heap file::", err)
			player.Message("Error encountered closing profile output file.")
			return
		}
		log.Commands.Println(player.Username() + " dumped memory profile of the game to rscgo.mprof")
		player.Message("Dumped memory profile.")
	}
	world.CommandHandlers["pprof"] = func(player *world.Player, args []string) {
		if len(args) < 1 {
			player.Message("Invalid args.  Usage: /pprof <start|stop>")
			return
		}
		switch args[0] {
		case "start":
			file, err := os.Create("rscgo.pprof")
			if err != nil {
				log.Warning.Println("Could not open file to dump CPU profile:", err)
				player.Message("Error encountered opening profile output file.")
				return
			}
			err = pprof.StartCPUProfile(file)
			if err != nil {
				log.Warning.Println("Could not start CPU profile:", err)
				player.Message("Error encountered starting CPU profile.")
				return
			}
			log.Commands.Println(player.Username() + " began profiling CPU time.")
			player.Message("CPU profiling started.")
		case "stop":
			pprof.StopCPUProfile()
			log.Commands.Println(player.Username() + " has finished profiling CPU time, output should be in rscgo.pprof")
			player.Message("CPU profiling finished.")
		default:
			player.Message("Invalid args.  Usage: /pprof <start|stop>")
		}
	}
	world.CommandHandlers["run"] = func(player *world.Player, args []string) {
		line := strings.Join(args, " ")
		env := world.ScriptEnv()
		env.Define("p", player)
		env.Define("target", player.TargetMob())
		env.Define("player", player)
		ret, err := vm.Execute(env, nil, "bind = import(\"bind\")\nworld = import(\"world\")\nlog = import(\"log\")\nids = import(\"ids\")\n\n"+line)
		if err != nil {
			player.Message("Error: " + err.Error())
			log.Info.Println("Anko Error: " + err.Error())
			return
		}
		switch ret.(type) {
		case string:
			player.Message(ret.(string))
		case int64:
			player.Message("int(" + strconv.Itoa(int(ret.(int64))) + ")")
		case int:
			player.Message("int(" + strconv.Itoa(ret.(int)) + ")")
		case bool:
			if ret.(bool) {
				player.Message("TRUE")
			} else {
				player.Message("FALSE")
			}
		default:
			player.Message(fmt.Sprintf("%v", ret))
		}
		log.Info.Println(ret)
	}
	world.CommandHandlers["reload"] = func(player *world.Player, args []string) {
		world.Clear()
		world.RunScripts()
		player.Message("Reloaded ./scripts/**.ank from working directory.")
		player.Message(fmt.Sprintf("Bind[%d item, %d obj, %d bound, %d npc, %d invBound, %d invObject, %d npcAtk, %d npcKill]", len(world.ItemTriggers), len(world.ObjectTriggers), len(world.BoundaryTriggers), len(world.NpcTriggers), len(world.InvOnBoundaryTriggers), len(world.InvOnObjectTriggers), len(world.NpcAtkTriggers), len(world.NpcDeathTriggers)))
		log.Info.Printf("Bind[%d item, %d obj, %d bound, %d npc, %d invBound, %d invObject, %d npcAtk, %d npcKill] loaded\n", len(world.ItemTriggers), len(world.ObjectTriggers), len(world.BoundaryTriggers), len(world.NpcTriggers), len(world.InvOnBoundaryTriggers), len(world.InvOnObjectTriggers), len(world.NpcAtkTriggers), len(world.NpcDeathTriggers))
	}
}

func notYetImplemented(player *world.Player) {
	player.Message("@que@@ora@Not yet implemented")
}
