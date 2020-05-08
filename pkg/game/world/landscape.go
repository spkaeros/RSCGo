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
	"fmt"
	"os"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/jag"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//TileData Represents a single tile in the game's landscape.
/*
	DiagonalWalls int
	HorizontalWalls byte
	VerticalWalls byte
	GroundElevation byte
	Roofs byte
	GroundTexture byte
GroundOverlay byte
*/
type TileData struct {
	CollisionMask int16
}

//Sector Represents a sector of 48x48(2304) tiles in the game's landscape.
type Sector struct {
	Tiles []TileData
}

//Sectors A map to store landscape sectors by their hashed file name.
var Sectors = make(map[int]*Sector)

//LoadCollisionData Loads the JAG archive './data/landscape.jag', decodes it, and stores the map sectors it holds in
// memory for quick access.
func LoadCollisionData() {
	archive := jag.New(config.DataDir() + string(os.PathSeparator) + "landscape.jag")

	entryFileCaret := 0
	metaDataCaret := 0
	// Sectors begin at: offsetX=48, offsetY=96
	for i := 0; i < archive.FileCount; i++ {
		id := int(uint32(archive.MetaData[metaDataCaret]&0xFF)<<24 | uint32(archive.MetaData[metaDataCaret+1]&0xFF)<<16 | uint32(archive.MetaData[metaDataCaret+2]&0xFF)<<8 | uint32(archive.MetaData[metaDataCaret+3]&0xFF))
		startCaret := entryFileCaret
		entryFileCaret += int(uint32(archive.MetaData[metaDataCaret+7]&0xFF)<<16 | uint32(archive.MetaData[metaDataCaret+8]&0xFF)<<8 | uint32(archive.MetaData[metaDataCaret+9]&0xFF))
		Sectors[id] = loadSector(archive.FileData[startCaret:entryFileCaret])
		metaDataCaret += 10
	}
}

type TileDefinition struct {
	Color   int
	Visible int
	Blocked int
}

var TileDefs []TileDefinition

//BoundaryDefs This holds the defining characteristics for all of the game's boundary scene objects, ordered by ID.
var BoundaryDefs []BoundaryDefinition

//BoundaryDefinition This represents a single definition for a single boundary object in the game.
type BoundaryDefinition struct {
	ID          int
	Name        string
	Commands    [2]string
	Description string
	Dynamic     bool
	Solid       bool
}

const (
	OverlayBlank = iota
	//OverlayGravel Used for roads, ID 1
	OverlayGravel
	//OverlayWater Used for regular water, ID 2
	OverlayWater
	//OverlayWoodFloor Used for the floors of buildings, ID 3
	OverlayWoodFloor
	//OverlayBridge Used for bridges, suspends wood floor over water, ID 4
	OverlayBridge
	//OverlayStoneFloor Used for the floors of buildings, ID 5
	OverlayStoneFloor
	//OverlayRedCarpet Used for the floors of buildings, ID 6
	OverlayRedCarpet
	//OverlayDarkWater Used for dark, swampy water, ID 7
	OverlayDarkWater
	//OverlayBlack Used for empty parts of upper planes, ID 8
	OverlayBlack
	//OverlayWhite Used as a separator, e.g for edge of water, mountains, etc.  ID 9
	OverlayWhite
	//OverlayBlack2 Not sure where it is used, ID 10
	OverlayBlack2
	//OverlayLava Used in dungeons and on Karamja/Crandor as lava, ID 11
	OverlayLava
	//OverlayBridge2 Used for a specific type of bridge, ID 12
	OverlayBridge2
	//OverlayBlueCarpet Used for the floors of buildings, ID 13
	OverlayBlueCarpet
	//OverlayPentagram Used for certain questing purposes, ID 14
	OverlayPentagram
	//OverlayPurpleCarpet Used for the floors of buildings, ID 15
	OverlayPurpleCarpet
	//OverlayBlack3 Not sure what it is used for, ID 16, traversable
	OverlayBlack3
	//OverlayStoneFloorLight Used for the entrance to temple of ikov, ID 17
	OverlayStoneFloorLight
	//OverlayUnknown Not sure what this is yet, ID 18
	OverlayUnknown
	//OverlayBlack4 Not sure what it is used for, ID 19
	OverlayBlack4
	//OverlayAgilityLog Blank suspended tile over blackness for agility challenged, ID 20
	OverlayAgilityLog
	//OverlayAgilityLog Blank suspended tile over blackness for agility challenged, ID 21
	OverlayAgilityLog2
	//OverlayUnknown2 Not sure what this is yet, ID 22
	OverlayUnknown2
	//OverlaySandFloor Used for sand floor, ID 23
	OverlaySandFloor
	//OverlayMudFloor Used for mud floor, ID 24
	OverlayMudFloor
	//OverlaySandFloor Used for water floor, ID 25
	OverlayWaterFloor
)

const (
	//ClipNorth Bitmask to represent a wall to the north.
	ClipNorth = 1 << iota
	//ClipEast Bitmask to represent a wall to the west.
	ClipEast
	//ClipSouth Bitmask to represent a wall to the south.
	ClipSouth
	//ClipWest Bitmask to represent a wall to the east.
	ClipWest
	ClipCanProjectile
	//ClipDiag1 Bitmask to represent a diagonal wall.
	ClipSwNe
	//ClipDiag2 Bitmask to represent a diagonal wall facing the opposite way.
	ClipSeNw
	//ClipFullBlock Bitmask to represent an object blocking an entire tile.
	ClipFullBlock
	// TODO: handle projectiles properly, after I define what properly means in this context
)

func ClipBit(direction int) int {
	var mask int
	if direction == North || direction == NorthEast || direction == NorthWest {
		mask |= ClipNorth
	}
	if direction == South || direction == SouthEast || direction == SouthWest {
		mask |= ClipSouth
	}
	if direction == East || direction == SouthEast || direction == NorthEast {
		mask |= ClipEast
	}
	if direction == West || direction == SouthWest || direction == NorthWest {
		mask |= ClipWest
	}
	return mask
}

//Masks returns appropriate collision bitmasks to check for obstacles on when traversing from this location
// toward the given x,y coordinates.
// Returns: [2]byte {verticalMasks, horizontalMasks}
func (l Location) Masks(x, y int) (masks [2]byte) {
	if x > l.X() {
		masks[0] |= ClipSouth
	} else {
		masks[0] |= ClipNorth
	}
	if y > l.Y() {
		masks[1] |= ClipEast
	} else {
		masks[1] |= ClipWest
	}
	// diags and solid objects are checked for automatically in the functions that you'd use this with, so
	return
}

//
func (l Location) Mask(toward Location) byte {
	masks := l.Masks(toward.X(), toward.Y())
	return masks[0] | masks[1]
}

/*
var blockedOverlays = [...]int{OverlayWater, OverlayDarkWater, OverlayBlack, OverlayWhite, OverlayLava, OverlayBlack2, OverlayBlack3, OverlayBlack4}

func isOverlayBlocked(overlay int) bool {
	for _, v := range blockedOverlays {
		if v == overlay {
			return true
		}
	}
	return false
}
*/
func IsTileBlocking(x, y int, bit byte, current bool) bool {
	return CollisionData(x, y).blocked(bit, current)
}

func (t TileData) blocked(bit byte, current bool) bool {
	if t.CollisionMask&int16(bit) != 0 {
		return true
	}
	// Diagonal walls (/, \) and impassable scenary objects (|=|) both effectively disable the occupied location
	// TODO: Is overlay clipping finished?
	return !current && (t.CollisionMask & (ClipSwNe | ClipSeNw | ClipFullBlock)) != 0
}

func sectorName(x, y int) string {
	regionX := (2304 + x) / RegionSize
	regionY := (1776 + y - (944 * ((y + 100) / 944))) / RegionSize
	return fmt.Sprintf("h%dx%dy%d", (y+100)/944, regionX, regionY)
}

func sectorFromCoords(x, y int) *Sector {
	if s, ok := Sectors[strutil.JagHash(sectorName(x, y))]; ok && s != nil {
		return s
	}
	// Default to returning a blank sector filled with zero-value tiles.
	return &Sector{Tiles: make([]TileData, 2304)}
}

func (s *Sector) tile(x, y int) TileData {
	areaX := (2304 + x) % RegionSize
	areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
	if len(s.Tiles) <= 0 {
		return TileData{}
	}
	return s.Tiles[areaX*RegionSize+areaY]
}

func CollisionData(x, y int) TileData {
	//sector := sectorFromCoords(x, y)
	//if sector == nil {
	//	return TileData{CollisionMask: ClipFullBlock}
	//}
	return sectorFromCoords(x, y).tile(x, y)
}

//loadSector Parses raw data into data structures that make up a 48x48 map sector.
func loadSector(data []byte) (s *Sector) {
	// If we were given less than the length of a decompressed, raw map sector
	if len(data) < 23040 {
		log.Warning.Printf("Too short sector data: %d\n", len(data))
		return nil
	}
	s = &Sector{Tiles: make([]TileData, 2304)}
	offset := 0

	blankCount := 0
	for x := 0; x < RegionSize; x++ {
		for y := 0; y < RegionSize; y++ {
			groundTexture := data[offset+1] & 0xFF
			groundOverlay := data[offset+2] & 0xFF
			//			roofTexture := data[offset+3] & 0xFF
			horizontalWalls := data[offset+4] & 0xFF
			verticalWalls := data[offset+5] & 0xFF
			diagonalWalls := int(uint32(data[offset+6]&0xFF)<<24 | uint32(data[offset+7]&0xFF)<<16 | uint32(data[offset+8]&0xFF)<<8 | uint32(data[offset+9]&0xFF))
			if groundOverlay == 250 {
				// -6 overflows to 250, and is water tile
				groundOverlay = OverlayWater
			}
			if (groundOverlay == 0 && (groundTexture) == 0) || groundOverlay == OverlayWater || groundOverlay == OverlayBlack {
				blankCount++
			}
			tileIdx := x*RegionSize + y
			if groundOverlay > 0 && int(groundOverlay) < len(TileDefs) && TileDefs[groundOverlay-1].Blocked != 0 {
				s.Tiles[tileIdx].CollisionMask |= ClipFullBlock
			}
			walls := [] []int {
				[]int{ int(verticalWalls)-1, ClipNorth, y },
				[]int{ int(horizontalWalls)-1, ClipEast,x },
			}
			for i := 0; i < 2; i++ {
				if walls[i][0] < 0 {
					continue
				}
				if boundary := walls[i][0]; boundary > len(BoundaryDefs)-1 {
					log.Debugf("Out of bounds indexing attempted into BoundaryDefs[%d]; while upper bound is currently %d\n", boundary, len(BoundaryDefs)-1)
					continue
				}
				if wall := BoundaryDefs[walls[i][0]]; !wall.Dynamic && wall.Solid {
					s.Tiles[x*RegionSize+y].CollisionMask |= int16(walls[i][1])
					if walls[i][2] > 0 {
						s.Tiles[(x-i)*RegionSize+((y-1)+i)].CollisionMask |= int16(walls[i][1]<<2)
					}
				}
			}
			// TODO: Affect adjacent tiles in an intelligent way to determine which are solid and which are not
			if diagonalWalls < 24000 && diagonalWalls > 0 {
				if diagonalWalls > 12000 {
					diagonalWalls -= 12000
				}
				diagonalWalls -= 1
				if wall := BoundaryDefs[diagonalWalls]; !wall.Dynamic && wall.Solid {
					if diagonalWalls > 12000 {
						s.Tiles[tileIdx].CollisionMask |= ClipSwNe
					} else {
						// diagonal that blocks: SE<->NW (/ aka |‾ or _|)
						s.Tiles[tileIdx].CollisionMask |= ClipSeNw
					}
				}
			}

			offset += 10
		}
	}
	if blankCount >= 2304 {
		return nil
	}

	return
}
