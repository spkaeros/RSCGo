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
	"github.com/spkaeros/rscgo/pkg/rand"
	"go.uber.org/atomic"
	"sync"
	"time"
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
	ID          int
	Boundaries  [2]Location
	StartPoint  Location
	ChatMessage string
	ChatTarget  int
}

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id int, startX int, startY int, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{ID: id, Mob: &Mob{Entity: &Entity{Index: int(NpcCounter.Swap(NpcCounter.Load() + 1)), Location: NewLocation(startX, startY)}, TransAttrs: &AttributeList{set: make(map[string]interface{})}}, ChatTarget: -1, ChatMessage: ""}
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
					go n.WalkTo(NewRandomLocation(n.Boundaries))
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
		n.ResetChanged()
		n.ResetMoved()
		n.ResetRemoved()
	}
	npcsLock.RUnlock()
}

//TraversePath If the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (n *NPC) TraversePath() {
	path := n.Path()
	if path == nil {
		return
	}
	if n.AtLocation(path.nextTile()) {
		path.CurrentWaypoint++
	}
	if n.FinishedPath() {
		n.ResetPath()
		return
	}
	dst := path.nextTile()
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
		n.ResetPath()
		return
	}
	if (newXBlocked && newYBlocked) || (newXBlocked && x != next.X() && y == next.Y()) || (newYBlocked && y != next.Y() && x == next.X()) {
		n.ResetPath()
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
		n.ResetPath()
		return
	}

	n.SetLocation(next, false)
}

func (n *NPC) Chat(target *Player, msgs ...string) {
	for _, msg := range msgs {
		for _, player := range target.NearbyPlayers() {
			player.SendPacket(NpcMessage(n, msg, target))
		}
		target.SendPacket(NpcMessage(n, msg, target))

		//		if i < len(msgs)-1 {
		time.Sleep(time.Millisecond * 1800)
		// TODO: is 3 ticks right?
		//		}
	}
}
