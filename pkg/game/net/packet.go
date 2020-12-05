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
	//	"math"
	"encoding/binary"
	"fmt"
	"io"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//WriteFlusher 
type WriteFlusher interface {
	io.Writer
	Flush() error
}

//Packet The definition of a game handlers.  Generally, these are commands, indexed by their Opcode(0-255), with
//  a 5000-byte buffer for arguments, stored in FrameBuffer.  If the handlers is bare, raw data is intended to be
//  transmitted when writing the handlers structure to a socket, otherwise we put a 2-byte unsigned short for the
//  length of the arguments buffer(plus one because the opcode is included in the payload size), and the 1-byte
//  opcode at the start of the handlers, as a header for the client to easily parse the information for each frame.
type Packet struct {
	Opcode      byte
	FrameBuffer []byte
	ReadIndex   int
	bitIndex    int
	Bare        bool
}

//NewPacket Creates a new handlers instance.
func NewPacket(opcode byte, payload []byte) *Packet {
	return &Packet{Opcode: opcode, FrameBuffer: payload, Bare: false}
}

//NewEmptyPacket Creates a new handlers instance intended for sending formatted data to the client.
func NewEmptyPacket(opcode byte) *Packet {
	return &Packet{Opcode: opcode, FrameBuffer: []byte{opcode}, Bare: false}
}

//NewReplyPacket Creates a new handlers instance intended for sending raw data to the client.
func NewReplyPacket(src []byte) *Packet {
	return &Packet{FrameBuffer: src, Bare: true}
}

//ReadUint128 Read the next 128-bit integer from the handlers payload, returns it as 2 uint64 words, 
func (p *Packet) ReadUint128() (lsb uint64, msb uint64) {
	return p.ReadUint64(), p.ReadUint64()
}

//ReadUint64 Read the next 64-bit integer from the handlers payload.
func (p *Packet) ReadUint64() uint64 {
	if checkError(p.Skip(8)) {
		return 0
	}
	// return uint64((p.ReadUint32() << 32) | p.ReadUint32())
	return binary.BigEndian.Uint64(p.FrameBuffer[p.ReadIndex-8:])
}

//ReadUint32 Read the next 32-bit integer from the handlers payload.
func (p *Packet) ReadUint32() int {
	if checkError(p.Skip(4)) {
		return 0
	}
	// return int((p.ReadUint16() << 16) | p.ReadUint16())
	return int(binary.BigEndian.Uint32(p.FrameBuffer[p.ReadIndex-4:]))
}

//ReadUint16 Read the next 16-bit integer from the handlers payload.
func (p *Packet) ReadUint16() int {
	// return int((p.ReadUint8() << 8) | p.ReadUint8())
	if checkError(p.Skip(2)) {
		return 0
	}
	return int(binary.BigEndian.Uint16(p.FrameBuffer[p.ReadIndex-2:]))
}

func checkError(err error) bool {
	if err != nil {
		debug.PrintStack()
		log.Warn("Packet error:", err)
		return true
	}
	return false
}

//ReadUint8 Read the next 8-bit integer from the handlers payload.
func (p *Packet) ReadUint8() uint8 {
	if checkError(p.Skip(1)) {
		return 0
	}
	return p.FrameBuffer[p.ReadIndex-1] & 0xFF
}

func (p *Packet) ReadUByte() byte {
	return p.ReadUint8()
}

//ReadInt8 returns the signed interpretation of the next payload byte.
func (p *Packet) ReadInt8() int8 {
	if checkError(p.Skip(1)) {
		return 0
	}
	return int8(p.FrameBuffer[p.ReadIndex-1])
}

//ReadBoolean Returns true if the next payload byte isn't 0
func (p *Packet) ReadBoolean() bool {
	if checkError(p.Skip(1)) {
		return false
	}
	return p.FrameBuffer[p.ReadIndex-1] != 0
}

func (p *Packet) Read(buf []byte) int {
	n := len(buf)
	if p.Available() < n {
		log.Warning.Println("PacketBufferError[OutOfBounds:Read] Tried to read too many bytes (" + strconv.Itoa(n) + ") from read buffer (length " + strconv.Itoa(p.Available()) + ")")
		return -1
	}
	copy(buf, p.FrameBuffer[p.ReadIndex:])
	if checkError(p.Skip(n)) {
		return -1
	}
	return n
}

//Flip Resets the read buffer caret to zero.
func (p *Packet) Flip() {
	p.ReadIndex = 0
}

//Rewind rewinds the reader index by n bytes
func (p *Packet) Rewind(n int) error {
	if n < 0 {
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Rewinding the buffer by less than 0 bytes is not permitted.  Perhaps you need *Packet.Skip ?", false)
	}
	if n > p.ReadIndex {
		p.ReadIndex = 0
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Tried to rewind reader caret ("+strconv.Itoa(p.ReadIndex)+") passed the start of the buffer (0)", false)
	}
	p.ReadIndex -= n
	return nil
}

//Skip skips the reader index by n bytes
func (p *Packet) Skip(n int) error {
	if n < 0 {
		return errors.NewNetworkError("BufferOutOfBounds; Skipping the buffer by less than 0 bytes is not permitted.  Perhaps you need *Packet.Rewind ?", false)
	}
	if p.Available() < n {
		p.ReadIndex = p.Length()
		return errors.NewNetworkError("Packet.Skip,BufferOutOfBounds; Tried to skip reader caret ("+strconv.Itoa(p.ReadIndex)+") passed the length of the buffer ("+strconv.Itoa(p.Length())+")", false)
	}
	p.ReadIndex += n
	return nil
}

//ReadStringN Reads the next n bytes from the payload and returns it as a UTF-8 string, regardless of payload contents.
func (p *Packet) ReadStringN(n int) (val string) {
	buf := make([]byte, n)
	readLen := p.Read(buf)
	if readLen < 0 {
		p.ReadIndex = p.Length()
		return string(p.FrameBuffer[p.ReadIndex:])
	}
	return string(buf)
}

//ReadString Read the next variable-length C-string from the handlers payload and return it as a Go-string.
// This will keep reading data until it reaches a string termination byte.
// String termination bytes in order of precedence:
// NULL (hex:0x00, escape-seq:'\x00', decimal:0);
// NewLine/LineFeed (hex:0xA, escape-seq:'\n', decimal:10)
func (p *Packet) ReadString() string {
	for _, separator := range []byte{0x00, 0x0A} {
		end := strings.IndexByte(string(p.FrameBuffer[p.ReadIndex:]), separator)
		if end <= 0 {
			continue
		}
		s := p.FrameBuffer[p.ReadIndex:p.ReadIndex+end]
		if checkError(p.Skip(end+1)) {
			log.Debug("Error caused by skipping bytes in ReadString")
		}
		return string(s)
	}
	s := p.FrameBuffer[p.ReadIndex:]
	if checkError(p.Skip(len(p.FrameBuffer) - p.ReadIndex)) {
		log.Debug("Error caused by skipping bytes in ReadString")
	}
	return string(s)
}

func (p *Packet) AddFramedString(s string) *Packet {
	return p.AddBytes(append(append([]byte{0}, append([]byte(s), 0)...)))
}

func (p *Packet) AddEncryptedString(s string) *Packet {
	log.Debug(s, strutil.ChatFilter.Format(s))
	buf, n := strutil.Encipher(strutil.ChatFilter.Format(s))
	if n < 128 {
		p.AddUint8(uint8(n))
	} else {
		p.AddUint8((uint8(n) >> 8) + 128)
		p.AddUint8(uint8(n))
	}
	return p.AddBytes(buf)
}

func (p *Packet) AddString(s string) *Packet {
	return p.AddUint8(0).AddBytes([]byte(s)).AddUint8(0)
}

func (p *Packet) EnsureCapacity(l int) {
	p.FrameBuffer = append(p.FrameBuffer, make([]byte, l)...)
}

//AddUint64 Adds a 64-bit integer to the handlers payload.
func (p *Packet) AddUint64(l uint64) *Packet {
	p.EnsureCapacity(8)
	binary.BigEndian.PutUint64(p.FrameBuffer[len(p.FrameBuffer)-8:], l)
	return p
}

//AddUint32 Adds a 32-bit integer to the handlers payload.
func (p *Packet) AddUint32(i uint32) *Packet {
	p.EnsureCapacity(4)
	binary.BigEndian.PutUint32(p.FrameBuffer[len(p.FrameBuffer)-4:], i)
	return p
}

func (p *Packet) AddSmart0832(i int) *Packet {
	if i >= 0x80 {
		return p.AddUint32(0x80<<24 | uint32(i))
	}
	return p.AddUint8(uint8(i))
}

func (p *Packet) AddSmart1632(i int) *Packet {
	i &= 0xFFFFFFFF
	if i <= 0xFFFF {
		return p.AddUint16(uint16(i))
	}
	return p.AddUint32(uint32(-2147483648 + i))
}

func (p *Packet) AddSmart0816(i int) *Packet {
	if i >= 0x80 {
		return p.AddUint16(0x80<<8 | uint16(i))
	}
	return p.AddUint8(uint8(i))
}

//AddUint8or32 Adds a 32-bit integer or an 8-byte integer to the handlers payload, depending on value.
// TODO: Deprecate and remove this in favor of above, improved name
func (p *Packet) AddUint8or32(i uint32) *Packet {
	return p.AddSmart0832(int(i))
}

//AddUint16 Adds a 16-bit integer to the handlers payload.
func (p *Packet) AddUint16(s uint16) *Packet {
	p.EnsureCapacity(2)
	binary.BigEndian.PutUint16(p.FrameBuffer[len(p.FrameBuffer)-2:], s)
	return p
}

//AddBoolean Adds a single byte to the payload, with the value 1 if b is true, and 0 if b is false.
func (p *Packet) AddBoolean(b bool) *Packet {
	p.EnsureCapacity(1)
	if b {
		p.FrameBuffer[len(p.FrameBuffer)-1] = 1
		return p
	}
	p.FrameBuffer[len(p.FrameBuffer)-1] = 0
	return p
}

//AddUint8 Adds an 8-bit integer to the handlers payload.
func (p *Packet) AddUint8(b uint8) *Packet {
	p.EnsureCapacity(1)
	p.FrameBuffer[len(p.FrameBuffer)-1] = b
	return p
}

//AddInt8 Adds an 8-bit signed integer to the handlers payload.
func (p *Packet) AddInt8(b int8) *Packet {
	p.EnsureCapacity(1)
	p.FrameBuffer[len(p.FrameBuffer)-1] = uint8(b)
	return p
}

//AddBytes Adds byte array to handlers payload
func (p *Packet) AddBytes(b []byte) *Packet {
	for _, v := range b {
		//		p.EnsureCapacity(1)
		p.AddUint8(v)
	}
	return p
}

var bitmasks [66]int32

func init() {
	for i := 0; i <= 64; i++ {
		bitmasks[i] = (1 << uint(i)) - 1
	}
	bitmasks[65] = -1
}

//AddSignedBits adds the value with the first bit masked off
func (p *Packet) AddSignedBits(value int, numBits int) *Packet {
	return p.AddBitmask(value&int(bitmasks[numBits]), numBits)
}

//AddBitmask Packs value into the numBits next bits of the packetbuilders byte buffer.
// Note: This method only keeps track of the data that it has written to the buffer; it will
// overwrite any non-bitmasked values in the buffer starting at the beginning.
func (p *Packet) AddBitmask(value int, numBits int) *Packet {
	// determine what byte we can start safely write this value into
	byteOffset := (p.bitIndex >> 3) + 1
	// determine what bit within that byte we can safely write this value into
	bitOffset := 8 - (p.bitIndex & 7)
	// increment written bits count
	p.bitIndex += numBits
	for ((p.bitIndex+7)>>3)+1 > len(p.FrameBuffer) {
		p.FrameBuffer = append(p.FrameBuffer, 0)
	}
	// Write our value, using some bitwise tricks to flip on the correct bits
	for numBits > bitOffset {
		// prepare the byte we're writing into for the new data, by masking away the bits we will not need
		p.FrameBuffer[byteOffset] &= byte(^bitmasks[bitOffset])
		// append bits of our value that fit onto byte
		p.FrameBuffer[byteOffset] |= byte(value >> uint(numBits-bitOffset&int(bitmasks[bitOffset])))
		// increment written bytes (maybe isn't necessary, we do not mix bitwise data with normal bytes ever)
		byteOffset++
		// decrease number of bits left to write
		numBits -= bitOffset
		bitOffset = 8
	}

	if numBits == bitOffset {
		// we reached the end of the last byte
		p.FrameBuffer[byteOffset] &= byte(^bitmasks[bitOffset])
		p.FrameBuffer[byteOffset] |= byte(value & int(bitmasks[bitOffset]))
	} else {
		// we were done encoding our value mid-byte, must reposition our bits
		p.FrameBuffer[byteOffset] &= byte(^(int(bitmasks[numBits]) << uint(bitOffset-numBits)))
		p.FrameBuffer[byteOffset] |= byte((value & int(bitmasks[numBits])) << uint(bitOffset-numBits))
	}
	return p

}

//Length returns length of byte buffer.
func (p *Packet) Length() int {
	return len(p.FrameBuffer)
}

//Available returns available read buffer bytes count.
func (p *Packet) Available() int {
	return p.Length() - p.ReadIndex
}

//Capacity returns the byte capacity left for this buffer
func (p *Packet) Capacity() int {
	return 5000 - p.Length()
}

func (p *Packet) String() string {
	return fmt.Sprintf("Packet{opcode:%d,available:%d,capacity:%d,payload:%v}", p.Opcode, p.Available(), p.Capacity(), p.FrameBuffer)
}

func (p *Packet) WriteIndex() int {
	return len(p.FrameBuffer)
}
