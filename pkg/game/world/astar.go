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
	"container/heap"
	"math"
)

type tileNode struct {
	parent              *tileNode
	loc                 Location
	hCost, gCost, nCost float64
	index               int
	open, closed        bool
}

func (n *tileNode) gCostFrom(neighbor *tileNode) float64 {
	stepPrice := 1.0
	if n.loc.DeltaX(neighbor.loc)+n.loc.DeltaY(neighbor.loc) > 1 {
		stepPrice = math.Sqrt2
	}
	return n.gCost + stepPrice
}

type tileQueue []*tileNode

func (q tileQueue) Len() int {
	return len(q)
}

func (q tileQueue) Less(i, j int) bool {
	return q[i].nCost < q[j].nCost
}

func (q tileQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *tileQueue) Push(x interface{}) {
	n := len(*q)
	node := x.(*tileNode)
	node.index = n
	*q = append(*q, node)
}

func (q *tileQueue) Pop() interface{} {
	old := *q
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*q = old[0 : n-1]
	return node
}

type Pathfinder struct {
	tileQueue
	activeTiles map[int]*tileNode
	last        Location
	start       Location
	end         Location
}

// This produces a unique hash for each tile possible within the current game world; in this sense, it is a perfect hashing algorithm.
// If the world ever expands a great deal, this may need to change to accomodate larger values.
func (l Location) Hash() int {
	return (l.X() << 14) | l.Y()
}

//NewPathfinder Returns a new A* pathfinder instance to derive an optimal path from start to end.
func NewPathfinder(start, end Location) *Pathfinder {
	startNode := &tileNode{loc: start, open: true}
	p := &Pathfinder{start: start, end: end, tileQueue: tileQueue{startNode}, activeTiles: map[int]*tileNode{start.Hash(): startNode, end.Hash(): {loc: end}}}
	heap.Init(&p.tileQueue)
	return p
}

func (p *Pathfinder) node(l Location) *tileNode {
	hash := (l.X() << 16) | l.Y()
	if v, ok := p.activeTiles[hash]; !ok || v == nil {
		p.activeTiles[hash] = &tileNode{loc: l}
	}
	return p.activeTiles[hash]
}

func (p *Pathfinder) MakePath() *Pathway {
	if IsTileBlocking(p.end.X(), p.end.Y(), 0, false) {
		return NewPathwayToLocation(p.end)
	}
	defer func() {
		for _, v := range p.activeTiles {
			v.index = 0
			v.hCost, v.gCost, v.nCost = 0, 0, 0
		}
	}()
	for p.tileQueue.Len() > 0 {
		active := heap.Pop(&p.tileQueue).(*tileNode)
		active.closed = true
		makePath := func(active *tileNode) *Pathway {
			path := &Pathway{StartX: 0, StartY: 0}
			for active.parent != nil && !active.parent.loc.Equals(p.start) {
				path.addFirstWaypoint(active.loc.X(), active.loc.Y())
				active = active.parent
			}
			return path
		}
		position := active.loc
		// if p.last.LongestDelta(position) == 0 /*|| p.tileQueue.Len() > 512*/ {
			// DoS prevention measures; astar will run forever if you let it
			//			return makePath(active)
			// return nil
		// }
		if position.Equals(p.end) {
			// We made it!
			return makePath(active)
		}
		p.last = active.loc

		// OrderedDirections is ordered as orthogonal then diagonals.
		// Direction precedent: E,W,N,S,SW,SE,NW,NE
		for _, direction := range OrderedDirections {
			//			node := &tileNode{loc: active.loc.Step(direction), open: false, closed: false}
			neighbor := p.node(active.loc.Step(direction))
			if IsTileBlocking(neighbor.loc.X(), neighbor.loc.Y(), active.loc.Mask(neighbor.loc), false) {
//			if !active.loc.Reachable(neighbor.loc) {
				continue
			}
			gCost := active.gCostFrom(neighbor)
			if !neighbor.open || gCost < neighbor.gCost {
				if neighbor.hCost == 0 {
					neighbor.hCost = neighbor.loc.EuclideanDistance(p.end)
				}
				neighbor.gCost = gCost
				neighbor.nCost = gCost + neighbor.hCost
				neighbor.parent = active
				if !neighbor.open || neighbor.closed {
					neighbor.open = true
					heap.Push(&p.tileQueue, neighbor)
					continue
				}
				heap.Fix(&p.tileQueue, neighbor.index)
			}
		}
	}
	return nil
}
