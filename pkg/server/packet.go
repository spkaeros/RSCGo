package server

import (
	"fmt"
)

type Packet struct {
	Opcode  byte
	Payload []byte
	Length  int
	bare    bool
	offset  int
}

func NewPacket(opcode byte, payload []byte, length int) *Packet {
	return &Packet{opcode, payload, length, false, 0}
}

func (p *Packet) ReadLong() (val uint64) {
	for i := 7; i >= 0; i-- {
		val |= uint64(p.ReadByte()) << uint(i*8)
	}
	return
}
func (p *Packet) ReadInt() (val uint32) {
	for i := 3; i >= 0; i-- {
		val |= uint32(p.ReadByte()) << uint(i*8)
	}
	return
}
func (p *Packet) ReadShort() (val uint16) {
	for i := 1; i >= 0; i-- {
		val |= uint16(p.ReadByte()) << uint(i*8)
	}
	return
}
func (p *Packet) ReadByte() (val uint8) {
	if p.offset+1 >= p.Length {
		fmt.Println("WARNING: Trying to read into packet with empty buffer!")
		return 0
	}
	defer func() {
		p.offset++
	}()
	return p.Payload[p.offset] & 0xFF
}

func (p *Packet) ReadString() (val string) {
	for c := p.ReadByte(); c != 0xA; c = p.ReadByte() {
		val += string(c)
	}
	return
}

func (p *Packet) AddLong(l uint64) {
	defer func() {
		p.Length += 8
	}()
	p.Payload = append(p.Payload, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
}

func (p *Packet) AddInt(i uint32) {
	defer func() {
		p.Length += 4
	}()
	p.Payload = append(p.Payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

func (p *Packet) AddShort(s uint16) {
	defer func() {
		p.Length += 2
	}()
	p.Payload = append(p.Payload, byte(s>>8), byte(s))
}

func (p *Packet) AddByte(b uint8) {
	defer func() {
		p.Length++
	}()
	p.Payload = append(p.Payload, b)
}

func (c *Client) ReadPacket() (*Packet, error) {
	buf, err := c.Read(3)
	if err != nil {
		return nil, err
	}
	length := int(int16(buf[0])<<8 | int16(buf[1]))
	opcode := buf[2] & 0xFF

	payloadBuffer, err := c.Read(length)
	if err != nil {
		return nil, err
	}

	return NewPacket(opcode, payloadBuffer, length), nil
}

func (p *Packet) prependHeader() {
	dataLen := len(p.Payload) + 1 // opcode
	p.Payload = append([]byte{byte((dataLen >> 8) & 0xFF), byte(dataLen & 0xFF), p.Opcode}, p.Payload...)
}

func (c *Client) WritePacket(p *Packet) {
	if !p.bare {
		p.prependHeader()
	}

	c.Write(p.Payload)
}
