package server

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(*Client, []string))

func init() {
	if f, err := os.OpenFile("logs"+string(os.PathSeparator)+"cmd.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Error.Println("Could not open commands log file for writing:", err)
	} else {
		log.Commands.SetOutput(f)
	}
	PacketHandlers["command"] = func(c *Client, p *packets.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		handler, ok := CommandHandlers[args[0]]
		if !ok {
			c.Message("@que@Invalid command.")
			log.Commands.Printf("[COMMAND] %v sent invalid command: /%v\n", c.player.Username, string(p.Payload))
			return
		}
		log.Commands.Printf("%v: /%v\n", c.player.Username, string(p.Payload))
		handler(c, args[1:])
	}
	CommandHandlers["dobj"] = func(c *Client, args []string) {
		if len(args) == 0 {
			args = []string{strconv.Itoa(c.player.X), strconv.Itoa(c.player.Y)}
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

		log.Commands.Printf("'%v' deleted object{id: %v; dir:%v} at %v,%v\n", c.player.Username, object.ID, object.Direction, x, y)
		world.RemoveObject(object)
	}
	CommandHandlers["kick"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /kick <player>")
			return
		}
		if pID, err := strconv.Atoi(args[0]); err == nil {
			affectedClient, ok := Clients.FromIndex(pID)
			if !ok {
				c.Message("@que@Could not find player.")
				return
			}
			log.Commands.Printf("'%v' kicked other player '%v'\n", c.player.Username, affectedClient.player.Username)
			c.Message("@que@Kicked: '" + affectedClient.player.Username + "'")
			affectedClient.outgoingPackets <- packets.Logout
			affectedClient.Destroy()
		} else {
			var name string
			for _, arg := range args {
				name += arg + " "
			}
			name = strings.TrimSpace(name)

			affectedClient, ok := Clients.FromUserHash(strutil.Base37(name))
			if !ok {
				c.Message("@que@Could not find player: '" + name + "'")
				return
			}

			log.Commands.Printf("'%v' kicked other player '%v'\n", c.player.Username, affectedClient.player.Username)
			c.Message("@que@Kicked: '" + affectedClient.player.Username + "'")
			affectedClient.outgoingPackets <- packets.Logout
			affectedClient.Destroy()
		}
	}
	CommandHandlers["memdump"] = func(c *Client, args []string) {
		file, err := os.Create("rscgo.mprof")
		if err != nil {
			log.Warning.Println("Could not open file to dump memory profile:", err)
			c.Message("Error encountered opening profile output file.")
			return
		}
		pprof.WriteHeapProfile(file)
		file.Close()
		log.Commands.Println(c.player.Username + " dumped memory profile of the server to rscgo.mprof")
		c.Message("Dumped memory profile.")
	}
	CommandHandlers["pprof"] = func(c *Client, args []string) {
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
			pprof.StartCPUProfile(file)
			log.Commands.Println(c.player.Username + " began profiling CPU time.")
			c.Message("CPU profiling started.")
		case "stop":
			pprof.StopCPUProfile()
			log.Commands.Println(c.player.Username + " has finished profiling CPU time, output should be in rscgo.pprof")
			c.Message("CPU profiling finished.")
		default:
			c.Message("Invalid args.  Usage: /pprof <start|stop>")
		}
	}
	CommandHandlers["object"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /object <id> <dir>, eg: /object 1154 north")
			return
		}
		if world.GetObject(c.player.X, c.player.Y) != nil {
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
		log.Commands.Printf("'%v' spawned new object{id: %v; dir:%v} at %v,%v\n", c.player.Username, id, direction, c.player.X, c.player.Y)
		world.AddObject(world.NewObject(id, direction, c.player.X, c.player.Y, false))
	}
	CommandHandlers["boundary"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /boundary <id> <dir>, eg: /boundary 1 north")
			return
		}
		if world.GetObject(c.player.X, c.player.Y) != nil {
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
		log.Commands.Printf("'%v' spawned new boundary{id: %v; dir:%v} at %v,%v\n", c.player.Username, id, direction, c.player.X, c.player.Y)
		world.AddObject(world.NewObject(id, direction, c.player.X, c.player.Y, true))
	}
	CommandHandlers["saveobjects"] = func(c *Client, args []string) {
		go func() {
			if count := SaveObjectLocations(); count > 0 {
				c.Message("Saved " + strconv.Itoa(count) + " game objects to world.db")
				log.Commands.Println(c.player.Username + " saved " + strconv.Itoa(count) + " game objects to world.db")
			} else {
				c.Message("Appears to have been an issue saving game objects to world.db.  Check server logs.")
				log.Commands.Println(c.player.Username + " failed to save game objects; count=" + strconv.Itoa(count))
			}
		}()
	}
	CommandHandlers["item"] = func(c *Client, args []string) {
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
		c.player.Items.Put(id, amount)
		c.outgoingPackets <- packets.InventoryItems(c.player)
	}
	CommandHandlers["goup"] = func(c *Client, args []string) {
		if nextLocation := c.player.Above(); !nextLocation.Equals(c.player.Location) {
			c.player.SetLocation(nextLocation)
			c.UpdatePlane()
		}
	}
	CommandHandlers["godown"] = func(c *Client, args []string) {
		if nextLocation := c.player.Below(); !nextLocation.Equals(c.player.Location) {
			c.player.SetLocation(nextLocation)
			c.UpdatePlane()
		}
	}
	CommandHandlers["npc"] = notYetImplemented
	CommandHandlers["summon"] = summon
	CommandHandlers["goto"] = gotoTeleport
	CommandHandlers["say"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.Message("@que@Invalid args.  Usage: /say <msg>")
			return
		}
		msg := "@whi@[@cya@GLOBAL@whi@] "
		switch c.player.Rank {
		case 2:
			msg += "@red@~"
		case 1:
			msg += "@blu@@"
		default:
			msg += "@yel@"
		}
		msg += c.player.Username + "@yel@:"
		for _, word := range args {
			msg += " " + word
		}
		Clients.Broadcast(func(c1 *Client) {
			c1.Message("@que@" + msg)
		})
	}
	CommandHandlers["tele"] = teleport
	CommandHandlers["teleport"] = teleport
	CommandHandlers["death"] = func(c *Client, args []string) {
		c.outgoingPackets <- packets.Death
	}
}

func teleport(c *Client, args []string) {
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
	log.Commands.Printf("Teleporting %v from %v,%v to %v,%v\n", c.player.Username, c.player.X, c.player.Y, x, y)
	c.Teleport(x, y)
}

func summon(c *Client, args []string) {
	if len(args) < 1 {
		c.Message("@que@Invalid args.  Usage: /summon <player_name>")
		return
	}
	var name string
	for _, arg := range args {
		name += arg + " "
	}
	name = strings.TrimSpace(name)

	c1, ok := Clients.FromUserHash(strutil.Base37(name))
	if !ok {
		c.Message("@que@@whi@[@cya@SERVER@whi@]: @gre@Could not find player: '" + name + "'")
		return
	}

	log.Commands.Printf("Summoning '%v' from %v,%v to '%v' at %v,%v\n", c1.player.Username, c1.player.X, c1.player.Y, c.player.Username, c.player.X, c.player.Y)
	c1.Teleport(c.player.X, c.player.Y)
}

func gotoTeleport(c *Client, args []string) {
	if len(args) < 1 {
		c.Message("@que@Invalid args.  Usage: /goto <player_name>")
		return
	}
	var name string
	for _, arg := range args {
		name += arg + " "
	}
	name = strings.TrimSpace(name)

	c1, ok := Clients.FromUserHash(strutil.Base37(name))
	if !ok {
		c.Message("@que@@whi@[@cya@SERVER@whi@]: @gre@Could not find player: '" + name + "'")
		return
	}

	log.Commands.Printf("Teleporting '%v' from %v,%v to '%v' at %v,%v\n", c.player.Username, c.player.X, c.player.Y, c1.player.Username, c1.player.X, c1.player.Y)
	c.Teleport(c1.player.X, c1.player.Y)
}

func notYetImplemented(c *Client, args []string) {
	c.Message("@que@@ora@Not yet implemented")
}
