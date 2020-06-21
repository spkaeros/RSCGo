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
	"github.com/spkaeros/rscgo/pkg/tasks"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/rand"
)

//Npcs A collection of every NPC in the game, sorted by index
var Npcs = NewMobList()

//NPC Represents a single non-playable character within the game world.
type NPC struct {
	Mob
	ID               int
	pathSteps        int
	lastMoved        time.Time
	StartPoint       Location
	Boundaries       [2]Location
	meleeRangeDamage damages
}

type (
	damageTable = map[uint64]int
	damages     = struct {
		damageTable
		sync.RWMutex
	}
)

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id int, startX int, startY int, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{
		ID: id,
		Mob: Mob{
			skills: entity.SkillTable{},
			Entity: Entity{
				Index:    Npcs.Size(),
				Location: NewLocation(startX, startY),
			},
			AttributeList: entity.NewAttributeList(),
		},
		meleeRangeDamage: damages{
			damageTable: make(damageTable),
		},
		Boundaries: [2]Location{NewLocation(minX, minY), NewLocation(maxX, maxY)},
	}
	n.StartPoint = n.Location.Clone()
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
func UpdateNPCPositions() {
	wait := sync.WaitGroup{}
	Npcs.RangeNpcs(func(n *NPC) bool {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if n.Busy() || n.IsFighting() || n.Equals(DeathPoint) {
				return
			}
			if n.pathSteps == 0 && (n.lastMoved.IsZero() || time.Now().After(n.lastMoved)) {
				// schedule when to start wandering again
				n.lastMoved = time.Now().Add(time.Second * time.Duration(rand.Rng.Intn(15)+5))
				// set how many steps we should wander for before taking a break
				n.pathSteps = rand.Rng.Intn(15)
			}
			// wander aimlessly until we run out of scheduled steps
			n.TraversePath()
		}()
		return false
	})
	wait.Wait()
}

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

//NpcActionPredicate callback to a function defined in the Anko scripts loaded at runtime, to be run when certain
// events occur.  If it returns true, it will block the event that triggered it from occurring
type NpcBlockingTrigger struct {
	// Check returns true if this handler should run.
	Check func(*Player, *NPC) bool
	// Action is the function that will run if Check returned true.
	Action func(*Player, *NPC)
}

//NpcDeathTriggers List of script callbacks to run when you kill an NPC
var NpcDeathTriggers []NpcBlockingTrigger

func (n *NPC) Damage(dmg int) {
	for _, r := range Region(n.X(), n.Y()).neighbors() {
		r.Players.RangePlayers(func(p1 *Player) bool {
			if !n.WithinRange(p1.Location, 16) {
				return false
			}
			p1.QueueNpcSplat(n, dmg)
			return false
		})
	}
}

func (n *NPC) DamageMelee(atk *Player, dmg int) {
	n.Damage(dmg)
	n.meleeRangeDamage.Lock()
	defer n.meleeRangeDamage.Unlock()
	n.meleeRangeDamage.damageTable[atk.UsernameHash()] += dmg
}

// Loops through the list of players that had dealt damage to this NPC, and
// adds up all the damage that was dealt by an active player, to help fairly distribute
// all of the NPCs EXP reward to all of the players that helped kill it, proprortional
// to how much they helped.
// Returns: the total number of hitpoints depleted by player melee damage.
func (n *NPC) TotalDamage() (total int) {
	n.meleeRangeDamage.RLock()
	defer n.meleeRangeDamage.RUnlock()
	for userHash, dmg := range n.meleeRangeDamage.damageTable {
		if Players.ContainsHash(userHash) {
			total += dmg
		}
	}

	return
}

func (n *NPC) Killed(killer entity.MobileEntity) {
	if killer, ok := killer.(*Player); ok {
		for _, t := range NpcDeathTriggers {
			if t.Check(killer, n) {
				go t.Action(killer, n)
			}
		}
	}
	// first pass is to find the total so we can split up the exp properly
	// this is because the total is not guaranteed to match max hitpoints since
	// the NPC can heal after damage has been dealt, among other things
	var dropPlayer *Player
	var mostDamage int
	totalDamage := n.TotalDamage()
	totalExp := n.ExperienceReward() ^ 3
	n.meleeRangeDamage.RLock()
	for usernameHash, damage := range n.meleeRangeDamage.damageTable {
		player, ok := Players.FindHash(usernameHash)
		if ok {
			exp := float64(totalExp) / float64(totalDamage)
			player.DistributeMeleeExp(int(exp) * damage)
			if damage > mostDamage || dropPlayer == nil {
				dropPlayer = player
				mostDamage = damage
			}
		}
	}
	n.meleeRangeDamage.RUnlock()

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
	n.SetCoords(0, 0, true)
}

func (n *NPC) Respawn() {
	n.meleeRangeDamage.Lock()
	n.meleeRangeDamage.damageTable = make(damageTable)
	n.meleeRangeDamage.Unlock()
	for i := 0; i <= 3; i++ {
		n.Skills().SetCur(i, n.Skills().Maximum(i))
	}
	n.SetLocation(n.StartPoint, true)
	n.UnsetVar("removed")
}

//TraversePath If the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (n *NPC) TraversePath() {
	dst := n.Location.Clone()
	dir := n.Direction()
	if Chance(25) {
		dir = rand.Rng.Intn(8)
	}
	if dir == East || dir == SouthEast || dir == NorthEast {
		dst.x.Dec()
	} else if dir == West || dir == SouthWest || dir == NorthWest {
		dst.x.Inc()
	}
	if dir == North || dir == NorthWest || dir == NorthEast {
		dst.y.Dec()
	} else if dir == South || dir == SouthWest || dir == SouthEast {
		dst.y.Inc()
	}

	if !n.Reachable(dst) || !dst.WithinArea(n.Boundaries) {
		return
	}

	n.pathSteps -= 1
	n.SetLocation(dst, false)
	//	}
}

//ChatIndirect sends a chat message to target and all of target's view area players, without any delay.
func (n *NPC) ChatIndirect(target *Player, msg string) {
	for _, player := range target.NearbyPlayers() {
		player.QueueNpcChat(n, target, msg)
	}
	target.QueueNpcChat(n, target, msg)
}

//Chat sends chat messages to target and all of target's view area players, with a 1800ms(3 tick) delay between each
// message.
func (n *NPC) Chat(target *Player, msgs ...string) {
	if len(msgs) <= 0 {
		return
	}

	for _, msg := range msgs {
		sleep := 3
		if len(msg) >= 83 {
			sleep = 4
		}
		n.ChatIndirect(target, msg)
		time.Sleep(time.Millisecond*640*time.Duration(sleep))
	}
}
