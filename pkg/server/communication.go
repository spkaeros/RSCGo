package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	Handlers[174] = func(c *Client, p *packets.Packet) {
		LogInfo.Printf("%v: '%v'", c.player.Username, strutil.FormatChatMessage(strutil.UnpackChatMessage(p.Payload)))
	}
}
