package packets

import (
	"bitbucket.org/zlacki/rscgo/pkg/rand"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

var epoch = uint64(time.Now().UnixNano() / int64(time.Millisecond))

//WelcomeMessage Welcome to the game on login
var WelcomeMessage = ServerMessage("Welcome to RuneScape")

//DefaultActionMessage This is a message to inform the player that the action they were trying to perform didn't do anything.
var DefaultActionMessage = ServerMessage("Nothing interesting happens.")

//ServerMessage Builds a packet containing a server message to display in the chat box.
func ServerMessage(msg string) (p *Packet) {
	p = NewOutgoingPacket(131)
	p.AddBytes([]byte(msg))
	return
}

//TeleBubble Builds a packet to draw a teleport bubble at the specified offsets.
func TeleBubble(offsetX, offsetY int) (p *Packet) {
	p = NewOutgoingPacket(36)
	p.AddByte(0) // type, 0 is mobs, 1 is stationary entities, e.g telegrab
	p.AddByte(uint8(offsetX))
	p.AddByte(uint8(offsetY))
	return
}

//InventoryItems Builds a packet containing the players inventory items.
func InventoryItems(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(53)
	player.Items.Lock.RLock()
	p.AddByte(uint8(len(player.Items.List)))
	for _, item := range player.Items.List {
		p.AddShort(uint16(item.ID)) // TODO: + 32768 if wielded.
		if db.Items[item.ID].Stackable {
			p.AddInt2(uint32(item.Amount))
		}
	}
	player.Items.Lock.RUnlock()
	return
}

//ServerInfo Builds a packet with the server information in it.
func ServerInfo(onlineCount int) (p *Packet) {
	// TODO: Real 204 RSC doesn't have this?
	p = NewOutgoingPacket(110)
	p.AddLong(epoch)
	p.AddInt(1337)
	p.AddShort(uint16(onlineCount))
	p.AddBytes([]byte("USA"))
	return p
}

//LoginBox Builds a packet to create a welcome box on the client with the inactiveDays since login, and lastIP connected from.
func LoginBox(inactiveDays int, lastIP string) (p *Packet) {
	p = NewOutgoingPacket(182)
	p.AddInt(uint32(strutil.IPToInteger(lastIP))) // IP
	p.AddShort(uint16(inactiveDays))              // Last logged in
	p.AddByte(0) // recovery questions set days, 200 = unset, 201 = set
	p.AddShort(1) // Unread messages, number minus one, 0 does not render anything
	p.AddBytes([]byte(lastIP))
	return p
}

//FightMode Builds a packet with the players fight mode information in it.
func FightMode(player *world.Player) (p *Packet) {
	// TODO: 204
	p = NewOutgoingPacket(132)
	p.AddByte(byte(player.FightMode()))
	return p
}

//Fatigue Builds a packet with the players fatigue percentage in it.
func Fatigue(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(114)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the server to match RSC
	p.AddShort(uint16(player.Fatigue() / 100))
	return p
}

//FriendList Builds a packet with the players friend list information in it.
func FriendList(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(71)
	p.AddByte(byte(len(player.FriendList)))
	for hash, online := range player.FriendList {
		p.AddLong(hash)
		// TODO: Online status
		status := 0
		if online {
			status = 0xFF
		}
		p.AddByte(byte(status)) // 255 for online, 0 for offline.
	}
	return p
}

//PrivateMessage Builds a packet with a private message from hash with content msg.
func PrivateMessage(hash uint64, msg string) (p *Packet) {
	p = NewOutgoingPacket(120)
	p.AddLong(hash)
	p.AddInt(rand.Uint32()) // unique Message ID to prevent duplicate messages somehow arriving or something idk
	for _, c := range strutil.PackChatMessage(msg) {
		p.AddByte(c)
	}
	return p
}

//IgnoreList Builds a packet with the players ignore list information in it.
func IgnoreList(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(109)
	p.AddByte(byte(len(player.IgnoreList)))
	for _, hash := range player.IgnoreList {
		p.AddLong(hash)
	}
	return p
}

//FriendUpdate Builds a packet with an online status update for the player with the specified hash
func FriendUpdate(hash uint64, online bool) (p *Packet) {
	p = NewOutgoingPacket(149)
	p.AddLong(hash)
	if online {
		p.AddByte(0xFF)
	} else {
		p.AddByte(0)
	}
	return
}

//TradeUpdate Builds a packet to update a trade offer
func TradeUpdate(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(97)
	player.TradeOffer.Lock.RLock()
	p.AddByte(uint8(len(player.TradeOffer.List)))
	for _, item := range player.TradeOffer.List {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
	}
	player.TradeOffer.Lock.RUnlock()
	return
}

//TradeOpen Builds a packet to open a trade window
func TradeOpen(player *world.Player) *Packet {
	return NewOutgoingPacket(92).AddShort(uint16(player.TradeTarget()))
}

//TradeTargetAccept Builds a packet to change trade targets accepted status
func TradeTargetAccept(accepted bool) *Packet {
	if accepted {
		return NewOutgoingPacket(162).AddByte(1)
	}
	return NewOutgoingPacket(162).AddByte(0)
}

//TradeAccept Builds a packet to change trade targets accepted status
func TradeAccept(accepted bool) *Packet {
	if accepted {
		return NewOutgoingPacket(15).AddByte(1)
	}
	return NewOutgoingPacket(15).AddByte(0)
}

//TradeConfirmationOpen Builds a packet to open the trade confirmation page
func TradeConfirmationOpen(player, other *world.Player) *Packet {
	p := NewOutgoingPacket(20)
	p.AddLong(other.UserBase37)

	other.TradeOffer.Lock.RLock()
	p.AddByte(uint8(len(other.TradeOffer.List)))
	for _, item := range other.TradeOffer.List {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
	}
	other.TradeOffer.Lock.RUnlock()

	player.TradeOffer.Lock.RLock()
	p.AddByte(uint8(len(player.TradeOffer.List)))
	for _, item := range player.TradeOffer.List {
		p.AddShort(uint16(item.ID))
		p.AddInt(uint32(item.Amount))
	}
	player.TradeOffer.Lock.RUnlock()
	return p
}

//ClientSettings Builds a packet containing the players client settings, e.g camera mode, mouse mode, sound fx...
func ClientSettings(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(240)
	// TODO: Right IDs?
	if player.GetClientSetting(0) {
		p.AddByte(1)
	} else {
		p.AddByte(0)
	}
	if player.GetClientSetting(2) {
		p.AddByte(1)
	} else {
		p.AddByte(0)
	}
	if player.GetClientSetting(3) {
		p.AddByte(1)
	} else {
		p.AddByte(0)
	}

	//	p.AddByte(0) // Camera auto/manual?
	//	p.AddByte(0) // Mouse buttons 1 or 2?
	//	p.AddByte(1) // Sound effects on/off?
	return
}

//BigInformationBox Builds a packet to trigger the opening of a large black text window with msg as its contents
func BigInformationBox(msg string) (p *Packet) {
	p = NewOutgoingPacket(222)
	p.AddBytes([]byte(msg))
	return p
}

//PlayerChat Builds a packet containing a view-area chat message from the player with the index sender and returns it.
func PlayerChat(sender int, msg string) *Packet {
	p := NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender))
	p.AddByte(1)
	p.AddByte(uint8(len(msg)))
	p.AddBytes([]byte(msg))
	return p
}

//PlayerStats Builds a packet containing all the player's stat information and returns it.
func PlayerStats(player *world.Player) *Packet {
	p := NewOutgoingPacket(156)
	for i := 0; i < 18; i++ {
		p.AddByte(uint8(player.Skillset.Current[i]))
	}

	for i := 0; i < 18; i++ {
		p.AddByte(uint8(player.Skillset.Maximum[i]))
	}

	for i := 0; i < 18; i++ {
		p.AddInt(uint32(player.Skillset.Experience[i]))
	}
	return p
}

//PlayerStat Builds a packet containing player's stat information for skill at idx and returns it.
func PlayerStat(player *world.Player, idx int) *Packet {
	p := NewOutgoingPacket(159)
	p.AddByte(byte(idx))
	p.AddInt(uint32(player.Skillset.Experience[idx]))
	return p
}

//NPCPositions Builds a packet containing view area NPC position and sprite information
func NPCPositions(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(79)
	counter := 0
	p.AddBits(len(player.LocalNPCs.List), 8)
	for _, n := range player.LocalNPCs.List {
		if n, ok := n.(*world.NPC); ok {
			if n.LongestDelta(player.Location) > 15 || n.TransAttrs.VarBool("remove", false) {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				player.LocalNPCs.Remove(n)
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
	for _, p1 := range player.LocalPlayers.List {
		if p1, ok := p1.(*world.Player); ok {
			if p1.LongestDelta(player.Location) > 15 || p1.TransAttrs.VarBool("remove", false) {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				player.LocalPlayers.Remove(p1)
				delete(player.KnownAppearances, p1.Index)
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
		if ticket, ok := player.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket {
			p.AddBits(0, 1)
		} else {
			p.AddBits(1, 1)
		}
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
	ourPlayer.AppearanceReqLock.Lock()
	appearanceList = append(appearanceList, ourPlayer.AppearanceReq...)
	ourPlayer.AppearanceReq = ourPlayer.AppearanceReq[:0]
	ourPlayer.AppearanceReqLock.Unlock()
	for _, p1 := range ourPlayer.LocalPlayers.List {
		if p1, ok := p1.(*world.Player); ok {
			if ticket, ok := ourPlayer.KnownAppearances[p1.Index]; !ok || ticket != p1.AppearanceTicket {
				appearanceList = append(appearanceList, p1)
			}
		}
	}
	if len(appearanceList) <= 0 {
		return nil
	}
	p.AddShort(uint16(len(appearanceList))) // Update size
	for _, player := range appearanceList {
		ourPlayer.KnownAppearances[player.Index] = player.AppearanceTicket
		p.AddShort(uint16(player.Index))
		p.AddByte(5) // Player appearances
		p.AddShort(uint16(player.AppearanceTicket))
		p.AddLong(strutil.Base37(player.Username))
		p.AddByte(12) // worn items length
		p.AddByte(1)  // head
		p.AddByte(2)  // body
		p.AddByte(3)  // unknown, always 3
		for i := 0; i < 9; i++ {
			p.AddByte(0)
		}
		p.AddByte(2)  // Hair
		p.AddByte(8)  // Top
		p.AddByte(14) // Bottom
		p.AddByte(0)  // Skin
		p.AddByte(3)  // Combat lvl
		p.AddByte(0)  // skulled
		//		p.AddByte(byte(player.Rank)) // Rank 2=admin,1=mod,0=normal
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

//EquipmentStats Builds a packet with the players equipment statistics in it.
func EquipmentStats(player *world.Player) (p *Packet) {
	p = NewOutgoingPacket(153)
	p.AddByte(uint8(player.ArmourPoints()))
	p.AddByte(uint8(player.AimPoints()))
	p.AddByte(uint8(player.PowerPoints()))
	p.AddByte(uint8(player.MagicPoints()))
	p.AddByte(uint8(player.PrayerPoints()))
	p.AddByte(uint8(player.RangedPoints()))
	return
}

//LoginResponse Builds a bare packet with the login response code.
func LoginResponse(v int) *Packet {
	return NewBarePacket([]byte{byte(v)})
}

//PlaneInfo Builds a packet to update information about the clients environment, e.g height, player index...
func PlaneInfo(player *world.Player) *Packet {
	playerInfo := NewOutgoingPacket(25)
	playerInfo.AddShort(uint16(player.Index))
	playerInfo.AddShort(2304) // alleged width, tiles per sector also...
	playerInfo.AddShort(1776) // alleged height

	playerInfo.AddShort(uint16(player.Plane())) // plane

	playerInfo.AddShort(944) // REAL plane height
	return playerInfo
}
