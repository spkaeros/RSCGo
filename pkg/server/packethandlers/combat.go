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
	"time"
)

func init() {
	PacketHandlers["attacknpc"] = func(player *world.Player, p *packet.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("player[%v] tried to attack nil NPC\n", player)
			return
		}
		if player.Busy() {
			return
		}
		if !world.NpcDefs[npc.ID].Attackable {
			log.Info.Println("Player attacked not attackable NPC!", world.NpcDefs[npc.ID])
			return
		}
		log.Info.Println(npc.ID)
		player.SetDistancedAction(func() bool {
			if player.NextTo(npc.Location) && player.WithinRange(npc.Location, 1) {
				go func() {
					for _, trigger := range script.NpcAtkTriggers {
						if trigger(player, npc) {
							return
						}
					}
					if time.Since(npc.TransAttrs.VarTime("lastFight")) <= time.Second*2 || npc.Busy() {
						return
					}
					player.ResetPath()
					npc.ResetPath()
					player.SetLocation(npc.Location, true)
					player.SetRegionRemoved()
					player.AddState(world.MSFighting)
					npc.AddState(world.MSFighting)
					player.SetDirection(world.LeftFighting)
					npc.SetDirection(world.RightFighting)
					player.SetFightTarget(npc)
					npc.SetFightTarget(player)
					ticker := time.NewTicker(time.Millisecond * 1200)
					defer ticker.Stop()
					curRound := 0
					for range ticker.C {
						if !player.HasState(world.MSFighting) || !player.Connected() {
							if npc.HasState(world.MSFighting) {
								npc.ResetFighting()
							}
							return
						}
						var attacker, defender world.MobileEntity
						var nextHit int
						if curRound%2 == 0 {
							attacker = player
							defender = npc
						} else {
							attacker = npc
							defender = player
						}
						attacker.Transients().SetVar("fightRound", attacker.Transients().VarInt("fightRound", 0)+1)
						curRound++
						if player.PrayerActivated(12) && attacker == npc {
							continue
						}
						nextHit = attacker.MeleeDamage(defender)
						if curHits := defender.Skills().Current(world.StatHits); nextHit > curHits {
							nextHit = curHits
						}
						defender.Skills().DecreaseCur(world.StatHits, nextHit)
						if defender.Skills().Current(world.StatHits) <= 0 {
							if defenderNpc, ok := defender.(*world.NPC); ok {
								for _, trigger := range script.NpcDeathTriggers {
									if trigger.Check(player, npc) {
										// TODO: maybe blocking for special corner cases?
										go trigger.Action(player, npc)
									}
								}

								if attackerPlayer, ok := attacker.(*world.Player); ok {
									attackerPlayer.PlaySound("victory")
									world.AddItem(world.NewGroundItemFor(attackerPlayer.UsernameHash(), 20, 1, defender.X(), defender.Y()))
								} else {
									world.AddItem(world.NewGroundItem(20, 1, defender.X(), defender.Y()))
								}

								defenderNpc.Skills().SetCur(world.StatHits, defenderNpc.Skills().Maximum(world.StatHits))
								defenderNpc.SetLocation(world.DeathPoint, true)
								attacker.ResetFighting()

								go func() {
									time.Sleep(time.Second * 10)
									defenderNpc.SetLocation(defenderNpc.StartPoint, true)
								}()
							} else if defenderPlayer, ok := defender.(*world.Player); ok {
								defenderPlayer.Killed(nil)
							}
							return
						}

						if defenderNpc, ok := defender.(*world.NPC); ok {
							hitUpdate := world.NpcDamage(defenderNpc, nextHit)
							player.SendPacket(hitUpdate)
							for _, p1 := range player.NearbyPlayers() {
								p1.SendPacket(hitUpdate)
							}
						} else if defenderPlayer, ok := defender.(*world.Player); ok {
							defenderPlayer.Damage(nextHit)
						}
					}
				}()
				return true
			} else {
				player.SetPath(world.MakePath(player.Location, npc.Location))
			}
			return false
		})
	}
	PacketHandlers["attackplayer"] = func(player *world.Player, p *packet.Packet) {
		affectedClient, ok := world.Players.FromIndex(p.ReadShort())
		if affectedClient == nil || !ok {
			log.Suspicious.Printf("player[%v] tried to attack nil player\n", player)
			return
		}
		if player.Busy() {
			return
		}
		if affectedClient.Busy() {
			log.Info.Printf("Target player busy during attack request  State: %d\n", affectedClient.State)
			return
		}
		affectedPlayer := affectedClient
		player.SetDistancedAction(func() bool {
			if player.NextTo(affectedPlayer.Location) && player.WithinRange(affectedPlayer.Location, 2) {
				player.ResetPath()
				if time.Since(affectedPlayer.TransAttrs.VarTime("lastRetreat")) <= time.Second*3 || affectedPlayer.IsFighting() {
					return true
				}
				affectedPlayer.ResetPath()
				affectedPlayer.PlaySound("underattack")
				player.SetLocation(affectedPlayer.Location, true)
				player.AddState(world.MSFighting)
				affectedPlayer.AddState(world.MSFighting)
				player.SetDirection(world.LeftFighting)
				affectedPlayer.SetDirection(world.RightFighting)
				player.Transients().SetVar("fightTarget", affectedPlayer)
				affectedPlayer.Transients().SetVar("fightTarget", player)
				go func() {
					ticker := time.NewTicker(time.Millisecond * 1200)
					defer ticker.Stop()
					curRound := 0
					for range ticker.C {
						if !affectedPlayer.HasState(world.MSFighting) || !player.HasState(world.MSFighting) || !player.Connected() || !affectedPlayer.Connected() {
							if affectedPlayer.HasState(world.MSFighting) {
								affectedPlayer.ResetFighting()
							}
							if player.HasState(world.MSFighting) {
								player.ResetFighting()
							}
							return
						}
						var attacker, defender *world.Player
						if curRound%2 == 0 {
							attacker = player
							defender = affectedPlayer
						} else {
							attacker = affectedClient
							defender = player
						}
						nextHit := attacker.MeleeDamage(defender)
						if nextHit > defender.Skills().Current(world.StatHits) {
							nextHit = defender.Skills().Current(world.StatHits)
						}
						defender.Skills().DecreaseCur(world.StatHits, nextHit)
						if defender.Skills().Current(world.StatHits) <= 0 {
							defender.Killed(attacker)
							return
						}
						defender.Damage(nextHit)

						attacker.Transients().SetVar("fightRound", attacker.Transients().VarInt("fightRound", 0)+1)
						curRound++
					}
				}()
				return true
			}
			return player.FinishedPath()
		})
	}
	PacketHandlers["fightmode"] = func(player *world.Player, p *packet.Packet) {
		mode := p.ReadByte()
		if mode < 0 || mode > 3 {
			log.Suspicious.Printf("Invalid fightmode selected (%v) by %v", mode, player.String())
			return
		}
		player.SetFightMode(int(mode))
	}
}
