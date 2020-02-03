package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/errors"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packet/handlers"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

//client Represents a single connecting client.
type client struct {
	player      *world.Player
	socket      net.Conn
	destroyer   sync.Once
	readWriter  *bufio.ReadWriter
	wsReader    io.Reader
	wsHeader    ws.Header
	wsLength    int64
}

//startNetworking Starts up 3 new goroutines; one for reading incoming data from the socket, one for writing outgoing data to the socket, and one for client state updates and parsing plus handling incoming world.  When the client kill signal is sent through the kill channel, the state update and packet handling goroutine will wait for both the reader and writer goroutines to complete their operations before unregistering the client.
func (c *client) startNetworking() {
	incomingPackets := make(chan *packet.Packet, 20)
	awaitDeath := sync.WaitGroup{}

	go func() {
		defer awaitDeath.Done()
		defer c.player.Destroy()
		awaitDeath.Add(1)
		for {
			select {
			case p := <-c.player.OutgoingPackets:
				if p == nil {
					return
				}
				c.writePacket(*p)
			case <-c.player.KillC:
				return
			}
		}
	}()
	go func() {
		defer awaitDeath.Done()
		defer c.player.Destroy()
		awaitDeath.Add(1)
		for {
			select {
			default:
				p, err := c.readPacket()
				if err != nil {
					if err, ok := err.(errors.NetError); ok && err.Error() != "Connection closed." && err.Error() != "Connection timed out." {
						if err.Error() != "SHORT_DATA" {
							log.Warning.Printf("Rejected Packet from: %s\n", c.player.String())
							log.Warning.Println(err)
						}
						continue
					}
					return
				}
				if !c.player.Connected() && p.Opcode != 32 && p.Opcode != 0 && p.Opcode != 2 && p.Opcode != 220 {
					log.Warning.Printf("Unauthorized packet[opcode:%v,len:%v] rejected from: %v\n", p.Opcode, len(p.Payload), c)
					return
				}
				incomingPackets <- p
			case <-c.player.KillC:
				return
			}
		}
	}()
	go func() {
		defer c.destroy()
		defer close(incomingPackets)
		defer awaitDeath.Wait()
		defer c.player.Destroy()
		for {
			select {
			case p := <-incomingPackets:
				if p == nil {
					log.Warning.Println("Tried processing nil packet!")
					continue
				}
				c.handlePacket(p)
			case <-c.player.KillC:
				return
			}
		}
	}()
}

//destroy Safely tears down a client, saves it to the database, and removes it from server-wide player list.
func (c *client) destroy() {
	c.destroyer.Do(func() {
		close(c.player.OutgoingPackets)
		if err := c.socket.Close(); err != nil {
			log.Error.Println("Couldn't close socket:", err)
		}
		c.player.SetConnected(false)
		c.player.SetRegionRemoved()
		if player, ok := world.Players.FromIndex(c.player.Index); !ok || player != c.player {
			log.Warning.Println("Destroying Player did not match player that is assigned index in map!")
			return
		}
		go db.SavePlayer(c.player)
		log.Info.Printf("Unregistered: %v\n", c.player.String())
		world.RemovePlayer(c.player)
		world.Players.BroadcastLogin(c.player, false)
		world.Players.Remove(c.player)
	})
}

//handlePacket Finds the mapped handler function for the specified packet, and calls it with the specified parameters.
func (c *client) handlePacket(p *packet.Packet) {
	handler := handlers.Handler(p.Opcode)
	if handler == nil {
		log.Info.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, len(p.Payload))
		fmt.Printf("CONTENT: %v\n", p.Payload)
		return
	}

	handler(c.player, p)
}

//newClient Creates a new instance of a client, launches goroutines to handle I/O for it, and returns a reference to it.
func newClient(socket net.Conn, ws2 bool) *client {
	c := &client{socket: socket}
	c.player = world.NewPlayer(world.Players.NextIndex(), strings.Split(socket.RemoteAddr().String(), ":")[0])
	c.readWriter = bufio.NewReadWriter(bufio.NewReader(socket), bufio.NewWriter(socket))
	if ws2 {
		c.wsHeader, c.wsReader, _ = wsutil.NextReader(socket, ws.StateServerSide)
	}
	c.startNetworking()
	return c
}

//Write Writes data to the client's socket from `b`.  Returns the length of the written bytes.
func (c *client) Write(src []byte) int {
	var err error
	var dataLen int
	if c.wsReader != nil {
		err = wsutil.WriteServerBinary(c.socket, src)
		dataLen = len(src)
	} else {
		dataLen, err = c.socket.Write(src)
	}
	if err != nil {
		log.Error.Println("Problem writing to websocket client:", err)
		c.player.Destroy()
		return -1
	}
	return dataLen
}

//Read Reads data off of the client's socket into 'dst'.  Returns length read into dst upon success.  Otherwise, returns -1 with a meaningful error message.
func (c *client) Read(dst []byte) (int, error) {
	// set the read deadline for the socket to 10 seconds from now.
	err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return -1, errors.ConnDeadline
	}

	if c.wsReader != nil {
		if c.wsHeader.Length <= c.wsLength {
			c.wsHeader, c.wsReader, err = wsutil.NextReader(c.readWriter.Reader, ws.StateServerSide)
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
					return -1, errors.ConnClosed
				} else if e, ok := err.(net.Error); ok && e.Timeout() {
					return -1, errors.ConnTimedOut
				} else {
					log.Warning.Println("Problem creating reader for next websocket frame:", err)
				}
				
				c.player.Destroy()
				return -1, err
			}
			// reset current read index
			c.wsLength = 0
		}
		n, err := c.wsReader.Read(dst)
		if err != nil {
			if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
				return -1, errors.ConnClosed
			} else if e, ok := err.(net.Error); ok && e.Timeout() {
				return -1, errors.ConnTimedOut
			} else if err == io.EOF {
				if !c.wsHeader.Fin {
					return -1, errors.ConnClosed
				}
				// EOF on fin means end of frame
				c.wsLength += int64(n)
				return n, nil
			} else {
				log.Warning.Println(err)
			}
			return -1, err
		}
		c.wsLength += int64(n)
		return n, nil
	}

	n, err := c.readWriter.Read(dst)
	if err != nil {
		log.Info.Println(err)
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return -1, errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return -1, errors.ConnTimedOut
		}
		return -1, err
	}
	return n, nil
}

//readPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *client) readPacket() (p *packet.Packet, err error) {
	header := make([]byte, 2)
	l, err := c.Read(header)
	if err != nil {
		return nil, err
	}
	if l < 2 {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}
	length := int(header[0]) - 1
	bigLength := length >= 160
	if bigLength {
		length = (length-160)<<8 + int(header[1])
	}

	if length+2 >= 5000 || length+2 < 2 {
		log.Suspicious.Printf("Invalid packet length from [%v]: %d\n", c, length)
		log.Warning.Printf("Packet from [%v] length out of bounds; got %d, expected between 0 and 5000\n", c, length)
		return nil, errors.NewNetworkError("Packet length out of bounds; must be between 0 and 5000.")
	}

	payload := make([]byte, length)

	if length > 0 {
		if l, err := c.Read(payload); err != nil {
			return nil, err
		} else if l < length {
			return nil, errors.NewNetworkError("SHORT_DATA")
		}
	}

	if !bigLength {
		// If the length in the packet header used 1 byte, the 2nd byte in the header is the final byte of frame data
		payload = append(payload, header[1])
	}

	return packet.NewPacket(payload[0], payload[1:]), nil
}

//writePacket This is a method to send a packet to the client.  If this is a bare packet, the packet payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this packet.
func (c *client) writePacket(p packet.Packet) {
	if p.Bare {
		c.Write(p.Payload)
		return
	}
	frameLength := len(p.Payload)
	header := make([]byte, 2)
	if frameLength >= 160 {
		header[0] = byte(frameLength>>8 + 160)
		header[1] = byte(frameLength)
	} else {
		header[0] = byte(frameLength)
		header[1] = p.Payload[frameLength-1]
		p.Payload = p.Payload[:frameLength-1]
	}
	c.Write(append(header, p.Payload...))
	return
}
