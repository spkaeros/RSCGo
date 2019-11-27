package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

var (
	crewHead           = 1
	metalHead          = 4
	downsHead          = 6
	beardHead          = 7
	baldHead           = 8
	validHeads         = []int{crewHead, metalHead, downsHead, beardHead, baldHead}
	validFemaleHeads   = []int{crewHead, metalHead, downsHead, baldHead}
	maleBody           = 2
	femaleBody         = 5
	validBodys         = []int{maleBody, femaleBody}
	validSkinColors    = []int{0xecded0, 0xccb366, 0xb38c40, 0x997326, 0x906020}
	validHeadColors    = []int{0xffc030, 0xffa040, 0x805030, 0x604020, 0x303030, 0xff6020, 0xff4000, 0xffffff, 65280, 65535}
	validBodyLegColors = []int{0xff0000, 0xff8000, 0xffe000, 0xa0e000, 57344, 32768, 41088, 45311, 33023, 12528, 0xe000e0, 0x303030, 0x604000, 0x805000, 0xffffff}
)

func inArray(a []int, i int) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func init() {
	PacketHandlers["changeappearance"] = func(c *world.Player, p *packet.Packet) {
		if !c.HasState(world.MSChangingAppearance) {
			// Make sure the player either has never logged in before, or talked to the makeover mage to get here.
			return
		}
		isMale := p.ReadBool()
		headType := int(p.ReadByte() + 1)
		bodyType := int(p.ReadByte() + 1)
		legType := int(p.ReadByte() + 1) // appearance2Colour, seems to be a client const, value seems to remain 2.  ofc, legs never change
		hairColor := int(p.ReadByte())
		topColor := int(p.ReadByte())
		legColor := int(p.ReadByte())
		skinColor := int(p.ReadByte())
		/*		if !inArray(validHeads, int(headType)) || !inArray(validBodys, int(bodyType)) || !inArray(validBodyLegColors, int(topColor)) ||
				!inArray(validBodyLegColors, int(legColor)) || !inArray(validHeadColors, int(hairColor))  ||
				!inArray(validSkinColors, int(skinColor)) || legType != 2 {*/
		if hairColor >= len(validHeadColors) || !inArray(validHeads, headType) || topColor >= len(validBodyLegColors) ||
			legColor >= len(validBodyLegColors) || skinColor >= len(validSkinColors) || !inArray(validBodys, bodyType) || legType != 3 || legColor >= len(validBodyLegColors) {
			log.Info.Printf("Invalid appearance data provided by %v: (headType:%v, bodyType:%v, legType:%v, hairColor:%v, topColor:%v, legColor:%v, skinColor:%v, gender:%v)\n", c.String(), headType, bodyType, legType, hairColor, topColor, legColor, skinColor, isMale)
			return
		}
		log.Info.Printf("(headType:%v, bodyType:%v, legType:%v, hairColor:%v, topColor:%v, legColor:%v, skinColor:%v, gender:%v)\n", headType, bodyType, legType, hairColor, topColor, legColor, skinColor, isMale)
		if !isMale {
			if bodyType != femaleBody {
				log.Info.Println("Correcting invalid packet data: female asked for male body type; setting to female body type, from", c)
				bodyType = femaleBody
			}
			if headType == beardHead {
				log.Info.Println("Correcting invalid packet data: female asked for male head type; setting to female head type, from", c)
				headType = metalHead
			}
		}
		c.AppearanceLock.Lock()
		{
			if c.Equips[0] == c.Appearance.Head {
				c.Equips[0] = headType
			}
			if c.Equips[1] == c.Appearance.Body {
				c.Equips[1] = bodyType
			}
			c.Appearance.Body = bodyType
			c.Appearance.Head = headType
			c.Appearance.Male = isMale
			c.Appearance.HeadColor = hairColor
			c.Appearance.SkinColor = skinColor
			c.Appearance.BodyColor = topColor
			c.Appearance.LegsColor = legColor
			c.ResetNeedsSelf()
			c.AppearanceTicket++
		}
		c.AppearanceLock.Unlock()
		c.RemoveState(world.MSChangingAppearance)
	}
}
