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
	"context"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"reflect"
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
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
					player.ResetPath()
					objectAction(player, object)
					return true
				}
				return false
			}
			if player.NextTo(object.Location) && (player.WithinRange(bounds[0], 1) || player.WithinRange(bounds[1], 1)) {
				player.ResetPath()
				objectAction(player, object)
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
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if world.Objects[object.ID].Type == 2 || world.Objects[object.ID].Type == 3 {
				if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
					player.ResetPath()
					objectAction(player, object)
					return true
				}
				return false
			}
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && (player.WithinRange(bounds[0], 1) || player.WithinRange(bounds[1], 1)) {
				player.ResetPath()
				objectAction(player, object)
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
		if object == nil {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", player, x, y)
			return
		}
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				boundaryAction(player, object)
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
		if object == nil {
			log.Info.Println("Boundary not found.")
			log.Suspicious.Printf("Player %v attempted to use a non-existant boundary at %d,%d\n", player, x, y)
			return
		}
		bounds := object.Boundaries()
		player.SetDistancedAction(func() bool {
			if (player.NextTo(bounds[1]) || player.NextTo(bounds[0])) && player.X() >= bounds[0].X() && player.Y() >= bounds[0].Y() && player.X() <= bounds[1].X() && player.Y() <= bounds[1].Y() {
				player.ResetPath()
				boundaryAction(player, object)
				return true
			}
			return false
		})
	}
	PacketHandlers["talktonpc"] = func(player *world.Player, p *packet.Packet) {
		idx := p.ReadShort()
		npc := world.GetNpc(idx)
		if npc == nil {
			return
		}
		if player.Busy() {
			return
		}
		player.SetDistancedAction(func() bool {
			if player.NextTo(npc.Location) && player.WithinRange(npc.Location, 1) && !npc.Busy() {
				player.ResetPath()
				npc.ResetPath()
				player.AddState(world.MSChatting)
				npc.AddState(world.MSChatting)
				if player.Location.Equals(npc.Location) {
				outer:
					for offX := -1; offX <= 1; offX++ {
						for offY := -1; offY <= 1; offY++ {
							if offX == 0 && offY == 0 {
								continue
							}
							newLoc := world.NewLocation(player.X() + offX, player.Y() + offY)
							if player.NextTo(newLoc) {
								player.SetLocation(newLoc, true)
								break outer
							}
						}
					}
					if player.Location.Equals(npc.Location) {
						return false
					}
				}
				player.SetDirection(player.DirectionTo(npc.X(), npc.Y()))
				npc.SetDirection(npc.DirectionTo(player.X(), player.Y()))
				go func() {
					defer func() {
						player.RemoveState(world.MSChatting)
						npc.RemoveState(world.MSChatting)
					}()
					for _, fn := range script.NpcTriggers {
						ran, err := fn(context.Background(), reflect.ValueOf(player), reflect.ValueOf(npc))
						if !ran.IsValid() {
							continue
						}
						if !err.IsNil() {
							log.Info.Println(err)
							continue
						}
						if ran.Bool() {
							return
						}
					}
					player.Message("The " + world.NpcDefs[npc.ID].Name + " does not appear interested in talking")
				}()
				return true
			} else {
				player.SetPath(world.MakePath(player.Location, npc.Location))
			}
			return false
		})
	}
}

func objectAction(player *world.Player, object *world.Object) {
	if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	player.AddState(world.MSBusy)

	go func() {
		defer func() {
			player.RemoveState(world.MSBusy)
		}()
		for _, fn := range script.ObjectTriggers {
			ran, err := fn(context.Background(), reflect.ValueOf(player), reflect.ValueOf(object))
			if !ran.IsValid() {
				continue
			}
			if !err.IsNil() {
				log.Info.Println(err)
				continue
			}
			if ran.Bool() {
				return
			}
		}
		player.SendPacket(world.DefaultActionMessage)
		//		if !script.Run("objectAction", player, "object", object) {
		//			player.SendPacket(world.DefaultActionMessage)
		//		}
	}()
}

func boundaryAction(player *world.Player, object *world.Object) {
	if player.Busy() || world.GetObject(object.X(), object.Y()) != object {
		// If somehow we became busy, the object changed before arriving, we do nothing.
		return
	}
	player.AddState(world.MSBusy)
	go func() {
		defer func() {
			player.RemoveState(world.MSBusy)
		}()
		for _, fn := range script.BoundaryTriggers {
			ran, err := fn(context.Background(), reflect.ValueOf(player), reflect.ValueOf(object))
			if !ran.IsValid() {
				continue
			}
			if !err.IsNil() {
				log.Info.Println(err)
				continue
			}
			if ran.Bool() {
				return
			}
		}
		player.SendPacket(world.DefaultActionMessage)
		//		if !script.Run("boundaryAction", player, "object", object) {
		//			player.SendPacket(world.DefaultActionMessage)
		//		}
	}()
}
