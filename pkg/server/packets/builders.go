package packets

import "time"

var epoch = uint64(time.Now().Unix())

func ServerMessage(msg string) (p *Packet) {
	p = NewOutgoingPacket(48)
	p.AddBytes([]byte(msg))
	return
}

func TeleBubble(offsetX, offsetY int) (p *Packet) {
	p = NewOutgoingPacket(23)
	p.AddByte(0) // type proly
	p.AddByte(uint8(offsetX))
	p.AddByte(uint8(offsetY))
	return
}

func ServerInfo(onlineCount int) (p *Packet) {
	p = NewOutgoingPacket(110)
	p.AddLong(epoch)
	p.AddInt(1337)
	p.AddShort(uint16(onlineCount))
	p.AddBytes([]byte("United States of America"))
	return p
}

func LoginBox(inactiveDays int, lastIP string) (p *Packet) {
	p = NewOutgoingPacket(248)
	p.AddShort(uint16(inactiveDays))
	p.AddBytes([]byte(lastIP))
	return p
}

func FightMode(mode int) (p *Packet) {
	p = NewOutgoingPacket(132)
	p.AddByte(byte(mode))
	return p
}

func Fatigue(percent int) (p *Packet) {
	p = NewOutgoingPacket(244)
	p.AddShort(uint16(percent))
	return p
}

func BigInformationBox(msg string) (p *Packet) {
	p = NewOutgoingPacket(64)
	p.AddBytes([]byte(msg))
	return p
}

func PlayerPositions(x, y, direction int) (p *Packet) {
	p = NewOutgoingPacket(145)
	p.AddBits(x, 11)
	p.AddBits(y, 13)
	p.AddBits(direction, 4)
	p.AddBits(0, 8)
	return
}

func PlayerAppearances(index int, userHash uint64) (p *Packet) {
	p = NewOutgoingPacket(53)
	p.AddShort(1) // Update size
	p.AddShort(uint16(index))
	p.AddByte(5)  // Player appearances
	p.AddShort(0) // Appearance ID wtf is it, changes every time we change appearance!
	p.AddLong(userHash)
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
	p.AddShort(3) // Combat lvl
	p.AddByte(0)  // skulled
	p.AddByte(2)  // Rank 2=admin,1=mod,0=normal
	return
}

//PlayerInfo Builds a packet to update information about the clients environment, e.g height, player index...
func PlayerInfo(index int, height int) *Packet {
	playerInfo := NewOutgoingPacket(131)
	playerInfo.AddShort(uint16(index))
	playerInfo.AddShort(2304)
	playerInfo.AddShort(1776)

	// getY + 100 / 1000
	playerInfo.AddShort(uint16(height))

	playerInfo.AddShort(944)
	return playerInfo
}
