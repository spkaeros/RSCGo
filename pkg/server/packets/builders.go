/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-22-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package packets

import (
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
)

var epoch = uint64(time.Now().UnixNano() / int64(time.Millisecond))

func ServerMessage(msg string) (p *Packet) {
	p = NewOutgoingPacket(48)
	p.AddBytes([]byte(msg))
	return
}

func TeleBubble(offsetX, offsetY int) (p *Packet) {
	p = NewOutgoingPacket(23)
	p.AddByte(0) // type, 0 is mobs, 1 is stationary entities, e.g telegrab
	p.AddByte(uint8(offsetX))
	p.AddByte(uint8(offsetY))
	return
}

func ServerInfo(onlineCount int) (p *Packet) {
	p = NewOutgoingPacket(110)
	p.AddLong(epoch)
	p.AddInt(1337)
	p.AddShort(uint16(onlineCount))
	p.AddBytes([]byte("USA"))
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

func Fatigue(player *entity.Player) (p *Packet) {
	p = NewOutgoingPacket(244)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the server to match RSC
	p.AddShort(uint16(player.Fatigue() / 100))
	return p
}

func BigInformationBox(msg string) (p *Packet) {
	p = NewOutgoingPacket(64)
	p.AddBytes([]byte(msg))
	return p
}

func PlayerPositions(player *entity.Player, local []*entity.Player) (p *Packet) {
	p = NewOutgoingPacket(145)
	// Note: X coords can be held in 10 bits and Y can be held in 12 bits
	//  Presumably, Jagex used 11 and 13 to evenly fill 3 bytes of data?
	p.AddBits(player.X(), 11)
	p.AddBits(player.Y(), 13)
	p.AddBits(int(player.Direction()), 4)
	p.AddBits(0, 8)
	for _, p1 := range local {
		p.AddBits(p1.Index, 11)
		offsetX := (p1.X() - player.X())
		if offsetX < 0 {
			offsetX += 32
		}
		offsetY := (p1.Y() - player.Y())
		if offsetY < 0 {
			offsetY += 32
		}
		p.AddBits(offsetX, 5)
		p.AddBits(offsetY, 5)
		p.AddBits(int(p1.Direction()), 4)
		p.AddBits(0, 1)
	}
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

func ObjectLocations(player *entity.Player, newObjects []*entity.Object) (p *Packet) {
	p = NewOutgoingPacket(27)
	for _, o := range newObjects {
		if o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		p.AddByte(byte(o.Direction))
	}
	return
}

func EquipmentStats(player *entity.Player) (p *Packet) {
	p = NewOutgoingPacket(177)
	p.AddShort(uint16(player.ArmourPoints()))
	p.AddShort(uint16(player.AimPoints()))
	p.AddShort(uint16(player.PowerPoints()))
	p.AddShort(uint16(player.MagicPoints()))
	p.AddShort(uint16(player.PrayerPoints()))
	p.AddShort(uint16(player.RangedPoints()))
	return
}

//LoginResponse Builds a bare packet with the login response code.
func LoginResponse(v int) *Packet {
	return NewBarePacket([]byte{byte(v)})
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
