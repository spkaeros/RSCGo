package packethandlers

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
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
			player.SendPacket(packetbuilders.ServerMessage("@que@Invalid command."))
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
			player.SendPacket(packetbuilders.ServerMessage("Error encountered opening profile output file."))
			return
		}
		err = pprof.WriteHeapProfile(file)
		if err != nil {
			log.Warning.Println("Could not write heap profile to file::", err)
			player.SendPacket(packetbuilders.ServerMessage("Error encountered writing profile output file."))
			return
		}
		err = file.Close()
		if err != nil {
			log.Warning.Println("Could not close heap file::", err)
			player.SendPacket(packetbuilders.ServerMessage("Error encountered closing profile output file."))
			return
		}
		log.Commands.Println(player.Username + " dumped memory profile of the server to rscgo.mprof")
		player.SendPacket(packetbuilders.ServerMessage("Dumped memory profile."))
	}
	script.CommandHandlers["pprof"] = func(player *world.Player, args []string) {
		if len(args) < 1 {
			player.SendPacket(packetbuilders.ServerMessage("Invalid args.  Usage: /pprof <start|stop>"))
			return
		}
		switch args[0] {
		case "start":
			file, err := os.Create("rscgo.pprof")
			if err != nil {
				log.Warning.Println("Could not open file to dump CPU profile:", err)
				player.SendPacket(packetbuilders.ServerMessage("Error encountered opening profile output file."))
				return
			}
			err = pprof.StartCPUProfile(file)
			if err != nil {
				log.Warning.Println("Could not start CPU profile:", err)
				player.SendPacket(packetbuilders.ServerMessage("Error encountered starting CPU profile."))
				return
			}
			log.Commands.Println(player.Username + " began profiling CPU time.")
			player.SendPacket(packetbuilders.ServerMessage("CPU profiling started."))
		case "stop":
			pprof.StopCPUProfile()
			log.Commands.Println(player.Username + " has finished profiling CPU time, output should be in rscgo.pprof")
			player.SendPacket(packetbuilders.ServerMessage("CPU profiling finished."))
		default:
			player.SendPacket(packetbuilders.ServerMessage("Invalid args.  Usage: /pprof <start|stop>"))
		}
	}
	script.CommandHandlers["saveobjects"] = func(player *world.Player, args []string) {
		go func() {
			if count := db.SaveObjectLocations(); count > 0 {
				player.SendPacket(packetbuilders.ServerMessage("Saved " + strconv.Itoa(count) + " game objects to world.db"))
				log.Commands.Println(player.Username + " saved " + strconv.Itoa(count) + " game objects to world.db")
			} else {
				player.SendPacket(packetbuilders.ServerMessage("Appears to have been an issue saving game objects to world.db.  Check server logs."))
				log.Commands.Println(player.Username + " failed to save game objects; count=" + strconv.Itoa(count))
			}
		}()
	}
	script.CommandHandlers["npc"] = func(player *world.Player, args []string) {
		if len(args) < 1 {
			player.SendPacket(packetbuilders.ServerMessage("@que@Invalid args.  Usage: /npc <id>"))
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id > 793 || id < 0 {
			player.SendPacket(packetbuilders.ServerMessage("@que@Invalid args.  Usage: /npc <id>"))
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
		env.Execute(line)
	}
	script.CommandHandlers["reloadscripts"] = func(player *world.Player, args []string) {
		script.Clear()
		script.Load()
		log.Info.Printf("Loaded %d inventory, %d object, %d boundary, and %d NPC action triggers.\n", len(script.InvTriggers), len(script.ObjectTriggers), len(script.BoundaryTriggers), len(script.NpcTriggers))
	}
}

func notYetImplemented(player *world.Player, args []string) {
	player.SendPacket(packetbuilders.ServerMessage("@que@@ora@Not yet implemented"))
}
