package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["changeappearance"] = func(c clients.Client, p *packetbuilders.Packet) {
		headGender := p.ReadBool()
		headType := p.ReadByte()
		bodyType := p.ReadByte()
		legType := p.ReadByte() // appearance2Colour, seems to be a client const, value seems to remain 2.  ofc, legs never change
		hairColor := p.ReadByte()
		topColor := p.ReadByte()
		legColor := p.ReadByte()
		skinColor := p.ReadByte()
		c.Player().Appearance = world.AppearanceTable{
			Head:      int(headType + 1),
			Body:      int(bodyType + 1),
			Legs:      int(legType + 1),
			Male:      headGender,
			HeadColor: int(hairColor),
			BodyColor: int(topColor),
			LegsColor: int(legColor),
			SkinColor: int(skinColor),
		}
		c.Player().AppearanceTicket++
		c.Player().TransAttrs.SetVar("self", false)
	}
}
