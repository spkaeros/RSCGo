package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["chatmsg"] = func(c *world.Player, p *packet.Packet) {
		for _, p1 := range c.NearbyPlayers() {
			if !p1.ChatBlocked() || p1.Friends(c.UserBase37) {
				p1.SendPacket(packetbuilders.PlayerChat(c.Index, string(p.Payload)))
			}
		}
	}
	PacketHandlers["addfriend"] = func(c *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.FriendList(c))
		}()
		if c.Friends(hash) {
			c.SendPacket(packetbuilders.ServerMessage(("@que@You are already friends with that person!")))
			return
		}
		if c.Ignoring(hash) {
			c.SendPacket(packetbuilders.ServerMessage("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your ignore list before friending them."))
			return
		}
		if c1, ok := players.FromUserHash(hash); ok {
			if c1.Friends(c.UserBase37) && c.FriendBlocked() {
				c1.SendPacket(packetbuilders.FriendUpdate(c.UserBase37, true))
			}
			if !c1.FriendBlocked() || c1.Friends(c.UserBase37) {
				c.FriendList[hash] = true
				return
			}
		}
		c.FriendList[hash] = false
	}
	PacketHandlers["privmsg"] = func(c *world.Player, p *packet.Packet) {
		if c1, ok := players.FromUserHash(p.ReadLong()); ok {
			if !c1.FriendBlocked() || c1.Friends(c.UserBase37) {
				c1.SendPacket(packetbuilders.PrivateMessage(c.UserBase37, strutil.ChatFilter.Format(strutil.ChatFilter.Unpack(p.Payload[8:]))))
			}
		}
	}
	PacketHandlers["removefriend"] = func(c *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.FriendList(c))
		}()
		if !c.Friends(hash) {
			c.SendPacket(packetbuilders.ServerMessage("@que@You are not friends with that person!"))
			return
		}
		if c1, ok := players.FromUserHash(hash); ok && c1.Friends(c.UserBase37) && c.FriendBlocked() {
			c1.SendPacket(packetbuilders.FriendUpdate(c.UserBase37, false))
		}
		delete(c.FriendList, hash)
	}
	PacketHandlers["addignore"] = func(c *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.IgnoreList(c))
		}()
		if c.Friends(hash) {
			c.SendPacket(packetbuilders.ServerMessage("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your friend list before ignoring them."))
			return
		}
		if c.Ignoring(hash) {
			c.SendPacket(packetbuilders.ServerMessage("@que@You are already ignoring that person!"))
			return
		}
		c.IgnoreList = append(c.IgnoreList, hash)
	}
	PacketHandlers["removeignore"] = func(c *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.IgnoreList(c))
		}()
		if !c.Ignoring(hash) {
			c.SendPacket(packetbuilders.ServerMessage("@que@You are not ignoring that person!"))
			return
		}
		for i, v := range c.IgnoreList {
			if v == hash {
				newSize := len(c.IgnoreList) - 1
				c.IgnoreList[i] = c.IgnoreList[newSize]
				c.IgnoreList = c.IgnoreList[:newSize]
				return
			}
		}
	}
	PacketHandlers["chooseoption"] = func(c *world.Player, p *packet.Packet) {
		choice := p.ReadByte()
		if !c.HasState(world.MSMenuChoosing) {
			return
		}
		if choice < 0 {
			return
		}
		c.RemoveState(world.MSMenuChoosing)
		c.OptionMenuC <- int8(choice)
	}
}
