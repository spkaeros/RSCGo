package server

import (
	"fmt"
	"github.com/gobwas/ws/wsutil"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/errors"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packethandlers"
	"github.com/spkaeros/rscgo/pkg/server/players"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

//Client Represents a single connecting client.
type Client struct {
	player          *world.Player
	IncomingPackets chan *packet.Packet
	CacheBuffer     []byte
	Socket          net.Conn
	DataBuffer      []byte
	DataLock        sync.RWMutex
	destroyer       sync.Once
}

//startReader Starts the client Socket reader goroutine.  Takes a waitgroup as an argument to facilitate synchronous destruction.
func (c *Client) startReader() {
	defer c.player.Destroy()
	for {
		select {
		default:
			p, err := c.readPacket()
			if err != nil {
				if err, ok := err.(errors.NetError); ok && err.Error() != "Connection closed." && err.Error() != "Connection timed out." {
					if err.Error() != "SHORT_DATA" {
						log.Warning.Printf("Rejected Packet from: %s\n", c)
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
			c.IncomingPackets <- p
		case <-c.player.Kill:
			return
		}
	}
}

//startWriter Starts the client Socket writer goroutine.
func (c *Client) startWriter() {
	defer c.player.Destroy()
	for {
		select {
		case p := <-c.player.OutgoingPackets:
			if p == nil {
				return
			}
			c.writePacket(*p)
		case <-c.player.Kill:
			return
		}
	}
}

//destroy Safely tears down a client, saves it to the database, and removes it from server-wide players.
func (c *Client) destroy(wg *sync.WaitGroup) {
	// Wait for network goroutines to finish.
	c.destroyer.Do(func() {
		(*wg).Wait()
		c.player.TransAttrs.UnsetVar("connected")
		close(c.player.OutgoingPackets)
		close(c.player.OptionMenuC)
		close(c.IncomingPackets)
		if err := c.Socket.Close(); err != nil {
			log.Error.Println("Couldn't close Socket:", err)
		}
		if _, ok := players.FromUserHash(c.player.UserBase37); ok {
			// Always try to launch I/O-heavy functions in their own goroutine.
			// Goroutines are light-weight and made for this kind of thing.
			go db.SavePlayer(c.player)
			world.RemovePlayer(c.player)
			c.player.TransAttrs.SetVar("remove", true)
			players.BroadcastLogin(c.player, false)
			players.Remove(c.player)
			log.Info.Printf("Unregistered: %v\n", c.player.String())
		}
	})
}

//startNetworking Starts up 3 new goroutines; one for reading incoming data from the Socket, one for writing outgoing data to the Socket, and one for client state updates and parsing plus handling incoming packetbuilders.  When the client kill signal is sent through the kill channel, the state update and packet handling goroutine will wait for both the reader and writer goroutines to complete their operations before unregistering the client.
func (c *Client) startNetworking() {
	var nwg sync.WaitGroup
	nwg.Add(2)
	go func() {
		defer nwg.Done()
		c.startReader()
	}()
	go func() {
		defer nwg.Done()
		c.startWriter()
	}()
	go func() {
		defer c.destroy(&nwg)
		for {
			select {
			case p := <-c.IncomingPackets:
				if p == nil {
					return
				}
				c.handlePacket(p)
			case <-c.player.Kill:
				return
			}
		}
	}()
}

//handlePacket Finds the mapped handler function for the specified packet, and calls it with the specified parameters.
func (c *Client) handlePacket(p *packet.Packet) {
	handler := packethandlers.Get(p.Opcode)
	if handler == nil {
		log.Info.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, len(p.Payload))
		fmt.Printf("CONTENT: %v\n", p.Payload)
		return
	}

	handler(c.player, p)
}

//newClient Creates a new instance of a Client, launches goroutines to handle I/O for it, and returns a reference to it.
func newClient(socket net.Conn, ws bool) *Client {
	c := &Client{Socket: socket, IncomingPackets: make(chan *packet.Packet, 20), DataBuffer: make([]byte, 5000)}
	c.player = world.NewPlayer(players.NextIndex(), strings.Split(socket.RemoteAddr().String(), ":")[0])
	c.player.Websocket = ws
	c.startNetworking()
	return c
}

//Write Writes data to the client's Socket from `b`.  Returns the length of the written bytes.
func (c *Client) Write(src []byte) int {
	var err error
	var dataLen int
	if c.player.Websocket {
		err = wsutil.WriteServerBinary(c.Socket, src)
		dataLen = len(src)
	} else {
		dataLen, err = c.Socket.Write(src)
	}
	if err != nil {
		log.Error.Println("Problem writing to websocket client:", err)
		c.player.Destroy()
		return -1
	}
	return dataLen
}

//Read Reads data off of the client's Socket into 'dst'.  Returns length read into dst upon success.  Otherwise, returns -1 with a meaningful error message.
func (c *Client) Read(dst []byte) (int, error) {
	// Set the read deadline for the socket to 10 seconds from now.
	err := c.Socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return -1, errors.ConnDeadline
	}

	expectedLen := len(dst)
	// Unstash any overflow data from previous read calls.
	cacheLen := len(c.CacheBuffer)
	if cacheLen > 0 {
		copy(dst, c.CacheBuffer)
		if cacheLen > expectedLen {
			c.CacheBuffer = c.CacheBuffer[expectedLen:]
			return expectedLen, nil
		} else {
			c.CacheBuffer = []byte{}
			if cacheLen == expectedLen {
				return expectedLen, nil
			}
		}
	}

	// Mark length of data left to read from socket after unstashing anything from the buffer
	reqDataLen := expectedLen - cacheLen

	var dataLen int
	var data []byte
	if !c.player.Websocket {
		dataLen, err = c.Socket.Read(dst[cacheLen:])
	} else {
		data, err = wsutil.ReadClientBinary(c.Socket)
		dataLen = len(data)
	}
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return -1, errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return -1, errors.ConnTimedOut
		}
		return -1, err
	}
	if c.player.Websocket {
		copy(dst[cacheLen:], data)
	}

	if dataLen < reqDataLen {
		// We didn't have enough data.  In practice, this produces an error I believe, but just in case!
		c.CacheBuffer = dst[:dataLen+cacheLen]
	} else if dataLen > reqDataLen {
		// We read too much data.  Stash what is not required.
		if c.player.Websocket {
			// Cache the recv'd data starting right after the last needed byte, next Read will unstash as if it were new data
			c.CacheBuffer = data[reqDataLen:]
		} else {
			// I don't think this can happen with TCP sockets.  We have finer control over what we read with them.
			// Just in case, I'll handle it in a semantically correct way, but I doubt it will ever run.
			c.CacheBuffer = dst[cacheLen+dataLen:]
		}
	}
	return dataLen + cacheLen, nil
}

//readPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *Client) readPacket() (*packet.Packet, error) {
	header := make([]byte, 2)
	if l, err := c.Read(header); err != nil {
		return nil, err
	} else if l < 2 {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}
	length := int(header[0])
	bigLength := length >= 160
	if bigLength {
		// length = (length-160)*256 + int(header[1])
		length = (length-160)<<8 + int(header[1])
	} else {
		// We have the final byte of frame data already, stored at header[1]
		length--
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
func (c *Client) writePacket(p packet.Packet) {
	if p.Bare {
		c.Write(p.Payload)
		return
	}
	frameLength := len(p.Payload)
	c.DataLock.Lock()
	header := c.DataBuffer[0:2]
	defer c.DataLock.Unlock()
	if frameLength >= 160 {
		//		header[0] = byte(frameLength/256+160)
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
