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
	Attackable  bool
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
}

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id int, startX int, startY int, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{ID: id, Mob: &Mob{Entity: &Entity{Index: Npcs.Size(), Location: NewLocation(startX, startY)}, AttributeList: entity.NewAttributeList()}}
	n.Transients().SetVar("skills", &entity.SkillTable{})
	n.Boundaries[0] = NewLocation(minX, minY)
	n.Boundaries[1] = NewLocation(maxX, maxY)
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
	n.StartPoint = n.Location.Clone()
	Npcs.Add(n)
	return n
}

func (n *NPC) Attackable() bool {
	if n.ID <= len(NpcDefs)-1 {
		return NpcDefs[n.ID].Attackable
	}
	return false
}

func (n *NPC) Name() string {
	if n.ID <= len(NpcDefs)-1 {
		return NpcDefs[n.ID].Name
	}
	return "nil"
}

func (n *NPC) Command() string {
	if n.ID <= len(NpcDefs)-1 {
		return NpcDefs[n.ID].Command
	}
	return "nil"
}

//UpdateNPCPositions Loops through the global NPC entityList and, if they are by a player, updates their path to a new path every so often,
// within their boundaries, and traverses each NPC along said path if necessary.
func UpdateNPCPositions() {
	Npcs.RangeNpcs(func(n *NPC) bool {
		if n.Busy() || n.IsFighting() || n.Equals(DeathPoint) {
			return false
		}

		if moveTime := n.VarTime("nextMove"); moveTime.IsZero() || time.Now().After(moveTime) {
			n.SetVar("nextMove", time.Now().Add(time.Second*time.Duration(rand.Int31N(10, 20))))
			n.SetVar("pathLength", rand.Int31N(5, 15))
		}

		if n.VarInt("pathLength", 0) > 0 {
			n.TraversePath()
		}
		return false
	})
	//npcsLock.RLock()
	//for _, n := range Npcs {
	//	if n.Busy() || n.IsFighting() || n.Equals(DeathPoint) {
	//		continue
	//	}
	//	if n.VarTime("nextMove").Before(time.Now()) {
	//		for _, r := range surroundingRegions(n.X(), n.Y()) {
	//			if r.Players.Size() > 0 {
	//				n.SetVar("nextMove", time.Now().Add(time.Second*time.Duration(rand.Int31N(10, 20))))
	//				n.SetVar("pathLength", rand.Int31N(5, 15))
	//				break
	//			}
	//		}
	//	}
	//	n.TraversePath()
	//}
	//npcsLock.RUnlock()
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

//MeleeExperience returns how much combat experience to award for killing an opponent with melee.
func (n *NPC) MeleeExperience(up bool) float64 {
	e := float64((n.Skills().CombatLevel()*2.0)+10.0) * 1.5
	if up {
		return math.Ceil(e)
	}
	return math.Floor(e)
}

func (n *NPC) Killed(killer entity.MobileEntity) {
	if killer, ok := killer.(*Player); ok {
		for _, t := range NpcDeathTriggers {
			if t.Check(killer, n) {
				go t.Action(killer, n)
			}
		}
	}
	AddItem(NewGroundItem(DefaultDrop, 1, n.X(), n.Y()))
	if killer, ok := killer.(*Player); ok {
		killer.DistributeMeleeExp(int(math.Ceil(n.MeleeExperience(true)/ 4.0)))
	}
	n.Skills().SetCur(entity.StatHits, n.Skills().Maximum(entity.StatHits))
	n.SetLocation(DeathPoint, true)
	killer.ResetFighting()
	n.ResetFighting()
	go func() {
		time.Sleep(time.Second * 10)
		n.SetLocation(n.StartPoint, true)
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
		
		if !n.Reachable(dst.X(), dst.Y()) ||
				dst.X() < n.Boundaries[0].X() || dst.X() > n.Boundaries[1].X() ||
				dst.Y() < n.Boundaries[0].Y() || dst.Y() > n.Boundaries[1].Y() {
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
