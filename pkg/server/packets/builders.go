package packets

func PlayerPositions(x, y, direction int) (p *Packet) {
	p = NewOutgoingPacket(145)
	p.AddBits(x, 11)
	p.AddBits(y, 13)
	p.AddBits(direction, 4)
	p.AddBits(0, 8)
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
