package packethandlers

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/collections"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["chatmsg"] = func(c collections.Client, p *packetbuilders.Packet) {
		for _, p1 := range c.Player().NearbyPlayers() {
			if !p1.ChatBlocked() || p1.Friends(c.Player().UserBase37) {
				if c1, ok := collections.Clients.FromIndex(p1.Index); ok && !p1.Ignoring(c.Player().UserBase37) {
					c1.SendPacket(packetbuilders.PlayerChat(c.Player().Index, string(p.Payload)))
				}
			}
		}
	}
	PacketHandlers["addfriend"] = func(c collections.Client, p *packetbuilders.Packet) {
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
		if c1, ok := collections.Clients.FromUserHash(hash); ok {
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
	PacketHandlers["privmsg"] = func(c collections.Client, p *packetbuilders.Packet) {
		if c1, ok := collections.Clients.FromUserHash(p.ReadLong()); ok {
			if !c1.Player().FriendBlocked() || c1.Player().Friends(c.Player().UserBase37) {
				c1.SendPacket( packetbuilders.PrivateMessage(c.Player().UserBase37, strutil.ChatFilter.Format(strutil.ChatFilter.Unpack(p.Payload[8:]))))
			}
		}
	}
	PacketHandlers["removefriend"] = func(c collections.Client, p *packetbuilders.Packet) {
		hash := p.ReadLong()
		defer func() {
			c.SendPacket(packetbuilders.FriendList(c.Player()))
		}()
		if !c.Player().Friends(hash) {
			c.Message("@que@You are not friends with that person!")
			return
		}
		if c1, ok := collections.Clients.FromUserHash(hash); ok && c1.Player().Friends(c.Player().UserBase37) && c.Player().FriendBlocked() {
			c1.SendPacket( packetbuilders.FriendUpdate(c.Player().UserBase37, false))
		}
		delete(c.Player().FriendList, hash)
	}
	PacketHandlers["addignore"] = func(c collections.Client, p *packetbuilders.Packet) {
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
	PacketHandlers["removeignore"] = func(c collections.Client, p *packetbuilders.Packet) {
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
