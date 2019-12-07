/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package packethandlers

import (
	"fmt"
	"github.com/mattn/anko/vm"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

func init() {
	PacketHandlers["command"] = func(player *world.Player, p *packet.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		handler, ok := script.CommandHandlers[args[0]]
		if !ok {
			player.Message("@que@Invalid command.")
			log.Commands.Printf("%v sent invalid command: /%v\n", player.Username, string(p.Payload))
			return
		}
		log.Commands.Printf("%v: /%v\n", player.Username, string(p.Payload))
		handler(player, args[1:])
	}
	script.CommandHandlers["memdump"] = func(player *world.Player, args []string) {
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
		log.Commands.Println(player.Username + " dumped memory profile of the server to rscgo.mprof")
		player.Message("Dumped memory profile.")
	}
	script.CommandHandlers["pprof"] = func(player *world.Player, args []string) {
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
			log.Commands.Println(player.Username + " began profiling CPU time.")
			player.Message("CPU profiling started.")
		case "stop":
			pprof.StopCPUProfile()
			log.Commands.Println(player.Username + " has finished profiling CPU time, output should be in rscgo.pprof")
			player.Message("CPU profiling finished.")
		default:
			player.Message("Invalid args.  Usage: /pprof <start|stop>")
		}
	}
	script.CommandHandlers["saveobjects"] = func(player *world.Player, args []string) {
		go func() {
			if count := db.SaveObjectLocations(); count > 0 {
				player.Message("Saved " + strconv.Itoa(count) + " game objects to world.db")
				log.Commands.Println(player.Username + " saved " + strconv.Itoa(count) + " game objects to world.db")
			} else {
				player.Message("Appears to have been an issue saving game objects to world.db.  Check server logs.")
				log.Commands.Println(player.Username + " failed to save game objects; count=" + strconv.Itoa(count))
			}
		}()
	}
	script.CommandHandlers["npc"] = func(player *world.Player, args []string) {
		if len(args) < 1 {
			player.Message("@que@Invalid args.  Usage: /npc <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id > 793 || id < 0 {
			player.Message("@que@Invalid args.  Usage: /npc <id>")
			return
		}

		x := player.X()
		y := player.Y()

		world.AddNpc(world.NewNpc(id, x, y, x-5, x+5, y-5, y+5))
	}
	script.CommandHandlers["anko"] = func(player *world.Player, args []string) {
		line := strings.Join(args, " ")
		env := script.WorldModule()
		env.Define("println", fmt.Println)
		env.Define("player", player)
		vm.Execute(env, nil, line)
	}
	script.CommandHandlers["reloadscripts"] = func(player *world.Player, args []string) {
		script.Clear()
		script.Load()
		player.Message(fmt.Sprintf("Bind[%d item, %d obj, %d bound, %d npc, %d invBound, %d invObject, %d npcAtk, %d npcKill]", len(script.ItemTriggers), len(script.ObjectTriggers), len(script.BoundaryTriggers), len(script.NpcTriggers), len(script.InvOnBoundaryTriggers), len(script.InvOnObjectTriggers), len(script.NpcAtkTriggers), len(script.NpcDeathTriggers)))
		log.Info.Printf("Bind[%d item, %d obj, %d bound, %d npc, %d invBound, %d invObject, %d npcAtk, %d npcKill] loaded\n", len(script.ItemTriggers), len(script.ObjectTriggers), len(script.BoundaryTriggers), len(script.NpcTriggers), len(script.InvOnBoundaryTriggers), len(script.InvOnObjectTriggers), len(script.NpcAtkTriggers), len(script.NpcDeathTriggers))
	}
}

func notYetImplemented(player *world.Player) {
	player.Message("@que@@ora@Not yet implemented")
}
