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
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
)

type TileData struct {
	DiagonalWalls   int
	HorizontalWalls byte
	VerticalWalls   byte
	Roofs           byte
	GroundElevation byte
	GroundOverlay   byte
	GroundTexture   byte
}

type Sector struct {
	Tiles [2304]*TileData
}

type SectorArchive struct {
	Sectors []*Sector
	stream  *Stream
}

type Stream struct {
	offset int
	data   []byte
}


func UnmarshalWorldMaps() {
	var jagArchive = &SectorArchive{stream: &Stream{offset: 0}}
	var bz2Header = [...]byte{'B', 'Z', 'h', '1'}
	var buf []byte
	jagBuffer, err := ioutil.ReadFile("./data/landscape.jag")
	if err != nil {
		log.Warning.Println("Problem occurred attempting to open the landscape file at './jagBuffer/landscape.jag':", err)
		return
	}
	jagArchive.stream.data = jagBuffer
	// decompressed_len medium int, compressed_len medium int, total of 6 bytes header.
	compressedAsWhole := jagArchive.stream.ReadU24BitInt() != jagArchive.stream.ReadU24BitInt()
	start := time.Now()
	if compressedAsWhole {
		// All these temp vars to load this shit sucks, I'm impressed that it even works given that Jagex mutilated the
		//  BZ2 format when making JAG archives, imo, and I'm using the standard library bzip2 implementation...
		//  In order to get the standard library to decode this strange antiquated file format, I had to manually insert
		//  a BZ2 header('B','Z','h','1') in between the JAG 6-byte size header, and the bz2 payload.  Drop the 6-byte
		//  JAG header from the start of the data, and voila, JAG will decode just like a BZ2.
		jagScanner := bufio.NewScanner(bzip2.NewReader(bytes.NewReader(append(bz2Header[:], jagBuffer[jagArchive.stream.offset:]...))))
		jagScanner.Split(bufio.ScanBytes)
		for jagScanner.Scan() {
			buf = append(buf, jagScanner.Bytes()...)
		}

		jagArchive.stream.data = buf
		jagArchive.stream.offset = 0
	}
	done := make(chan struct{})
	go func() {
		totalFiles := int(jagArchive.stream.ReadUShort())
		endOffset := jagArchive.stream.offset + totalFiles*10 // file info is 10 bytes per file, entry file data immediately follows
		for i := 0; i < totalFiles; i++ {
			startOffset := endOffset
			jagArchive.stream.offset += 7
			endOffset += int(jagArchive.stream.ReadU24BitInt())
			if compressedAsWhole {
				r, err := gzip.NewReader(bytes.NewReader(buf[startOffset:endOffset]))
				entryScanner := bufio.NewScanner(r)
				if err != nil {
					// If we error on attempting to gunzip, it's not compressed.
					entryScanner = bufio.NewScanner(bytes.NewReader(buf[startOffset:endOffset]))
				}
				entryScanner.Split(bufio.ScanBytes)
				var entryData []byte
				for entryScanner.Scan() {
					entryData = append(entryData, entryScanner.Bytes()...)
				}
				if sector := LoadSector(entryData); sector != nil {
					jagArchive.Sectors = append(jagArchive.Sectors, sector)
				}
				entryData = []byte{}
			} else {
				log.Warning.Println("Attempted to read JAG archive with individually compressed entries, no support for it yet")
				//			Files[Identifiers[i]] = jagArchive.stream.jagBuffer[startOffsets[i]:startOffsets[i]+decompressedSizes[i]]
				// TODO: Decompress individual archived files
			}
		}

		buf = []byte{}
		jagArchive.stream.data = []byte{}
		done <- struct{}{}
	}()
	go func() {
		select {
		case <-done:
			defer close(done)
			defer runtime.GC()
			log.Info.Println("Finished loading landscape data for clipping in", time.Since(start))
		}
	}()
}

func LoadSector(data []byte) (s *Sector) {
	if len(data) < 10*2304 {
//		log.Warning.Printf("Too short sector data: %d\n", len(data))
		return nil
	}
	s = &Sector{}
	blankCount := 0
	for _, tile := range s.Tiles {
		if tile == nil {
			tile = new(TileData)
		}
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

func (s *Stream) ReadUByte() uint8 {
	defer func() {
		s.offset++
	}()
	return s.data[s.offset] & 0xFF
}

func (s *Stream) ReadByte() int8 {
	defer func() {
		s.offset++
	}()
	return int8(s.data[s.offset])
}

func (s *Stream) ReadUShort() uint16 {
	s.offset += 2
	return uint16(s.data[s.offset-2]&0xFF)<<8 | uint16(s.data[s.offset-1]&0xFF)
}

func (s *Stream) ReadU24BitInt() uint32 {
	s.offset += 3
	return uint32(s.data[s.offset-3]&0xFF)<<16 + uint32(s.data[s.offset-2]&0xFF)<<8 + uint32(s.data[s.offset-1]&0xFF)
}

func (s *Stream) ReadSmart2() uint16 {
	i := s.data[s.offset] & 0xFF
	if i < 128 {
		return uint16(s.ReadUByte())
	}

	return s.ReadUShort() - 32768
}

func (s *Stream) ReadSmart() uint16 {
	i := s.data[s.offset] & 0xFF
	if i < 128 {
		return uint16(s.ReadUByte()) - 64
	}

	return s.ReadUShort() - 49152
}

func (s *Stream) ReadUInt() uint32 {
	defer func() {
		s.offset += 4
	}()
	return uint32(s.data[s.offset]&0xFF)<<24 + uint32(s.data[s.offset+1]&0xFF)<<16 + uint32(s.data[s.offset+2]&0xFF)<<8 + uint32(s.data[s.offset+3]&0xFF)
}

func (s *Stream) ReadLong() uint64 {
	return uint64((s.ReadUInt()&0xFFFFFFFF)<<32) + uint64(s.ReadUInt()&0xFFFFFFFF)
}

func (s *Stream) ReadShort() int16 {
	defer func() {
		s.offset += 2
	}()
	j := int(s.ReadByte())*256 + int(s.ReadByte())
	if j > 0x7FFF {
		j -= 0x10000
	}
	return int16(j)
}

func (s *Stream) ReadStringIndex() string {
	length := int(0)
	for s.data[s.offset+length] != 0 {
		length++
	}
	if length > 0 && length < len(s.data)-s.offset {
		builder := strings.Builder{}
		builder.Write(s.data[s.offset : s.offset+length])
		return builder.String()
	}
	return ""
}

func (s *Stream) ReadString() string {
	length := int(0)
	for s.data[s.offset+length] != 0xA {
		length++
	}
	if length > 0 && length < len(s.data)-s.offset {
		builder := strings.Builder{}
		builder.Write(s.data[s.offset : s.offset+length])
		return builder.String()
	}
	return ""
}
