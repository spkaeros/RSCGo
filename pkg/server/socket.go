package server

import (
	"fmt"
	"io"
	"net"
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

//Write Writes data to the client's socket from `b`.  Returns the length of the written bytes.
func (c *Client) Write(b []byte) int {
	l, err := c.socket.Write(b)
	if err != nil {
		LogError.Println("Could not write to client socket.", err)
		c.Destroy()
	} else if l != len(b) {
		// Possibly non-fatal?
		LogError.Printf("Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
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
	// I'm using a pre-allocated buffer for incoming packet data, to avoid allocation overhead.
	header := c.buffer[:3]
	if l, err := c.Read(header); err != nil || l != 3 {
		// This could happen legitimately, under certain strange circumstances.  Not proof of malicious intent.
		return nil, err
	}
	// TODO: Should I fit this into a byte?  Maybe use bit-packing for a custom header protocol?
	length := int(int16(header[0])<<8 | int16(header[1]))
	opcode := header[2]

	if length+3 >= 5000 || length+3 < 3 {
		// This should only happen if someone is either editing their outgoing network data, or using a modified client.
		if len(Flags.Verbose) > 0 {
			LogWarning.Printf("Packet length out of bounds; got %d, expected between 4 and 5000\n", length+3)
		}
		return nil, errors.NewNetworkError("Packet length out of bounds; must be between 4 and 5000.")
	}

	if !c.player.Connected && opcode != 32 && opcode != 0 {
		// This should only happen if someone is either editing their outgoing network data, or using a modified client.
		if len(Flags.Verbose) > 0 {
			LogWarning.Printf("Unauthorized packet{opcode:%v,len:%v] rejected from: %v\n", opcode, length+3, c)
		}
		return nil, errors.NewNetworkError("Unauthorized packet received.")
	}

	payload := c.buffer[3 : length+3]

	if l, err := c.Read(payload); err != nil || l != length {
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
