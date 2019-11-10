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
	"github.com/spkaeros/rscgo/pkg/jag"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"io/ioutil"
	"runtime"
	"sync"
)

//TileData Represents a single tile in the game's landscape.
type TileData struct {
	DiagonalWalls   int
	HorizontalWalls byte
	VerticalWalls   byte
	Roofs           byte
	GroundElevation byte
	GroundOverlay   byte
	GroundTexture   byte
}

//Sector Represents a sector of 48x48(2304) tiles in the game's landscape.
type Sector struct {
	Tiles []TileData
}

//Sectors A slice filled with map sector data.
var Sectors []*Sector

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

	decodeFile := func(startCaret, fileSize int) {
		defer wg.Done()
		gzLock.Lock()
		err := gzReader.Reset(bytes.NewBuffer(archive.FileData[startCaret:startCaret+fileSize]))
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
			Sectors = append(Sectors, sector)
			sectorLock.Unlock()
		}
		runtime.GC()
	}

	entryFileCaret := 0
	metaDataCaret := 0
	for i := 0; i < archive.FileCount; i++ {
		metaDataCaret += 10
		fileSize := readU24BitInt(archive.MetaData, metaDataCaret)
		go decodeFile(entryFileCaret, fileSize)
		entryFileCaret += fileSize
	}
	wg.Wait()
}

//LoadSector Parses raw data into data structures that make up a 48x48 map sector.
func LoadSector(data []byte) (s *Sector) {
	// If we were given less than the length of a decompressed, raw map sector
	if len(data) < 23040 {
		log.Warning.Printf("Too short sector data: %d\n", len(data))
		return nil
	}
	s = &Sector{Tiles: make([]TileData, 2304)}
	blankCount := 0
	for _, tile := range s.Tiles {
		tile.GroundElevation = data[0] & 0xFF
		tile.GroundTexture = data[1] & 0xFF
		tile.GroundOverlay = data[2] & 0xFF
		if tile.GroundOverlay == 0 {
			blankCount++
		}
		tile.Roofs = data[3] & 0xFF
		tile.HorizontalWalls = data[4] & 0xFF
		tile.VerticalWalls = data[5] & 0xFF
		tile.DiagonalWalls = int(uint32(data[6]&0xFF<<24) | uint32(data[7]&0xFF<<16) |
			uint32(data[8]&0xFF<<8) | uint32(data[9]&0xFF))
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
