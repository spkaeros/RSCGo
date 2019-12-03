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
	if !IsTileBlocking(x, y-1, ClipSouth, false) {
		if neighbor := p.nodes[x<<13|(y-1)]; neighbor == nil {
			p.nodes[x<<13|(y-1)] = &pNode{loc: NewLocation(x, y-1)}
		}
		neighbors = append(neighbors, p.nodes[x<<13|(y-1)])
	}
	if !IsTileBlocking(x+1, y, ClipEast, false) {
		if neighbor := p.nodes[(x+1)<<13|y]; neighbor == nil {
			p.nodes[(x+1)<<13|y] = &pNode{loc: NewLocation(x+1, y)}
		}
		neighbors = append(neighbors, p.nodes[(x+1)<<13|y])
	}
	if !IsTileBlocking(x, y+1, ClipNorth, false) {
		if neighbor := p.nodes[x<<13|(y+1)]; neighbor == nil {
			p.nodes[x<<13|(y+1)] = &pNode{loc: NewLocation(x, y+1)}
		}
		neighbors = append(neighbors, p.nodes[x<<13|(y+1)])
	}
	if !IsTileBlocking(x-1, y, ClipWest, false) {
		if neighbor := p.nodes[(x-1)<<13|y]; neighbor == nil {
			p.nodes[(x-1)<<13|y] = &pNode{loc: NewLocation(x-1, y)}
		}
		neighbors = append(neighbors, p.nodes[(x-1)<<13|y])
	}

	if !IsTileBlocking(x-1, y-1, ClipSouth|ClipWest, false) {
		if !IsTileBlocking(x-1, y, ClipWest, false) && !IsTileBlocking(x, y-1, ClipSouth, false) {
			if neighbor := p.nodes[(x-1)<<13|(y-1)]; neighbor == nil {
				p.nodes[(x-1)<<13|(y-1)] = &pNode{loc: NewLocation(x-1, y-1)}
			}
			neighbors = append(neighbors, p.nodes[(x-1)<<13|(y-1)])
		}
	}
	if !IsTileBlocking(x+1, y-1, ClipSouth|ClipEast, false) {
		if !IsTileBlocking(x+1, y, ClipEast, false) && !IsTileBlocking(x, y-1, ClipSouth, false) {
			if neighbor := p.nodes[(x+1)<<13|(y-1)]; neighbor == nil {
				p.nodes[(x+1)<<13|(y-1)] = &pNode{loc: NewLocation(x+1, y-1)}
			}
			neighbors = append(neighbors, p.nodes[(x+1)<<13|(y-1)])
		}
	}
	if !IsTileBlocking(x+1, y+1, ClipNorth|ClipEast, false) {
		if !IsTileBlocking(x+1, y, ClipEast, false) && !IsTileBlocking(x, y+1, ClipNorth, false) {
			if neighbor := p.nodes[(x+1)<<13|(y+1)]; neighbor == nil {
				p.nodes[(x+1)<<13|(y+1)] = &pNode{loc: NewLocation(x+1, y+1)}
			}
			neighbors = append(neighbors, p.nodes[(x+1)<<13|(y+1)])
		}
	}
	if !IsTileBlocking(x-1, y+1, ClipNorth|ClipWest, false) {
		if !IsTileBlocking(x-1, y, ClipWest, false) && !IsTileBlocking(x, y+1, ClipNorth, false) {
			if neighbor := p.nodes[(x-1)<<13|(y+1)]; neighbor == nil {
				p.nodes[(x-1)<<13|(y+1)] = &pNode{loc: NewLocation(x-1, y+1)}
			}
			neighbors = append(neighbors, p.nodes[(x-1)<<13|(y+1)])
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
	if parent.loc.DeltaX(neighbor.loc) == 0 || parent.loc.DeltaY(neighbor.loc) == 0 {
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
			other.hCost = 1 * hCost(other, p.end)
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
	if IsTileBlocking(p.end.X(), p.end.Y(), ClipDiag1|ClipDiag2|ClipFullBlock, false) {
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
	if active.parent != nil {
		position := active.loc
		for !p.start.Equals(position) {
			path.addFirstWaypoint(position.X(), position.Y())
			active = active.parent
			position = active.loc
		}
	}
	return path
}
