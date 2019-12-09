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
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

func init() {
	PacketHandlers["objectaction"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Object not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant object at %d,%d\n", player, x, y)
			return
		}
		player.SetDistancedAction(func() bool {
			if player.AtObject(object) {
				player.ResetPath()
				player.AddState(world.MSBusy)

				go func() {
					defer func() {
						player.RemoveState(world.MSBusy)
					}()

					for _, trigger := range script.ObjectTriggers {
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
	}
	PacketHandlers["objectaction2"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil {
			log.Info.Println("Object not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant object at %d,%d\n", player, x, y)
			return
		}
		player.SetDistancedAction(func() bool {
			if player.AtObject(object) {
				player.ResetPath()
				player.AddState(world.MSBusy)

				go func() {
					defer func() {
						player.RemoveState(world.MSBusy)
					}()

					for _, trigger := range script.ObjectTriggers {
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
	}
	PacketHandlers["boundaryaction2"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil || !object.Boundary {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", player, x, y)
			return
		}
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				boundaryAction(player, object, 1)
				return true
			}
			return false
		})
	}
	PacketHandlers["boundaryaction"] = func(player *world.Player, p *packet.Packet) {
		if player.Busy() {
			return
		}
		x := p.ReadShort()
		y := p.ReadShort()
		object := world.GetObject(x, y)
		if object == nil || !object.Boundary {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", player, x, y)
			return
		}
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				boundaryAction(player, object, 0)
				return true
			}
			return false
		})
	}
	PacketHandlers["talktonpc"] = func(player *world.Player, p *packet.Packet) {
		idx := p.ReadShort()
		npc := world.GetNpc(idx)
		if npc == nil{
			return
		}
		player.SetDistancedAction(func() bool {
			if player.Busy() {
				player.ResetPath()
				return true
			}
			if npc.Busy() {
				player.Message(npc.Name() + " is busy at the moment")
				player.ResetPath()
				return true
			}
			if player.WithinRange(npc.Location, 1) && player.NextTo(npc.Location) {
				player.ResetPath()
				for _, triggerDef := range script.NpcTriggers {
					if triggerDef.Check(npc) {
						npc.ResetPath()
						if player.Location.Equals(npc.Location) {
						outer:
							for offX := -1; offX <= 1; offX++ {
								for offY := -1; offY <= 1; offY++ {
									if offX == 0 && offY == 0 {
										continue
									}
									newLoc := world.NewLocation(player.X()+offX, player.Y()+offY)
									switch player.DirectionTo(newLoc.X(), newLoc.Y()) {
									case world.North:
										if world.IsTileBlocking(newLoc.X(), newLoc.Y(), world.ClipSouth, false) {
											continue
										}
									case world.South:
										if world.IsTileBlocking(newLoc.X(), newLoc.Y(), world.ClipNorth, false) {
											continue
										}
									case world.East:
										if world.IsTileBlocking(newLoc.X(), newLoc.Y(), world.ClipWest, false) {
											continue
										}
									case world.West:
										if world.IsTileBlocking(newLoc.X(), newLoc.Y(), world.ClipEast, false) {
											continue
										}
									case world.NorthWest:
										if world.IsTileBlocking(player.X(), player.Y()-1, world.ClipSouth, false) {
											continue
										}
										if world.IsTileBlocking(player.X()+1, player.Y(), world.ClipEast, false) {
											continue
										}
									case world.NorthEast:
										if world.IsTileBlocking(player.X(), player.Y()-1, world.ClipSouth, false) {
											continue
										}
										if world.IsTileBlocking(player.X()-1, player.Y(), world.ClipWest, false) {
											continue
										}
									case world.SouthWest:
										if world.IsTileBlocking(player.X(), player.Y()+1, world.ClipNorth, false) {
											continue
										}
										if world.IsTileBlocking(player.X()+1, player.Y(), world.ClipEast, false) {
											continue
										}
									case world.SouthEast:
										if world.IsTileBlocking(player.X(), player.Y()+1, world.ClipNorth, false) {
											continue
										}
										if world.IsTileBlocking(player.X()-1, player.Y(), world.ClipWest, false) {
											continue
										}
									}
									if player.NextTo(newLoc) {
										player.SetLocation(newLoc, false)
										break outer
									}
								}
							}
						}

						if !player.Location.Equals(npc.Location) {
							player.SetDirection(player.DirectionTo(npc.X(), npc.Y()))
							npc.SetDirection(npc.DirectionTo(player.X(), player.Y()))
						}
						player.AddState(world.MSChatting)
						npc.AddState(world.MSChatting)
						go func() {
							defer func() {
								player.RemoveState(world.MSChatting)
								npc.RemoveState(world.MSChatting)
							}()
							triggerDef.Action(player, npc)
						}()
						return true
					}
				}
				player.Message("The " + npc.Name() + " does not appear interested in talking")
				return true
			}

			player.WalkTo(npc.Location)
			return false
		})
	}
	PacketHandlers["invonboundary"] = func(player *world.Player, p *packet.Packet) {
		targetX := p.ReadShort()
		targetY := p.ReadShort()
		p.ReadByte() // dir, useful?
		invIndex := p.ReadShort()

		object := world.GetObject(targetX, targetY)
		if object == nil || !object.Boundary {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", player, targetX, targetY)
			return
		}
		if invIndex >= player.Inventory.Size() {
			log.Suspicious.Printf("Player %v attempted to use a non-existant item(idx:%v, cap:%v) on a boundary at %d,%d\n", player, invIndex, player.Inventory.Size()-1, targetX, targetY)
			return
		}
		invItem := player.Inventory.Get(invIndex)
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				player.AddState(world.MSBusy)
				go func() {
					defer func() {
						player.RemoveState(world.MSBusy)
					}()
					for _, fn := range script.InvOnBoundaryTriggers {
						if fn(player, object, invItem) {
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			return false
		})
	}

	PacketHandlers["invonobject"] = func(player *world.Player, p *packet.Packet) {
		targetX := p.ReadShort()
		targetY := p.ReadShort()
		invIndex := p.ReadShort()

		object := world.GetObject(targetX, targetY)
		if object == nil || object.Boundary {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", player, targetX, targetY)
			return
		}
		if invIndex >= player.Inventory.Size() {
			log.Suspicious.Printf("Player %v attempted to use a non-existant item(idx:%v, cap:%v) on a boundary at %d,%d\n", player, invIndex, player.Inventory.Size()-1, targetX, targetY)
			return
		}
		invItem := player.Inventory.Get(invIndex)
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if world.ObjectDefs[object.ID].Type == 2 || world.ObjectDefs[object.ID].Type == 3 {
				if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
					player.ResetPath()
					player.AddState(world.MSBusy)
					go func() {
						defer func() {
							player.RemoveState(world.MSBusy)
						}()
						for _, fn := range script.InvOnObjectTriggers {
							if fn(player, object, invItem) {
								return
							}
						}
						player.SendPacket(world.DefaultActionMessage)
					}()
					return true
				}
				return false
			}
			if player.AtObject(object) {
				player.ResetPath()
				player.AddState(world.MSBusy)
				go func() {
					defer func() {
						player.RemoveState(world.MSBusy)
					}()
					for _, fn := range script.InvOnObjectTriggers {
						if fn(player, object, invItem) {
							return
						}
					}
					player.SendPacket(world.DefaultActionMessage)
				}()
				return true
			}
			return false
		})
	}
}


func boundaryAction(player *world.Player, object *world.Object, click int) {
	if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	player.AddState(world.MSBusy)
	go func() {
		defer func() {
			player.RemoveState(world.MSBusy)
		}()

		for _, trigger := range script.BoundaryTriggers {
			if trigger.Check(object, click) {
				trigger.Action(player, object, click)
				return
			}
		}
		player.SendPacket(world.DefaultActionMessage)
	}()
}
