/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package net

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/log"
)

//Packet The definition of a game net.  Generally, these are commands, indexed by their Opcode(0-255), with
//  a 5000-byte buffer for arguments, stored in FrameBuffer.  If the net is bare, raw data is intended to be
//  transmitted when writing the net structure to a socket, otherwise we put a 2-byte unsigned short for the
//  length of the arguments buffer(plus one because the opcode is included in the payload size), and the 1-byte
//  opcode at the start of the net, as a header for the client to easily parse the information for each frame.
type Packet struct {
	Opcode       byte
	FrameBuffer  []byte
	HeaderBuffer []byte
	readIndex    int
	bitIndex     int
}

//NewPacket Creates a new net instance.
func NewPacket(opcode byte, payload []byte) *Packet {
	return &Packet{Opcode: opcode, FrameBuffer: payload, HeaderBuffer: make([]byte, 2)}
}

//NewEmptyPacket Creates a new net instance intended for sending formatted data to the client.
func NewEmptyPacket(opcode byte) *Packet {
	return &Packet{Opcode: opcode, FrameBuffer: []byte{opcode}, HeaderBuffer: make([]byte, 2)}
}

//NewReplyPacket Creates a new net instance intended for sending raw data to the client.
func NewReplyPacket(src []byte) *Packet {
	return &Packet{FrameBuffer: src}
}

func (p *Packet) readNSizeUints(n int) []uint64 {
	read := func(numBytes int) uint64 {
		var val uint64
		buf := make([]byte, numBytes)
		_ = p.Read(buf)
		for idx, b := range buf {
			val |= uint64(b) << uint((numBytes-1-idx)<<3)
		}
		return val
	}

	var set []uint64
	for ; n > 0; n -= 8 {
		set = append(set, read(int(math.Min(float64(n), 8))))
	}
	return set
}

//ReadUint128 Read the next 128-bit integer from the net payload.
func (p *Packet) ReadUint128() (msb uint64, lsb uint64) {
	buf := p.readNSizeUints(16)
	return buf[0], buf[1]
}

//ReadUint64 Read the next 64-bit integer from the net payload.
func (p *Packet) ReadUint64() uint64 {
	return p.readNSizeUints(8)[0]
}

//ReadUint32 Read the next 32-bit integer from the net payload.
func (p *Packet) ReadUint32() int {
	return int(p.readNSizeUints(4)[0])
}

//ReadUint16 Read the next 16-bit integer from the net payload.
func (p *Packet) ReadUint16() int {
	return int(p.readNSizeUints(2)[0])
}

func checkError(err error) bool {
	if err != nil {
		log.Warning.Println(err)
		return true
	}
	return false
}

//ReadUint8 Read the next 8-bit integer from the net payload.
func (p *Packet) ReadUint8() byte {
	if checkError(p.Skip(1)) {
		return 0
	}
	return p.FrameBuffer[p.readIndex-1] & 0xFF
}

//ReadInt8 returns the signed interpretation of the next payload byte.
func (p *Packet) ReadInt8() int8 {
	return int8(p.ReadUint8())
}

//ReadBoolean Returns true if the next payload byte isn't 0
func (p *Packet) ReadBoolean() bool {
	return p.ReadUint8() != 0
}

func (p *Packet) Read(buf []byte) int {
	n := len(buf)
	if p.Available() < n {
		log.Warning.Println("PacketBufferError[OutOfBounds:Read] Tried to read too many bytes (" + strconv.Itoa(n) + ") from read buffer (length " + strconv.Itoa(p.Available()) + ")")
		return -1
	}
	copy(buf, p.FrameBuffer[p.readIndex:])
	p.Skip(n)
	return n
}

func (p *Packet) Flip() {
	p.readIndex = 0
}

//Rewind rewinds the reader index by n bytes
func (p *Packet) Rewind(n int) error {
	if n < 0 {
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Rewinding the buffer by less than 0 bytes is not permitted.  Perhaps you need *Packet.Skip ?", false)
	}
	if n > p.readIndex {
		p.readIndex = 0
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Tried to rewind reader caret ("+strconv.Itoa(p.readIndex)+") passed the start of the buffer (0)", false)
	}
	p.readIndex -= n
	return nil
}

//Skip skips the reader index by n bytes
func (p *Packet) Skip(n int) error {
	if n < 0 {
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Skipping the buffer by less than 0 bytes is not permitted.  Perhaps you need *Packet.Rewind ?", false)
	}
	if p.Available() < n {
		p.readIndex = p.Length()
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Tried to skip reader caret ("+strconv.Itoa(p.readIndex)+") passed the length of the buffer ("+strconv.Itoa(p.Length())+")", false)
	}
	p.readIndex += n
	return nil
}

//ReadStringN Reads the next n bytes from the payload and returns it as a UTF-8 string, regardless of payload contents.
func (p *Packet) ReadStringN(n int) (val string) {
	buf := make([]byte, n)
	readLen := p.Read(buf)
	if readLen < 0 {
		p.readIndex = p.Length()
		return string(p.FrameBuffer[p.readIndex:])
	}
	return string(buf)
}

//ReadString Read the next variable-length C-string from the net payload and return it as a Go-string.
// This will keep reading data until it reaches a null-byte or a new-line character ( '\0', 0xA, 0, 10 ).
func (p *Packet) ReadString() string {
	start := p.readIndex
	s := string(p.FrameBuffer[start:])
	end := strings.IndexByte(s, 0)+1
	if end < 0 {
		end = strings.IndexByte(s, '\n')+1
		if end < 0 {
			end = p.Length()
		}
	}
	p.readIndex += end
	return s[:end-1]
}

//AddUint64 Adds a 64-bit integer to the net payload.
func (p *Packet) AddUint64(l uint64) *Packet {
	p.FrameBuffer = append(p.FrameBuffer, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	return p
}

//AddUint32 Adds a 32-bit integer to the net payload.
func (p *Packet) AddUint32(i uint32) *Packet {
	p.FrameBuffer = append(p.FrameBuffer, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return p
}

//AddUint8or32 Adds a 32-bit integer or an 8-byte integer to the net payload, depending on value.
func (p *Packet) AddUint8or32(i uint32) *Packet {
	if i < 128 {
		p.FrameBuffer = append(p.FrameBuffer, uint8(i))
		return p
	}
	p.FrameBuffer = append(p.FrameBuffer, byte((i>>24)+128), byte(i>>16), byte(i>>8), byte(i))
	return p
}

//SetUint16At Rewrites the data at offset to the provided short uint value
func (p *Packet) SetUint16At(offset int, s uint16) *Packet {
	if offset >= len(p.FrameBuffer) || offset < 0 {
		log.Warning.Println("Attempted out of bounds Packet.SetUint16At: ", offset, s)
		return p
	}
	p.FrameBuffer[offset+1] = byte(s >> 8)
	p.FrameBuffer[offset+2] = byte(s)
	return p
}

//AddUint16 Adds a 16-bit integer to the net payload.
func (p *Packet) AddUint16(s uint16) *Packet {
	p.FrameBuffer = append(p.FrameBuffer, byte(s>>8), byte(s))
	return p
}

//AddBoolean Adds a single byte to the payload, with the value 1 if b is true, and 0 if b is false.
func (p *Packet) AddBoolean(b bool) *Packet {
	if b {
		p.FrameBuffer = append(p.FrameBuffer, 1)
		return p
	}
	p.FrameBuffer = append(p.FrameBuffer, 0)
	return p
}

//AddUint8 Adds an 8-bit integer to the net payload.
func (p *Packet) AddUint8(b uint8) *Packet {
	p.FrameBuffer = append(p.FrameBuffer, b)
	return p
}

//AddInt8 Adds an 8-bit signed integer to the net payload.
func (p *Packet) AddInt8(b int8) *Packet {
	p.FrameBuffer = append(p.FrameBuffer, uint8(b))
	return p
}

//AddBytes Adds byte array to net payload
func (p *Packet) AddBytes(b []byte) *Packet {
	p.FrameBuffer = append(p.FrameBuffer, b...)
	return p
}

//AddBitmask Packs value into the numBits next bits of the packetbuilders byte buffer.
func (p *Packet) AddBitmask(value int, numBits int) *Packet {
	masks := []int32{0, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f, 0xff, 0x1ff, 0x3ff, 0x7ff, 0xfff, 0x1fff,
		0x3fff, 0x7fff, 0xffff, 0x1ffff, 0x3ffff, 0x7ffff, 0xfffff, 0x1fffff, 0x3fffff, 0x7fffff, 0xffffff,
		0x1ffffff, 0x3ffffff, 0x7ffffff, 0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff, -1}
	byteOffset := (p.bitIndex >> 3) + 1
	bitOffset := 8 - (p.bitIndex & 7)
	p.bitIndex += numBits
	length := ((p.bitIndex + 7) / 8) + 1
	for length > len(p.FrameBuffer) {
		p.FrameBuffer = append(p.FrameBuffer, 0)
	}
	for ; numBits > bitOffset; bitOffset = 8 {
		p.FrameBuffer[byteOffset] &= byte(^masks[bitOffset])
		p.FrameBuffer[byteOffset] |= byte(value >> uint(numBits-bitOffset&int(masks[bitOffset])))
		byteOffset++
		numBits -= bitOffset
	}
	if numBits == bitOffset {
		p.FrameBuffer[byteOffset] &= byte(^masks[bitOffset])
		p.FrameBuffer[byteOffset] |= byte(value & int(masks[bitOffset]))
	} else {
		p.FrameBuffer[byteOffset] &= byte(^(int(masks[numBits]) << uint(bitOffset-numBits)))
		p.FrameBuffer[byteOffset] |= byte((value & int(masks[numBits])) << uint(bitOffset-numBits))
	}

	return p
}

//Length returns length of byte buffer.
func (p *Packet) Length() int {
	return len(p.FrameBuffer)
}

//Available returns available read buffer bytes count.
func (p *Packet) Available() int {
	return p.Length() - p.readIndex
}

//Capacity returns the byte capacity left for this buffer
func (p *Packet) Capacity() int {
	return 5000 - p.Length()
}

func (p *Packet) String() string {
	return fmt.Sprintf("Packet{opcode:%d,available:%d,capacity:%d,payload:%v}", p.Opcode, p.Available(), p.Capacity(), p.FrameBuffer)
}

func (p *Packet) ReadIndex() int {
	return p.readIndex
}
