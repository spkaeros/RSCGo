package client

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"net"
	"strings"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
)
/*
func init() {
	PacketHandlers["pingreq"] = func(c *Client, p *packetbuilders.Packet) {
		c.OutgoingPackets <- packetbuilders.ResponsePong
	}
}
*/
//Write Writes data to the client's Socket from `b`.  Returns the length of the written bytes.
func (c *Client) Write(b []byte) int {
	if !c.player.Websocket {
		l, err := c.Socket.Write(b)
		if err != nil {
			log.Error.Println("Could not write to client Socket.", err)
			c.Destroy()
		} else if l != len(b) {
			// Possibly non-fatal?
			log.Error.Printf("Wrong number of bytes written to Client Socket.  Expected %d, got %d.\n", len(b), l)
		}
		return l
	}
	w := wsutil.NewWriter(c.Socket, ws.StateServerSide, ws.OpBinary)
	l, err := w.Write(b)
	if err != nil {
		log.Error.Println("Could not write to client Socket.", err)
		c.Destroy()
	} else if l != len(b) {
		// Possibly non-fatal?
		log.Error.Printf("Wrong number of bytes written to Client Socket.  Expected %d, got %d.\n", len(b), l)
	}
	if err := w.Flush(); err != nil {
		log.Warning.Println("Error writing to Websocket:", err)
	}
	return l
}

//Read Reads data off of the client's Socket into 'dst'.  Returns length read into dst upon success.  Otherwise, returns -1 with a meaningful error message.
func (c *Client) Read(dst []byte) (int, error) {
	err := c.Socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return -1, errors.ConnDeadline
	}

	if len(c.PacketData) >= len(dst) {
		// If we have enough data to fill dst, fill it, stash the remaining leftovers
		copy(dst, c.PacketData)
		if len(c.PacketData) == len(dst) {
			c.PacketData = c.PacketData[:0]
		} else {
			c.PacketData = c.PacketData[len(dst):]
		}
		return len(dst), nil
	}

	var data []byte
	if c.player.Websocket {
		data, _, err = wsutil.ReadData(c.Socket, ws.StateServerSide)
	} else {
		_, err = c.Socket.Read(dst)
		data = dst
	}

	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return -1, errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return -1, errors.ConnTimedOut
		}
		return -1, err
	}

	if len(c.PacketData) > 0 {
		// unstash extra data
		if c.player.Websocket {
			data = append(c.PacketData, data...)
		} else {
			dst = append(c.PacketData, dst...)
		}
		c.PacketData = c.PacketData[:0]
	}

	if c.player.Websocket {
		copy(dst, data)

		if len(data) > len(dst) {
			// stash extra data
			c.PacketData = data[len(dst):]
		}
	}
	return len(dst), nil
}

//ReadPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *Client) ReadPacket() (*packetbuilders.Packet, error) {
	// TODO: Is allocation overhead more expensive than mutex locks?  If so, I must change this back to pre-allocated, and guard it with a RWMutex
	header := make([]byte, 2)
	if l, err := c.Read(header); err != nil {
		return nil, err
	} else if l < 2 {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}
	length := int(header[0])
	if length >= 160 {
		length = (length-160)*256 + int(header[1])
	} else {
		length--
	}

	if length >= 5000 || length < 0 {
		log.Suspicious.Printf("Invalid packet length from [%v]: %d\n", c, length)
		log.Warning.Printf("Packet length out of bounds; got %d, expected between 4 and 5000\n", length+3)
		return nil, errors.NewNetworkError("Packet length out of bounds; must be between 4 and 5000.")
	}

	payload := make([]byte, length)

	if l, err := c.Read(payload); err != nil {
		return nil, err
	} else if l < length {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}

	if length < 160 {
		payload = append(payload, header[1])
	}

	return packetbuilders.NewPacket(payload[0], payload[1:]), nil
}

//WritePacket This is a method to send a packet to the client.  If this is a bare packet, the packet payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this packet.
func (c *Client) WritePacket(p packetbuilders.Packet) {
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
