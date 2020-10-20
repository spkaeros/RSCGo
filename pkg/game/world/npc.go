/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"sync"
	"time"

	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/tasks"
)

//Npcs A collection of every NPC in the game, sorted by index
var Npcs = NewMobList()

//NPC Represents a single non-playable character within the game world.
type NPC struct {
	Mob
	ID, MoveTick, PathSteps       int
	StartPoint                    entity.Location
	Boundaries                    [2]entity.Location
	meleeRangeDamage, magicDamage damages
}

type (
	damageTable = map[uint64]int
	damages     struct {
		total int
		damageTable
		sync.RWMutex
	}
)

func (e *Entity) Type() entity.Type {
	return entity.Type(-1)
}

func (e *Entity) SetServerIndex(i int) {
	e.Index = i
}

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id, startX, startY, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{
		ID: id,
		Mob: Mob{
			skills: entity.SkillTable{},
			Entity: &Entity{
				Index:    Npcs.Size(),
				Location: NewLocation(startX, startY),
			},
			AttributeList: entity.NewAttributeList(),
		},
		meleeRangeDamage: damages{
			damageTable: make(damageTable),
		},
		magicDamage: damages{
			damageTable: make(damageTable),
		},
		Boundaries: [2]entity.Location{NewLocation(minX, minY), NewLocation(maxX, maxY)},
	}
	n.StartPoint = n.Clone()
	if id < len(definitions.Npcs)-1 {
		n.Skills().SetCur(0, definitions.Npcs[id].Attack)
		n.Skills().SetCur(1, definitions.Npcs[id].Defense)
		n.Skills().SetCur(2, definitions.Npcs[id].Strength)
		n.Skills().SetCur(3, definitions.Npcs[id].Hits)

		n.Skills().SetMax(0, definitions.Npcs[id].Attack)
		n.Skills().SetMax(1, definitions.Npcs[id].Defense)
		n.Skills().SetMax(2, definitions.Npcs[id].Strength)
		n.Skills().SetMax(3, definitions.Npcs[id].Hits)
	}
	for i := 4; i < 18; i++ {
		n.Skills().SetCur(i, 0)
		n.Skills().SetMax(i, 0)
	}

	Npcs.Add(n)
	return n
}

func (d damages) Put(username uint64, dmg int) {
	d.Lock()
	defer d.Unlock()
	d.total += dmg
	d.damageTable[username] += dmg
}

func (n *NPC) CacheDamage(hash uint64, dmg int) {
	n.meleeRangeDamage.Lock()
	defer n.meleeRangeDamage.Unlock()
	n.meleeRangeDamage.damageTable[hash] += dmg
}

// Returns true if this NPCs definition has the attackable hostility bit set.
func (n *NPC) Attackable() bool {
	if n.ID > len(definitions.Npcs)-1 {
		return false
	}

	return definitions.Npcs[n.ID].Hostility&1 == 1
}

// Returns true if this NPCs definition has the retreat near death hostility bit set.
func (n *NPC) Retreats() bool {
	if n.ID > len(definitions.Npcs)-1 {
		return false
	}

	return definitions.Npcs[n.ID].Hostility&2 == 2
}

// Returns true if this NPCs definition has the aggressive hostility bit set.
func (n *NPC) Aggressive() bool {
	if n.ID > len(definitions.Npcs)-1 {
		return false
	}

	return definitions.Npcs[n.ID].Hostility&4 == 4
}

func (n *NPC) Name() string {
	if n.ID > len(definitions.Npcs)-1 {
		return "nil"
	}
	return definitions.Npcs[n.ID].Name
}

func (n *NPC) Command() string {
	if n.ID > len(definitions.Npcs)-1 {
		return "nil"
	}
	return definitions.Npcs[n.ID].Command
}

//UpdateNPCPositions Loops through the global NPC entityList and, if they are by a player, updates their path to a new path every so often,
// within their boundaries, and traverses each NPC along said path if necessary.
// func UpdateNPCPositions() {
// wait := sync.WaitGroup{}
// Npcs.RangeNpcs(func(n *NPC) bool {
// wait.Add(1)
// go func() {
// defer wait.Done()
// if n.Busy() || n.IsFighting() || n.Equals(DeathPoint) {
// return
// }
// if n.MoveTick > 0 {
// n.MoveTick--
// }
// if Chance(25) && n.pathSteps == 0 && n.moveTick <= 0 {
// // move some amount between 2-15 tiles, moving 1 tile per tick
// n.PathSteps = int(rand.Rng.Float64() * 15 - 2) + 2
// // wait some amount between 25-50 ticks before doing this again
// n.MoveTick = int(rand.Rng.Float64() * 50 - 25) + 25
// }
// // wander aimlessly until we run out of scheduled steps
// n.TraversePath()
// }()
// return false
// })
// wait.Wait()
// }

func (n *NPC) UpdateRegion(x, y int) {
	curArea := Region(n.X(), n.Y())
	newArea := Region(x, y)
	if newArea != curArea {
		if curArea.NPCs.Contains(n) {
			curArea.NPCs.Remove(n)
		}
		newArea.NPCs.Add(n)
	}
}

//ResetNpcUpdateFlags Resets the synchronization update flags for all NPCs in the game world.
func ResetNpcUpdateFlags() {
	wait := sync.WaitGroup{}
	Npcs.RangeNpcs(func(n *NPC) bool {
		wait.Add(1)
		go func() {
			defer wait.Done()
			n.ResetRegionRemoved()
			n.ResetRegionMoved()
			n.ResetSpriteUpdated()
			n.ResetAppearanceChanged()
		}()
		return false
	})
	wait.Wait()
}

type NpcBlockingTrigger struct {
	Check  func(*Player, *NPC) bool
	Action func(*Player, *NPC)
}

//NpcDeathTriggers List of script callbacks to run when you kill an NPC
var NpcDeathTriggers []NpcBlockingTrigger

// func (n *NPC) Damage(dmg int) {
// n.enqueueArea(npcSplatHandle, HitSplat{n, dmg})
// n.Skills().DecreaseCur(entity.StatHits, dmg)
// }
//
// func (n *NPC) DamageFrom(m entity.MobileEntity, dmg, dmgKind int) {
// if dmgKind == 0 {
// n.meleeRangeDamage.Put(AsPlayer(m).UsernameHash(), dmg)
// } else if dmgKind == 1 {
// n.magicDamage.Put(AsPlayer(m).UsernameHash(), dmg)
// }
// n.Damage(dmg)
// 	}

func (n *NPC) rewardKillers() (winner *Player) {
	n.meleeRangeDamage.RLock()
	defer n.meleeRangeDamage.RUnlock()
	totalDamage, totalExp, amount := n.meleeRangeDamage.total, n.ExperienceReward() & ^3, 0
	for username, damage := range n.meleeRangeDamage.damageTable {
		if player, ok := Players.FindHash(username); ok {
			player.DistributeMeleeExp(int(float64(totalExp) / float64(totalDamage) * float64(damage)))
			if damage > amount || winner == nil {
				winner = player
				amount = damage
			}
		}
	}
	return winner
}

func (n *NPC) Killed(killer entity.MobileEntity) {
	if killer, ok := killer.(*Player); ok {
		for _, t := range NpcDeathTriggers {
			if t.Check(killer, n) {
				go t.Action(killer, n)
			}
		}
	}
	dropPlayer := n.rewardKillers()
	// first pass is to find the total so we can split up the exp properly
	// this is because the total is not guaranteed to match max hitpoints since
	// the NPC can heal after damage has been dealt, among other things
	if dropPlayer != nil {
		AddItem(NewGroundItemFor(dropPlayer.UsernameHash(), DefaultDrop, 1, n.X(), n.Y()))
	} else {
		AddItem(NewGroundItem(DefaultDrop, 1, n.X(), n.Y()))
	}

	killer.ResetFighting()
	n.ResetFighting()
	n.Remove()
	tasks.Schedule(16, func() bool {
		n.Respawn()
		return true
	})
	return
}

func (n *NPC) Remove() {
	n.SetVar("removed", true)
	n.SetLocation(DeathPoint, true)
}

func (n *NPC) Respawn() {
	for i := 0; i <= 3; i++ {
		n.Skills().SetCur(i, n.Skills().Maximum(i))
	}
	n.SetVar("removed", false)
	n.SetLocation(n.StartPoint.Clone(), true)
	n.meleeRangeDamage.Lock()
	defer n.meleeRangeDamage.Unlock()
	n.meleeRangeDamage.damageTable = make(damageTable)
	n.magicDamage.Lock()
	defer n.magicDamage.Unlock()
	n.magicDamage.damageTable = make(damageTable)
}

//TraversePath If the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (n *NPC) TraversePath() {
	dst := n.Clone()
	dir := n.Direction()
	if Chance(25) {
		dir = rand.Rng.Intn(8)
	}
	if dir == East || dir == SouthEast || dir == NorthEast {
		dst.SetX(dst.X() - 1)
	} else if dir == West || dir == SouthWest || dir == NorthWest {
		dst.SetX(dst.X() + 1)
	}
	if dir == North || dir == NorthWest || dir == NorthEast {
		dst.SetY(dst.Y() - 1)
	} else if dir == South || dir == SouthWest || dir == SouthEast {
		dst.SetY(dst.Y() + 1)
	}

	if !n.Reachable(dst) || !dst.WithinArea(n.Boundaries) {
		return
	}

	n.PathSteps--
	n.SetLocation(dst, false)
}

//ChatIndirect sends a chat message to target and all of target's view area players, without any delay.
func (n *NPC) ChatIndirect(target *Player, msg string) {
	// for _, player := range target.NearbyPlayers() {
	// player.QueueNpcChat(n, target, msg)
	// }
	// target.QueueNpcChat(n, target, msg)
	n.enqueueArea(npcChatHandle, NewTargetedMessage(n, target, msg))
}

func (n *NPC) enqueueArea(handle string, e interface{}) {
	updated := NewMobList()
	for _, region := range Region(n.X(), n.Y()).neighbors() {
		region.Players.RangePlayers(func(p *Player) bool {
			if !updated.Contains(p) && p.Near(n, 15) {
				p.enqueue(handle, e)
				updated.Add(p)
			}
			return false
		})
	}
}

//Chat sends chat messages to target and all of target's view area players, with a 1800ms(3 tick) delay between each
// message.
func (n *NPC) Chat(target *Player, msgs ...string) {
	if len(msgs) <= 0 {
		return
	}

	n.enqueueArea(npcChatHandle, NewTargetedMessage(n, target, msgs[0]))
	wait := time.Duration(0)
	for _, v := range msgs {
		// tasks.TickList.Schedule(wait, func() bool {
		n.enqueueArea(npcChatHandle, NewTargetedMessage(n, target, v))
		// return true
		// })
		wait += 3
		if len(msgs[0]) >= 84 {
			wait++
		}
		time.Sleep(TickMillis * wait)
	}
	// if len(msgs) > 1 {
	// wait := 3
	// if len(msgs[0]) >= 84 {
	// wait++
	// }
	// // time.Sleep(TickMillis*wait)
	// tasks.TickList.Schedule(wait, func() bool {
	// n.Chat(target, (msgs[1:])...)
	// return true
	// })
	// }
}
