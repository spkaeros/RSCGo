package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

func init() {
	PacketHandlers["pingreq"] = func(c *Client, p *packets.Packet) {
		c.outgoingPackets <- packets.ResponsePong
	}
}

//Write Writes data to the client's socket from `b`.  Returns the length of the written bytes.
func (c *Client) Write(b []byte) int {
	l, err := c.socket.Write(b)
	if err != nil {
		log.Error.Println("Could not write to client socket.", err)
		c.Destroy()
	} else if l != len(b) {
		// Possibly non-fatal?
		log.Error.Printf("Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
	return l
}

//Read Reads data off of the client's socket into 'dst'.  Returns length read into dst upon success.  Otherwise, returns -1 with a meaningful error message.
func (c *Client) Read(dst []byte) (int, error) {
	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		// This shouldn't happen
		return -1, errors.ConnDeadline
	}
	length, err := c.socket.Read(dst)
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return -1, errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return -1, errors.ConnTimedOut
		}
	} else if length != len(dst) {
		return length, errors.NewNetworkError(fmt.Sprintf("Client.Read: unexpected length.  Expected %d, got %d.\n", len(dst), length))
	}

	return length, nil
}

//ReadPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *Client) ReadPacket() (*packets.Packet, error) {
	// TODO: Is allocation overhead more expensive than mutex locks?  If so, I must change this back to pre-allocated, and guard it with a RWMutex
	header := make([]byte, 2)
	if l, err := c.Read(header); err != nil || l != 2 {
		// This could happen legitimately, under certain strange circumstances.  Not proof of malicious intent.
		return nil, err
	}
	length := int(header[0])
	bigLength := length >= 160
	if bigLength {
		length = (length-160)*256 + int(header[1])
	}

	if length >= 5000 || length < 0 {
		log.Suspicious.Printf("Invalid packet length from [%v]: %d\n", c, length)
		if len(Flags.Verbose) > 0 {
			log.Warning.Printf("Packet length out of bounds; got %d, expected between 4 and 5000\n", length+3)
		}
		return nil, errors.NewNetworkError("Packet length out of bounds; must be between 4 and 5000.")
	}

	if bigLength {
		payload := make([]byte, length)

		if l, err := c.Read(payload); err != nil || l != length {
			return nil, err
		}

		return packets.NewPacket(payload[0], payload[1:]), nil
	}
	payload := make([]byte, length-1)

	if l, err := c.Read(payload); err != nil || l != length-1 {
		return nil, err
	}
	payload = append(payload, header[1])

	return packets.NewPacket(payload[0], payload[1:]), nil
}

//WritePacket This is a method to send a packet to the client.  If this is a bare packet, the packet payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this packet.
func (c *Client) WritePacket(p packets.Packet) {
	var buf []byte
	if !p.Bare {
		frameLength := len(p.Payload)
		if frameLength >= 160 {
			buf = append(buf, byte(160+frameLength/256))
			buf = append(buf, byte(frameLength))
		} else {
			buf = append(buf, byte(frameLength))
			buf = append(buf, p.Payload[frameLength-1])
			p.Payload = p.Payload[:frameLength-1]
		}
	}
	buf = append(buf, p.Payload...)

	c.Write(buf)
}
