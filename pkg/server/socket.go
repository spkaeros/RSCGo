package server

import (
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

func init() {
	PacketHandlers["pingreq"] = func(c *Client, p *packets.Packet) {
		c.outgoingPackets <- packets.ResponsePong
	}
}

//Write Writes data to the client's socket from the slice of bytes `b`
// Returns the length of the written bytes.
func (c *Client) Write(b []byte) int {
	l, err := c.socket.Write(b)
	if err != nil {
		LogError.Println("Could not write to client socket.", err)
		//c.kill <- struct{}{}
		close(c.kill)
	} else if l != len(b) {
		// Possibly non-fatal?
		LogError.Printf("Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
	return l
}

//Read Reads data off of the client's socket into the slice of bytes 'dst'
// Returns nil upon successful read.  Otherwise, returns a meaningful error message.
func (c *Client) Read(dst []byte) error {
	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		// This shouldn't happen
		return errors.ConnDeadline
	}
	length, err := c.socket.Read(dst)
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return errors.ConnTimedOut
		}
	} else if length != len(dst) {
		return errors.NewNetworkError("Client.Read: unexpected length.  Expected " + strconv.Itoa(len(dst)) + ", got " + strconv.Itoa(length) + ".")
	}

	return nil
}

//ReadPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *Client) ReadPacket() (*packets.Packet, error) {
	header := c.buffer[:3]
	if err := c.Read(header); err != nil {
		return nil, err
	}
	length := int(int16(header[0])<<8 | int16(header[1]))
	opcode := header[2]

	payload := c.buffer[3 : length+3]

	if err := c.Read(payload); err != nil {
		return nil, err
	}

	return packets.NewPacket(opcode, payload), nil
}

//WritePacket This is a method to send a packet to the client.  If this is a bare packet, the packet payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this packet.
func (c *Client) WritePacket(p *packets.Packet) {
	if !p.Bare {
		l := len(p.Payload) - 2
		p.Payload[0] = byte(l >> 8)
		p.Payload[1] = byte(l)
		//		if c.isaacStream != nil {
		//			p.Payload[2] ^= c.isaacStream.decoder.Uint8()
		//		}
	}

	c.Write(p.Payload)
}
