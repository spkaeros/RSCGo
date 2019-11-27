package packethandlers

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/clients"
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
	PacketHandlers["command"] = func(c clients.Client, p *packet.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		handler, ok := script.CommandHandlers[args[0]]
		if !ok {
			c.Message("@que@Invalid command.")
			log.Commands.Printf("%v sent invalid command: /%v\n", c.Player().Username, string(p.Payload))
			return
		}
		log.Commands.Printf("%v: /%v\n", c.Player().Username, string(p.Payload))
		handler(c, args[1:])
	}
	script.CommandHandlers["memdump"] = func(c clients.Client, args []string) {
		file, err := os.Create("rscgo.mprof")
		if err != nil {
			log.Warning.Println("Could not open file to dump memory profile:", err)
			c.Message("Error encountered opening profile output file.")
			return
		}
		err = pprof.WriteHeapProfile(file)
		if err != nil {
			log.Warning.Println("Could not write heap profile to file::", err)
			c.Message("Error encountered writing profile output file.")
			return
		}
		err = file.Close()
		if err != nil {
			log.Warning.Println("Could not close heap file::", err)
			c.Message("Error encountered closing profile output file.")
			return
		}
		log.Commands.Println(c.Player().Username + " dumped memory profile of the server to rscgo.mprof")
		c.Message("Dumped memory profile.")
	}
	script.CommandHandlers["pprof"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("Invalid args.  Usage: /pprof <start|stop>")
			return
		}
		switch args[0] {
		case "start":
			file, err := os.Create("rscgo.pprof")
			if err != nil {
				log.Warning.Println("Could not open file to dump CPU profile:", err)
				c.Message("Error encountered opening profile output file.")
				return
			}
			err = pprof.StartCPUProfile(file)
			if err != nil {
				log.Warning.Println("Could not start CPU profile:", err)
				c.Message("Error encountered starting CPU profile.")
				return
			}
			log.Commands.Println(c.Player().Username + " began profiling CPU time.")
			c.Message("CPU profiling started.")
		case "stop":
			pprof.StopCPUProfile()
			log.Commands.Println(c.Player().Username + " has finished profiling CPU time, output should be in rscgo.pprof")
			c.Message("CPU profiling finished.")
		default:
			c.Message("Invalid args.  Usage: /pprof <start|stop>")
		}
	}
	script.CommandHandlers["saveobjects"] = func(c clients.Client, args []string) {
		go func() {
			if count := db.SaveObjectLocations(); count > 0 {
				c.Message("Saved " + strconv.Itoa(count) + " game objects to world.db")
				log.Commands.Println(c.Player().Username + " saved " + strconv.Itoa(count) + " game objects to world.db")
			} else {
				c.Message("Appears to have been an issue saving game objects to world.db.  Check server logs.")
				log.Commands.Println(c.Player().Username + " failed to save game objects; count=" + strconv.Itoa(count))
			}
		}()
	}
	script.CommandHandlers["npc"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /npc <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id > 793 || id < 0 {
			c.Message("@que@Invalid args.  Usage: /npc <id>")
			return
		}

		x := c.Player().X()
		y := c.Player().Y()

		world.AddNpc(world.NewNpc(id, x, y, x-5, x+5, y-5, y+5))
	}
	script.CommandHandlers["anko"] = func(c clients.Client, args []string) {
		line := strings.Join(args, " ")
		env := script.WorldModule()
		env.Define("println", fmt.Println)
		env.Define("player", c.Player())
		env.Execute(line)
	}
	script.CommandHandlers["reloadscripts"] = func(c clients.Client, args []string) {
		script.Clear()
		script.Load()
		log.Info.Printf("Loaded %d inventory, %d object, %d boundary, and %d NPC action triggers.\n", len(script.InvTriggers), len(script.ObjectTriggers), len(script.BoundaryTriggers), len(script.NpcTriggers))
	}
}

func notYetImplemented(c clients.Client, args []string) {
	c.Message("@que@@ora@Not yet implemented")
}
