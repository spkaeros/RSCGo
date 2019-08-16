package packets

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
)

//Packet The definition of a game packet.  Generally, these are commands, indexed by their Opcode(0-255), with
//  a 5000-byte buffer for arguments, stored in Payload.  If the packet is bare, raw data is intended to be
//  transmitted when writing the packet stucture to a socket, otherwise we put a 2-byte unsigned short for the
//  length of the arguments buffer(plus one because the opcode is included in the payload size), and the 1-byte
//  opcode at the start of the packet, as a header for the client to easily parse the information for each frame.
type Packet struct {
	Opcode      byte
	Payload     []byte
	Bare        bool
	readIndex   int
	bitPosition int
	length      int
}

var ResponsePong = NewOutgoingPacket(3)

var LogWarning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)

//NewPacket Creates a new packet instance.
func NewPacket(opcode byte, payload []byte) *Packet {
	return &Packet{opcode, payload, false, 0, 0, 0}
}

//NewOutgoingPacket Creates a new packet instance intended for sending formatted data to the client.
func NewOutgoingPacket(opcode byte) *Packet {
	buf := []byte{0xA5, 0xA5, opcode}
	return &Packet{opcode, buf, false, 0, 0, 0}
}

//NewBarePacket Creates a new packet instance intended for sending raw data to the client.
func NewBarePacket(src []byte) *Packet {
	return &Packet{0, src, true, 0, 0, 0}
}

//ReadLong Read the next 64-bit integer from the packet payload.
func (p *Packet) ReadLong() (val uint64, err error) {
	for i := 7; i >= 0; i-- {
		b, err := p.ReadByte()
		if err == errors.BufferOverflow {
			LogWarning.Printf("WARNING: Tried to read data from empty packet in a long!  Rewinding offset...")
			p.readIndex -= 7 - i
			return 0, err
		}
		val |= uint64(b) << uint(i*8)
	}
	return val, nil
}

//ReadInt Read the next 32-bit integer from the packet payload.
func (p *Packet) ReadInt() (val uint32, err error) {
	for i := 3; i >= 0; i-- {
		b, err := p.ReadByte()
		if err == errors.BufferOverflow {
			LogWarning.Printf("WARNING: Tried to read data from empty packet in a int!  Rewinding offset...")
			p.readIndex -= 3 - i
			return 0, err
		}
		val |= uint32(b) << uint(i*8)
	}
	return val, nil
}

//ReadShort Read the next 16-bit integer from the packet payload.
func (p *Packet) ReadShort() (val uint16, err error) {
	for i := 1; i >= 0; i-- {
		b, err := p.ReadByte()
		if err == errors.BufferOverflow {
			LogWarning.Printf("WARNING: Tried to read data from empty packet in a short!  Rewinding offset...")
			p.readIndex -= 1 - i
			return 0, err
		}
		val |= uint16(b) << uint(i*8)
	}
	return val, nil
}

//ReadByte Read the next 8-bit integer from the packet payload.
func (p *Packet) ReadByte() (byte, error) {
	if p.readIndex+1 > len(p.Payload) {
		LogWarning.Printf("WARNING: Tried to read data from empty packet in a byte!  Rewinding offset...")
		return byte(0), errors.BufferOverflow
	}
	defer func() {
		p.readIndex++
	}()
	return p.Payload[p.readIndex] & 0xFF, nil
}

//ReadString Read the next variable-length C-string from the packet payload and return it as a Go-string.
//  This will keep reading data until it reaches a null-byte ( '\0', 0xA, 10 ).
func (p *Packet) ReadString() (val string, err error) {
	for c, err := p.ReadByte(); err == nil && c != 0xA; c, err = p.ReadByte() {
		if err == errors.BufferOverflow {
			return "", err
		}
		val += string(c)
	}
	return val, nil
}

//AddLong Adds a 64-bit integer to the packet payload.
func (p *Packet) AddLong(l uint64) *Packet {
	p.Payload = append(p.Payload, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	return p
}

//AddInt Adds a 32-bit integer to the packet payload.
func (p *Packet) AddInt(i uint32) *Packet {
	p.Payload = append(p.Payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return p
}

//AddShort Adds a 16-bit integer to the packet payload.
func (p *Packet) AddShort(s uint16) *Packet {
	p.Payload = append(p.Payload, byte(s>>8), byte(s))
	return p
}

//AddByte Adds an 8-bit integer to the packet payload.
func (p *Packet) AddByte(b uint8) *Packet {
	p.Payload = append(p.Payload, b)
	return p
}

func (p *Packet) String() string {
	return fmt.Sprintf("Packet{opcode='%d',len='%d',payload={ %v }}", p.Opcode, len(p.Payload), p.Payload)
}

var bitmasks = []int32{0, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f, 0xff, 0x1ff, 0x3ff, 0x7ff, 0xfff, 0x1fff,
	0x3fff, 0x7fff, 0xffff, 0x1ffff, 0x3ffff, 0x7ffff, 0xfffff, 0x1fffff, 0x3fffff, 0x7fffff, 0xffffff,
	0x1ffffff, 0x3ffffff, 0x7ffffff, 0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff, -1}
var bitPosition = 0

func (p *Packet) AddBits(value int, numBits int) *Packet {
	bytePos := (p.bitPosition >> 3) + 3
	bitOffset := 8 - (p.bitPosition & 7)
	p.bitPosition += numBits
	p.length = ((p.bitPosition + 7) / 8) + 3
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
