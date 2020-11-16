/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmai.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in a copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package entity

type Type int

const (
	TypePlayer     Type = 1 << iota // 1
	TypeNpc                         // 2
	TypeObject                      // 4
	TypeDoor                        // 8
	TypeItem                        // 16
	TypeGroundItem                  // 32

	TypeMob    = TypePlayer | TypeNpc
	TypeEntity = TypeObject | TypeDoor | TypeItem | TypeGroundItem
)

type Entity interface {
	Location
	ServerIndex() int
	SetServerIndex(int)
	Type() Type
}

type Location interface {
	X() int
	Y() int
	SetX(int)
	SetY(int)
	Wilderness() int
	Above() Location
	Below() Location
	Collides(Location) bool
	PlaneY(bool) int
	Hash() int
	EuclideanDistance(o Location) float64
	Mask(o Location) byte
	Point() Location
	Clone() Location
	Reachable(Location) bool
	NextTo(Location) bool
	ReachableCoords(int, int) bool
	CanReach([2]Location) bool
	Near(Location, int) bool
	IsValid() bool
	Step(int) Location
	NextTileToward(o Location) Location
	DirectionToward(o Location) int
	DirectionTo(int, int) int
	WithinArea([2]Location) bool
	WithinReach(Location) bool
	LongestDelta(o Location) int
	LongestDeltaCoords(int, int) int
	DeltaX(o Location) int
	TheirDeltaX(o Location) int
	Delta(o Location) int
	DeltaY(o Location) int
	TheirDeltaY(o Location) int
	Equals(o interface{}) bool
	WithinRange(o Location, rad int) bool
	Plane() int
	PivotTo(Location) [2][]int
}
