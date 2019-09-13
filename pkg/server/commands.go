package server

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
	"bitbucket.org/zlacki/rscgo/pkg/world"
)

//CommandHandlers A map to assign in-game commands to the functions they should execute.
var CommandHandlers = make(map[string]func(*Client, []string))

func init() {
	PacketHandlers["command"] = func(c *Client, p *packets.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		handler, ok := CommandHandlers[args[0]]
		if !ok {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid command.")
			LogInfo.Printf("[COMMAND] %v sent invalid command: /%v\n", c.player.Username, string(p.Payload))
			return
		}
		LogInfo.Printf("[COMMAND] %v: /%v\n", c.player.Username, string(p.Payload))
		handler(c, args[1:])
	}
	CommandHandlers["dobj"] = func(c *Client, args []string) {
		if len(args) != 2 {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /dobj <x> <y>")
			return
		}
		x, err := strconv.Atoi(args[0])
		if err != nil {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /dobj <x> <y>")
			return
		}
		y, err := strconv.Atoi(args[1])
		if err != nil {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /dobj <x> <y>")
			return
		}
		if !world.WithinWorld(x, y) {
			c.outgoingPackets <- packets.ServerMessage("@que@Coordinates out of world boundaries.")
			return
		}
		object := world.GetObject(x, y)
		if object == nil {
			c.outgoingPackets <- packets.ServerMessage(fmt.Sprintf("@que@Can not find object at coords %d,%d", x, y))
			return
		}

		world.RemoveObject(object)
	}
	CommandHandlers["kick"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /kick <player>")
			return
		}
		if pID, err := strconv.Atoi(args[0]); err == nil {
			affectedClient := ClientFromIndex(pID)
			if affectedClient == nil {
				c.outgoingPackets <- packets.ServerMessage("@que@Could not find player.")
				return
			}
			c.outgoingPackets <- packets.ServerMessage("@que@Kicked: '" + affectedClient.player.Username + "'")
			affectedClient.outgoingPackets <- packets.Logout
			affectedClient.Destroy()
		} else {
			var name string
			for _, arg := range args {
				name += arg + " "
			}
			name = strings.TrimSpace(name)

			affectedClient := ClientFromHash(strutil.Base37(name))
			if affectedClient == nil {
				c.outgoingPackets <- packets.ServerMessage("@que@Could not find player: '" + name + "'")
				return
			}

			c.outgoingPackets <- packets.ServerMessage("@que@Kicked: '" + affectedClient.player.Username + "'")
			affectedClient.outgoingPackets <- packets.Logout
			affectedClient.Destroy()
		}
	}
	CommandHandlers["object"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /object <id> <dir>, eg: /object 1154 north")
			return
		}
		if world.GetObject(c.player.X, c.player.Y) != nil {
			c.outgoingPackets <- packets.ServerMessage("@que@You must remove the old object at this location first!")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /object <id>")
			return
		}
		direction := world.North
		if len(args) > 1 {
			switch args[1] {
			case "northeast":
			case "ne":
				direction = world.NorthEast
			case "northwest":
			case "nw":
				direction = world.NorthWest
			case "east":
			case "e":
				direction = world.East
			case "west":
			case "w":
				direction = world.West
			case "south":
			case "s":
				direction = world.South
			case "southeast":
			case "se":
				direction = world.SouthEast
			case "southwest":
			case "sw":
				direction = world.SouthWest
			case "north":
			case "n":
			default:
				direction = world.North
			}
		}
		world.AddObject(world.NewObject(id, direction, c.player.X, c.player.Y, false))
	}
	CommandHandlers["item"] = notYetImplemented
	CommandHandlers["goup"] = notYetImplemented
	CommandHandlers["godown"] = notYetImplemented
	CommandHandlers["npc"] = notYetImplemented
	CommandHandlers["summon"] = summon
	CommandHandlers["goto"] = gotoTeleport
	CommandHandlers["say"] = func(c *Client, args []string) {
		if len(args) < 1 {
			c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /say <msg>")
			return
		}
		msg := "@whi@[@cya@GLOBAL@whi@] @yel@" + c.player.Username + "@whi@:@yel@"
		for _, arg := range args {
			msg += " " + arg
		}
		for _, c1 := range Clients {
			if c1.player.Connected {
				c1.outgoingPackets <- packets.ServerMessage(fmt.Sprintf("@que@%s", msg))
			}
		}
	}
	CommandHandlers["tele"] = teleport
	CommandHandlers["teleport"] = teleport
	CommandHandlers["death"] = func(c *Client, args []string) {
		c.outgoingPackets <- packets.Death
	}
}

func teleport(c *Client, args []string) {
	if len(args) < 2 {
		c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /tele <x> <y>")
		return
	}
	x, _ := strconv.Atoi(args[0])
	y, _ := strconv.Atoi(args[1])
	if x >= world.MaxX || y >= world.MaxY || x < 0 || y < 0 {
		c.outgoingPackets <- packets.ServerMessage(fmt.Sprintf("@que@Invalid coordinates.  Must be between 0,0 and %v,%v", world.MaxX, world.MaxY))
		return
	}
	newLocation := world.NewLocation(x, y)
	LogInfo.Printf("Teleporting %v from %v to %v\n", c.player.Username, c.player, newLocation)
	c.outgoingPackets <- packets.TeleBubble(0, 0)
	for _, p1 := range world.GetRegionFromLocation(c.player.Location).Players.NearbyPlayers(c.player) {
		diffX := c.player.X - p1.X
		diffY := c.player.Y - p1.Y
		if c1, ok := ClientsIdx[p1.Index]; ok {
			c1.outgoingPackets <- packets.TeleBubble(diffX, diffY)
		}
	}
	c.player.TransAttrs["plrremove"] = true
	c.player.SetLocation(*newLocation)
}

func summon(c *Client, args []string) {
	if len(args) < 1 {
		c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /summon <player_name>")
		return
	}
	var name string
	for _, arg := range args {
		name += arg + " "
	}
	name = strings.TrimSpace(name)

	if c1, ok := Clients[strutil.Base37(name)]; ok {
		c1.outgoingPackets <- packets.TeleBubble(0, 0)
		for _, p1 := range world.GetRegionFromLocation(c1.player.Location).Players.NearbyPlayers(c1.player) {
			diffX := c1.player.X - p1.X
			diffY := c1.player.Y - p1.Y
			if c2, ok := ClientsIdx[p1.Index]; ok {
				c2.outgoingPackets <- packets.TeleBubble(diffX, diffY)
			}
		}
		c1.player.TransAttrs["plrremove"] = true
		c1.player.SetLocation(c.player.Location)
		return
	}
	c.outgoingPackets <- packets.ServerMessage("@que@@whi@[@cya@SERVER@whi@]: @gre@Could not find player with username '" + name + "'")
}

func gotoTeleport(c *Client, args []string) {
	if len(args) < 1 {
		c.outgoingPackets <- packets.ServerMessage("@que@Invalid args.  Usage: /goto <player_name>")
		return
	}
	var name string
	for _, arg := range args {
		name += arg + " "
	}
	name = strings.TrimSpace(name)

	if c1, ok := Clients[strutil.Base37(name)]; ok {
		c.outgoingPackets <- packets.TeleBubble(0, 0)
		for _, p1 := range world.GetRegionFromLocation(c.player.Location).Players.NearbyPlayers(c.player) {
			diffX := c.player.X - p1.X
			diffY := c.player.Y - p1.Y
			if c2, ok := Clients[p1.UserBase37]; ok {
				c2.outgoingPackets <- packets.TeleBubble(diffX, diffY)
			}
		}
		c.player.TransAttrs["plrremove"] = true
		c.player.SetLocation(c1.player.Location)
		return
	}
	c.outgoingPackets <- packets.ServerMessage("@que@@whi@[@cya@SERVER@whi@]: @gre@Could not find player with username '" + name + "'")
}

func notYetImplemented(c *Client, args []string) {
	c.outgoingPackets <- packets.ServerMessage("@que@@ora@Not yet implemented")
}
