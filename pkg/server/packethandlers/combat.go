package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"time"
)

func init() {
	PacketHandlers["attacknpc"] = func(c *world.Player, p *packet.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("player[%v] tried to attack nil NPC\n", c)
			return
		}
		if c.Busy() {
			return
		}
		if !world.NpcDefs[npc.ID].Attackable {
			log.Info.Println("Player attacked not attackable NPC!", world.NpcDefs[npc.ID])
			return
		}
		log.Info.Println(npc.ID)
		c.SetDistancedAction(func() bool {
			if c.NextTo(npc.Location) && c.WithinRange(npc.Location, 2) {
				c.ResetPath()
				npc.ResetPath()
				c.SetLocation(npc.Location, true)
				c.AddState(world.MSFighting)
				npc.AddState(world.MSFighting)
				c.SetDirection(world.LeftFighting)
				npc.SetDirection(world.RightFighting)
				c.SetFightTarget(npc)
				npc.SetFightTarget(c)
				go func() {
					ticker := time.NewTicker(time.Millisecond * 1200)
					defer ticker.Stop()
					curRound := 0
					for range ticker.C {
						if !c.HasState(world.MSFighting) || !c.Connected() {
							if npc.HasState(world.MSFighting) {
								script.EngineChannel <- func() {
									npc.ResetFighting()
								}
							}
							return
						}
						var attacker, defender world.MobileEntity
						var nextHit int
						if curRound%2 == 0 {
							attacker = c
							defender = npc
						} else {
							attacker = npc
							defender = c
						}
						nextHit = attacker.MeleeDamage(defender)
						if curHits := defender.Skills().Current(world.StatHits); nextHit > curHits {
							nextHit = curHits
						}
						defender.Skills().DecreaseCur(world.StatHits, nextHit)
						if defender.Skills().Current(world.StatHits) <= 0 {
							if defenderNpc, ok := defender.(*world.NPC); ok {
								script.EngineChannel <- func() {
									if attackerPlayer, ok := attacker.(*world.Player); ok {
										attackerPlayer.SendPacket(packetbuilders.Sound("victory"))
										world.AddItem(world.NewGroundItemFor(attackerPlayer.UserBase37, 20, 1, defender.X(), defender.Y()))
									} else {
										world.AddItem(world.NewGroundItem(20, 1, defender.X(), defender.Y()))
									}
									attacker.ResetFighting()
									defenderNpc.Skills().SetCur(world.StatHits, defenderNpc.Skills().Maximum(world.StatHits))
									defenderNpc.SetLocation(world.DeathPoint, true)
								}

								go func() {
									time.Sleep(time.Second * 10)
									script.EngineChannel <- func() {
										defenderNpc.SetLocation(defenderNpc.StartPoint, true)
									}
								}()
							} else if defenderPlayer, ok := defender.(*world.Player); ok {
								script.EngineChannel <- func() {
									attacker.ResetFighting()
									world.AddItem(world.NewGroundItem(20, 1, defender.X(), defender.Y()))
									for i := 0; i < 18; i++ {
										defenderPlayer.Skills().SetCur(i, defenderPlayer.Skills().Maximum(i))
									}
									defenderPlayer.SendPacket(packetbuilders.PlayerStats(defenderPlayer))
									defenderPlayer.SendPacket(packetbuilders.Death)
									defenderPlayer.SendPacket(packetbuilders.Sound("death"))
									defenderPlayer.Transients().SetVar("deathTime", time.Now())
									// TODO: Keep 3 most valuable items
									defenderPlayer.Inventory().Range(func(item *world.Item) bool {
										if item.Worn {
											defenderPlayer.DequipItem(item)
										}
										world.AddItem(world.NewGroundItem(item.ID, item.Amount, defender.X(), defender.Y()))
										return true
									})
									defenderPlayer.Inventory().Clear()
									defenderPlayer.SendPacket(packetbuilders.InventoryItems(defenderPlayer))
									defenderPlayer.SendPacket(packetbuilders.EquipmentStats(defenderPlayer))
									plane := defenderPlayer.Plane()
									defenderPlayer.SetLocation(world.SpawnPoint, true)
									if defenderPlayer.Plane() != plane {
										defenderPlayer.SendPacket(packetbuilders.PlaneInfo(defenderPlayer))
									}
								}
							}
							return
						}

						if defenderNpc, ok := defender.(*world.NPC); ok {
							hitUpdate := packetbuilders.NpcDamage(defenderNpc, nextHit)
							c.SendPacket(hitUpdate)
							for _, p1 := range c.NearbyPlayers() {
								p1.SendPacket(hitUpdate)
							}
						} else if defenderPlayer, ok := defender.(*world.Player); ok {
							hitUpdate := packetbuilders.PlayerDamage(defenderPlayer, nextHit)
							c.SendPacket(hitUpdate)
							for _, p1 := range c.NearbyPlayers() {
								p1.SendPacket(hitUpdate)
							}
						}

						attacker.Transients().SetVar("fightRound", attacker.Transients().VarInt("fightRound", 0)+1)
						curRound++
					}
				}()
				return true
			} else {
				c.SetPath(world.MakePath(c.Location, npc.Location))
			}
			return false
		})
	}
	PacketHandlers["attackplayer"] = func(c *world.Player, p *packet.Packet) {
		affectedClient, ok := players.FromIndex(p.ReadShort())
		if affectedClient == nil || !ok {
			log.Suspicious.Printf("player[%v] tried to attack nil player\n", c)
			return
		}
		if c.Busy() {
			return
		}
		if affectedClient.Busy() {
			log.Info.Printf("Target player busy during attack request  State: %d\n", affectedClient.State)
			return
		}
		affectedPlayer := affectedClient
		c.SetDistancedAction(func() bool {
			if c.NextTo(affectedPlayer.Location) && c.WithinRange(affectedPlayer.Location, 2) {
				c.ResetPath()
				if time.Since(affectedPlayer.TransAttrs.VarTime("lastRetreat")) <= time.Second*3 {
					return false
				}
				affectedPlayer.ResetPath()
				affectedPlayer.SendPacket(packetbuilders.Sound("underattack"))
				c.SetLocation(affectedPlayer.Location, true)
				c.AddState(world.MSFighting)
				affectedPlayer.AddState(world.MSFighting)
				c.SetDirection(world.LeftFighting)
				affectedPlayer.SetDirection(world.RightFighting)
				c.Transients().SetVar("fightTarget", affectedPlayer)
				affectedPlayer.Transients().SetVar("fightTarget", c)
				go func() {
					ticker := time.NewTicker(time.Millisecond * 1200)
					defer ticker.Stop()
					curRound := 0
					for range ticker.C {
						if !affectedPlayer.HasState(world.MSFighting) || !c.HasState(world.MSFighting) || !c.Connected() || !affectedPlayer.Connected() {
							if affectedPlayer.HasState(world.MSFighting) {
								script.EngineChannel <- func() {
									affectedPlayer.ResetFighting()
								}
							}
							if c.HasState(world.MSFighting) {
								script.EngineChannel <- func() {
									c.ResetFighting()
								}
							}
							return
						}
						var attacker, defender *world.Player
						if curRound%2 == 0 {
							attacker = c
							defender = affectedPlayer
						} else {
							attacker = affectedClient
							defender = c
						}
						nextHit := attacker.MeleeDamage(defender)
						if nextHit > defender.Skills().Current(world.StatHits) {
							nextHit = defender.Skills().Current(world.StatHits)
						}
						defender.Skills().DecreaseCur(world.StatHits, nextHit)
						if defender.Skills().Current(world.StatHits) <= 0 {
							script.EngineChannel <- func() {
								attacker.ResetFighting()
								world.AddItem(world.NewGroundItem(20, 1, defender.X(), defender.Y()))
								attacker.SendPacket(packetbuilders.Sound("victory"))
								defender.SendPacket(packetbuilders.Sound("death"))
								// TODO: Keep 3 most valuable items
								defender.Inventory().Range(func(item *world.Item) bool {
									world.AddItem(world.NewGroundItemFor(attacker.UserBase37, item.ID, item.Amount, defender.X(), defender.Y()))
									return true
								})
								defender.Inventory().Clear()
								attacker.SendPacket(packetbuilders.ServerMessage("You have defeated " + defender.Username + "!"))
								defender.Skills().SetCur(world.StatHits, defender.Skills().Maximum(world.StatHits))
								defender.SendPacket(packetbuilders.PlayerStats(defender))
								defender.Transients().SetVar("deathTime", time.Now())
								defender.SendPacket(packetbuilders.Death)
								defender.SetLocation(world.SpawnPoint, true)
								if defender.Plane() != world.SpawnPoint.Plane() {
									defender.SendPacket(packetbuilders.PlaneInfo(defender))
								}
							}
							return
						}
						hitUpdate := packetbuilders.PlayerDamage(defender, nextHit)
						c.SendPacket(hitUpdate)
						for _, p1 := range c.NearbyPlayers() {
							p1.SendPacket(hitUpdate)
						}

						attacker.Transients().SetVar("fightRound", attacker.Transients().VarInt("fightRound", 0)+1)
						curRound++
					}
				}()
				return true
			}
			return c.FinishedPath()
		})
	}
	PacketHandlers["fightmode"] = func(c *world.Player, p *packet.Packet) {
		mode := p.ReadByte()
		if mode < 0 || mode > 3 {
			log.Suspicious.Printf("Invalid fightmode selected (%v) by %v", mode, c.String())
			return
		}
		c.SetFightMode(int(mode))
	}
}
