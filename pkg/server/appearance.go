package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["changeappearance"] = func(c *Client, p *packets.Packet) {
		headGender := p.ReadBool()
		headType := p.ReadByte()
		bodyType := p.ReadByte()
		legType := p.ReadByte() // appearance2Colour, seems to be a client const, value seems to remain 2.  ofc, legs never change
		hairColor := p.ReadByte()
		topColor := p.ReadByte()
		legColor := p.ReadByte()
		skinColor := p.ReadByte()
		log.Info.Printf("Player appearance modification requested: headGender:%d,headType:%d,bodyType:%d,legType:%d,hairColor:%d,topColor:%d,legColor:%d,skinColor:%d\n", headGender, headType,
			bodyType, legType, hairColor, topColor, legColor, skinColor)
		c.player.Appearance = world.AppearanceTable{
			Head:      int(headType+1),
			Body:      int(bodyType+1),
			Legs:      int(legType+1),
			Male:      headGender,
			HeadColor: int(hairColor),
			BodyColor: int(topColor),
			LegsColor: int(legColor),
			SkinColor: int(skinColor),
		}
		c.player.AppearanceTicket++
		c.player.TransAttrs.SetVar("self", false)
	}
}
