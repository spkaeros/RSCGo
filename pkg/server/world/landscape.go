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
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/spkaeros/rscgo/pkg/jag"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"io/ioutil"
	"runtime"
	"sync"
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

//Sectors A slice filled with map sector data.
//var Sectors []*Sector
var Sectors = make(map[int]*Sector)

//sectorLock Mutexes to safely load the game's landscape data concurrently.
// TODO: Would a semaphore work better for this?
var sectorLock sync.RWMutex

//LoadMapData Loads the JAG archive './data/landscape.jag', decodes it, and stores the map sectors it holds in
// memory for quick access.
func LoadMapData() {
	archive := jag.New("./data/landscape.jag")
	var gzLock sync.Mutex
	var gzReader = new(gzip.Reader)
	defer gzReader.Close()
	var wg sync.WaitGroup
	wg.Add(archive.FileCount)
	Boundarys = append(Boundarys, BoundaryDefinition{})

	decodeFile := func(data []byte, id int) {
		defer wg.Done()
		gzLock.Lock()
		err := gzReader.Reset(bytes.NewBuffer(data))
		if err != nil {
			log.Warning.Println("Ran into some sort of problem with jag entry gzReader:", err)
			gzLock.Unlock()
			return
		}
		sectorData, err := ioutil.ReadAll(gzReader)
		gzLock.Unlock()
		if err != nil {
			log.Warning.Println("Ran into some sort of problem with gunzip on jag archive entry:", err)
			return
		}
		if sector := LoadSector(sectorData); sector != nil {
			sectorLock.Lock()
			Sectors[id] = sector
			sectorLock.Unlock()
		}
		runtime.GC()
	}

	entryFileCaret := 0
	metaDataCaret := 0
	for i := 0; i < archive.FileCount; i++ {
		metaDataCaret += 4
		id := readInt(archive.MetaData, metaDataCaret)
		metaDataCaret += 6
		startCaret := entryFileCaret
		entryFileCaret += readU24BitInt(archive.MetaData, metaDataCaret)
		go decodeFile(archive.FileData[startCaret:entryFileCaret], id)
	}
	wg.Wait()
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
	OverlayGravel = 1
	OverlayWater = 2
	OverlayWood = 3
	OverlayStoneFloor = 5
	OverlayRedCarpet = 6
	OverlayDarkWater = 7
	OverlayBlack = 8
	OverlayWhite = 9
	OverlayLava = 11
	OverlayPentagram = 14
)

var BlockedOverlays = [...]int{OverlayWater, OverlayDarkWater, OverlayBlack, OverlayWhite, OverlayLava}

func isOverlayBlocked(overlay int) bool {
	for _, v := range BlockedOverlays {
		if v == overlay {
			return true
		}
	}
	return false
}

func isTileBlocking(x, y int, bit byte) bool {
	return ClipData(x, y).blocked(bit)
}

func (t TileData) blocked(bit byte) bool {
	if t.CollisionMask & int(bit) != 0 {
		return true
	}
	// Diag
	if t.CollisionMask & 0x10 != 0 {
		return true
	}
	// oppososite diag
	if t.CollisionMask & 0x20 != 0 {
		return true
	}
	// tile entirely blocked
	if t.CollisionMask & 0x40 != 0 {
		return true
	}
	// if it's not a traversable ground type
	return isOverlayBlocked(int(t.GroundOverlay))
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
		return TileData{}
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
//			s.Tiles[x*RegionSize+y].GroundElevation = data[offset+0] & 0xFF
//			s.Tiles[x*RegionSize+y].GroundTexture = data[offset+1] & 0xFF
//			s.Tiles[x*RegionSize+y].GroundOverlay = data[offset+2] & 0xFF
//			s.Tiles[x*RegionSize+y].Roofs = data[offset+3] & 0xFF
//			s.Tiles[x*RegionSize+y].HorizontalWalls = data[offset+4] & 0xFF
//			s.Tiles[x*RegionSize+y].VerticalWalls = data[offset+5] & 0xFF
//			s.Tiles[x*RegionSize+y].DiagonalWalls = int(uint32(data[offset+6]&0xFF) << 24 + uint32(data[offset+7]&0xFF) << 16 +
//				uint32(data[offset+8]&0xFF) << 8 + uint32(data[offset+9]&0xFF))
			groundOverlay := data[offset+2] & 0xFF
			horizontalWalls := data[offset+4] & 0xFF
			verticalWalls := data[offset+5] & 0xFF
			diagonalWalls := int(uint32(data[offset+6]&0xFF) << 24 + uint32(data[offset+7]&0xFF) << 16 + uint32(data[offset+8]&0xFF) << 8 + uint32(data[offset+9]&0xFF))
			if groundOverlay == 250 {
				// -6 overflows to 250, and is water tile
				groundOverlay = 2
			}
			if (groundOverlay == 0 && (data[offset+1] & 0xFF) == 0) || groundOverlay == OverlayWater || groundOverlay == OverlayBlack {
				blankCount++
			}
			tileIdx := x*RegionSize+y
			s.Tiles[tileIdx].GroundOverlay = groundOverlay
			if groundOverlay > 0 && Tiles[groundOverlay-1].ObjectType != 0 {
				s.Tiles[tileIdx].CollisionMask |= 0x40
			}
			if verticalWalls > 0 && Boundarys[verticalWalls-1].Unknown == 0 && Boundarys[verticalWalls-1].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= 1
				if x > 0 || y > 0 {
					// -1 is tile x,y-1
					s.Tiles[tileIdx-1].CollisionMask |= 4
				}
			}
			if horizontalWalls > 0 && Boundarys[horizontalWalls-1].Unknown == 0 && Boundarys[horizontalWalls-1].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= 2
				if x >= 1 || y >= RegionSize {
					// -48 is tile x-1,y
					s.Tiles[tileIdx-48].CollisionMask |= 8
				}
			}
			if diagonalWalls > 0 && diagonalWalls < 12000 && Boundarys[diagonalWalls-1].Unknown == 0 && Boundarys[diagonalWalls-1].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= 0x20
			}
			if diagonalWalls >= 12000 && diagonalWalls < 24000 && Boundarys[diagonalWalls-12001].Unknown == 0 && Boundarys[diagonalWalls-12001].Traversable != 0 {
				s.Tiles[tileIdx].CollisionMask |= 0x10
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
