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
	"strconv"
	"strings"
	"time"
)

type TileData struct {
	DiagonalWalls int
	HorizontalWalls byte
	VerticalWalls byte
	Roofs byte
	GroundElevation byte
	GroundOverlay byte
	GroundTexture byte
}

type Sector struct {
	Tiles [2304]*TileData
	Name string // e.g, h0x50y50
}

type SectorArchive struct {
	Sectors []*Sector
	stream *Stream
}

type Stream struct {
	offset int
	data []byte
}

var Identifiers []int
var decompressedSizes []int
var compressedSizes []int
var startOffsets []int
var Files [][]byte
var jagArchive = &SectorArchive{stream: &Stream{data: []byte{}, offset: 0}}
var bz2Header = [...]byte {'B', 'Z', 'h', '1'}
var activeSectors = [...]int { 343253403, 343253404, 343253405, 343253457, 343253458, 343253459, 343253460, 343253461, 343253462, 343253463, 343253466, 343253518, 343253519, 343253520, 343253521, 343253522, 343253523, 343253524, 343253525, 343480384, 343480385, 343480386, 343480438, 343480503, 343480504, 343480505, 343480506, 355283396, 355283397, 355283458, 355283512, 355283513, 355283515, 355283516, 355283517, 355283518, 355510377, 355510378, 355510439, 355510494, 355510495, 355510496, 355510497, 355510498, 355510499, 355737358, 355737359, 355737475, 355737476, 355737477, 355737478, 355737479, 355737480, 355964339, 355964340, 355964394, 355964455, 355964456, 355964457, 355964458, 355964459, 355964460, 355964461, 356191320, 356191321, 356191382, 356191437, 356191438, 356191439, 356191440, 356191441, 356191442, 356418301, 356418302, 356418303, 356418355, 356418356, 356418357, 356418358, 356418359, 356418360, 356418418, 356418419, 356418422, 356418423, 356645283, 356645284, 356645337, 356645338, 356645339, 356645340, 356645341, 356645344, 356645345, 356645397, 356645398, 356645399, 356645403, 356645404, 356872263, 356872264, 356872265, 356872317, 356872318, 356872319, 356872320, 356872321, 356872322, 356872324, 356872325, 356872326, 356872378, 356872380, 356872384, 356872385, 357099245, 357099246, 357099298, 357099299, 357099300, 357099301, 357099302, 357099307, 357099359, 357099361, 357099362, 357099363, 357099364, 357099365, 357099366, 357326225, 357326226, 357326227, 357326279, 357326280, 357326281, 357326282, 357326283, 357326284, 357326285, 357326340, 357326341, 357326342, 357326343, 357326344, 357326345, 357326346, 357326347, 369129237, 369129238, 369129239, 369129291, 369129292, 369129293, 369129294, 369129295, 369129296, 369129297, 369129298, 369129355, 369129357, 369129358, 369129359, 369356218, 369356219, 369356220, 369356272, 369356273, 369356274, 369356275, 369356276, 369356277, 369356278, 369356338, 369356339, 369356340, 369583199, 369583200, 369583201, 369583253, 369583254, 369583255, 369583256, 369583257, 369583258, 369583316, 369583317, 369583318, 369583319, 369583320, 369583321, 369810180, 369810181, 369810182, 369810234, 369810235, 369810236, 369810237, 369810238, 369810239, 369810296, 369810297, 369810298, 369810299, 369810300, 369810301, 369810302, 370037161, 370037162, 370037163, 370037215, 370037216, 370037217, 370037218, 370037219, 370037220, 370037221, 370037222, 370037223, 370037224, 370037276, 370037277, 370037278, 370037279, 370037280, 370037281, 370037282, 370037283, 370264142, 370264143, 370264144, 370264196, 370264197, 370264198, 370264199, 370264200, 370264201, 370264202, 370264203, 370264204, 370264205, 370264257, 370264258, 370264259, 370264260, 370264261, 370264262, 370264263, 370264264, 370491123, 370491124, 370491125, 370491177, 370491178, 370491179, 370491180, 370491181, 370491182, 370491183, 370491184, 370491185, 370491186, 370491238, 370491239, 370491240, 370491241, 370491242, 370491243, 370491244, 370491245, 370718104, 370718105, 370718106, 370718158, 370718159, 370718160, 370718161, 370718162, 370718163, 370718164, 370718165, 370718166, 370718167, 370718219, 370718220, 370718221, 370718222, 370718223, 370718224, 370718225, 370718226, 370945085, 370945086, 370945087, 370945139, 370945140, 370945141, 370945142, 370945143, 370945144, 370945145, 370945146, 370945147, 370945148, 370945200, 370945201, 370945202, 370945203, 370945204, 370945205, 370945206, 370945207, 336277248, 337639133, 285553830, 285553831, 285553832, 285553884, 285553885, 285553886, 285553887, 285553888, 285553889, 285553890, 285553891, 285553892, 285553893, 285553945, 285553946, 285553947, 285553948, 285553949, 285553950, 285553951, 285553952, 285780811, 285780812, 285780813, 285780865, 285780866, 285780867, 285780868, 285780869, 285780870, 285780871, 285780872, 285780873, 285780874, 285780926, 285780927, 285780928, 285780930, 285780931, 285780932, 285780933, 297583823, 297583824, 297583825, 297583877, 297583878, 297583879, 297583880, 297583881, 297583882, 297583883, 297583884, 297583885, 297583886, 297583938, 297583939, 297583940, 297583941, 297583942, 297583943, 297583944, 297583945, 297810804, 297810805, 297810806, 297810858, 297810859, 297810860, 297810861, 297810862, 297810863, 297810864, 297810865, 297810866, 297810867, 297810919, 297810920, 297810921, 297810922, 297810923, 297810924, 297810925, 297810926, 298037785, 298037786, 298037787, 298037839, 298037840, 298037841, 298037842, 298037843, 298037844, 298037845, 298037846, 298037847, 298037848, 298037900, 298037901, 298037902, 298037903, 298037904, 298037905, 298037906, 298037907, 298264766, 298264767, 298264768, 298264820, 298264821, 298264822, 298264823, 298264824, 298264825, 298264826, 298264827, 298264828, 298264829, 298264881, 298264882, 298264883, 298264884, 298264885, 298264886, 298264887, 298264888, 298491747, 298491748, 298491749, 298491801, 298491802, 298491803, 298491804, 298491805, 298491806, 298491807, 298491808, 298491809, 298491810, 298491862, 298491863, 298491864, 298491865, 298491866, 298491867, 298491868, 298491869, 298718728, 298718729, 298718730, 298718782, 298718783, 298718784, 298718785, 298718786, 298718787, 298718788, 298718789, 298718790, 298718791, 298718843, 298718844, 298718845, 298718846, 298718847, 298718848, 298718849, 298718850, 298945709, 298945710, 298945711, 298945763, 298945764, 298945765, 298945766, 298945767, 298945768, 298945769, 298945770, 298945771, 298945772, 298945824, 298945825, 298945826, 298945827, 298945828, 298945829, 298945830, 298945831, 299172690, 299172691, 299172692, 299172744, 299172745, 299172746, 299172747, 299172748, 299172749, 299172750, 299172751, 299172753, 299172805, 299172806, 299172807, 299172808, 299172809, 299172810, 299172811, 299172812, 299399671, 299399672, 299399673, 299399725, 299399726, 299399727, 299399728, 299399729, 299399730, 299399731, 299399732, 299399733, 299399734, 299399786, 299399787, 299399788, 299399789, 299399790, 299399791, 299399792, 299399793, 299626652, 299626653, 299626654, 299626706, 299626707, 299626708, 299626709, 299626710, 299626711, 299626712, 299626713, 299626714, 299626715, 299626767, 299626768, 299626769, 299626770, 299626771, 299626772, 299626773, 299626774, 311429664, 311429665, 311429666, 311429718, 311429719, 311429720, 311429721, 311429722, 311429723, 311429724, 311429725, 311429726, 311429727, 311429779, 311429780, 311429781, 311429782, 311429783, 311429784, 311429785, 311429786, 311656645, 311656646, 311656647, 311656699, 311656700, 311656701, 311656702, 311656703, 311656704, 311656705, 311656706, 311656707, 311656708, 311656760, 311656761, 311656762, 311656763, 311656764, 311656765, 311656766, 311656767, 311883626, 311883627, 311883628, 311883680, 311883681, 311883682, 311883683, 311883684, 311883685, 311883686, 311883687, 311883688, 311883742, 311883743, 311883744, 311883745, 311883746, 311883747, 311883748, 312110607, 312110608, 312110609, 312110661, 312110662, 312110663, 312110664, 312110665, 312110666, 312110667, 312110668, 312110669, 312110723, 312110724, 312110725, 312110726, 312110727, 312110728, 312110729, 312337588, 312337589, 312337590, 312337642, 312337643, 312337644, 312337645, 312337646, 312337647, 312337648, 312337649, 312337650, 312337704, 312337705, 312337706, 312337707, 312337708, 312337709, 312337710, 312564569, 312564570, 312564571, 312564623, 312564624, 312564625, 312564626, 312564627, 312564628, 312564629, 312564630, 312564631, 312564632, 312564684, 312564685, 312564686, 312564687, 312564688, 312564689, 312564690, 312564691, 312791550, 312791551, 312791552, 312791604, 312791605, 312791606, 312791607, 312791608, 312791609, 312791610, 312791611, 312791612, 312791613, 312791665, 312791666, 312791667, 312791668, 312791669, 312791670, 312791671, 312791672, 313018531, 313018532, 313018533, 313018585, 313018586, 313018587, 313018588, 313018589, 313018590, 313018591, 313018592, 313018593, 313018594, 313018646, 313018647, 313018648, 313018649, 313018650, 313018651, 313018652, 313018653, 313245512, 313245513, 313245514, 313245566, 313245567, 313245568, 313245569, 313245570, 313245571, 313245572, 313245573, 313245574, 313245575, 313245627, 313245628, 313245629, 313245630, 313245631, 313245632, 313245633, 313245634 }


func UnmarshalWorldMaps() {
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
		var buf []byte
		// All these temp vars to load this shit sucks, I'm impressed that it even works given that Jagex mutilated the
		//  BZ2 format when making JAG archives, imo, and I'm using the standard library bzip2 implementation...
		//  In order to get the standard library to decode this strange antiquated file format, I had to manually insert
		//  a BZ2 header('B','Z','h','1') in between the JAG 6-byte size header, and the bz2 payload.  Drop the 6-byte
		//  JAG header from the start of the data, and voila, JAG will decode just like a BZ2.
		jagScanner := bufio.NewScanner(bzip2.NewReader(bytes.NewReader(append(bz2Header[:], jagBuffer[jagArchive.stream.offset:]...))))
		jagScanner.Split(bufio.ScanBytes)
		for jagScanner.Scan()  {
			buf = append(buf, jagScanner.Bytes()...)
		}

		jagBuffer = buf
		jagArchive.stream.data = buf
		jagArchive.stream.offset = 0
	}
	totalFiles := int(jagArchive.stream.ReadUShort())
	endOffset := jagArchive.stream.offset+totalFiles*10 // file info is 10 bytes per file, entry file data immediately follows
	done := make(chan struct{})
	go func() {
		for i := 0; i < totalFiles; i++ {
			Identifiers = append(Identifiers, int(jagArchive.stream.ReadUInt()))
			decompressedSizes = append(decompressedSizes, int(jagArchive.stream.ReadU24BitInt()))
			compressedSizes = append(compressedSizes, int(jagArchive.stream.ReadU24BitInt()))
			startOffsets = append(startOffsets, endOffset)
			endOffset += compressedSizes[i]
			if compressedAsWhole {
				Files = append(Files, []byte{})
				buffer := bytes.NewReader(jagBuffer[startOffsets[i]:endOffset])
				r, err := gzip.NewReader(buffer)
				entryScanner := bufio.NewScanner(r)
				entryScanner.Split(bufio.ScanBytes)
				if err != nil {
					entryScanner = bufio.NewScanner(buffer)
					entryScanner.Split(bufio.ScanBytes)
				}
				for entryScanner.Scan() {
					Files[i] = append(Files[i], entryScanner.Bytes()...)
				}

				/*
					sector := LoadSector(fmt.Sprintf("%d", Identifiers[i]), Files[i])
					if sector != nil {
						jagArchive.Sectors = append(jagArchive.Sectors, sector)
					}
				*/
			} else {
				log.Warning.Println("Attempted to read JAG archive with individually compressed entries, no support for it yet")
				//			Files[Identifiers[i]] = jagArchive.stream.jagBuffer[startOffsets[i]:startOffsets[i]+decompressedSizes[i]]
				// TODO: Decompress individual archived files
			}
		}

		for _, name := range activeSectors {
			jagArchive.Sectors = append(jagArchive.Sectors, LoadSector(strconv.Itoa(name), Files[indexFromHash(name)]))
		}

		done <- struct{}{}
	}()
	go func() {
		select {
		case <-done:
			defer close(done)
			log.Info.Println("Finished loading landscape data for clipping in", time.Since(start))
		}
	}()
}

func init() {
	log.Info.Println(runtime.GOMAXPROCS(runtime.NumCPU() * 4), runtime.NumCPU())
}

func LoadSector(name string, data []byte) (s *Sector) {
	if len(data) < 10*2304 {
		log.Warning.Printf("Too short sector data: %d\n", len(data))
		return nil
	}
	s = &Sector{Name: name}
	for _, tile := range s.Tiles {
		if tile == nil {
			tile = new(TileData)
		}
		tile.GroundElevation = data[0] & 0xFF
		tile.GroundTexture = data[1] & 0xFF
		tile.GroundOverlay = data[2] & 0xFF
		tile.Roofs = data[3] & 0xFF
		tile.HorizontalWalls = data[4] & 0xFF
		tile.VerticalWalls = data[5] & 0xFF
		tile.DiagonalWalls = int(uint32(data[6]&0xFF<<24) | uint32(data[7]&0xFF<<16) |
			uint32(data[8]&0xFF<<8) | uint32(data[9]&0xFF))
	}
	return
}

func indexFromHash(hash int) int {
	for i, v := range Identifiers {
		if v == hash {
			return i
		}
	}
	return 0
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
	return uint16(s.data[s.offset - 2] & 0xFF) << 8 | uint16(s.data[s.offset-1]&0xFF)
}

func (s *Stream) ReadU24BitInt() uint32 {
	s.offset += 3
	return uint32(s.data[s.offset - 3] & 0xFF) << 16 + uint32(s.data[s.offset - 2] & 0xFF) << 8 + uint32(s.data[s.offset - 1]&0xFF)
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
	return uint32(s.data[s.offset] & 0xFF) << 24 + uint32(s.data[s.offset + 1] & 0xFF) << 16 + uint32(s.data[s.offset + 2] & 0xFF) << 8 + uint32(s.data[s.offset+3]&0xFF)
}

func (s *Stream) ReadLong() uint64 {
	return uint64((s.ReadUInt() & 0xFFFFFFFF) << 32) + uint64(s.ReadUInt() & 0xFFFFFFFF)
}

func (s *Stream) ReadShort() int16 {
	defer func() {
		s.offset += 2
	}()
	j := int(s.ReadByte()) * 256 + int(s.ReadByte())
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
	if length > 0 && length < len(s.data) - s.offset {
		builder := strings.Builder{}
		builder.Write(s.data[s.offset:s.offset+length])
		return builder.String()
	}
	return ""
}

func (s *Stream) ReadString() string {
	length := int(0)
	for s.data[s.offset+length] != 0xA {
		length++
	}
	if length > 0 && length < len(s.data) - s.offset {
		builder := strings.Builder{}
		builder.Write(s.data[s.offset:s.offset+length])
		return builder.String()
	}
	return ""
}