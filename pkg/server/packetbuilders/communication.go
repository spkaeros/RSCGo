package packetbuilders

import (
	"bitbucket.org/zlacki/rscgo/pkg/rand"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

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
	for _, c := range strutil.ChatFilter.Pack(msg) {
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

//PrivacySettings Builds a packet containing the players privacy settings for display in the settings menu.
func PrivacySettings(player *world.Player) *Packet {
	p := NewOutgoingPacket(51)
	p.AddBool(player.ChatBlocked())
	p.AddBool(player.FriendBlocked())
	p.AddBool(player.TradeBlocked())
	p.AddBool(player.DuelBlocked())
	return p
}
