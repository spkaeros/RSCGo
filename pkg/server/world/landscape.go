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
	"github.com/spkaeros/rscgo/pkg/jag"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//TileData Represents a single tile in the game's landscape.
type TileData struct {
	/*
	DiagonalWalls int
	HorizontalWalls byte
	VerticalWalls byte
	GroundElevation byte
	Roofs byte
	GroundTexture byte
	 */
	GroundOverlay   byte
	CollisionMask   int
}

//Sector Represents a sector of 48x48(2304) tiles in the game's landscape.
type Sector struct {
	Tiles []TileData
}

//Sectors A map to store landscape sectors by their hashed file name.
var Sectors = make(map[int]*Sector)

//LoadMapData Loads the JAG archive './data/landscape.jag', decodes it, and stores the map sectors it holds in
// memory for quick access.
func LoadMapData() {
	archive := jag.New("./data/landscape.jag")

	entryFileCaret := 0
	metaDataCaret := 0
	for i := 0; i < archive.FileCount; i++ {
		metaDataCaret += 4
		id := readInt(archive.MetaData, metaDataCaret)
		metaDataCaret += 6
		startCaret := entryFileCaret
		entryFileCaret += readU24BitInt(archive.MetaData, metaDataCaret)
		Sectors[id] = LoadSector(archive.FileData[startCaret:entryFileCaret])
	}
}

type TileDefinition struct {
	Color int
	Visible int
	ObjectType int
}

var Tiles []TileDefinition

//Boundarys This holds the defining characteristics for all of the game's boundary scene objects, ordered by ID.
var Boundarys []BoundaryDefinition

//BoundaryDefinition This represents a single definition for a single boundary object in the game.
type BoundaryDefinition struct {
	ID          int
	Name        string
	Commands    []string
	Description string
	Unknown     int
	Traversable int
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
	//OverlayBlack3 Not sure what it is used for, ID 16
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
	//WallNorth Bitmask to represent a wall to the north.
	WallNorth = 1
	//WallSouth Bitmask to represent a wall to the south.
	WallSouth = 4
	//WallEast Bitmask to represent a wall to the west.
	WallEast = 2
	//WallWest Bitmask to represent a wall to the east.
	WallWest = 8
	//WallEast Bitmask to represent a diagonal wall.
	WallDiag1 = 0x10
	//WallEast Bitmask to represent a diagonal wall facing the opposite way.
	WallDiag2 = 0x20
	//WallEast Bitmask to represent an object occupying an entire tile.
	WallObject = 0x40
)

var BlockedOverlays = [...]int{OverlayWater, OverlayDarkWater, OverlayBlack, OverlayWhite, OverlayLava, OverlayBlack2, OverlayBlack3, OverlayBlack4}

func isOverlayBlocked(overlay int) bool {
	for _, v := range BlockedOverlays {
		if v == overlay {
			return true
		}
	}
	return false
}

func IsTileBlocking(x, y int, bit byte, current bool) bool {
	return ClipData(x, y).blocked(bit, current)
}

func (t TileData) blocked(bit byte, current bool) bool {
	if t.CollisionMask & int(bit) != 0 {
		return true
	}
	// Diag
	if !current && (t.CollisionMask & WallDiag1) != 0 {
		return true
	}
	// oppososite diag
	if !current && (t.CollisionMask & WallDiag2) != 0 {
		return true
	}
	// tile entirely blocked
	if !current && (t.CollisionMask & WallObject) != 0 {
		return true
	}
	// if it's not a traversable ground type
	return !current && isOverlayBlocked(int(t.GroundOverlay))
}

func SectorName(x, y int) string {
	regionX := (2304+x)/RegionSize
	regionY := (1776+y-(944*((y+100)/944)))/RegionSize
	return fmt.Sprintf("h%dx%dy%d", (y+100)/944, regionX, regionY)
}

func SectorFromCoords(x, y int) *Sector {
	return Sectors[strutil.JagHash(SectorName(x, y))]
}

func (s *Sector) Tile(x, y int) TileData {
	areaX := (2304+x) % RegionSize
	areaY := (1776+y-(944*((y+100)/944))) % RegionSize
	return s.Tiles[areaX * RegionSize + areaY]
}

func ClipData(x, y int) TileData {
	sector := SectorFromCoords(x, y)
	if sector == nil {
		return TileData{GroundOverlay:0, CollisionMask:0x40}
	}
	return sector.Tile(x, y)
}

//LoadSector Parses raw data into data structures that make up a 48x48 map sector.
func LoadSector(data []byte) (s *Sector) {
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
			diagonalWalls := int(uint32(data[offset+6]&0xFF) << 24 + uint32(data[offset+7]&0xFF) << 16 + uint32(data[offset+8]&0xFF) << 8 + uint32(data[offset+9]&0xFF))
			if groundOverlay == 250 {
				// -6 overflows to 250, and is water tile
				groundOverlay = 2
			}
			if (groundOverlay == 0 && (groundTexture) == 0) || groundOverlay == OverlayWater || groundOverlay == OverlayBlack {
				blankCount++
			}
			tileIdx := x*RegionSize+y
			s.Tiles[tileIdx].GroundOverlay = groundOverlay
			if groundOverlay > 0 && Tiles[groundOverlay-1].ObjectType != 0 {
				s.Tiles[tileIdx].CollisionMask |= 0x40
			}
			if verticalWalls > 0 && Boundarys[verticalWalls-1].Unknown == 0 && Boundarys[verticalWalls-1].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= WallNorth
				if tileIdx >= 1 {
					// -1 is tile x,y-1
					s.Tiles[tileIdx-1].CollisionMask |= WallSouth
				}
			}
			if horizontalWalls > 0 && Boundarys[horizontalWalls-1].Unknown == 0 && Boundarys[horizontalWalls-1].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= WallEast
				if tileIdx >= 48 {
					// -48 is tile x-1,y
					s.Tiles[tileIdx-48].CollisionMask |= WallWest
				}
			}
			if diagonalWalls > 0 && diagonalWalls < 12000 && Boundarys[diagonalWalls-1].Unknown == 0 && Boundarys[diagonalWalls-1].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= WallDiag2
			}
			if diagonalWalls >= 12000 && diagonalWalls < 24000 && Boundarys[diagonalWalls-12001].Unknown == 0 && Boundarys[diagonalWalls-12001].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= WallDiag1
			}
			offset += 10
		}
	}
	if blankCount >= 2304 {
		return nil
	}

	return
}

//readU24BitInt Reads an unsigned 3-byte int from data, starting at caret-3
func readU24BitInt(data []byte, caret int) int {
	return int(uint32(data[caret-3]&0xFF)<<16 + uint32(data[caret-2]&0xFF)<<8 + uint32(data[caret-1]&0xFF))
}

//readInt Reads an unsigned 3-byte int from data, starting at caret-3
func readInt(data []byte, caret int) int {
	return int(uint32(data[caret-4]&0xFF)<<24 + uint32(data[caret-3]&0xFF)<<16 + uint32(data[caret-2]&0xFF)<<8 + uint32(data[caret-1]&0xFF))
}
