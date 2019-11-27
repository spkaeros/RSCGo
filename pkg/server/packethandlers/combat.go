package packethandlers

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"time"
)

func init() {
	PacketHandlers["attacknpc"] = func(c clients.Client, p *packet.Packet) {
		npc := world.GetNpc(p.ReadShort())
		if npc == nil {
			log.Suspicious.Printf("player[%v] tried to attack nil NPC\n", c)
			return
		}
		if c.Player().State != world.MSIdle {
			return
		}
		if !world.NpcDefs[npc.ID].Attackable {
			log.Info.Println("Player attacked not attackable NPC!", world.NpcDefs[npc.ID])
			return
		}
		c.Player().SetDistancedAction(func() bool {
			if c.Player().NextTo(npc.Location) && c.Player().WithinRange(npc.Location, 2) {
				c.Player().ResetPath()
				npc.ResetPath()
				c.Player().SetLocation(npc.Location, true)
				c.Player().State = world.MSFighting
				npc.State = world.MSFighting
				c.Player().SetDirection(world.LeftFighting)
				npc.SetDirection(world.RightFighting)
				c.Player().Transients().SetVar("fighting", true)
				c.Player().Transients().SetVar("fightTarget", npc)
				npc.Transients().SetVar("fighting", true)
				npc.Transients().SetVar("fightTarget", c.Player())
				go func() {
					ticker := time.NewTicker(time.Millisecond * 1200)
					defer ticker.Stop()
					curRound := 0
					for range ticker.C {
						if !c.Player().Transients().VarBool("fighting", false) || !c.Player().Transients().VarBool("connected", false) {
							if npc.Transients().VarBool("fighting", false) {
								script.EngineChannel <- func() {
									npc.ResetFighting()
								}
							}
							return
						}
						var attacker, defender world.MobileEntity
						var nextHit int
						if curRound%2 == 0 {
							attacker = c.Player()
							defender = npc
						} else {
							attacker = npc
							defender = c.Player()
						}
						nextHit = attacker.MeleeDamage(defender)
						if curHits := defender.Stats().Current(world.StatHits); nextHit > curHits {
							nextHit = curHits
						}
						defender.Stats().DecreaseCur(world.StatHits, nextHit)
						if defender.Stats().Current(world.StatHits) <= 0 {
							if defenderNpc, ok := defender.(*world.NPC); ok {
								script.EngineChannel <- func() {
									if attackerPlayer, ok := attacker.(*world.Player); ok {
										attackerPlayer.SendPacket(packetbuilders.Sound("victory"))
										world.AddItem(world.NewGroundItemFor(attackerPlayer.UserBase37, 20, 1, defender.X(), defender.Y()))
									} else {
										world.AddItem(world.NewGroundItem(20, 1, defender.X(), defender.Y()))
									}
									attacker.ResetFighting()
									defenderNpc.Stats().SetCur(world.StatHits, defenderNpc.Stats().Maximum(world.StatHits))
									defenderNpc.SetLocation(world.DeathSpot, true)
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
									defenderPlayer.Stats().SetCur(world.StatHits, defenderPlayer.Stats().Maximum(world.StatHits))
									defenderPlayer.SendPacket(packetbuilders.PlayerStats(defenderPlayer))
									defenderPlayer.SendPacket(packetbuilders.Sound("death"))
									// TODO: Keep 3 most valuable items
									defenderPlayer.Inventory().Range(func(item *world.Item) bool {
										world.AddItem(world.NewGroundItem(item.ID, item.Amount, defender.X(), defender.Y()))
										return true
									})
									defenderPlayer.Inventory().Clear()
									defenderPlayer.Transients().SetVar("deathTime", time.Now())
									defenderPlayer.SendPacket(packetbuilders.Death)
									defenderPlayer.SetLocation(world.SpawnPoint, true)
									defenderPlayer.SendPacket(packetbuilders.PlaneInfo(defenderPlayer))
								}
							}
							return
						}

						if defenderNpc, ok := defender.(*world.NPC); ok {
							hitUpdate := packetbuilders.NpcDamage(defenderNpc, nextHit)
							c.SendPacket(hitUpdate)
							for _, p1 := range c.Player().NearbyPlayers() {
								p1.SendPacket(hitUpdate)
							}
						} else if defenderPlayer, ok := defender.(*world.Player); ok {
							hitUpdate := packetbuilders.PlayerDamage(defenderPlayer, nextHit)
							c.SendPacket(hitUpdate)
							for _, p1 := range c.Player().NearbyPlayers() {
								p1.SendPacket(hitUpdate)
							}
						}

						attacker.Transients().SetVar("fightRound", attacker.Transients().VarInt("fightRound", 0)+1)
						curRound++
					}
				}()
				return true
			} else {
				c.Player().SetPath(world.MakePath(c.Player().Location, npc.Location))
			}
			return false
		})
	}
	PacketHandlers["attackplayer"] = func(c clients.Client, p *packet.Packet) {
		affectedClient, ok := clients.FromIndex(p.ReadShort())
		if affectedClient == nil || !ok {
			log.Suspicious.Printf("player[%v] tried to attack nil player\n", c)
			return
		}
		if c.Player().State != world.MSIdle {
			return
		}
		if affectedClient.Player().State != world.MSIdle {
			log.Info.Printf("Target player busy during attack request  State: %d\n", affectedClient.Player().State)
			return
		}
		affectedPlayer := affectedClient.Player()
		c.Player().SetDistancedAction(func() bool {
			if c.Player().NextTo(affectedPlayer.Location) && c.Player().WithinRange(affectedPlayer.Location, 2) {
				affectedPlayer.SendPacket(packetbuilders.Sound("underattack"))
				c.Player().ResetPath()
				affectedPlayer.ResetPath()
				c.Player().SetLocation(affectedPlayer.Location, true)
				c.Player().State = world.MSFighting
				affectedPlayer.State = world.MSFighting
				c.Player().SetDirection(world.LeftFighting)
				affectedPlayer.SetDirection(world.RightFighting)
				c.Player().Transients().SetVar("fighting", true)
				c.Player().Transients().SetVar("fightTarget", affectedPlayer)
				affectedPlayer.Transients().SetVar("fighting", true)
				affectedPlayer.Transients().SetVar("fightTarget", c.Player())
				go func() {
					ticker := time.NewTicker(time.Millisecond * 1200)
					defer ticker.Stop()
					curRound := 0
					for range ticker.C {
						if !affectedPlayer.Transients().VarBool("fighting", false) || !c.Player().Transients().VarBool("fighting", false) || !c.Player().Transients().VarBool("connected", false) || !affectedPlayer.Transients().VarBool("connected", false) {
							if affectedPlayer.Transients().VarBool("fighting", false) {
								script.EngineChannel <- func() {
									affectedPlayer.ResetFighting()
								}
							}
							if c.Player().Transients().VarBool("fighting", false) {
								script.EngineChannel <- func() {
									c.Player().ResetFighting()
								}
							}
							return
						}
						var attacker, defender *world.Player
						if curRound%2 == 0 {
							attacker = c.Player()
							defender = affectedPlayer
						} else {
							attacker = affectedClient.Player()
							defender = c.Player()
						}
						nextHit := attacker.MeleeDamage(defender)
						if nextHit > defender.Stats().Current(world.StatHits) {
							nextHit = defender.Stats().Current(world.StatHits)
						}
						defender.Stats().DecreaseCur(world.StatHits, nextHit)
						if defender.Stats().Current(world.StatHits) <= 0 {
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
								defender.Stats().SetCur(world.StatHits, defender.Stats().Maximum(world.StatHits))
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
						for _, p1 := range c.Player().NearbyPlayers() {
							p1.SendPacket(hitUpdate)
						}

						attacker.Transients().SetVar("fightRound", attacker.Transients().VarInt("fightRound", 0)+1)
						curRound++
					}
				}()
				return true
			} else {
				c.Player().SetPath(world.MakePath(c.Player().Location, affectedPlayer.Location))
			}
			return false
		})
	}
	PacketHandlers["fightmode"] = func(c clients.Client, p *packet.Packet) {
		mode := p.ReadByte()
		if mode < 0 || mode > 3 {
			log.Suspicious.Printf("Invalid fightmode selected (%v) by %v", mode, c.Player().String())
			return
		}
		c.Player().SetFightMode(int(mode))
	}
}
