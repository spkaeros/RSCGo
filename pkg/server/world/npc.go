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

	"github.com/spkaeros/rscgo/pkg/rand"
	"go.uber.org/atomic"
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
var Npcs []*NPC
var npcsLock sync.RWMutex

//NPC Represents a single non-playable character within the game world.
type NPC struct {
	*Mob
	ID         int
	Boundaries [2]Location
	StartPoint Location
}

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id int, startX int, startY int, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{ID: id, Mob: &Mob{Entity: &Entity{Index: int(NpcCounter.Swap(NpcCounter.Load() + 1)), Location: NewLocation(startX, startY)}, TransAttrs: &AttributeList{set: make(map[string]interface{})}}}
	n.Transients().SetVar("skills", &SkillTable{})
	n.Boundaries[0] = NewLocation(minX, minY)
	n.Boundaries[1] = NewLocation(maxX, maxY)
	n.StartPoint = NewLocation(startX, startY)
	if id < 794 {
		n.Skills().current[0] = NpcDefs[id].Attack
		n.Skills().current[1] = NpcDefs[id].Defense
		n.Skills().current[2] = NpcDefs[id].Strength
		n.Skills().current[3] = NpcDefs[id].Hits
		n.Skills().maximum[0] = NpcDefs[id].Attack
		n.Skills().maximum[1] = NpcDefs[id].Defense
		n.Skills().maximum[2] = NpcDefs[id].Strength
		n.Skills().maximum[3] = NpcDefs[id].Hits
	}
	npcsLock.Lock()
	Npcs = append(Npcs, n)
	npcsLock.Unlock()
	return n
}

func (n *NPC) Name() string {
	if n.ID > 793 || n.ID < 0 {
		return "nil"
	}
	return NpcDefs[n.ID].Name
}

func (n *NPC) Command() string {
	if n.ID > 793 || n.ID < 0 {
		return "nil"
	}
	return NpcDefs[n.ID].Command
}

//UpdateNPCPositions Loops through the global NPC entityList and, if they are by a player, updates their path to a new path every so often,
// within their boundaries, and traverses each NPC along said path if necessary.
func UpdateNPCPositions() {
	npcsLock.RLock()
	for _, n := range Npcs {
		if n.Busy() || n.IsFighting() || n.Equals(DeathPoint) {
			continue
		}
		if n.TransAttrs.VarTime("nextMove").Before(time.Now()) {
			for _, r := range surroundingRegions(n.X(), n.Y()) {
				r.Players.lock.RLock()
				if len(r.Players.set) > 0 {
					r.Players.lock.RUnlock()
					n.TransAttrs.SetVar("nextMove", time.Now().Add(time.Second*time.Duration(rand.Int31N(5, 15))))
					//					go n.WalkTo(NewRandomLocation(n.Boundaries))
					n.TransAttrs.SetVar("pathLength", rand.Int31N(5, 15))
					break
				}
				r.Players.lock.RUnlock()
			}
		}
		n.TraversePath()
	}
	npcsLock.RUnlock()
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
	npcsLock.RLock()
	for _, n := range Npcs {
		for _, fn := range n.ResetTickables {
			fn()
		}
		n.ResetTickables = n.ResetTickables[:0]
	}
	npcsLock.RUnlock()
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
		r.Players.lock.RLock()
		for _, p1 := range r.Players.set {
			if p1, ok := p1.(*Player); ok {
				p1.SendPacket(NpcDamage(n, dmg))
			}
		}
		r.Players.lock.RUnlock()
	}
}

func (n *NPC) Killed(killer MobileEntity) {
	if killer, ok := killer.(*Player); ok {
		for _, t := range NpcDeathTriggers {
			if t.Check(killer, n) {
				go t.Action(killer, n)
			}
		}
	}
	AddItem(NewGroundItem(20, 1, n.X(), n.Y()))
	n.Skills().SetCur(StatHits, n.Skills().Maximum(StatHits))
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
	/*	path := n.Path()
		if path == nil {
			return
		}
		if n.AtLocation(path.nextTile()) {
			path.CurrentWaypoint++
		}
		if n.FinishedPath() {
			n.ResetPath()
			return
		}*/
	//dst := path.nextTile()
	if n.TransAttrs.VarInt("pathLength", 0) <= 0 {
		return
	}

	for tries := 0; tries < 10; tries++ {
		dst := NewLocation(n.X(), n.Y())
		if rand.Uint8n(4) == 3 {
			n.TransAttrs.SetVar("pathDir", int(rand.Uint8n(8)))
		}
		switch n.TransAttrs.VarInt("pathDir", North) {
		case North:
			if IsTileBlocking(dst.X(), dst.Y()-1, ClipSouth, false) ||
				dst.X() > n.Boundaries[1].X() || dst.Y()-1 > n.Boundaries[1].Y() ||
				dst.X() < n.Boundaries[0].X() || dst.Y()-1 < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", South)
				continue
			}
			dst.y.Dec()
		case South:
			if IsTileBlocking(dst.X(), dst.Y()+1, ClipNorth, false) ||
				dst.X() > n.Boundaries[1].X() || dst.Y()+1 > n.Boundaries[1].Y() ||
				dst.X() < n.Boundaries[0].X() || dst.Y()+1 < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", North)
				continue
			}
			dst.y.Inc()
		case East:
			if IsTileBlocking(dst.X()-1, dst.Y(), ClipWest, false) ||
				dst.X()-1 > n.Boundaries[1].X() || dst.Y() > n.Boundaries[1].Y() ||
				dst.X()-1 < n.Boundaries[0].X() || dst.Y() < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", West)
				continue
			}
			dst.x.Dec()
		case West:
			if IsTileBlocking(dst.X()+1, dst.Y(), ClipEast, false) ||
				dst.X()+1 > n.Boundaries[1].X() || dst.Y() > n.Boundaries[1].Y() ||
				dst.X()+1 < n.Boundaries[0].X() || dst.Y() < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", East)
				continue
			}
			dst.x.Inc()
		case NorthEast:
			if IsTileBlocking(dst.X()-1, dst.Y()-1, ClipSouth|ClipWest, false) ||
				IsTileBlocking(dst.X(), dst.Y()-1, ClipSouth, false) ||
				IsTileBlocking(dst.X()-1, dst.Y(), ClipWest, false) ||
				dst.X()-1 > n.Boundaries[1].X() || dst.Y()-1 > n.Boundaries[1].Y() ||
				dst.X()-1 < n.Boundaries[0].X() || dst.Y()-1 < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", SouthWest)
				continue
			}
			dst.y.Dec()
			dst.x.Dec()
		case SouthEast:
			if IsTileBlocking(dst.X()-1, dst.Y()+1, ClipNorth|ClipWest, false) ||
				IsTileBlocking(dst.X(), dst.Y()+1, ClipNorth, false) ||
				IsTileBlocking(dst.X()-1, dst.Y(), ClipWest, false) ||
				dst.X()-1 > n.Boundaries[1].X() || dst.Y()+1 > n.Boundaries[1].Y() ||
				dst.X()-1 < n.Boundaries[0].X() || dst.Y()+1 < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", NorthWest)
				continue
			}
			dst.y.Inc()
			dst.x.Dec()
		case NorthWest:
			if IsTileBlocking(dst.X()+1, dst.Y()-1, ClipSouth|ClipEast, false) ||
				IsTileBlocking(dst.X(), dst.Y()-1, ClipSouth, false) ||
				IsTileBlocking(dst.X()+1, dst.Y(), ClipEast, false) ||
				dst.X()+1 > n.Boundaries[1].X() || dst.Y()-1 > n.Boundaries[1].Y() ||
				dst.X()+1 < n.Boundaries[0].X() || dst.Y()-1 < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", SouthEast)
				continue
			}
			dst.y.Dec()
			dst.x.Inc()
		case SouthWest:
			if IsTileBlocking(dst.X()+1, dst.Y()+1, ClipNorth|ClipEast, false) ||
				IsTileBlocking(dst.X(), dst.Y()+1, ClipNorth, false) ||
				IsTileBlocking(dst.X()+1, dst.Y(), ClipEast, false) ||
				dst.X()+1 > n.Boundaries[1].X() || dst.Y()+1 > n.Boundaries[1].Y() ||
				dst.X()+1 < n.Boundaries[0].X() || dst.Y()+1 < n.Boundaries[0].Y() {
				n.TransAttrs.SetVar("pathDir", NorthEast)
				continue
			}
			dst.y.Inc()
			dst.x.Inc()
		}

		x, y := n.X(), n.Y()
		next := NewLocation(x, y)
		xBlocked, yBlocked := false, false
		newXBlocked, newYBlocked := false, false
		if y > dst.Y() {
			yBlocked = IsTileBlocking(x, y, ClipNorth, true)
			newYBlocked = IsTileBlocking(x, y-1, ClipSouth, false)
			if !newYBlocked {
				next.y.Dec()
			}
		} else if y < dst.Y() {
			yBlocked = IsTileBlocking(x, y, ClipSouth, true)
			newYBlocked = IsTileBlocking(x, y+1, ClipNorth, false)
			if !newYBlocked {
				next.y.Inc()
			}
		}
		if x > dst.X() {
			xBlocked = IsTileBlocking(x, next.Y(), ClipEast, true)
			newXBlocked = IsTileBlocking(x-1, next.Y(), ClipWest, false)
			if !newXBlocked {
				next.x.Dec()
			}
		} else if x < dst.X() {
			xBlocked = IsTileBlocking(x, next.Y(), ClipWest, true)
			newXBlocked = IsTileBlocking(x+1, next.Y(), ClipEast, false)
			if !newXBlocked {
				next.x.Inc()
			}
		}

		if (xBlocked && yBlocked) || (xBlocked && y == dst.Y()) || (yBlocked && x == dst.X()) {
			return
		}
		if (newXBlocked && newYBlocked) || (newXBlocked && x != next.X() && y == next.Y()) || (newYBlocked && y != next.Y() && x == next.X()) {
			return
		}

		if next.X() > x {
			newXBlocked = IsTileBlocking(next.X(), next.Y(), ClipEast, false)
		} else if next.X() < x {
			newXBlocked = IsTileBlocking(next.X(), next.Y(), ClipWest, false)
		}

		if next.Y() > y {
			newYBlocked = IsTileBlocking(next.X(), next.Y(), ClipNorth, false)
		} else if next.Y() < y {
			newYBlocked = IsTileBlocking(next.X(), next.Y(), ClipSouth, false)
		}

		if (newXBlocked && newYBlocked) || (newXBlocked && y == next.Y()) || (newYBlocked && x == next.X()) {
			return
		}

		n.TransAttrs.DecVar("pathLength", 1)

		n.SetLocation(next, false)
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
