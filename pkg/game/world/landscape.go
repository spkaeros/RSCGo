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
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/jag"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//CollisionMask Represents a single tile in the game's landscape.
type CollisionMask int16

//Sector Represents a sector of 48x48(2304) tiles in the game's landscape.
type Sector struct {
	Tiles []CollisionMask
}

//Sectors A map to store landscape sectors by their hashed file name.
var Sectors = make(map[int]*Sector)
var SectorsLock sync.RWMutex

//LoadCollisionData Loads the JAG archive './data/landscape.jag', decodes it, and stores the map sectors it holds in
// memory for quick access.
func LoadCollisionData() {
	archive := jag.New(config.DataDir() + string(os.PathSeparator) + "landscape.jag")

	fileOffset := 0
	metaDataOffset := 0
	// Sectors begin at: offsetX=48, offsetY=96
	for i := 0; i < archive.FileCount; i++ {
		id := int(binary.BigEndian.Uint32(archive.MetaData[metaDataOffset:]))
		compSz :=  int(uint32(archive.MetaData[metaDataOffset+7]&0xFF)<<16 | uint32(archive.MetaData[metaDataOffset+8]&0xFF)<<8 | uint32(archive.MetaData[metaDataOffset+9]&0xFF))
		SectorsLock.Lock()
		Sectors[id] = loadSector(archive.FileData[fileOffset:fileOffset+compSz])
		SectorsLock.Unlock()
		metaDataOffset += 10
		fileOffset += compSz
	}
}

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
	if y > l.Y() {
		masks[0] |= ClipNorth
	} else if y < l.Y() {
		masks[0] |= ClipSouth
	}
	if x > l.X() {
		masks[1] |= ClipEast
	} else if x < l.X() {
		masks[1] |= ClipWest
	}
	// diags and solid objects are checked for automatically in the functions that you'd use this with, so
	return masks
}

//
func (l Location) Mask(toward Location) byte {
	masks := l.Masks(toward.X(), toward.Y())
	return masks[0] | masks[1]
}

/*
var blockedOverlays = [...]int{OverlayWater, definitions.OverlayDarkWater, definitions.OverlayBlack, definitions.OverlayWhite, definitions.OverlayLava, definitions.OverlayBlack2, definitions.OverlayBlack3, definitions.OverlayBlack4}

func isOverlayBlocked(overlay int) bool {
	for _, v := range blockedOverlays {
		if v == definitions.Overlay {
			return true
		}
	}
	return false
}
*/
func IsTileBlocking(x, y int, bit byte, current bool) bool {
	return CollisionData(x, y).blocked(bit, current)
}

func (t CollisionMask) blocked(bit byte, current bool) bool {
	// Diagonal walls (/, \) and impassable scenary objects (|=|) both effectively disable the occupied location
	// TODO: Is overlay clipping finished?
	return t&CollisionMask(bit) != 0 || (!current && (t&(ClipSwNe|ClipSeNw|ClipFullBlock)) != 0)
}

func sectorName(x, y int) string {
	regionX := (2304 + x) / RegionSize
	regionY := (1776 + y - (944 * ((y + 100) / 944))) / RegionSize
	return fmt.Sprintf("h%dx%dy%d", (y+100)/944, regionX, regionY)
}

func sectorFromCoords(x, y int) *Sector {
	SectorsLock.RLock()
	defer SectorsLock.RUnlock()
	if s, ok := Sectors[strutil.JagHash(sectorName(x, y))]; ok && s != nil {
		return s
	}
	// Default to returning a blank sector filled with zero-value tiles.
	return &Sector{Tiles: make([]CollisionMask, 2304)}
}

func (s *Sector) tile(x, y int) CollisionMask {
	areaX := (2304 + x) % RegionSize
	areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
	if len(s.Tiles) <= 0 {
		return 0
	}
	return s.Tiles[areaX*RegionSize+areaY]
}

func CollisionData(x, y int) CollisionMask {
	return sectorFromCoords(x, y).tile(x, y)
}

//loadSector Parses raw data into data structures that make up a 48x48 map sector.
func loadSector(data []byte) (s *Sector) {
	// 48*48=2304 tiles per sector and 10 bytes per tile makes each sector 23040 byte
	if len(data) < 23040 {
		log.Warning.Printf("Too short sector data: %d\n", len(data))
		return nil
	}
	s = &Sector{Tiles: make([]CollisionMask, 2304)}
	offset := 0

	blankCount := 0
	for x := 0; x < RegionSize; x++ {
		for y := 0; y < RegionSize; y++ {
			groundTexture := data[offset+1] & 0xFF
			groundOverlay := data[offset+2] & 0xFF
			//roofTexture := data[offset+3] & 0xFF
			horizontalWalls := data[offset+4] & 0xFF
			verticalWalls := data[offset+5] & 0xFF
			diagonalWalls := binary.BigEndian.Uint32(data[offset+6:])
			if groundOverlay == 250 {
				// -6 overflows to 250, and is water tile
				groundOverlay = definitions.OverlayWater
			}
			if (groundOverlay == 0 && (groundTexture) == 0) || groundOverlay == definitions.OverlayWater || groundOverlay == definitions.OverlayBlack {
				blankCount++
			}
			tileIdx := x*RegionSize + y
			if groundOverlay > 0 && int(groundOverlay) < len(definitions.TileOverlays) && definitions.TileOverlays[groundOverlay-1].Blocked != 0 {
				s.Tiles[tileIdx] |= ClipFullBlock
			}
			walls := [][]int{
				{int(verticalWalls) - 1, ClipNorth, y},
				{int(horizontalWalls) - 1, ClipEast, x},
			}
			for i := 0; i < 2; i++ {
				if walls[i][0] < 0 {
					continue
				}
				if boundary := walls[i][0]; boundary > len(definitions.BoundaryObjects)-1 {
					log.Debugf("Out of bounds indexing attempted into definitions.BoundaryObjects[%d]; while upper bound is currently %d\n", boundary, len(definitions.BoundaryObjects)-1)
					continue
				}
				if wall := definitions.BoundaryObjects[walls[i][0]]; !wall.Dynamic && wall.Solid {
					s.Tiles[x*RegionSize+y] |= CollisionMask(walls[i][1])
					if walls[i][2] > 0 {
						s.Tiles[(x-i)*RegionSize+((y-1)+i)] |= CollisionMask(walls[i][1] << 2)
					}
				}
			}
			// TODO: Affect adjacent tiles in an intelligent way to determine which are solid and which are not
			if diagonalWalls < 24000 && diagonalWalls > 0 {
				if diagonalWalls > 12000 {
					diagonalWalls -= 12000
				}
				diagonalWalls -= 1
				if wall := definitions.BoundaryObjects[diagonalWalls]; !wall.Dynamic && wall.Solid {
					if diagonalWalls > 12000 {
						s.Tiles[tileIdx] |= ClipSwNe
					} else {
						// diagonal that blocks: SE<->NW (/ aka |‾ or _|)
						s.Tiles[tileIdx] |= ClipSeNw
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
