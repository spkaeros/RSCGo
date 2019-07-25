package server

import (
	rscrand "bitbucket.org/zlacki/rscgo/rand"
	"fmt"
	"net"
	"strings"
)

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
	send   chan *packet
}

func (c channel) write(b []byte) {
	l, err := c.socket.Write(b)
	if err != nil {
		fmt.Println("ERROR: Could not write to client socket.")
		fmt.Println(err)
	}
	if l != len(b) {
		fmt.Printf("WARNING: Wrong number of bytes written to client socket.  Expected %d, got %d.\n", len(b), l)
	}
}

func (c channel) sendPacket(p *packet) {
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

	for i := 0; i < dataLen; i++ {
		buf = append(buf, p.payload[i])
	}

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

var handlers = make(map[byte]func(*client, *packet))

func sessionIDRequest(c *client, p *packet) {
	c.uID = p.payload[0]
	p1 := newPacket(0, []byte{}, 0)
	p1.bare = true
	p1.addLong(rscrand.GetSecureRandomLong())
	c.send <- p1
}

func loginRequest(c *client, p *packet) {
	// TODO: RSA decryption, blabla.
	// Currently returns an invalid username or password response
	p1 := newPacket(0, []byte{}, 0)
	p1.bare = true
	p1.addByte(3)
	c.send <- p1
}

func init() {
	handlers[32] = sessionIDRequest
	handlers[0] = loginRequest
}

func (c *client) handlePacket(p *packet) {
	handler, ok := handlers[p.opcode]
	if !ok {
		fmt.Printf("Unhandled packet: {opcode: %d; length: %d};\n", p.opcode, p.length)
		return
	}
	handler(c, p)
}

func getIPFromConn(c net.Conn) string {
	parts := strings.Split(c.RemoteAddr().String(), ":")
	if len(parts) < 1 {
		return "nil"
	}
	return parts[0]
}