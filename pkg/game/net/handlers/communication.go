/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package handlers

import (
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

func init() {
	game.AddHandler("chatmsg", func(player *world.Player, p *net.Packet) {
		for _, p1 := range player.NearbyPlayers() {
			if !p1.ChatBlocked() || p1.FriendsWith(player.UsernameHash()) {
				//p1.SendPacket(world.PlayerChat(player.Index, string(p.FrameBuffer)))
				p1.QueuePublicChat(player, string(p.FrameBuffer))
			}
		}
	})
	game.AddHandler("addfriend", func(player *world.Player, p *net.Packet) {
		hash := p.ReadUint64()
		defer func() {
			player.SendPacket(world.FriendList(player))
		}()
		if player.FriendsWith(hash) {
			player.Message("@que@You are already friends with that person!")
			return
		}
		if player.Ignoring(hash) {
			player.Message("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your ignore list before friending them.")
			return
		}
		if p1, ok := world.Players.FindHash(hash); ok && p1 != nil &&
			(!p1.FriendBlocked() || p1.FriendList.ContainsHash(hash)) {
			player.FriendList.Add(strutil.Base37.Decode(hash))
			p1.SendPacket(world.FriendUpdate(player.UsernameHash(), true))
		}
	})
	game.AddHandler("privmsg", func(player *world.Player, p *net.Packet) {
		hash := p.ReadUint64()
		if p1, ok := world.Players.FindHash(hash); ok && p1 != nil &&
			(!p1.FriendBlocked() || p1.FriendList.ContainsHash(hash)) {
			// c1.SendPacket(world.PrivateMessage(player.UsernameHash(), strutil.ChatFilter.Format(strutil.ChatFilter.Unpack(p.FrameBuffer[8:]))))
			p1.SendPacket(world.PrivateMessage(player.UsernameHash(), strutil.ChatFilter.Format(string(p.FrameBuffer[8:]))))
		}
	})
	game.AddHandler("removefriend", func(player *world.Player, p *net.Packet) {
		hash := p.ReadUint64()
		defer func() {
			player.SendPacket(world.FriendList(player))
		}()
		if !player.FriendsWith(hash) {
			player.Message("@que@You are not friends with that person!")
			return
		}
		player.FriendList.Remove(strutil.Base37.Decode(hash))
		if player.FriendBlocked() {
			if p1, ok := world.Players.FindHash(hash); ok && p1 != nil &&
				p1.FriendList.ContainsHash(hash) {
				p1.FriendList.ToggleStatus(player.Username())
			}
		}
	})
	game.AddHandler("addignore", func(player *world.Player, p *net.Packet) {
		hash := p.ReadUint64()
		defer func() {
			player.SendPacket(world.IgnoreList(player))
		}()
		if player.FriendsWith(hash) {
			player.Message("@que@Please remove '" + strutil.Base37.Decode(hash) + "' from your friend list before ignoring them.")
			return
		}
		if player.Ignoring(hash) {
			player.Message("@que@You are already ignoring that person!")
			return
		}
		player.IgnoreList = append(player.IgnoreList, hash)
	})
	game.AddHandler("removeignore", func(player *world.Player, p *net.Packet) {
		hash := p.ReadUint64()
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
	})
	game.AddHandler("chooseoption", func(player *world.Player, p *net.Packet) {
		choice := p.ReadUint8()
		if player.VarInt("state", 0)&world.StateChatChoosing&^world.MSItemAction == 0 {
			return
		}
		if choice < 0 {
			return
		}
		player.ReplyMenuC <- int8(choice)
	})
}
