package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["chatmsg"] = func(player *world.Player, p *packet.Packet) {
		for _, p1 := range player.NearbyPlayers() {
			if !p1.ChatBlocked() || p1.Friends(player.UserBase37) {
				p1.SendPacket(packetbuilders.PlayerChat(player.Index, string(p.Payload)))
			}
		}
	}
	PacketHandlers["addfriend"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(packetbuilders.FriendList(player))
		}()
		if player.Friends(hash) {
			player.SendPacket(packetbuilders.ServerMessage(("@que@You are already friends with that person!")))
			return
		}
		if player.Ignoring(hash) {
			player.SendPacket(packetbuilders.ServerMessage("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your ignore list before friending them."))
			return
		}
		if c1, ok := players.FromUserHash(hash); ok {
			if c1.Friends(player.UserBase37) && player.FriendBlocked() {
				c1.SendPacket(packetbuilders.FriendUpdate(player.UserBase37, true))
			}
			if !c1.FriendBlocked() || c1.Friends(player.UserBase37) {
				player.FriendList[hash] = true
				return
			}
		}
		player.FriendList[hash] = false
	}
	PacketHandlers["privmsg"] = func(player *world.Player, p *packet.Packet) {
		if c1, ok := players.FromUserHash(p.ReadLong()); ok {
			if !c1.FriendBlocked() || c1.Friends(player.UserBase37) {
				c1.SendPacket(packetbuilders.PrivateMessage(player.UserBase37, strutil.ChatFilter.Format(strutil.ChatFilter.Unpack(p.Payload[8:]))))
			}
		}
	}
	PacketHandlers["removefriend"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(packetbuilders.FriendList(player))
		}()
		if !player.Friends(hash) {
			player.SendPacket(packetbuilders.ServerMessage("@que@You are not friends with that person!"))
			return
		}
		if c1, ok := players.FromUserHash(hash); ok && c1.Friends(player.UserBase37) && player.FriendBlocked() {
			c1.SendPacket(packetbuilders.FriendUpdate(player.UserBase37, false))
		}
		delete(player.FriendList, hash)
	}
	PacketHandlers["addignore"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(packetbuilders.IgnoreList(player))
		}()
		if player.Friends(hash) {
			player.SendPacket(packetbuilders.ServerMessage("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your friend list before ignoring them."))
			return
		}
		if player.Ignoring(hash) {
			player.SendPacket(packetbuilders.ServerMessage("@que@You are already ignoring that person!"))
			return
		}
		player.IgnoreList = append(player.IgnoreList, hash)
	}
	PacketHandlers["removeignore"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(packetbuilders.IgnoreList(player))
		}()
		if !player.Ignoring(hash) {
			player.SendPacket(packetbuilders.ServerMessage("@que@You are not ignoring that person!"))
			return
		}
		for i, v := range player.IgnoreList {
			if v == hash {
				newSize := len(player.IgnoreList) - 1
				player.IgnoreList[i] = player.IgnoreList[newSize]
				player.IgnoreList = player.IgnoreList[:newSize]
				return
			}
		}
	}
	PacketHandlers["chooseoption"] = func(player *world.Player, p *packet.Packet) {
		choice := p.ReadByte()
		if !player.HasState(world.MSMenuChoosing) {
			return
		}
		if choice < 0 {
			return
		}
		player.RemoveState(world.MSMenuChoosing)
		player.OptionMenuC <- int8(choice)
	}
}
