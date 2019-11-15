package packethandlers

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(clients.Client, []string))

func init() {
	PacketHandlers["command"] = func(c clients.Client, p *packetbuilders.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		handler, ok := CommandHandlers[args[0]]
		if !ok {
			c.Message("@que@Invalid command.")
			log.Commands.Printf("%v sent invalid command: /%v\n", c.Player().Username, string(p.Payload))
			return
		}
		log.Commands.Printf("%v: /%v\n", c.Player().Username, string(p.Payload))
		handler(c, args[1:])
	}
	CommandHandlers["dobj"] = func(c clients.Client, args []string) {
		if len(args) == 0 {
			args = []string{strconv.Itoa(int(c.Player().X.Load())), strconv.Itoa(int(c.Player().Y.Load()))}
		}
		if len(args) < 2 {
			c.Message("@que@Invalid args.  Usage: /dobj <x> <y>")
			return
		}
		x, err := strconv.Atoi(args[0])
		if err != nil {
			c.Message("@que@Invalid args.  Usage: /dobj <x> <y>")
			return
		}
		y, err := strconv.Atoi(args[1])
		if err != nil {
			c.Message("@que@Invalid args.  Usage: /dobj <x> <y>")
			return
		}
		if !world.WithinWorld(x, y) {
			c.Message("@que@Coordinates out of world boundaries.")
			return
		}
		object := world.GetObject(x, y)
		if object == nil {
			c.Message(fmt.Sprintf("@que@Can not find object at coords %d,%d", x, y))
			return
		}

		log.Commands.Printf("'%v' deleted object{id: %v; dir:%v} at %v,%v\n", c.Player().Username, object.ID, object.Direction, x, y)
		world.RemoveObject(object)
	}
	CommandHandlers["kick"] = func(c clients.Client, args []string) {
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
	CommandHandlers["memdump"] = func(c clients.Client, args []string) {
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
	CommandHandlers["pprof"] = func(c clients.Client, args []string) {
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
	CommandHandlers["object"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /object <id> <dir>, eg: /object 1154 north")
			return
		}
		x := int(c.Player().X.Load())
		y := int(c.Player().Y.Load())
		if world.GetObject(x, y) != nil {
			c.Message("@que@You must remove the old object at this location first!")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			c.Message("@que@Invalid args.  Usage: /object <id> <dir>, eg: /object 1154 north")
			return
		}
		direction := world.North
		if len(args) > 1 {
			if d, err := strconv.Atoi(args[1]); err == nil {
				if d < world.North || d > world.NorthEast {
					c.Message("@que@Invalid direction; must be between 0 and 8, or simply spell out the direction or its initials.")
					return
				}
				direction = d
			} else {
				direction = world.ParseDirection(args[1])
			}
		}
		log.Commands.Printf("'%v' spawned new object{id: %v; dir:%v} at %v,%v\n", c.Player().Username, id, direction, x, y)
		world.AddObject(world.NewObject(id, direction, x, y, false))
	}
	CommandHandlers["boundary"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /boundary <id> <dir>, eg: /boundary 1 north")
			return
		}
		x := int(c.Player().X.Load())
		y := int(c.Player().Y.Load())
		if world.GetObject(x, y) != nil {
			c.Message("@que@You must remove the old boundary at this location first!")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			c.Message("@que@Invalid args.  Usage: /boundary <id> <dir>, eg: /boundary 1154 north")
			return
		}
		direction := world.North
		if len(args) > 1 {
			if d, err := strconv.Atoi(args[1]); err == nil {
				if d < world.North || d > world.NorthEast {
					c.Message("@que@Invalid direction; must be between 0 and 8, or simply spell out the direction or its initials.")
					return
				}
				direction = d
			} else {
				direction = world.ParseDirection(args[1])
			}
		}
		log.Commands.Printf("'%v' spawned new boundary{id: %v; dir:%v} at %v,%v\n", c.Player().Username, id, direction, x, y)
		world.AddObject(world.NewObject(id, direction, x, y, true))
	}
	CommandHandlers["saveobjects"] = func(c clients.Client, args []string) {
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
	CommandHandlers["item"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /item <id> <quantity>")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil || id > 1289 || id < 0 {
			c.Message("@que@Invalid args.  Usage: /item <id> <quantity>")
			return
		}
		amount := 1
		if len(args) > 1 {
			amount, err = strconv.Atoi(args[1])
			if err != nil || amount <= 0 {
				c.Message("@que@Invalid args.  Usage: /item <id> <quantity>")
				return
			}
		}
		if !db.Items[id].Stackable && amount > 1 {
			for i := 0; i < amount; i++ {
				c.Player().Items.Add(id, 1)
				if c.Player().Items.Size() >= 30 {
					break
				}
			}
		} else {
			c.Player().Items.Add(id, amount)
		}
		c.SendPacket(packetbuilders.InventoryItems(c.Player()))
	}
	CommandHandlers["goup"] = func(c clients.Client, args []string) {
		if nextLocation := c.Player().Above(); !nextLocation.Equals(&c.Player().Location) {
			c.Player().SetLocation(nextLocation)
			c.UpdatePlane()
		}
	}
	CommandHandlers["godown"] = func(c clients.Client, args []string) {
		if nextLocation := c.Player().Below(); !nextLocation.Equals(&c.Player().Location) {
			c.Player().SetLocation(nextLocation)
			c.UpdatePlane()
		}
	}
	CommandHandlers["npc"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /npc <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id > 793 || id < 0 {
			c.Message("@que@Invalid args.  Usage: /npc <id>")
			return
		}

		x := int(c.Player().X.Load())
		y := int(c.Player().Y.Load())

		world.AddNpc(world.NewNpc(id, x, y, x-5, x+5, y-5, y+5))
	}
	CommandHandlers["summon"] = summon
	CommandHandlers["goto"] = gotoTeleport
	CommandHandlers["say"] = func(c clients.Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /say <msg>")
			return
		}
		msg := "@whi@[@cya@GLOBAL@whi@] "
		switch c.Player().Rank {
		case 2:
			msg += "@red@~"
		case 1:
			msg += "@blu@@"
		default:
			msg += "@yel@"
		}
		msg += c.Player().Username + "@yel@:"
		for _, word := range args {
			msg += " " + word
		}
		clients.Range(func(c1 clients.Client) {
			c1.Message("@que@" + msg)
		})
	}
	CommandHandlers["tele"] = teleport
	CommandHandlers["teleport"] = teleport
	CommandHandlers["death"] = func(c clients.Client, args []string) {
		c.SendPacket(packetbuilders.Death)
	}
	CommandHandlers["anko"] = func(c clients.Client, args []string) {
		line := strings.Join(args, " ")
		env := script.WorldModule()
		env.Define("println", fmt.Println)
		env.Define("player", c.Player())
		env.Execute(line)
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
	log.Commands.Printf("Teleporting %v from %v,%v to %v,%v\n", c.Player().Username, c.Player().X, c.Player().Y, x, y)
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

	log.Commands.Printf("Summoning '%v' from %v,%v to '%v' at %v,%v\n", c1.Player().Username, c1.Player().X.Load(), c1.Player().Y.Load(), c.Player().Username, c.Player().X.Load(), c.Player().Y.Load())
	c1.Teleport(int(c.Player().X.Load()), int(c.Player().Y.Load()))
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

	log.Commands.Printf("Teleporting '%v' from %v,%v to '%v' at %v,%v\n", c.Player().Username, c.Player().X.Load(), c.Player().Y.Load(), c1.Player().Username, c1.Player().X.Load(), c1.Player().Y.Load())
	c.Teleport(int(c1.Player().X.Load()), int(c1.Player().Y.Load()))
}

func notYetImplemented(c clients.Client, args []string) {
	c.Message("@que@@ora@Not yet implemented")
}
