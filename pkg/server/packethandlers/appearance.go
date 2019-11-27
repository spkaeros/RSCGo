package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["changeappearance"] = func(c clients.Client, p *packet.Packet) {
		if !c.Player().HasState(world.MSChangingAppearance) {
			// Make sure the player either has never logged in before, or talked to the makeover mage to get here.
			return
		}
		headGender := p.ReadBool()
		headType := p.ReadByte()
		bodyType := p.ReadByte()
		legType := p.ReadByte() // appearance2Colour, seems to be a client const, value seems to remain 2.  ofc, legs never change
		hairColor := p.ReadByte()
		topColor := p.ReadByte()
		legColor := p.ReadByte()
		skinColor := p.ReadByte()
		if c.Player().Equips[0] == c.Player().Appearance.Head {
			c.Player().Equips[0] = int(headType + 1)
		}
		if c.Player().Equips[1] == c.Player().Appearance.Body {
			c.Player().Equips[1] = int(bodyType + 1)
		}
		if c.Player().Equips[2] == c.Player().Appearance.Legs {
			c.Player().Equips[2] = int(legType + 1)
		}
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
		c.Player().RemoveState(world.MSChangingAppearance)
		c.Player().TransAttrs.SetVar("self", false)
	}
}
