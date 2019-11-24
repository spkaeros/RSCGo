package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["chatmsg"] = func(c clients.Client, p *packet.Packet) {
		for _, p1 := range c.Player().NearbyPlayers() {
			if !p1.ChatBlocked() || p1.Friends(c.Player().UserBase37) {
				p1.SendPacket(packetbuilders.PlayerChat(c.Player().Index, string(p.Payload)))
			}
		}
	}
	PacketHandlers["addfriend"] = func(c clients.Client, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.FriendList(c.Player()))
		}()
		if c.Player().Friends(hash) {
			c.Message("@que@You are already friends with that person!")
			return
		}
		if c.Player().Ignoring(hash) {
			c.Message("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your ignore list before friending them.")
			return
		}
		if c1, ok := clients.FromUserHash(hash); ok {
			if c1.Player().Friends(c.Player().UserBase37) && c.Player().FriendBlocked() {
				c1.SendPacket(packetbuilders.FriendUpdate(c.Player().UserBase37, true))
			}
			if !c1.Player().FriendBlocked() || c1.Player().Friends(c.Player().UserBase37) {
				c.Player().FriendList[hash] = true
				return
			}
		}
		c.Player().FriendList[hash] = false
	}
	PacketHandlers["privmsg"] = func(c clients.Client, p *packet.Packet) {
		if c1, ok := clients.FromUserHash(p.ReadLong()); ok {
			if !c1.Player().FriendBlocked() || c1.Player().Friends(c.Player().UserBase37) {
				c1.SendPacket(packetbuilders.PrivateMessage(c.Player().UserBase37, strutil.ChatFilter.Format(strutil.ChatFilter.Unpack(p.Payload[8:]))))
			}
		}
	}
	PacketHandlers["removefriend"] = func(c clients.Client, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.FriendList(c.Player()))
		}()
		if !c.Player().Friends(hash) {
			c.Message("@que@You are not friends with that person!")
			return
		}
		if c1, ok := clients.FromUserHash(hash); ok && c1.Player().Friends(c.Player().UserBase37) && c.Player().FriendBlocked() {
			c1.SendPacket(packetbuilders.FriendUpdate(c.Player().UserBase37, false))
		}
		delete(c.Player().FriendList, hash)
	}
	PacketHandlers["addignore"] = func(c clients.Client, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.IgnoreList(c.Player()))
		}()
		if c.Player().Friends(hash) {
			c.Message("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your friend list before ignoring them.")
			return
		}
		if c.Player().Ignoring(hash) {
			c.Message("@que@You are already ignoring that person!")
			return
		}
		c.Player().IgnoreList = append(c.Player().IgnoreList, hash)
	}
	PacketHandlers["removeignore"] = func(c clients.Client, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.IgnoreList(c.Player()))
		}()
		if !c.Player().Ignoring(hash) {
			c.Message("@que@You are not ignoring that person!")
			return
		}
		for i, v := range c.Player().IgnoreList {
			if v == hash {
				newSize := len(c.Player().IgnoreList) - 1
				c.Player().IgnoreList[i] = c.Player().IgnoreList[newSize]
				c.Player().IgnoreList = c.Player().IgnoreList[:newSize]
				return
			}
		}
	}
}
