package server

import (
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
}

func NewPacket(opcode byte, payload []byte, length int) *Packet {
	return &Packet{opcode, payload, length, false}
}

type Channel struct {
	socket net.Conn
}

func (c Channel) Write(b []byte) {
	l, err := c.socket.Write(b)
	if err != nil {
		fmt.Println("ERROR: Could not Write to Client socket.")
		fmt.Println(err)
	}
	if l != len(b) {
		fmt.Printf("WARNING: Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
}

func (c Channel) NextPacket() (*Packet, *NetError) {
	headerBuffer := make([]byte, 3)

	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		fmt.Printf("Rejected Packet from: '%s'\n", getIPFromConn(c.socket))
		fmt.Println(err)
		return nil, Deadline()
	}
	headerLength, err := c.socket.Read(headerBuffer)

	if err == io.EOF {
		return nil, Closed()
	} else if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, Timeout()
	} else if err != nil {
		if strings.Contains(err.Error(), "use of closed") {
			return nil, &NetError{"Trying to read a closed socket.", true, true}
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			return nil, Closed()
		}
		fmt.Printf("Rejected Packet from: '%s'\n", getIPFromConn(c.socket))
		fmt.Println(err)
		return nil, &NetError{msg: "Unexpected I/O error encountered while reading Packet header."}
	} else if headerLength != 3 {
		fmt.Printf("Rejected Packet from: '%s'\n", getIPFromConn(c.socket))
		return nil, &NetError{msg: "Packet header unexpected length.  Expected 3 bytes, got " + strconv.Itoa(headerLength) + " bytes."}
	}

	length := int(headerBuffer[0] & 0xFF)
	if length >= 160 {
		length = (length-160)*256 + int(headerBuffer[1]&0xFF)
	} else {
		length--
	}

	opcode := headerBuffer[2] & 0xFF
	// Opcode is part of the length variable sent from Jagex Client.
	// IMO, opcode is a part of the header, so I read it into the header.
	// TODO: Check Jagex Client for any cases that would break this code, e.g bare opcode-free context-based data?
	length--

	payloadBuffer := make([]byte, length)

	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		fmt.Printf("Rejected Packet[opcode:%d, len:%d] from: '%s'\n", opcode, length, getIPFromConn(c.socket))
		fmt.Println(err)
		return nil, Deadline()
	}
	payloadLength, err := c.socket.Read(payloadBuffer)

	if err == io.EOF {
		return nil, Closed()
	} else if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, Timeout()
	} else if err != nil {
		if strings.Contains(err.Error(), "use of closed") {
			return nil, &NetError{"Trying to read a closed socket.", true, true}
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			return nil, Closed()
		}
		fmt.Printf("Rejected Packet[opcode:%d, len:%d] from: '%s'\n", opcode, length, getIPFromConn(c.socket))
		return nil, &NetError{msg: "Unexpected I/O error encountered while reading Packet header."}
	} else if payloadLength != length {
		fmt.Printf("Rejected Packet[opcode:%d, len:%d] from: '%s'\n", opcode, length, getIPFromConn(c.socket))
		return nil, &NetError{msg: "Packet frame unexpected length.  Expected " + strconv.Itoa(length) + " bytes, got " + strconv.Itoa(payloadLength) + " bytes."}
	}

	if length < 160 {
		payloadBuffer = append(payloadBuffer, headerBuffer[1])
		length++
	}

	return NewPacket(opcode, payloadBuffer, length), nil
}

func (c Channel) WritePacket(p *Packet) {
	buf := make([]byte, p.length)
	dataLen := len(p.payload)
	packetLen := dataLen + 1
	if !p.bare {
		if packetLen >= 160 {
			buf = append(buf, byte(160+packetLen/256), byte(packetLen&0xFF))
		} else {
			buf = append(buf, byte(packetLen&0xFF))
			if dataLen > 0 {
				dataLen--
				buf = append(buf, p.payload[dataLen])
			}
			buf = append(buf, p.opcode)
		}
	}

	buf = append(buf, p.payload[:dataLen]...)

	c.Write(buf)
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
