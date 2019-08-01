package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type NetError struct {
	msg    string
	ping   bool
	closed bool
}

func (e *NetError) Error() string {
	return e.msg
}

func Closed() *NetError {
	return &NetError{msg: "Connection reset by peer.", closed: true}
}

func Timeout() *NetError {
	return &NetError{msg: "Connection timed out.", ping: true}
}

func Deadline() *NetError {
	return &NetError{msg: "Could not set read deadline for Client listener.", closed: true}
}

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
	buf, err :=rsa.DecryptPKCS1v15(rand.Reader, RsaKey, p.payload)
	if err != nil {
		fmt.Println("WARNING: Could not decrypt RSA login block")
		return false
	}
	p.payload = buf
	p.length = len(buf)
	return true
}

func (c *Client) Write(b []byte) {
	l, err := c.socket.Write(b)
	if err != nil {
		fmt.Println("ERROR: Could not Write to Client socket.")
		fmt.Println(err)
	}
	if l != len(b) {
		fmt.Printf("WARNING: Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
}

func (c *Client) Read(len int) ([]byte, error) {
	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		// This shouldn't happen
		return nil, Deadline()
	}
	buf := make([]byte, len)
	length, err := c.socket.Read(buf)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil, Timeout()
		}
		if strings.Contains(err.Error(), "use of closed") {
			return nil, &NetError{msg: "Trying to read a closed socket.", closed: true}
		}
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") {
			return nil, Closed()
		}
	}
	if length != len {
		return nil, &NetError{msg: "Client.Read: unexpected length.  Expected " + strconv.Itoa(len) + ", got " + strconv.Itoa(length) + "."}
	}

	return buf, nil
}

func (c *Client) NextPacket() (*Packet, error) {
	headerBuffer, err := c.Read(3)
	if err != nil {
		return nil, err
	}

	length := int(headerBuffer[0] & 0xFF)
	if length >= 160 {
		length = (length-160)*256 + int(headerBuffer[1]&0xFF)
	} else {
		// TODO: Should it be <= 160, and should it be >= 1?
		// If the payload length is less than 160 bytes, the 2nd byte in the header is used to store the last byte
		//  of payload data.  Subtract one from length so that we don't try to read it from the end of the payload.
		length--
	}

	// Opcode byte is included in the length variable, but we read it into the header buffer since it should be there.
	opcode := headerBuffer[2] & 0xFF
	length--

	payloadBuffer, err := c.Read(length)
	if err != nil {
		return nil, err
	}

	if length < 160 {
		payloadBuffer = append(payloadBuffer, headerBuffer[1])
		length++
	}

	return NewPacket(opcode, payloadBuffer, length), nil
}

func (c *Client) WritePacket(p *Packet) {
	dataLen := len(p.payload)
	packetLen := dataLen + 1
	buf := make([]byte, 0)
	if !p.bare {
		if packetLen >= 160 {
			buf = append(buf, byte(160+packetLen/256), byte(packetLen&0xFF))
		} else {
			buf = append(buf, byte(packetLen&0xFF))
			if dataLen > 0 {
				dataLen--
				buf = append(buf, p.payload[dataLen])
			}
		}
		buf = append(buf, p.opcode&0xFF)
	}
	buf = append(buf, p.payload[:dataLen]...)

	c.Write(buf)
}

func (p *Packet) ReadLong() int64 {
	l := int64(p.payload[p.offset] >> 56)
	l |= int64(p.payload[p.offset + 1] >> 48)
	l |= int64(p.payload[p.offset + 2] >> 40)
	l |= int64(p.payload[p.offset + 3] >> 32)
	l |= int64(p.payload[p.offset + 4] >> 24)
	l |= int64(p.payload[p.offset + 5] >> 16)
	l |= int64(p.payload[p.offset + 6] >> 8)
	l |= int64(p.payload[p.offset + 7])
	p.offset += 8
	return l
}
func (p *Packet) ReadInt() int32 {
	i := int32(p.payload[p.offset] >> 24)
	i |= int32(p.payload[p.offset + 1] >> 16)
	i |= int32(p.payload[p.offset + 2] >> 8)
	i |= int32(p.payload[p.offset + 3])
	p.offset += 4
	return i
}
func (p *Packet) ReadShort() int16 {
	i := int16(p.payload[p.offset] >> 8)
	i |= int16(p.payload[p.offset + 1])
	p.offset += 2
	return i
}
func (p *Packet) ReadByte() byte {
	p.offset++
	return p.payload[p.offset]
}

func (p *Packet) ReadString(len int) string {
	if p.offset + len > p.length {
		fmt.Printf("WARNING: Requested string length too long.  Requested %d, only %d left in buffer.", len, p.length-p.offset)
		return ""
	}
	defer func() {
		p.offset += len
	}()
	return string(p.payload[p.offset:p.offset+len])
}

func (p *Packet) AddLong(l uint64) {
	p.payload = append(p.payload, byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	p.length += 8
}

func (p *Packet) AddInt(i uint32) {
	p.payload = append(p.payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	p.length += 4
}

func (p *Packet) AddShort(s uint16) {
	p.payload = append(p.payload, byte(s>>8), byte(s))
	p.length += 2
}

func (p *Packet) AddByte(b uint8) {
	p.payload = append(p.payload, b)
	p.length++
}

func getIPFromConn(c net.Conn) string {
	parts := strings.Split(c.RemoteAddr().String(), ":")
	if len(parts) < 1 {
		return "nil"
	}
	return parts[0]
}
