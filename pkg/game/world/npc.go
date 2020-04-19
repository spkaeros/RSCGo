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
	"math"
	"time"
	
	"go.uber.org/atomic"
	
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/log"
)

//NpcDefinition This represents a single definition for a single NPC in the game.
type NpcDefinition struct {
	ID          int
	Name        string
	Description string
	Command     string
	Hits        int
	Attack      int
	Strength    int
	Defense     int
	Hostility   int
}

//NpcDefs This holds the defining characteristics for all of the game's NPCs, ordered by ID.
var NpcDefs []NpcDefinition

//NpcCounter Counts the number of total NPCs within the world.
var NpcCounter = atomic.NewUint32(0)
//Npcs A collection of every NPC in the game, sorted by index
//var Npcs []*NPC
var Npcs = NewMobList()

//NPC Represents a single non-playable character within the game world.
type NPC struct {
	*Mob
	ID         int
	Boundaries [2]Location
	StartPoint Location
	damageDeltas map[uint64]int
}

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id int, startX int, startY int, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{ID: id, Mob: &Mob{Entity: &Entity{Index: Npcs.Size(), Location: NewLocation(startX, startY)}, AttributeList: entity.NewAttributeList()},
			damageDeltas: make(map[uint64]int), Boundaries: [2]Location{ NewLocation(minX, minY), NewLocation(maxX, maxY) } }
	Npcs.Add(n)
	n.SetVar("skills", &entity.SkillTable{})
	n.SetVar("startPoint", n.Location.Clone())
	for i := 0; i < 18; i++ {
		n.Skills().SetCur(i, 1)
		n.Skills().SetMax(i, 1)
	}
	if id < 794 {
		n.Skills().SetCur(0, NpcDefs[id].Attack)
		n.Skills().SetCur(1, NpcDefs[id].Defense)
		n.Skills().SetCur(2, NpcDefs[id].Strength)
		n.Skills().SetCur(3, NpcDefs[id].Hits)
		
		n.Skills().SetMax(0, NpcDefs[id].Attack)
		n.Skills().SetMax(1, NpcDefs[id].Defense)
		n.Skills().SetMax(2, NpcDefs[id].Strength)
		n.Skills().SetMax(3, NpcDefs[id].Hits)

	}
	return n
}

// Returns true if this NPCs definition has the attackable hostility bit set.
func (n *NPC) Attackable() bool {
	if n.ID > len(NpcDefs)-1 {
		return false
	}

	return NpcDefs[n.ID].Hostility&1 == 1
}

// Returns true if this NPCs definition has the retreat near death hostility bit set.
func (n *NPC) Retreats() bool {
	if n.ID > len(NpcDefs)-1 {
		return false
	}

	return NpcDefs[n.ID].Hostility&2 == 2
}

// Returns true if this NPCs definition has the aggressive hostility bit set.
func (n *NPC) Aggressive() bool {
	if n.ID > len(NpcDefs)-1 {
		return false
	}

	return NpcDefs[n.ID].Hostility&4 == 4
}

func (n *NPC) Name() string {
	if n.ID > len(NpcDefs)-1 {
		return "nil"
	}
	return NpcDefs[n.ID].Name
}

func (n *NPC) Command() string {
	if n.ID > len(NpcDefs)-1 {
		return "nil"
	}
	return NpcDefs[n.ID].Command
}

//UpdateNPCPositions Loops through the global NPC entityList and, if they are by a player, updates their path to a new path every so often,
// within their boundaries, and traverses each NPC along said path if necessary.
func UpdateNPCPositions() {
	Npcs.RangeNpcs(func(n *NPC) bool {
		if n.Busy() || n.IsFighting() || n.Equals(DeathPoint) {
			return false
		}
		moveTime := n.VarTime("moveTime")
		if n.VarInt("pathLength", 0) == 0 && (moveTime.IsZero() || time.Now().After(moveTime)) {
			// schedule when to start wandering again
			n.SetVar("moveTime", time.Now().Add(time.Second*time.Duration(rand.Int31N(10, 15))))
			// set how many steps we should wander for before taking a break
			if n.VarInt("pathLength", 0) == 0 {
				n.SetVar("pathLength", rand.Int31N(5, 15))
			}
		} else {
			// wander aimlessly until we run out of scheduled steps 
			n.TraversePath()
		}
		return false
	})
}

func (n *NPC) UpdateRegion(x, y int) {
	curArea := getRegion(n.X(), n.Y())
	newArea := getRegion(x, y)
	if newArea != curArea {
		if curArea.NPCs.Contains(n) {
			curArea.NPCs.Remove(n)
		}
		newArea.NPCs.Add(n)
	}
}

//ResetNpcUpdateFlags Resets the synchronization update flags for all NPCs in the game world.
func ResetNpcUpdateFlags() {
	Npcs.RangeNpcs(func (n *NPC) bool {
		n.ResetRegionRemoved()
		n.ResetRegionMoved()
		n.ResetSpriteUpdated()
		n.ResetAppearanceChanged()
		return false
	})
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
	for _, r := range surroundingRegions(n.X(), n.Y()) {
		r.Players.RangePlayers(func(p1 *Player) bool {
			if !n.WithinRange(p1.Location, 16) {
				return false
			}
			p1.SendPacket(NpcDamage(n, dmg))
			return false
		})
	}
}

func (n *NPC) DamageMelee(atk *Player, dmg int) {
	n.Damage(dmg)
	if delta, ok := n.damageDeltas[atk.UsernameHash()]; ok {
		n.damageDeltas[atk.UsernameHash()] = delta+dmg
		return
	}
	n.damageDeltas[atk.UsernameHash()] = dmg
}

//MeleeExperience returns how much combat experience to award for killing an opponent with melee.
func (n *NPC) MeleeExperience(up bool) float64 {
	e := float64((n.Skills().CombatLevel()*2.0)+10.0) * 1.5
	if up {
		return math.Ceil(e)
	}
	return math.Floor(e)
}

// Loops through the list of players that had dealt damage to this NPC, and
// adds up all the damage that was dealt by an active player, to help fairly distribute
// all of the NPCs EXP reward to all of the players that helped kill it, proprortional
// to how much they helped.
// Returns: the total number of hitpoints depleted by player melee damage.
func(n *NPC) TotalDamage() (total int) {
	for userHash, dmg:= range n.damageDeltas {
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
	totalExp := int(n.MeleeExperience(true)) & 0xFFFFFFFFC
	for usernameHash, damage := range n.damageDeltas {
		player, ok := Players.FromUserHash(usernameHash)
		log.Info.Println(usernameHash,damage)
		if ok {
			exp := float64(totalExp)/float64(totalDamage)
			player.DistributeMeleeExp(int(exp)*damage / 4)
		}
		if damage > mostDamage || dropPlayer == nil {
			if ok {
				dropPlayer = player
			}
			mostDamage = damage
		}
	}
	n.damageDeltas = make(map[uint64]int)
	if dropPlayer != nil {
		AddItem(NewGroundItemFor(dropPlayer.UsernameHash(), DefaultDrop, 1, n.X(), n.Y()))
	} else {
		AddItem(NewGroundItem(DefaultDrop, 1, n.X(), n.Y()))
	}
	n.Skills().SetCur(entity.StatHits, n.Skills().Maximum(entity.StatHits))
	n.SetLocation(DeathPoint, true)
	killer.ResetFighting()
	n.ResetFighting()
	go func() {
		// TODO: npc definition entry for respawn time
		time.Sleep(time.Second * 10)
		n.SetLocation(n.VarChecked("startPoint").(Location), true)
	}()
	return
}

//TraversePath If the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (n *NPC) TraversePath() {
	if n.VarInt("pathLength", 0) <= 0 {
		return
	}
	
	for tries := 0; tries < 10; tries++ {
		if Chance(25) {
			n.SetVar("pathDir", int(rand.Uint8n(8)))
		}
		
		dst := n.Location.Clone()
		dir := n.VarInt("pathDir", North);
		if dir == West || dir == SouthWest || dir == NorthWest {
			dst.x.Inc()
		} else if dir == East || dir == SouthEast || dir == NorthEast {
			dst.x.Dec()
		}
		if dir == North || dir == NorthWest || dir == NorthEast {
			dst.y.Dec()
		} else if dir == South || dir == SouthWest || dir == SouthEast {
			dst.y.Inc()
		}
		
		if !n.Reachable(dst) || !dst.WithinArea(n.Boundaries) {
			n.SetVar("pathDir", int(rand.Uint8n(8)))
			continue
		}

		n.Dec("pathLength", 1)
		n.SetLocation(dst, false)
		break
	}
}

//ChatIndirect sends a chat message to target and all of target's view area players, without any delay.
func (n *NPC) ChatIndirect(target *Player, msg string) {
	for _, player := range target.NearbyPlayers() {
		player.SendPacket(NpcMessage(n, msg, target))
	}
	target.SendPacket(NpcMessage(n, msg, target))
}

//Chat sends chat messages to target and all of target's view area players, with a 1800ms(3 tick) delay between each
// message.
func (n *NPC) Chat(target *Player, msgs ...string) {
	for _, msg := range msgs {
		n.ChatIndirect(target, msg)

		//		if i < len(msgs)-1 {
		time.Sleep(time.Millisecond * 1800)
		// TODO: is 3 ticks right?
		//		}
	}
}
