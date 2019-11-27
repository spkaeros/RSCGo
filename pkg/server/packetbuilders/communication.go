package packetbuilders

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//FriendList Builds a packet with the players friend list information in it.
func FriendList(player *world.Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(71)
	p.AddByte(byte(len(player.FriendList)))
	for hash, online := range player.FriendList {
		p.AddLong(hash)
		status := 0
		if online {
			status = 0xFF
		}
		p.AddByte(byte(status)) // 255 for online, 0 for offline.
	}
	return p
}

//PrivateMessage Builds a packet with a private message from hash with content msg.
func PrivateMessage(hash uint64, msg string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(120)
	p.AddLong(hash)
	p.AddInt(rand.Uint32()) // unique Message ID to prevent duplicate messages somehow arriving or something idk
	for _, c := range strutil.ChatFilter.Pack(msg) {
		p.AddByte(c)
	}
	return p
}

//IgnoreList Builds a packet with the players ignore list information in it.
func IgnoreList(player *world.Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(109)
	p.AddByte(byte(len(player.IgnoreList)))
	for _, hash := range player.IgnoreList {
		p.AddLong(hash)
	}
	return p
}

//FriendUpdate Builds a packet with an online status update for the player with the specified hash
func FriendUpdate(hash uint64, online bool) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(149)
	p.AddLong(hash)
	if online {
		p.AddByte(0xFF)
	} else {
		p.AddByte(0)
	}
	return
}

//PlayerChat Builds a packet containing a view-area chat message from the player with the index sender and returns it.
func PlayerChat(sender int, msg string) *packet.Packet {
	p := packet.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender))
	p.AddByte(1)
	p.AddByte(uint8(len(msg)))
	p.AddBytes([]byte(msg))
	return p
}

//PlayerDamage Builds a packet containing a view-area damage display for this player
func PlayerDamage(victim *world.Player, damage int) *packet.Packet {
	p := packet.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(victim.Index))
	p.AddByte(2)
	p.AddByte(uint8(damage))
	p.AddByte(uint8(victim.Skillset.Current(world.StatHits)))
	p.AddByte(uint8(victim.Skillset.Maximum(world.StatHits)))
	return p
}

//NpcDamage Builds a packet containing a view-area damage display for this NPC
func NpcDamage(victim *world.NPC, damage int) *packet.Packet {
	p := packet.NewOutgoingPacket(104)
	p.AddShort(1)
	p.AddShort(uint16(victim.Index))
	p.AddByte(2)
	p.AddByte(uint8(damage))
	p.AddByte(uint8(victim.Skillset.Current(world.StatHits)))
	p.AddByte(uint8(victim.Skillset.Maximum(world.StatHits)))
	return p
}

func NpcMessage(sender *world.NPC, message string, target *world.Player) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(104)
	p.AddShort(1)
	p.AddShort(uint16(sender.Index))
	p.AddByte(1)
	p.AddShort(uint16(target.Index))
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	messageRaw := strutil.ChatFilter.Pack(message)
	p.AddByte(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddByte(c)
	}
	return
}

func PlayerMessage(sender *world.Player, message string) (p *packet.Packet) {
	p = packet.NewOutgoingPacket(234)
	p.AddShort(1)
	p.AddShort(uint16(sender.Index))
	p.AddByte(6)
	if len(message) > 255 {
		message = message[:255]
	}
	message = strutil.ChatFilter.Format(message)
	messageRaw := strutil.ChatFilter.Pack(message)
	p.AddByte(uint8(len(messageRaw)))
	for _, c := range messageRaw {
		p.AddByte(c)
	}
	return
}

//PrivacySettings Builds a packet containing the players privacy settings for display in the settings menu.
func PrivacySettings(player *world.Player) *packet.Packet {
	p := packet.NewOutgoingPacket(51)
	p.AddBool(player.ChatBlocked())
	p.AddBool(player.FriendBlocked())
	p.AddBool(player.TradeBlocked())
	p.AddBool(player.DuelBlocked())
	return p
}

func OptionMenuOpen(questions ...string) *packet.Packet {
	p := packet.NewOutgoingPacket(245)
	p.AddByte(uint8(len(questions)))
	for _, question := range questions {
		p.AddByte(uint8(len(question)))
		p.AddBytes([]byte(question))
	}
	return p
}

var OptionMenuClose = packet.NewOutgoingPacket(252)
