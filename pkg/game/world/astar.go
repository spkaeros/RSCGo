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

type pNode struct {
	parent              *pNode
	loc                 Location
	hCost, gCost, nCost float64
	index               int
	open, closed        bool
}

func (p *Pathfinder) neighbors(n *pNode) []*pNode {
	x, y := n.loc.X(), n.loc.Y()
	var neighbors []*pNode
	for deltaX := -1; deltaX <= 1; deltaX++ {
		for deltaY := -1; deltaY <= 1; deltaY++ {
			if deltaX|deltaY == 0 {
				continue
			}
			neighborX, neighborY := x+deltaX, y+deltaY
			neighborHash := neighborX<<13|neighborY
			if n.loc.ReachableCoords(neighborX, neighborY) {
				if neighbor := p.nodes[neighborHash];
						neighbor == nil {
					p.nodes[neighborHash] = &pNode{loc: NewLocation(neighborX, neighborY)}
				}
				neighbors = append(neighbors, p.nodes[neighborHash])
			}
		}
	}
	return neighbors
}

type Pathfinder struct {
	nodes map[int]*pNode
	open  queue
	start Location
	end   Location
}

type queue []*pNode

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].nCost < q[j].nCost
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	node := x.(*pNode)
	node.index = n
	*q = append(*q, node)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*q = old[0 : n-1]
	return node
}

//NewPathfinder Returns a new A* pathfinder instance to derive an optimal path from start to end.
func NewPathfinder(start, end Location) *Pathfinder {
	p := &Pathfinder{start: start, end: end, nodes: make(map[int]*pNode), open: make(queue, 1)}
	p.open[0] = &pNode{loc: start, open: true}
	p.nodes[start.X()<<13|start.Y()] = p.open[0]
	p.nodes[end.X()<<13|end.Y()] = &pNode{loc: end}
	heap.Init(&p.open)
	return p
}

func gCost(parent, neighbor *pNode) float64 {
	if parent.loc.LongestDelta(neighbor.loc) == 0 {
		return parent.gCost + 1
	}
	return parent.gCost + math.Sqrt2
}

func hCost(neighbor *pNode, end Location) float64 {
	return (math.Sqrt2 - 1.0) * float64(neighbor.loc.DeltaX(end)+neighbor.loc.DeltaY(end))
}

func (p *Pathfinder) compare(active, other *pNode) {
	gCost := gCost(active, other)
	if !other.open || gCost < other.gCost {
		other.gCost = gCost
		if other.hCost == 0 {
			other.hCost = hCost(other, p.end)*3
		}
		other.nCost = other.gCost + other.hCost
		other.parent = active

		if !other.open {
			other.open = true
			heap.Push(&p.open, other)
		} else {
			heap.Fix(&p.open, other.index)
		}
	}
}

func (p *Pathfinder) MakePath() *Pathway {
	if IsTileBlocking(p.end.X(), p.end.Y(), 0, false) {
		return NewPathwayToLocation(p.end)
	}
	for p.open.Len() > 0 {
		active := heap.Pop(&p.open).(*pNode)
		active.closed = true
		position := active.loc
		if position.Equals(p.end) {
			break
		}
		for _, neighbor := range p.neighbors(active) {
			if neighbor.closed {
				continue
			}
			p.compare(active, neighbor)
		}
	}

	path := &Pathway{StartX: 0, StartY: 0}

	active := p.nodes[p.end.X()<<13|p.end.Y()]
	if active != nil && active.parent != nil {
		position := active.loc
		for !p.start.Equals(position) {
			path.addFirstWaypoint(position.X(), position.Y())
			active = active.parent
			position = active.loc
		}
	}
	return path
}
