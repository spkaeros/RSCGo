package packetbuilders

import (
	"github.com/spkaeros/rscgo/pkg/server/world"
)

//NPCPositions Builds a packet containing view area NPC position and sprite information
func NPCPositions(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(79)
	counter := 0
	p.AddBits(len(player.LocalNPCs.List), 8)
	var removing = world.List{}
	for _, n := range player.LocalNPCs.List {
		if n, ok := n.(*world.NPC); ok {
			if n.LongestDelta(player.Location) > 15 || n.TransAttrs.VarBool("remove", false) {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				removing.List = append(removing.List, n)
				counter++
			} else if n.TransAttrs.VarBool("moved", false) || n.TransAttrs.VarBool("changed", false) {
				p.AddBits(1, 1)
				if n.TransAttrs.VarBool("moved", false) {
					p.AddBits(0, 1)
					p.AddBits(n.Direction(), 3)
				} else {
					p.AddBits(1, 1)
					p.AddBits(n.Direction(), 4)
				}
				counter++
			} else {
				p.AddBits(0, 1)
			}
		}
	}
	for _, n := range removing.List {
		player.LocalNPCs.Remove(n)
	}
	newCount := 0
	for _, n := range player.NewNPCs() {
		if len(player.LocalNPCs.List) >= 255 || newCount >= 25 {
			break
		}
		newCount++
		player.LocalNPCs.Add(n)
		p.AddBits(n.Index, 12)
		offsetX := int(n.X.Load()) - int(player.X.Load())
		if offsetX < 0 {
			offsetX += 32
		}
		offsetY := int(n.Y.Load()) - int(player.Y.Load())
		if offsetY < 0 {
			offsetY += 32
		}
		p.AddBits(offsetX, 5)
		p.AddBits(offsetY, 5)
		p.AddBits(n.Direction(), 4)
		p.AddBits(n.ID, 10)
		counter++
	}
	if counter <= 0 {
		return nil
	}
	return
}

//PlayerPositions Builds a packet containing view area player position and sprite information, including ones own information, and returns it.
// If no players need to be updated, returns nil.
func PlayerPositions(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(191)
	// Note: X coords can be held in 10 bits and Y can be held in 12 bits
	//  Presumably, Jagex used 11 and 13 to evenly fill 3 bytes of data?
	p.AddBits(int(player.X.Load()), 11)
	p.AddBits(int(player.Y.Load()), 13)
	p.AddBits(player.Direction(), 4)
	p.AddBits(len(player.LocalPlayers.List), 8)
	counter := 0
	if player.TransAttrs.VarBool("remove", false) || !player.TransAttrs.VarBool("self", false) || player.TransAttrs.VarBool("moved", false) || player.TransAttrs.VarBool("changed", false) {
		counter++
	}
	var removing = world.List{}
	for _, p1 := range player.LocalPlayers.List {
		if p1, ok := p1.(*world.Player); ok {
			if p1.LongestDelta(player.Location) > 15 || p1.TransAttrs.VarBool("remove", false) {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				removing.List = append(removing.List, p1)
				player.AppearanceLock.Lock()
				delete(player.KnownAppearances, p1.Index)
				player.AppearanceLock.Unlock()
				counter++
			} else if p1.TransAttrs.VarBool("moved", false) || p1.TransAttrs.VarBool("changed", false) {
				p.AddBits(1, 1)
				if p1.TransAttrs.VarBool("moved", false) {
					p.AddBits(0, 1)
					p.AddBits(p1.Direction(), 3)
				} else {
					p.AddBits(1, 1)
					p.AddBits(p1.Direction(), 4)
				}
				counter++
			} else {
				p.AddBits(0, 1)
			}
		}
	}
	for _, p1 := range removing.List {
		player.LocalPlayers.Remove(p1)
	}
	newPlayerCount := 0
	for _, p1 := range player.NewPlayers() {
		if len(player.LocalPlayers.List) >= 255 || newPlayerCount >= 25 {
			// No more than 255 players in view at once, no more than 25 new players at once.
			break
		}
		newPlayerCount++
		p.AddBits(p1.Index, 11)
		offsetX := int(p1.X.Load()) - int(player.X.Load())
		if offsetX < 0 {
			offsetX += 32
		}
		offsetY := int(p1.Y.Load()) - int(player.Y.Load())
		if offsetY < 0 {
			offsetY += 32
		}
		p.AddBits(offsetX, 5)
		p.AddBits(offsetY, 5)
		p.AddBits(p1.Direction(), 4)
		player.AppearanceLock.RLock()
		if ticket, ok := player.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket {
			p.AddBits(0, 1)
		} else {
			p.AddBits(1, 1)
		}
		player.AppearanceLock.RUnlock()
		player.LocalPlayers.Add(p1)
		counter++
	}
	if counter <= 0 {
		return nil
	}
	return
}

//PlayerAppearances Builds a packet with the view-area player appearance profiles in it.
func PlayerAppearances(ourPlayer *world.Player) (p *Packet) {
	p = NewOutgoingPacket(234)
	var appearanceList []*world.Player
	if !ourPlayer.TransAttrs.VarBool("self", false) {
		appearanceList = append(appearanceList, ourPlayer)
	}
	ourPlayer.AppearanceLock.Lock()
	appearanceList = append(appearanceList, ourPlayer.AppearanceReq...)
	ourPlayer.AppearanceReq = ourPlayer.AppearanceReq[:0]
	ourPlayer.AppearanceLock.Unlock()
	for _, p1 := range ourPlayer.LocalPlayers.List {
		if p1, ok := p1.(*world.Player); ok {
			ourPlayer.AppearanceLock.RLock()
			if ticket, ok := ourPlayer.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket {
				appearanceList = append(appearanceList, p1)
			}
			ourPlayer.AppearanceLock.RUnlock()
		}
	}
	if len(appearanceList) <= 0 {
		return nil
	}
	p.AddShort(uint16(len(appearanceList))) // Update size
	for _, player := range appearanceList {
		ourPlayer.AppearanceLock.Lock()
		ourPlayer.KnownAppearances[player.Index] = player.AppearanceTicket
		ourPlayer.AppearanceLock.Unlock()
		p.AddShort(uint16(player.Index))
		p.AddByte(5) // player appearances
		p.AddShort(uint16(player.AppearanceTicket))
		p.AddLong(player.UserBase37)
		p.AddByte(12) // length of sprites.  Anything less than 12 will get padded with 0s
		//		p.AddByte(uint8(player.Appearance.Head))
		//		p.AddByte(uint8(player.Appearance.Body))
		//		p.AddByte(uint8(player.Appearance.Legs))
		ourPlayer.AppearanceLock.RLock()
		for i := 0; i < 12; i++ {
			p.AddByte(uint8(player.Equips[i]))
		}
		ourPlayer.AppearanceLock.RUnlock()
		p.AddByte(uint8(player.Appearance.HeadColor))
		p.AddByte(uint8(player.Appearance.BodyColor))
		p.AddByte(uint8(player.Appearance.LegsColor))
		p.AddByte(uint8(player.Appearance.SkinColor))
		p.AddByte(uint8(player.Skillset.CombatLevel()))
		p.AddByte(0) // TODO: skulled
	}
	return
}

//ObjectLocations Builds a packet with the view-area object positions in it, relative to the player.
// If no new objects are available and no existing local objects are removed from area, returns nil.
func ObjectLocations(player *world.Player) (p *Packet) {
	counter := 0
	p = NewOutgoingPacket(48)
	for _, o := range player.LocalObjects.List {
		if o, ok := o.(*world.Object); ok {
			if o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, 21) || world.GetObject(int(o.X.Load()), int(o.Y.Load())) != o {
				p.AddShort(60000)
				p.AddByte(byte(o.X.Load() - player.X.Load()))
				p.AddByte(byte(o.Y.Load() - player.Y.Load()))
				//				p.AddByte(byte(o.Direction))
				player.LocalObjects.Remove(o)
				counter++
			}
		}
	}
	for _, o := range player.NewObjects() {
		if o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X.Load() - player.X.Load()))
		p.AddByte(byte(o.Y.Load() - player.Y.Load()))
		//		p.AddByte(byte(o.Direction))
		player.LocalObjects.Add(o)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

//BoundaryLocations Builds a packet with the view-area boundary positions in it, relative to the player.
// If no new objects are available and no existing local boundarys are removed from area, returns nil.
func BoundaryLocations(player *world.Player) (p *Packet) {
	counter := 0
	p = NewOutgoingPacket(91)
	for _, o := range player.LocalObjects.List {
		if o, ok := o.(*world.Object); ok {
			if !o.Boundary {
				continue
			}
			if !player.WithinRange(o.Location, 21) {
				//p.AddShort(65535)
				p.AddByte(255)
				p.AddByte(byte(o.X.Load() - player.X.Load()))
				p.AddByte(byte(o.Y.Load() - player.Y.Load()))
				//p.AddByte(byte(o.Direction))
				player.LocalObjects.Remove(o)
				counter++
			}
		}
	}
	for _, o := range player.NewObjects() {
		if !o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X.Load() - player.X.Load()))
		p.AddByte(byte(o.Y.Load() - player.Y.Load()))
		p.AddByte(byte(o.Direction))
		player.LocalObjects.Add(o)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}

//ItemLocations Builds a packet with the view-area item positions in it, relative to the player.
// If no new items are available and no existing items are removed from area, returns nil.
func ItemLocations(player *world.Player) (p *Packet) {
	counter := 0
	p = NewOutgoingPacket(99)
	for _, i := range player.LocalItems.List {
		if i, ok := i.(*world.GroundItem); ok {
			x, y := i.X.Load(), i.Y.Load()
			if !player.WithinRange(i.Location, 21) {
				p.AddByte(255)
				p.AddByte(byte(x - player.X.Load()))
				p.AddByte(byte(y - player.Y.Load()))
				player.LocalItems.Remove(i)
				counter++
			} else if !i.VisibleTo(player) || !world.GetRegion(int(x), int(y)).Items.Contains(i) {
				p.AddShort(uint16(i.ID + 0x8000)) // + 32768
				p.AddByte(byte(x - player.X.Load()))
				p.AddByte(byte(y - player.Y.Load()))
				player.LocalItems.Remove(i)
				counter++
			}
		}
	}
	for _, i := range player.NewItems() {
		p.AddShort(uint16(i.ID))
		p.AddByte(byte(i.X.Load() - player.X.Load()))
		p.AddByte(byte(i.Y.Load() - player.Y.Load()))
		player.LocalItems.Add(i)
		counter++
	}
	if counter == 0 {
		return nil
	}
	return
}
