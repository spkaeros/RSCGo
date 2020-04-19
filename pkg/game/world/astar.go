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

func (p *Pathfinder) neighbors(n *pNode) (neighbors []*pNode) {
//	for deltaX := -1; deltaX <= 1; deltaX++ {
	for _, direction := range OrderedDirections {
		node := &pNode{loc: n.loc.Step(direction)}
		if !n.loc.Reachable(node.loc) {
			continue
		}
		if p.nodes.hasNode(node) {
			continue
		}
		neighbors = append(neighbors, node)
	}
//			p.visited = append(p.visited, node)
//		for deltaY := -1; deltaY <= 1; deltaY++ {
//			if deltaX|deltaY == 0 {
//				continue
//			}
//			if !p.start.WithinRange(node.loc, 20) {
//				continue
//			}
//		}
		
	return neighbors
}

type Pathfinder struct {
	nodes queue
	visited queue
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

func (q queue) hasNode(n *pNode) bool {
	for _, n1 := range q {
		if n.loc.X() == n1.loc.X() && n.loc.Y() == n1.loc.Y() {
			return true
		}
	}
	return false
}

//NewPathfinder Returns a new A* pathfinder instance to derive an optimal path from start to end.
func NewPathfinder(start, end Location) *Pathfinder {
	p := &Pathfinder{start: start, end: end, nodes: queue{&pNode{loc: start, open: true}}}
	heap.Init(&p.nodes)
	return p
}

func (n *pNode) gCostTo(neighbor *pNode) float64 {
//	deltaX := parent.loc.X() - neighbor.loc.X()
//	deltaY := parent.loc.Y() - neighbor.loc.Y()
//	toNext := (math.Sqrt(float64(deltaX*deltaX) + float64(deltaY*deltaY)))
//	return toNext * (math.Sqrt2 - 1.0)
	stepPrice := 1.0
	if n.loc.DeltaX(neighbor.loc) + n.loc.DeltaY(neighbor.loc) > 1 {
		stepPrice = math.Sqrt2
	}
	return n.gCost + stepPrice
//	return parent.gCost + ((math.Sqrt2 - 1.0) * float64(neighbor.loc.DeltaX(end)+neighbor.loc.DeltaY(end)))
//	return parent.gCost + math.Sqrt(float64(neighbor.loc.DeltaX(parent.loc)) + float64(neighbor.loc.DeltaY(parent.loc)))
}

func hCost(neighbor *pNode, end Location) float64 {
	return math.Sqrt(float64(neighbor.loc.DeltaX(end)*neighbor.loc.DeltaX(end)) + float64(neighbor.loc.DeltaY(end)*neighbor.loc.DeltaY(end)))
}

func (p *Pathfinder) compare(active, other *pNode) {
	gCost := active.gCostTo(other)
	if !other.open || gCost < other.gCost {
		other.gCost = gCost
		if other.hCost == 0 {
			other.hCost = hCost(other, p.end)*3
		}
		other.nCost = other.gCost + other.hCost
		other.parent = active

		if !other.open {
			other.open = true
			heap.Push(&p.nodes, other)
		} else {
			heap.Fix(&p.nodes, other.index)
		}
	}
}

func (p *Pathfinder) MakePath() *Pathway {
	if IsTileBlocking(p.end.X(), p.end.Y(), 0, false) {
		return NewPathwayToLocation(p.end)
	}
	for p.nodes.Len() > 0 {
//		active := p.nodes[0]
//		p.nodes = p.nodes[1:]
		active := heap.Pop(&p.nodes).(*pNode)
		active.closed = true
		position := active.loc
		if position.Equals(p.end) {
			path := &Pathway{StartX: 0, StartY: 0}
			for !p.start.Equals(position) {
				path.addFirstWaypoint(position.X(), position.Y())
				active = active.parent
				position = active.loc
			}
			return path
		}
		neighbors := p.neighbors(active)
		for _, neighbor := range neighbors {
			if !neighbor.closed {
				p.compare(active, neighbor)
			}
		}
	}
/*
	active := p.nodes[p.end.X()<<13|p.end.Y()]
	if active != nil && active.parent != nil {
		position := active.loc
		for !p.start.Equals(position) {
			path.addFirstWaypoint(position.X(), position.Y())
			active = active.parent
			position = active.loc
		}
	}
*/
	return nil
}
