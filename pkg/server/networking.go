package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type netError struct {
	msg string
	ping bool
	closed bool
}

func (e *netError) Error() string {
	return e.msg
}

func connectionClosed() *netError {
	return &netError{msg: "Connection reset by peer.", closed: true}
}

func timedOut() *netError {
	return &netError{msg: "Connection timed out.", ping: true}
}

func deadlineError() *netError {
	return &netError{msg: "Could not set read deadline for Client listener.", closed: true}
}

type packet struct {
	opcode  byte
	payload []byte
	length  int
	bare    bool
}

func newPacket(opcode byte, payload []byte, length int) *packet {
	return &packet{opcode, payload, length, false}
}

type channel struct {
	socket net.Conn
}

func (c channel) write(b []byte) {
	l, err := c.socket.Write(b)
	if err != nil {
		fmt.Println("ERROR: Could not write to Client socket.")
		fmt.Println(err)
	}
	if l != len(b) {
		fmt.Printf("WARNING: Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
}

func (c channel) readPacket() (*packet, *netError) {
	headerBuffer := make([]byte, 3)

	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		fmt.Printf("Rejected packet from: '%s'\n", getIPFromConn(c.socket))
		fmt.Println(err)
		return nil, deadlineError()
	}
	headerLength, err := c.socket.Read(headerBuffer)

	if err == io.EOF {
		return nil, connectionClosed()
	} else if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, timedOut()
	} else if err != nil {
		if strings.Contains(err.Error(), "use of closed") {
			return nil, &netError{"Trying to read a closed socket.", true, true}
		}
		fmt.Printf("Rejected packet from: '%s'\n", getIPFromConn(c.socket))
		fmt.Println(err)
		return nil, &netError{msg: "Unexpected I/O error encountered while reading packet header."}
	} else if headerLength != 3 {
		fmt.Printf("Rejected packet from: '%s'\n", getIPFromConn(c.socket))
		return nil, &netError{msg: "Packet header unexpected length.  Expected 3 bytes, got " + strconv.Itoa(headerLength) + " bytes."}
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
	// TODO: Check Jagex Client for any cases that would break this code, e.g raw opcode-free context-based packets?
	length--

	payloadBuffer := make([]byte, length)

	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		fmt.Printf("Rejected packet[opcode: %d, len:%d] from: '%s'\n", opcode, length, getIPFromConn(c.socket))
		fmt.Println(err)
		return nil, deadlineError()
	}
	payloadLength, err := c.socket.Read(payloadBuffer)

	if err == io.EOF {
		return nil, connectionClosed()
	} else if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, timedOut()
	} else if err != nil {
		fmt.Printf("Rejected packet[opcode: %d, len:%d] from: '%s'\n", opcode, length, getIPFromConn(c.socket))
		return nil, &netError{msg: "Unexpected I/O error encountered while reading packet header."}
	} else if payloadLength != length {
		fmt.Printf("Rejected packet[opcode: %d, len:%d] from: '%s'\n", opcode, length, getIPFromConn(c.socket))
		return nil, &netError{msg: "Packet frame unexpected length.  Expected " + strconv.Itoa(length) + " bytes, got " + strconv.Itoa(payloadLength) + " bytes."}
	}

	if length < 160 {
		payloadBuffer = append(payloadBuffer, headerBuffer[1])
		length++
	}

	return newPacket(opcode, payloadBuffer, length), nil
}

func (c channel) writePacket(p *packet) {
	buf := make([]byte, 0)
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

	c.write(buf)
}

func (p *packet) addLong(l uint64) {
	p.payload = append(p.payload, byte(l >> 56), byte(l >> 48), byte(l >> 40), byte(l >> 32), byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l))
	p.length += 8
}

func (p *packet) addInt(i uint32) {
	p.payload = append(p.payload, byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i))
	p.length += 4
}

func (p *packet) addShort(s uint16) {
	p.payload = append(p.payload, byte(s >> 8), byte(s))
	p.length += 2
}

func (p *packet) addByte(b uint8) {
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
