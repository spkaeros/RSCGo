package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type Packet struct {
	Opcode    byte
	Payload   []byte
	bare      bool
	readIndex int
}

func NewPacket(opcode byte, payload []byte) *Packet {
	return &Packet{opcode, payload, false, 0}
}

func NewOutgoingPacket(opcode byte) *Packet {
	buf := []byte{ 0xA5, 0xA5, opcode }
	return &Packet{opcode, buf, false, 0}
}

func NewBarePacket(src []byte) *Packet {
	return &Packet{0, src, true, 0}
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
	if p.readIndex+1 > len(p.Payload) {
		fmt.Println("WARNING: Trying to read into packet with empty buffer!")
		return 0
	}
	defer func() {
		p.readIndex++
	}()
	return p.Payload[p.readIndex] & 0xFF
}

func (p *Packet) ReadString() (val string) {
	for c := p.ReadByte(); c != 0xA; c = p.ReadByte() {
		val += string(c)
	}
	return
}

func (p *Packet) AddLong(l uint64) *Packet {
	p.Payload = append(p.Payload, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	return p
}

func (p *Packet) AddInt(i uint32) *Packet {
	p.Payload = append(p.Payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return p
}

func (p *Packet) AddShort(s uint16) *Packet {
	p.Payload = append(p.Payload, byte(s>>8), byte(s))
	return p
}

func (p *Packet) AddByte(b uint8) *Packet {
	p.Payload = append(p.Payload, b)
	return p
}

func (c *Client) ReadPacket() (*Packet, error) {
	header := c.buffer[:3]
	if err := c.Read(header); err != nil {
		return nil, err
	}
	length := int(int16(header[0])<<8 | int16(header[1]))
	opcode := header[2] & 0xFF

	payload := c.buffer[3:length+3]

	if err := c.Read(payload); err != nil {
		return nil, err
	}

	if opcode == 0 {
		// Login block encrypted with block cipher using shared secret, to send/recv credentials and stream cipher key
		buf, err := rsa.DecryptPKCS1v15(rand.Reader, RsaKey, payload)
		if err != nil {
			LogDebug(1, "WARNING: Could not decrypt RSA login block: `%v`\n", err.Error())
			c.sendLoginResponse(9)
			return nil, err
		}
		payload = buf
	}

	return NewPacket(opcode, payload), nil
}

func (c *Client) WritePacket(p *Packet) {
	if !p.bare {
		l := len(p.Payload) - 2
		p.Payload[0] = byte(l >> 8)
		p.Payload[1] = byte(l)
	}

	c.Write(p.Payload)
}
