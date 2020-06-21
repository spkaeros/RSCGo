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
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

func init() {
	game.AddHandler("objectaction", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		x := p.ReadUint16()
		y := p.ReadUint16()
		object := world.GetObject(x, y)
		if object == nil || object.Boundary {
			log.Suspicious.Printf("Player %v attempted to use a non-existent object at %d,%d\n", player, x, y)
			return
		}
		player.SetTickAction(func() bool {
			if player.AtObject(object) {
				player.ResetPath()
				player.AddState(world.MSBatching)

				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
					}()

					for _, trigger := range world.ObjectTriggers {
						if trigger.Check(object, 0) {
							trigger.Action(player, object, 0)
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			return false

		})
	})
	game.AddHandler("objectaction2", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		x := p.ReadUint16()
		y := p.ReadUint16()
		object := world.GetObject(x, y)
		if object == nil || object.Boundary {
			log.Suspicious.Printf("Player %v attempted to use a non-existent object at %d,%d\n", player, x, y)
			return
		}
		player.SetTickAction(func() bool {
			if player.AtObject(object) {
				player.ResetPath()
				player.AddState(world.MSBatching)

				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
					}()

					for _, trigger := range world.ObjectTriggers {
						if trigger.Check(object, 1) {
							trigger.Action(player, object, 1)
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()

				return true
			}
			return false
		})
	})
	game.AddHandler("boundaryaction2", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		x := p.ReadUint16()
		y := p.ReadUint16()
		object := world.GetObject(x, y)
		if object == nil || !object.Boundary {
			log.Suspicious.Printf("Player %v attempted to use a non-existent boundary at %d,%d\n", player, x, y)
			return
		}
		bounds := object.Boundaries()
		player.SetTickAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
					// If somehow we became busy, the object changed before arriving, we do nothing.
					return true
				}
				player.AddState(world.MSBatching)
				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
					}()

					for _, trigger := range world.BoundaryTriggers {
						if trigger.Check(object, 1) {
							trigger.Action(player, object, 1)
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			return false
		})
	})
	game.AddHandler("boundaryaction", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		x := p.ReadUint16()
		y := p.ReadUint16()
		object := world.GetObject(x, y)
		if object == nil || !object.Boundary {
			log.Suspicious.Printf("%v attempted to use a non-existent boundary at %d,%d\n", player, x, y)
			return
		}
		bounds := object.Boundaries()
		player.SetTickAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
					// If somehow we became busy, the object changed before arriving, we do nothing.
					return true
				}
				player.AddState(world.MSBatching)
				go func() {
					defer func() {
						player.RemoveState(world.MSBatching)
					}()

					for _, trigger := range world.BoundaryTriggers {
						if trigger.Check(object, 0) {
							trigger.Action(player, object, 0)
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			return false
		})
	})
	game.AddHandler("talktonpc", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		idx := p.ReadUint16()
		npc := world.GetNpc(idx)
		if npc == nil {
			return
		}
		if player.IsFighting() {
			return
		}
		player.WalkingArrivalAction(npc, 1, func() {
			player.ResetPath()
			if npc.Busy() {
				player.Message(npc.Name() + " is busy at the moment")
				return
			}
			if player.Busy() {
				return
			}
			for _, triggerDef := range world.NpcTriggers {
				if triggerDef.Check(npc) {
					npc.ResetPath()
					if player.Location.Equals(npc.Location) {
						for _, direction := range world.OrderedDirections {
							neighbor := npc.Step(direction)
							if npc.Reachable(neighbor) {
								npc.SetLocation(neighbor, true)
								break
							}
						}
					}

					if !player.Location.Equals(npc.Location) {
						player.SetDirection(player.DirectionTo(npc.X(), npc.Y()))
						npc.SetDirection(npc.DirectionTo(player.X(), player.Y()))
					}
					go func() {
						player.SetVar("targetMob", npc)
						defer func() {
							player.UnsetVar("targetMob")
							player.RemoveState(world.StateChatting)
							npc.RemoveState(world.StateChatting)
						}()
						player.AddState(world.StateChatting)
						npc.AddState(world.StateChatting)
						triggerDef.Action(player, npc)
					}()
					return
				}
			}
			player.Message("The " + npc.Name() + " does not appear interested in talking")
		})
	})
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
			player.WalkTo(object.Location)
			return false
		})
	})
	game.AddHandler("invonplayer", func(player *world.Player, p *net.Packet) {
		if player.Busy() || player.IsFighting() {
			return
		}
		targetIndex := p.ReadUint16()
		invIndex := p.ReadUint16()

		if targetIndex == player.Index {
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
			if player.WithinRange(target.Location, 1) && player.NextTo(target.Location) {
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
			player.WalkTo(target.Location)
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
		player.WalkTo(object.Location)
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
				player.WalkTo(object.Location)
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
			player.WalkTo(object.Location)
			return false
		})
	})
}
