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
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("invonboundary", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		targetX := p.ReadUint16()
		targetY := p.ReadUint16()
		p.ReadUint8() // dir, useful?
		invIndex := p.ReadUint16()

		object := world.GetObject(targetX, targetY)
		if object == nil || !object.Boundary {
			log.Suspicious.Printf("%v attempted to use a non-existent boundary at %d,%d\n", player, targetX, targetY)
			return
		}
		if invIndex >= player.Inventory.Size() {
			log.Suspicious.Printf("%v attempted to use a non-existent item(idx:%v, cap:%v) on a boundary at %d,%d\n", player, invIndex, player.Inventory.Size()-1, targetX, targetY)
			return
		}
		invItem := player.Inventory.Get(invIndex)
		bounds := object.Boundaries()
		player.SetTickAction(func() bool {
			if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
				// If somehow we became busy, or the object changed before arriving, we do nothing.
				return true
			}
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				player.AddState(world.MSBatching)
				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
					}()
					for _, fn := range world.InvOnBoundaryTriggers {
						if fn(player, object, invItem) {
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			player.WalkTo(object.Point())
			return false
		})
	})
	game.AddHandler("invonplayer", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		targetIndex := p.ReadUint16()
		invIndex := p.ReadUint16()

		if targetIndex == player.ServerIndex() {
			log.Suspicious.Printf("%s attempted to use an inventory item on themself\n", player.String())
			return
		}

		target, ok := world.Players.FindIndex(targetIndex)
		if !ok || target == nil || !target.Connected() {
			log.Suspicious.Printf("%s attempted to use an inventory item on a player that doesn't exist\n", player.String())
			return
		}
		if invIndex >= player.Inventory.Size() {
			log.Suspicious.Printf("%s attempted to use a non-existent item(idx:%v, cap:%v)  on a player(%s)\n", player.String(), invIndex, player.Inventory.Size()-1, target.String())
			return
		}
		invItem := player.Inventory.Get(invIndex)
		player.SetTickAction(func() bool {
			if player.Busy() || !player.Connected() || target == nil || target.Busy() || !target.Connected() {
				return true
			}
			if player.Near(target.Point(), 1) && player.NextTo(target.Point()) {
				player.ResetPath()
				player.AddState(world.MSBatching)
				target.AddState(world.MSBatching)
				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
						target.RemoveState(world.MSBatching)
					}()
					for _, trigger := range world.InvOnPlayerTriggers {
						if trigger.Check(invItem) {
							trigger.Action(player, target, invItem)
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			player.WalkTo(target.Point())
			return false
		})
	})
	game.AddHandler("invonobject", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		targetX := p.ReadUint16()
		targetY := p.ReadUint16()
		invIndex := p.ReadUint16()

		object := world.GetObject(targetX, targetY)
		if object == nil || object.Boundary {
			log.Suspicious.Printf("%v attempted to use a non-existent boundary at %d,%d\n", player, targetX, targetY)
			return
		}
		if invIndex >= player.Inventory.Size() {
			log.Suspicious.Printf("%v attempted to use a non-existent item(idx:%v, cap:%v) on a boundary at %d,%d\n", player, invIndex, player.Inventory.Size()-1, targetX, targetY)
			return
		}
		invItem := player.Inventory.Get(invIndex)
		bounds := object.Boundaries()
		player.WalkTo(object.Point())
		player.SetTickAction(func() bool {
			if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
				// If somehow we became busy, the object changed before arriving, we do nothing.
				return true
			}
			if definitions.Scenary(object.ID).CollisionType == 2 || definitions.Scenary(object.ID).CollisionType == 3 {
				if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
					player.ResetPath()
					player.AddState(world.MSBatching)
					go func() {
						defer func() {
							player.RemoveState(world.MSBatching)
						}()
						for _, fn := range world.InvOnObjectTriggers {
							if fn(player, object, invItem) {
								return
							}
						}
						player.SendPacket(world.DefaultActionMessage)
					}()
					return true
				}
				player.WalkTo(object.Point())
				return false
			}
			if player.AtObject(object) {
				player.ResetPath()
				player.AddState(world.MSBatching)
				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
					}()
					for _, fn := range world.InvOnObjectTriggers {
						if fn(player, object, invItem) {
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			player.WalkTo(object.Point())
			return false
		})
	})
}
