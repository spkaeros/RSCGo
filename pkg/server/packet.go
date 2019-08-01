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
	return int64(p.ReadByte()<<56) | int64(p.ReadByte()<<48) | int64(p.ReadByte()<<40) | int64(p.ReadByte()<<32) |
		int64(p.ReadByte()<<24) | int64(p.ReadByte()<<16) | int64(p.ReadByte()<<8) | int64(p.ReadByte())
}
func (p *Packet) ReadInt() (val int32) {
	return int32(p.ReadByte()<<24) | int32(p.ReadByte()<<16) | int32(p.ReadByte()<<8) | int32(p.ReadByte())
}
func (p *Packet) ReadShort() (val int16) {
	return int16(p.ReadByte()<<8) | int16(p.ReadByte())
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
	length := (int(buf[0]&0xFF) << 8) | int(buf[1]&0xFF)
	opcode := buf[2] & 0xFF

	payloadBuffer, err := c.Read(length)
	if err != nil {
		return nil, err
	}

	return NewPacket(opcode, payloadBuffer, length), nil
}

func (p *Packet) prependHeader() {
	dataLen := p.length + 1 // opcode
	p.payload = append([]byte{byte(dataLen>>8) & 0xFF, byte(dataLen & 0xFF), p.opcode}, p.payload...)
}

func (c *Client) WritePacket(p *Packet) {
	if !p.bare {
		p.prependHeader()
	}

	fmt.Println(p.payload)
	c.Write(p.payload)
}
