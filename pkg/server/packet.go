package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type Packet struct {
	opcode  byte
	payload []byte
	length  int
	bare    bool
	offset  int
}

func NewPacket(opcode byte, payload []byte, length int) *Packet {
	return &Packet{opcode, payload, length, false, 0}
}

func (p *Packet) DecryptRSA() bool {
	buf, err := rsa.DecryptPKCS1v15(rand.Reader, RsaKey, p.payload)
	if err != nil {
		fmt.Println("WARNING: Could not decrypt RSA login block")
		return false
	}
	p.payload = buf
	p.length = len(buf)
	return true
}

func (p *Packet) ReadLong() (val int64) {
	for i := 7; i >= 0; i-- {
		val |= int64(p.ReadByte()) << uint(i*8)
	}
	return
}
func (p *Packet) ReadInt() (val int32) {
	for i := 3; i >= 0; i-- {
		val |= int32(p.ReadByte()) << uint(i*8)
	}
	return
}
func (p *Packet) ReadShort() (val int16) {
	for i := 1; i >= 0; i-- {
		val |= int16(p.ReadByte()) << uint(i*8)
	}
	return
}
func (p *Packet) ReadByte() (val byte) {
	if p.offset+1 >= p.length {
		fmt.Println("WARNING: Trying to read into packet with empty buffer!")
		return 0
	}
	defer func() {
		p.offset++
	}()
	return p.payload[p.offset] & 0xFF
}

func (p *Packet) ReadString(len int) string {
	if p.offset+len > p.length {
		fmt.Printf("WARNING: Requested string length too long.  Requested %d, only %d left in buffer.\n", len, p.length-p.offset)
		len = p.length - p.offset
	}
	defer func() {
		p.offset += len
	}()
	return string(p.payload[p.offset : p.offset+len])
}

func (p *Packet) AddLong(l uint64) {
	defer func() {
		p.length += 8
	}()
	p.payload = append(p.payload, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
}

func (p *Packet) AddInt(i uint32) {
	defer func() {
		p.length += 4
	}()
	p.payload = append(p.payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

func (p *Packet) AddShort(s uint16) {
	defer func() {
		p.length += 2
	}()
	p.payload = append(p.payload, byte(s>>8), byte(s))
}

func (p *Packet) AddByte(b uint8) {
	defer func() {
		p.length++
	}()
	p.payload = append(p.payload, b)
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
	dataLen := len(p.payload) + 1 // opcode
	p.payload = append([]byte{byte((dataLen >> 8) & 0xFF), byte(dataLen & 0xFF), p.opcode}, p.payload...)
}

func (c *Client) WritePacket(p *Packet) {
	if !p.bare {
		p.prependHeader()
	}

	c.Write(p.payload)
}
