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
	// "math/rand"
	"time"
	
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/tasks"
	"github.com/spkaeros/rscgo/pkg/rand"
)

//Npcs A collection of every NPC in the game, sorted by index
var Npcs = NewMobList()

//NPC Represents a single non-playable character within the game world.
type NPC struct {
	Mob
	ID       					  int
	StartPoint                    entity.Location
	Boundaries                    [2]entity.Location
	Steps, Ticks				  int
	meleeRangeDamage, magicDamage damages
}

type (
	damageTable = map[uint64]int
	damages     struct {
		damageTable map[uint64]int
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
		meleeRangeDamage: damages {
			damageTable: make(map[uint64]int),
		},
		magicDamage: damages {
			damageTable: make(map[uint64]int),
		},
		Boundaries: [2]entity.Location{NewLocation(minX, minY), NewLocation(maxX, maxY)},
	}
	defer Npcs.Add(n)
	n.StartPoint = n.Clone()
	if n.valid() {
		skills := [18] int {
			definitions.Npcs[id].Attack,
			definitions.Npcs[id].Defense,
			definitions.Npcs[id].Strength,
			definitions.Npcs[id].Hits,
		}
		for i := 0; i < 18; i += 1 {
			n.Skills().SetCur(i, skills[i])
			n.Skills().SetMax(i, skills[i])
		}
	}
	return n
}

func (d damages) Put(username uint64, dmg int) {
	d.Lock()
	defer d.Unlock()
	d.damageTable[username] += dmg
}

func (d damages) Reset() {
	d.Lock()
	defer d.Unlock()
	d.damageTable = make(map[uint64]int)
}

func (n *NPC) CacheDamage(hash uint64, dmg int) {
	n.meleeRangeDamage.Put(hash, dmg)
}

func (n *NPC) valid() bool {
	return n.ID < len(definitions.Npcs)
}

// Returns true if this NPCs definition has the attackable hostility bit set.
func (n *NPC) Attackable() bool {
	if !n.valid() {
		return false
	}

	return definitions.Npcs[n.ID].Hostility&1 == 1
}

// Returns true if this NPCs definition has the retreat near death hostility bit set.
func (n *NPC) Retreats() bool {
	if !n.valid() {
	// if n.ID > len(d efinitions.Npcs)-1 {
		return false
	}

	return definitions.Npcs[n.ID].Hostility&2 == 2
}

// Returns true if this NPCs definition has the aggressive hostility bit set.
func (n *NPC) Aggressive() bool {
	if !n.valid() {
		return false
	}

	return definitions.Npcs[n.ID].Hostility&4 == 4
}

func (n *NPC) Name() string {
	if !n.valid() {
		return "nil"
	}
	return definitions.Npcs[n.ID].Name
}

func (n *NPC) Command() string {
	if !n.valid() {
		return "nil"
	}
	return definitions.Npcs[n.ID].Command
}

type NpcBlockingTrigger struct {
	Check func(*Player, *NPC) bool
	Action func(*Player, *NPC)
}

//NpcDeathTriggers List of script callbacks to run when you kill an NPC
var NpcDeathTriggers []NpcBlockingTrigger

func (n *NPC) rewardKillers() (winner *Player) {
	n.meleeRangeDamage.RLock()
	defer n.meleeRangeDamage.RUnlock()
	totalExp := float64(n.ExperienceReward())
	amount := 0
	total := 0.0
	for _, damage := range n.meleeRangeDamage.damageTable {
		total += float64(damage)
	}
	for username, damage := range n.meleeRangeDamage.damageTable {
		if player, ok := Players.FindHash(username); ok {
			if damage > amount || winner == nil {
				winner = player
				amount = damage
			}
			player.DistributeMeleeExp(totalExp / float64(total) * float64(damage))
		}
	}
	return winner
}

func (n *NPC) Killed(killer entity.MobileEntity) {
	if killer, ok := killer.(*Player); ok {
		for _, t := range NpcDeathTriggers {
			if t.Check(killer, n) {
				t.Action(killer, n)
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
	for i := 0; i < 18; i++ {
		n.Skills().SetCur(i, n.Skills().Maximum(i))
	}
	n.UnsetVar("removed")
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
	n.Steps -= 1
	if p := n.VarPlayer("targetPlayer"); p != nil {
		if !n.Near(p, 6) {
			n.UnsetVar("targetPlayer")
		}
		if n.Aggressive() && n.Near(p, 1) && !n.Collides(p) && !p.Busy() && !p.IsFighting() {
			if t := p.SessionCache().VarTime("lastFight"); time.Since(t) < 1920*time.Millisecond {
				return
			}
			StartCombat(n, p)
			return
		}
		path := n.PivotTo(p)
		if len(path[0])|len(path[1]) == 0 {
			return
		}
		dst := NewLocation(path[0][0], path[1][0])
		if n.Collides(dst) || !dst.WithinArea(n.Boundaries) {
			return
		}
		n.SetLocation(dst, false)
		return
	}

	dir := n.Direction()
	dst := n.Step(dir)
	if Chance(15) {
		dir = rand.Intn(8)
		dst = n.Step(dir)
		for i := 0; i < 10 && n.Collides(dst); i += 1 {
			dir = rand.Intn(8)
			dst = n.Step(dir)
		}
	}
	
	if n.Collides(dst) || !dst.WithinArea(n.Boundaries) {
		return
	}

	n.SetLocation(dst, false)
}

//ChatIndirect sends a chat message to target and all of target's view area players, without any delay.
func (n *NPC) ChatIndirect(target *Player, msg string) {
	n.enqueueArea(npcEvents, NewTargetedMessage(n,target,msg))
}

func (n *NPC) Enqueue(handle string, e interface{}) {
	n.enqueueArea(handle, e)
}

func (n *NPC) enqueueArea(handle string, e interface{}) {
	updated := NewMobList()
	for _, region := range Region(n.X(), n.Y()).neighbors() {
		region.Players.RangePlayers(func(p *Player) bool {
			if !updated.Contains(p) && p.Near(n, p.ViewRadius()) {
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
	for _, v := range msgs {
		n.enqueueArea(npcEvents, NewTargetedMessage(n, target, v))
		wait := 3
		if len(v) >= 84 {
			wait += 1
		}
		tasks.Stall(wait)
	}
}

