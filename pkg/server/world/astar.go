/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import "math"

type Node struct {
	hCost, gCost int
	cost         int
	open         bool
	parent       *Node
	loc          Location
}

type Pathfinder struct {
	nodes map[int]*Node
	open []*Node
	start Location
	end Location
}

//NewPathfinder Returns a new A* pathfinder instance to derive an optimal path from start to end.
func NewPathfinder(start, end Location) *Pathfinder {
	p := &Pathfinder{start: start, end: end, nodes: make(map[int]*Node), open: []*Node{{loc: start, open: true}}}
	p.nodes[start.X() << 32 | start.Y()] = p.open[0]
	p.nodes[end.X() << 32 | end.Y()] = &Node{loc: end, open: true}
	return p
}

func (p *Pathfinder) hasOpen(node *Node) bool {
	for _, n := range p.open {
		if node == n {
			return true
		}
	}
	return false
}

func (p *Pathfinder) getCheapest() *Node {
	var node *Node
	min := math.MaxInt32
	for _, n := range p.open {
		if !n.open {
			continue
		}
		if n.cost < min {
			min = n.cost
			node = n
		}
	}
	return node
}

func travelCost(start, end Location) int {
	deltaX, deltaY := start.DeltaX(end), start.DeltaY(end)
	shortL, longL := deltaX, deltaY
	if deltaX > deltaY {
		shortL = deltaY
		longL = deltaX
	}
	return shortL * 14 + ((longL - shortL) * 10)
}

func (p *Pathfinder) removeOpen(node *Node) {
	for i, n := range p.open {
		if n == node {
			p.open = append(p.open[:i], p.open[i+1:]...)
			break
		}
	}
	node.open = false
}

func (p *Pathfinder) compare(active, other *Node) {
	gCost := active.gCost + travelCost(active.loc, other.loc)
	cost := travelCost(other.loc, p.end)
	fCost := gCost + cost
	if other.cost > fCost {
		p.removeOpen(other)
	} else if other.open && !p.hasOpen(other) {
		other.gCost = gCost
		other.hCost = cost
		other.cost = fCost
		other.parent = active
		p.open = append(p.open, other)
	}
}

func (p *Pathfinder) MakePath() *Pathway {
	if IsTileBlocking(p.end.X(), p.end.Y(), 0x40, false) {
		return NewPathwayToLocation(p.end)
	}
	for len(p.open) > 0 {
		active := p.getCheapest()
		position := active.loc
		if position.Equals(p.end) {
			break
		}
		p.removeOpen(active)

		x, y := position.X(), position.Y()
		for nextX := x - 1; nextX <= x + 1; nextX++ {
			for nextY := y - 1; nextY <= y + 1; nextY++ {
				if nextX == x && nextY == y {
					continue
				}

				adj := NewLocation(nextX, nextY)
				sprites := [3][3]int{{SouthWest, West, NorthWest}, {South, -1, North}, {SouthEast, East, NorthEast}}
				xIndex, yIndex := position.X()-adj.X()+1, position.Y()-adj.Y()+1
				nextTileMask := 4
				curTileMask := 1
				if xIndex < 0 || xIndex >= 3 {
					continue
				}
				if yIndex < 0 || yIndex >= 3 {
					continue
				}
				dir := sprites[xIndex][yIndex]
				switch dir {
				case North:
					nextTileMask = WallSouth
					curTileMask = WallNorth
				case South:
					nextTileMask = WallNorth
					curTileMask = WallSouth
				case East:
					nextTileMask = WallWest
					curTileMask = WallEast
				case West:
					nextTileMask = WallEast
					curTileMask = WallWest
				case NorthEast:
					nextTileMask = WallSouth | WallWest
					curTileMask = WallNorth | WallWest
				case NorthWest:
					nextTileMask = WallSouth | WallEast
					curTileMask = WallNorth | WallEast
				case SouthEast:
					nextTileMask = WallNorth | WallWest
					curTileMask = WallSouth | WallWest
				case SouthWest:
					nextTileMask = WallNorth | WallEast
					curTileMask = WallSouth | WallEast
				}
				if !IsTileBlocking(position.X(), position.Y(), byte(curTileMask), true) && !IsTileBlocking(adj.X(), adj.Y(), byte(nextTileMask), false) {
					switch dir {
					case NorthEast:
						if IsTileBlocking(position.X(), position.Y()-1, byte(nextTileMask), false) {
							continue
						}
						if IsTileBlocking(position.X()-1, position.Y(), byte(nextTileMask), false) {
							continue
						}
					case NorthWest:
						if IsTileBlocking(position.X(), position.Y()-1, byte(nextTileMask), false) {
							continue
						}
						if IsTileBlocking(position.X()+1, position.Y(), byte(nextTileMask), false) {
							continue
						}
					case SouthEast:
						if IsTileBlocking(position.X(), position.Y()+1, byte(nextTileMask), false) {
							continue
						}
						if IsTileBlocking(position.X()-1, position.Y(), byte(nextTileMask), false) {
							continue
						}
					case SouthWest:
						if IsTileBlocking(position.X(), position.Y()+1, byte(nextTileMask), false) {
							continue
						}
						if IsTileBlocking(position.X()+1, position.Y(), byte(nextTileMask), false) {
							continue
						}
					}
					node, ok := p.nodes[adj.X() << 32 | adj.Y()] //&Node{loc: adj, open: true}
					if !ok {
						node = &Node{loc:adj, open:true}
						p.nodes[adj.X() << 32 | adj.Y()] = node
					}
					p.compare(active, node)
				}
			}
		}
	}

	path := &Pathway{StartX: 0, StartY: 0}

	active := p.nodes[p.end.X() << 32 | p.end.Y()]
	if active.parent != nil {
		position := active.loc
		for !p.start.Equals(position) {
			path.AddWaypoint(position.X(), position.Y())
			active = active.parent
			position = active.loc
		}
	}
	return path
}
