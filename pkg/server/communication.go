package server

import (
	"strconv"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	Handlers[174] = func(c *Client, p *packets.Packet) {
		LogInfo.Printf("%v: '%v'", c.player.Username, strutil.FormatChatMessage(strutil.UnpackChatMessage(p.Payload)))
	}
	Handlers[120] = func(c *Client, p *packets.Packet) {
		args := strutil.ModalParse(string(p.Payload))
		if args[0] == "tele" {
		}
		switch args[0] {
		case "bubble":
			break
		case "death":
			c.WritePacket(packets.Death)
			break
		case "tele":
			if len(args) < 3 {
				c.WritePacket(packets.ServerMessage("@que@Invalid args.  Usage: /tele <x> <y>"))
				return
			}
			x, _ := strconv.Atoi(args[1])
			y, _ := strconv.Atoi(args[2])
			LogInfo.Printf("Teleporting %v from %v,%v to %v,%v\n", c, c.player.Location().X(), c.player.Location().Y(), x, y)
			c.player.Location().SetX(x)
			c.player.Location().SetY(y)
			c.WritePacket(packets.TeleBubble(0, 0))
			LogInfo.Printf("%v: '/%v'", c.player.Username, args)
			break
		default:
			c.WritePacket(packets.ServerMessage("@que@Invalid command."))
			break
		}
	}
}
