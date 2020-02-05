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
	"strconv"

	"github.com/spkaeros/rscgo/pkg/log"
)

//Packet The definition of a game net.  Generally, these are commands, indexed by their Opcode(0-255), with
//  a 5000-byte buffer for arguments, stored in Payload.  If the net is bare, raw data is intended to be
//  transmitted when writing the net stucture to a socket, otherwise we put a 2-byte unsigned short for the
//  length of the arguments buffer(plus one because the opcode is included in the payload size), and the 1-byte
//  opcode at the start of the net, as a header for the client to easily parse the information for each frame.
type Packet struct {
	Opcode      byte
	Payload     []byte
	Bare        bool
	readIndex   int
	bitPosition int
	length      int
}

//NewPacket Creates a new net instance.
func NewPacket(opcode byte, payload []byte) *Packet {
	return &Packet{opcode, payload, false, 0, 0, 0}
}

//NewOutgoingPacket Creates a new net instance intended for sending formatted data to the client.
func NewOutgoingPacket(opcode byte) *Packet {
	return &Packet{opcode, []byte{opcode}, false, 0, 0, 0}
}

//NewBarePacket Creates a new net instance intended for sending raw data to the client.
func NewBarePacket(src []byte) *Packet {
	return &Packet{0, src, true, 0, 0, 0}
}

func (p *Packet) readVarLengthInt(numBytes int) uint64 {
	var val uint64
	for i := numBytes - 1; i >= 0; i-- {
		val |= uint64(p.ReadByte()) << uint(i*8)
	}
	return val
}

//ReadLong Read the next 64-bit integer from the net payload.
func (p *Packet) ReadLong() uint64 {
	return p.readVarLengthInt(8)
}

//ReadInt Read the next 32-bit integer from the net payload.
func (p *Packet) ReadInt() int {
	return int(p.readVarLengthInt(4))
}

//ReadShort Read the next 16-bit integer from the net payload.
func (p *Packet) ReadShort() int {
	return int(p.readVarLengthInt(2))
}

//ReadBool Reads the next byte, if it is 1 returns true, else returns false.
func (p *Packet) ReadBool() bool {
	return p.ReadByte() == 1
}

func (p *Packet) checkError(err error) bool {
	if err != nil {
		return false
	}
	return true
}

//ReadByte Read the next 8-bit integer from the net payload.
func (p *Packet) ReadByte() byte {
	if p.readIndex+1 > len(p.Payload) {
		log.Warning.Println("Error parsing net arguments: { opcode=" + strconv.Itoa(int(p.Opcode)) + "; offset=" + strconv.Itoa(p.readIndex) + " };")
		return byte(0)
	}
	defer func() {
		p.readIndex++
	}()
	return p.Payload[p.readIndex] & 0xFF
}

//ReadSByte Read the next 8-bit integer from the net payload, as a signed byte.
func (p *Packet) ReadSByte() int8 {
	if p.readIndex+1 > len(p.Payload) {
		log.Warning.Println("Error parsing net arguments: { opcode=" + strconv.Itoa(int(p.Opcode)) + "; offset=" + strconv.Itoa(p.readIndex) + " };")
		return int8(0)
	}
	defer func() {
		p.readIndex++
	}()
	return int8(p.Payload[p.readIndex])
}

//ReadString Read the next n bytes from the net payload and return it as a Go-string.
func (p *Packet) ReadString(n int) (val string) {
	for i := 0; i < n; i++ {
		val += string(p.ReadByte())
	}
	return
}

// FIXME: Below is custom, non-jagex, better ReadString for use in modified client
/*
//ReadString Read the next variable-length C-string from the net payload and return it as a Go-string.
// This will keep reading data until it reaches a null-byte or a new-line character ( '\0', 0xA, 0, 10 ).
func (p *Packet) ReadString() string {
	var val string
	for c := p.ReadByte(); c != 0 && c != 0xA; c = p.ReadByte() {
		val += string(c)
	}
	return val
}
*/

//AddLong Adds a 64-bit integer to the net payload.
func (p *Packet) AddLong(l uint64) *Packet {
	p.Payload = append(p.Payload, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	return p
}

//AddInt Adds a 32-bit integer to the net payload.
func (p *Packet) AddInt(i uint32) *Packet {
	p.Payload = append(p.Payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return p
}

//AddInt2 Adds a 32-bit integer or an 8-byte integer to the net payload, depending on value.
func (p *Packet) AddInt2(i uint32) *Packet {
	if i < 128 {
		p.Payload = append(p.Payload, uint8(i))
		return p
	}
	p.Payload = append(p.Payload, byte((i>>24)+128), byte(i>>16), byte(i>>8), byte(i))
	return p
}

func (p *Packet) SetShort(offset int, s uint16) *Packet {
	if offset >= len(p.Payload) || offset < 0 {
		log.Warning.Println("Attempted out of bounds Packet.SetShort: ", offset, s)
		return p
	}
	p.Payload[offset+1] = byte(s >> 8)
	p.Payload[offset+2] = byte(s)
	return p
}

//AddShort Adds a 16-bit integer to the net payload.
func (p *Packet) AddShort(s uint16) *Packet {
	p.Payload = append(p.Payload, byte(s>>8), byte(s))
	return p
}

//AddBool Adds a single byte to the payload, with the value 1 if b is true, and 0 if b is false.
func (p *Packet) AddBool(b bool) *Packet {
	if b {
		p.Payload = append(p.Payload, 1)
		return p
	}
	p.Payload = append(p.Payload, 0)
	return p
}

//AddByte Adds an 8-bit integer to the net payload.
func (p *Packet) AddByte(b uint8) *Packet {
	p.Payload = append(p.Payload, b)
	return p
}

//AddSByte Adds an 8-bit signed integer to the net payload.
func (p *Packet) AddSByte(b int8) *Packet {
	p.Payload = append(p.Payload, uint8(b))
	return p
}

//AddBytes Adds byte array to net payload
func (p *Packet) AddBytes(b []byte) *Packet {
	p.Payload = append(p.Payload, b...)
	return p
}

func (p *Packet) String() string {
	return fmt.Sprintf("Packet{opcode='%d',len='%d',payload={ %v }}", p.Opcode, len(p.Payload), p.Payload)
}

//AddBits Packs value into the numBits next bits of the packetbuilders byte buffer.
func (p *Packet) AddBits(value int, numBits int) *Packet {
	bitmasks := []int32{0, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f, 0xff, 0x1ff, 0x3ff, 0x7ff, 0xfff, 0x1fff,
		0x3fff, 0x7fff, 0xffff, 0x1ffff, 0x3ffff, 0x7ffff, 0xfffff, 0x1fffff, 0x3fffff, 0x7fffff, 0xffffff,
		0x1ffffff, 0x3ffffff, 0x7ffffff, 0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff, -1}
	bytePos := (p.bitPosition >> 3) + 1
	bitOffset := 8 - (p.bitPosition & 7)
	p.bitPosition += numBits
	p.length = ((p.bitPosition + 7) / 8) + 1
	for p.length > len(p.Payload) {
		p.Payload = append(p.Payload, 0)
	}
	for ; numBits > bitOffset; bitOffset = 8 {
		p.Payload[bytePos] &= byte(^bitmasks[bitOffset])
		p.Payload[bytePos] |= byte(value >> uint(numBits-bitOffset&int(bitmasks[bitOffset])))
		bytePos++
		numBits -= bitOffset
	}
	if numBits == bitOffset {
		p.Payload[bytePos] &= byte(^bitmasks[bitOffset])
		p.Payload[bytePos] |= byte(value & int(bitmasks[bitOffset]))
	} else {
		p.Payload[bytePos] &= byte(^(bitmasks[numBits] << uint(bitOffset-numBits)))
		p.Payload[bytePos] |= byte((value & int(bitmasks[numBits])) << uint(bitOffset-numBits))
	}

	return p
}
