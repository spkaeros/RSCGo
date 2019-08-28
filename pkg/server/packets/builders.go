/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-22-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-27-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package packets

import (
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

var epoch = uint64(time.Now().UnixNano() / int64(time.Millisecond))

var WelcomeMessage = ServerMessage("Welcome to RuneScape")

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

func FightMode(player *entity.Player) (p *Packet) {
	p = NewOutgoingPacket(132)
	p.AddByte(byte(player.FightMode()))
	return p
}

func Fatigue(player *entity.Player) (p *Packet) {
	p = NewOutgoingPacket(244)
	// Fatigue is converted to percentage differently in the client.
	// 100% clientside is 750, serverside is 75000.  Needs the extra precision on the server to match RSC
	p.AddShort(uint16(player.Fatigue() / 100))
	return p
}

func FriendList(player *entity.Player) (p *Packet) {
	p = NewOutgoingPacket(249)
	p.AddByte(byte(len(player.FriendList)))
	for _, hash := range player.FriendList {
		p.AddLong(hash)
		// TODO: Online status
		p.AddByte(0) // 99 for online, 0 for offline.
	}
	return p
}

func FriendUpdate(player *entity.Player, hash uint64, online bool) (p *Packet) {
	p = NewOutgoingPacket(25)
	p.AddLong(hash)
	if online {
		p.AddByte(99)
	} else {
		p.AddByte(0)
	}
	return
}

func ClientSettings(player *entity.Player) (p *Packet) {
	p = NewOutgoingPacket(152)
	p.AddByte(0) // Camera auto/manual?
	p.AddByte(0) // Mouse buttons 1 or 2?
	p.AddByte(1) // Sound effects on/off?
	return
}

func BigInformationBox(msg string) (p *Packet) {
	p = NewOutgoingPacket(64)
	p.AddBytes([]byte(msg))
	return p
}

func PlayerChat(sender int, msg string) *Packet {
	p := NewOutgoingPacket(53)
	p.AddShort(1)
	p.AddShort(uint16(sender))
	p.AddByte(1)
	p.AddByte(uint8(len(msg)))
	p.AddBytes([]byte(msg))
	return p
}

func PlayerStats(player *entity.Player) *Packet {
	p := NewOutgoingPacket(180)
	for i := 0; i < 18; i++ {
		p.AddShort(uint16(player.Skillset.Current[i]))
	}

	for i := 0; i < 18; i++ {
		p.AddShort(uint16(player.Skillset.Maximum[i]))
	}

	for i := 0; i < 18; i++ {
		p.AddLong(uint64(player.Skillset.Experience[i]))
	}
	return p
}

func PlayerStat(player *entity.Player, idx int) *Packet {
	p := NewOutgoingPacket(208)
	p.AddByte(byte(idx))
	p.AddShort(uint16(player.Skillset.Current[idx]))
	p.AddShort(uint16(player.Skillset.Maximum[idx]))
	p.AddLong(uint64(player.Skillset.Experience[idx]))
	return p
}

func PlayerPositions(player *entity.Player, local []*entity.Player, removing []*entity.Player) (p *Packet) {
	p = NewOutgoingPacket(145)
	// Note: X coords can be held in 10 bits and Y can be held in 12 bits
	//  Presumably, Jagex used 11 and 13 to evenly fill 3 bytes of data?
	p.AddBits(player.X(), 11)
	p.AddBits(player.Y(), 13)
	p.AddBits(int(player.Direction()), 4)
	p.AddBits(len(player.LocalPlayers.List), 8)
	counter := 0
	if player.HasMoved || !player.HasSelf {
		counter++
	}
	//	if !player.HasSelf {
	//		counter++
	//	}
	for _, p1 := range removing {
		p.AddBits(1, 1)
		p.AddBits(1, 1)
		p.AddBits(3, 2)
		player.LocalPlayers.RemovePlayer(p1)
		counter++
	}
	for _, p1 := range player.LocalPlayers.List {
		p1, ok := p1.(*entity.Player)
		if ok {
			if p1.Location().LongestDelta(player.Location()) > 15 || p1.Removing {
				p.AddBits(1, 1)
				p.AddBits(1, 1)
				p.AddBits(3, 2)
				player.LocalPlayers.RemovePlayer(p1)
				counter++
			} else if p1.HasMoved {
				p.AddBits(1, 1)
				p.AddBits(0, 1)
				p.AddBits(int(p1.Direction()), 3)
				counter++
			} else {
				p.AddBits(0, 1)
			}
		}
	}
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
		p.AddBits(1, 1)
		player.LocalPlayers.AddPlayer(p1)
		counter++
	}
	if counter <= 0 {
		return nil
	}
	return
}

func PlayerAppearances(ourPlayer *entity.Player, local []*entity.Player) (p *Packet) {
	p = NewOutgoingPacket(53)
	if ourPlayer.AppearanceChanged {
		local = append(local, ourPlayer)
	}
	if len(local) <= 0 {
		return nil
	}
	p.AddShort(uint16(len(local))) // Update size
	for _, player := range local {
		p.AddShort(uint16(player.Index))
		p.AddByte(5)  // Player appearances
		p.AddShort(0) // Appearance ID wtf is it, changes every time we change appearance!
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
		p.AddShort(3) // Combat lvl
		p.AddByte(0)  // skulled
		p.AddByte(2)  // Rank 2=admin,1=mod,0=normal
	}
	return
}

func ObjectLocations(player *entity.Player, newObjects []*entity.Object, removingObjects []*entity.Object) (p *Packet) {
	p = NewOutgoingPacket(27)
	for _, o := range removingObjects {
		if o.Boundary {
			continue
		}
		p.AddShort(32767)
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		p.AddByte(byte(o.Direction))
		player.LocalObjects.RemoveObject(o)
	}
	for _, o := range newObjects {
		if o.Boundary {
			continue
		}
		p.AddShort(uint16(o.ID))
		p.AddByte(byte(o.X() - player.X()))
		p.AddByte(byte(o.Y() - player.Y()))
		p.AddByte(byte(o.Direction))
		player.LocalObjects.AddObject(o)
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
func PlayerInfo(player *entity.Player) *Packet {
	playerInfo := NewOutgoingPacket(131)
	playerInfo.AddShort(uint16(player.Index))
	playerInfo.AddShort(2304)
	playerInfo.AddShort(1776)

	playerInfo.AddShort(uint16((player.Y() + 100) / 1000))

	playerInfo.AddShort(944)
	return playerInfo
}
