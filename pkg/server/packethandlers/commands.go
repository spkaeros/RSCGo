package packethandlers

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/clients"
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
	"time"
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
	script.CommandHandlers["kick"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /kick <player>")
			return
		}
		var (
			affectedClient clients.Client
			ok             bool
		)
		if pID, err := strconv.Atoi(args[0]); err == nil {
			affectedClient, ok = clients.FromIndex(pID)
		} else {
			affectedClient, ok = clients.FromUserHash(strutil.Base37.Encode(strings.Join(args, " ")))
		}
		if affectedClient == nil || !ok {
			c.Message("@que@Could not find Player().")
			return
		}
		log.Commands.Printf("'%v' kicked other player '%v'\n", c.Player().Username, affectedClient.Player().Username)
		c.Message("@que@Kicked: '" + affectedClient.Player().Username + "'")
		affectedClient.SendPacket(packetbuilders.Logout)
		affectedClient.Destroy()
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

		x := c.Player().CurX()
		y := c.Player().CurY()

		world.AddNpc(world.NewNpc(id, x, y, x-5, x+5, y-5, y+5))
	}
	script.CommandHandlers["summon"] = summon
	script.CommandHandlers["goto"] = gotoTeleport
	script.CommandHandlers["tele"] = teleport
	script.CommandHandlers["teleport"] = teleport
	script.CommandHandlers["death"] = func(c clients.Client, args []string) {
		c.SendPacket(packetbuilders.Death)
		c.Player().Transients().SetVar("deathTime", time.Now())
		c.Player()
		c.Player().SetLocation(world.SpawnPoint)
	}
	script.CommandHandlers["anko"] = func(c clients.Client, args []string) {
		line := strings.Join(args, " ")
		env := script.WorldModule()
		env.Define("println", fmt.Println)
		env.Define("player", c.Player())
		env.Execute(line)
	}
	script.CommandHandlers["reloadscripts"] = func(c clients.Client, args []string) {
		script.InvTriggers = script.InvTriggers[:0]
		script.BoundaryTriggers = script.BoundaryTriggers[:0]
		script.ObjectTriggers = script.ObjectTriggers[:0]
		script.NpcTriggers = script.NpcTriggers[:0]
		script.Load()
		log.Info.Printf("Loaded %d inventory, %d object, %d boundary, and %d NPC action triggers.\n", len(script.InvTriggers), len(script.ObjectTriggers), len(script.BoundaryTriggers), len(script.NpcTriggers))
	}
	script.CommandHandlers["tile"] = func(c clients.Client, args []string) {
		regionX := (2304+c.Player().CurX())/world.RegionSize
		regionY := (1776+c.Player().CurY()-(944*c.Player().Plane()))/world.RegionSize
		mapSector := fmt.Sprintf("h%dx%dy%d", c.Player().Plane(), regionX, regionY)
		areaX := (2304+c.Player().CurX()) % 48
		areaY := (1776+c.Player().CurY()-(944*c.Player().Plane())) % 48
		tile := world.ClipData(c.Player().CurX(), c.Player().CurY())
//		c.Message(fmt.Sprintf("@que@%v sector(%v rel:(%v,%v)): V:%v, H:%v, D:%v, R:%v, O:%v, T:%v, E:%v, bitmask:%v", c.Player().Location.String(), mapSector, areaX, areaY, tile.VerticalWalls, tile.HorizontalWalls, tile.DiagonalWalls, tile.Roofs, tile.GroundOverlay, tile.GroundTexture, tile.GroundElevation, tile.CollisionMask))
		c.Message(fmt.Sprintf("@que@%v sector(%v rel:(%v,%v)): Overlay:%v, bitmask:%v", c.Player().Location.String(), mapSector, areaX, areaY, tile.GroundOverlay, tile.CollisionMask))
	}
	script.CommandHandlers["walk"] = func(c clients.Client, args []string) {
		x, _ := strconv.Atoi(args[0])
		y, _ := strconv.Atoi(args[1])
		c.Player().SetPath(world.MakePath(c.Player().Location, world.NewLocation(x, y)))
	}
}

func teleport(c clients.Client, args []string) {
	if len(args) != 2 {
		c.Message("@que@Invalid args.  Usage: /tele <x> <y>")
		return
	}
	x, err := strconv.Atoi(args[0])
	if err != nil {
		c.Message("@que@Invalid args.  Usage: /tele <x> <y>")
		return
	}
	y, err := strconv.Atoi(args[1])
	if err != nil {
		c.Message("@que@Invalid args.  Usage: /tele <x> <y>")
		return
	}
	if !world.WithinWorld(x, y) {
		c.Message("@que@Coordinates out of world boundaries.")
		return
	}
	log.Commands.Printf("Teleporting %v from %v,%v to %v,%v\n", c.Player().Username, c.Player().CurX(), c.Player().CurY(), x, y)
	c.Teleport(x, y)
}

func summon(c clients.Client, args []string) {
	if len(args) < 1 {
		c.Message("@que@Invalid args.  Usage: /summon <player_name>")
		return
	}
	var name string
	for _, arg := range args {
		name += arg + " "
	}
	name = strings.TrimSpace(name)

	c1, ok := clients.FromUserHash(strutil.Base37.Encode(name))
	if !ok {
		c.Message("@que@@whi@[@cya@SERVER@whi@]: @gre@Could not find player: '" + name + "'")
		return
	}

	log.Commands.Printf("Summoning '%v' from %v,%v to '%v' at %v,%v\n", c1.Player().Username, c1.Player().CurX(), c1.Player().CurY(), c.Player().Username, c.Player().CurX(), c.Player().CurY())
	c1.Teleport(c.Player().CurX(), c.Player().CurY())
}

func gotoTeleport(c clients.Client, args []string) {
	if len(args) < 1 {
		c.Message("@que@Invalid args.  Usage: /goto <player_name>")
		return
	}
	var name string
	for _, arg := range args {
		name += arg + " "
	}
	name = strings.TrimSpace(name)

	c1, ok := clients.FromUserHash(strutil.Base37.Encode(name))
	if !ok {
		c.Message("@que@@whi@[@cya@SERVER@whi@]: @gre@Could not find player: '" + name + "'")
		return
	}

	log.Commands.Printf("Teleporting '%v' from %v,%v to '%v' at %v,%v\n", c.Player().Username, c.Player().CurX(), c.Player().CurY(), c1.Player().Username, c1.Player().CurX(), c1.Player().CurY())
	c.Teleport(c1.Player().CurX(), c1.Player().CurY())
}

func notYetImplemented(c clients.Client, args []string) {
	c.Message("@que@@ora@Not yet implemented")
}
