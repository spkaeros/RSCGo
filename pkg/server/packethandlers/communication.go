/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	PacketHandlers["chatmsg"] = func(player *world.Player, p *packet.Packet) {
		for _, p1 := range player.NearbyPlayers() {
			if !p1.ChatBlocked() || p1.Friends(player.UserBase37) {
				p1.SendPacket(world.PlayerChat(player.Index, string(p.Payload)))
			}
		}
	}
	PacketHandlers["addfriend"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(world.FriendList(player))
		}()
		if player.Friends(hash) {
			player.Message("@que@You are already friends with that person!")
			return
		}
		if player.Ignoring(hash) {
			player.Message("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your ignore list before friending them.")
			return
		}
		if c1, ok := players.FromUserHash(hash); ok {
			if c1.Friends(player.UserBase37) && player.FriendBlocked() {
				c1.SendPacket(world.FriendUpdate(player.UserBase37, true))
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
				c1.SendPacket(world.PrivateMessage(player.UserBase37, strutil.ChatFilter.Format(strutil.ChatFilter.Unpack(p.Payload[8:]))))
			}
		}
	}
	PacketHandlers["removefriend"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(world.FriendList(player))
		}()
		if !player.Friends(hash) {
			player.Message("@que@You are not friends with that person!")
			return
		}
		if c1, ok := players.FromUserHash(hash); ok && c1.Friends(player.UserBase37) && player.FriendBlocked() {
			c1.SendPacket(world.FriendUpdate(player.UserBase37, false))
		}
		delete(player.FriendList, hash)
	}
	PacketHandlers["addignore"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(world.IgnoreList(player))
		}()
		if player.Friends(hash) {
			player.Message("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your friend list before ignoring them.")
			return
		}
		if player.Ignoring(hash) {
			player.Message("@que@You are already ignoring that person!")
			return
		}
		player.IgnoreList = append(player.IgnoreList, hash)
	}
	PacketHandlers["removeignore"] = func(player *world.Player, p *packet.Packet) {
		hash := p.ReadLong()
		defer func() {
			player.SendPacket(world.IgnoreList(player))
		}()
		if !player.Ignoring(hash) {
			player.Message("@que@You are not ignoring that person!")
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
		if !player.HasState(world.MSOptionMenu) {
			return
		}
		if choice < 0 {
			return
		}
		player.OptionMenuC <- int8(choice)
	}
}
