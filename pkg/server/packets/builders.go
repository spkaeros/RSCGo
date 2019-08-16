package packets

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
