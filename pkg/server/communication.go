package server

import (
	"strconv"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

var CommandHandlers = make(map[string]func(*Client, []string))

func init() {
	CommandHandlers["tele"] = func(c *Client, args []string) {
		if len(args) < 2 {
			c.WritePacket(packets.ServerMessage("@que@Invalid args.  Usage: /tele <x> <y>"))
			return
		}
		x, _ := strconv.Atoi(args[0])
		y, _ := strconv.Atoi(args[1])
		LogInfo.Printf("Teleporting %v from %v,%v to %v,%v\n", c, c.player.Location().X(), c.player.Location().Y(), x, y)
		c.player.Location().SetX(x)
		c.player.Location().SetY(y)
		c.WritePacket(packets.TeleBubble(0, 0))
	}
	CommandHandlers["death"] = func(c *Client, args []string) {
		c.WritePacket(packets.Death)
	}
	Handlers[174] = func(c *Client, p *packets.Packet) {
		LogInfo.Printf("%v: '%v'", c.player.Username, strutil.FormatChatMessage(strutil.UnpackChatMessage(p.Payload)))
	}
	Handlers[120] = func(c *Client, p *packets.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		handler, ok := CommandHandlers[args[0]]
		if !ok {
			c.WritePacket(packets.ServerMessage("@que@Invalid command."))
			return
		}
		handler(c, args[1:])
	}
}
