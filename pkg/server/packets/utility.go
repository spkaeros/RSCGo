package packets

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
	"time"
)

//Logout Resets client to login welcome screen
var Logout = NewOutgoingPacket(4)

//WelcomeMessage Welcome to the game on login
var WelcomeMessage = ServerMessage("Welcome to RuneScape")

//Death The 'Oh dear...You are dead' fade-to-black graphic effect when you die.
var Death = NewOutgoingPacket(83)

//ResponsePong Response to a RSC protocol ping packet
var ResponsePong = NewOutgoingPacket(9)

//CannotLogout Message that you can not logout right now.
var CannotLogout = NewOutgoingPacket(183)

//DefaultActionMessage This is a message to inform the player that the action they were trying to perform didn't do anything.
var DefaultActionMessage = ServerMessage("Nothing interesting happens.")

//ServerMessage Builds a packet containing a server message to display in the chat box.
func ServerMessage(msg string) (p *Packet) {
	p = NewOutgoingPacket(131)
	p.AddBytes([]byte(msg))
	return
}

var epoch = uint64(time.Now().UnixNano() / int64(time.Millisecond))

//TeleBubble Builds a packet to draw a teleport bubble at the specified offsets.
func TeleBubble(offsetX, offsetY int) (p *Packet) {
	p = NewOutgoingPacket(36)
	p.AddByte(0) // type, 0 is mobs, 1 is stationary entities, e.g telegrab
	p.AddByte(uint8(offsetX))
	p.AddByte(uint8(offsetY))
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

//BigInformationBox Builds a packet to trigger the opening of a large black text window with msg as its contents
func BigInformationBox(msg string) (p *Packet) {
	p = NewOutgoingPacket(222)
	p.AddBytes([]byte(msg))
	return p
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
