/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package engine

import (
	"bufio"
	"io"
	stdnet "net"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/net/handlers"
	"github.com/spkaeros/rscgo/pkg/game/net/handshake"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

//client Represents a single connecting client.
type client struct {
	player     *world.Player
	socket     stdnet.Conn
	destroyer  sync.Once
	readWriter *bufio.ReadWriter
	websocket  bool
	readSize  int
	readLimit  int
	frameFin bool
}

//startNetworking Starts up 3 new goroutines; one for reading incoming data from the socket, one for writing outgoing data to the socket, and one for client state updates and parsing plus handling incoming world.  When the client kill signal is sent through the kill channel, the state update and net handling goroutine will wait for both the reader and writer goroutines to complete their operations before unregistering the client.
func (c *client) startNetworking() {
	incomingPackets := make(chan *net.Packet, 20)
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
		read := func(c *client, data []byte) int {
			written := 0; 
			for written < len(data) {
				err := c.socket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(15)))
				if err != nil {
					return -1
				}
				if c.websocket && (c.readSize >= c.readLimit) {
					// reset buffer read index and create the next reader
					header, reader, err := wsutil.NextReader(c.socket, ws.StateServerSide)
					c.readLimit = int(header.Length)
					c.readSize = 0
					c.frameFin = header.Fin
					if err != nil {
						if err == io.EOF || err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
							return -1
						} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
							return -1
						}
						log.Warning.Println("Problem creating reader for next websocket frame:", err)
					}
					c.readWriter.Reader = bufio.NewReader(reader)
				}
				n, err := c.readWriter.Read(data[written:])
				c.readSize += n
				if err != nil {
					if err == io.EOF {
						if !c.frameFin {
//							log.Info.Println(err)
							continue
						}
					}
					if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
						return -1
					} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
						return -1
					}
					continue
				}
				written += n
			}
//			log.Info.Printf("data[%d]:%v\n", written, data)
			return written
		}
		awaitDeath.Add(1)
		for {
			select {
			default:
				header := make([]byte, 2)
				if read(c, header) < 2 {
					continue
				}
				frameSize := int(header[0] & 0xFF)
				if frameSize >= 160 {
					frameSize = ((frameSize-160)<<8) | int(header[1] & 0xFF)
				} else {
					frameSize -= 1
				}
			
				// Upper bound is an approximation of the max size of the clientside outgoing data buffer
				if frameSize >= 23768 || frameSize < 0 {
					log.Suspicious.Printf("Invalid packet length from [%v]: %d\n", c, frameSize)
					continue
				}
				localData := make([]byte, frameSize)
				if frameSize > 0 {
					if read(c, localData) == -1 {
						continue
					}
				}
				if frameSize < 160 {
					localData = append(localData, header[1])
				}
				if !c.player.Connected() && !handshake.EarlyOperation(int(localData[0])) {
					log.Warning.Printf("Unauthorized packet[opcode:%v,size:%v (expected:%v)] rejected from: %v\n", localData[0], len(localData), frameSize, c)
					continue
				}
				incomingPackets <- net.NewPacket(localData[0], localData[1:])

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

//destroy Safely tears down a client, saves it to the database, and removes it from game-wide player list.
func (c *client) destroy() {
	c.destroyer.Do(func() {
		go func() {
			c.player.UpdateWG.RLock()
			c.player.SetConnected(false)
			if err := c.socket.Close(); err != nil {
				log.Error.Println("Couldn't close socket:", err)
			}
			close(c.player.OutgoingPackets)
			if player, ok := world.Players.FromIndex(c.player.Index); c.player.Index == -1 || (ok && player != c.player) || !ok {
				log.Warning.Printf("Unregistered: Unauthenticated connection ('%v'@'%v')\n", c.player.Username(), c.player.CurrentIP())
				if ok {
					log.Suspicious.Printf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", c.player.Index, player)
				}
				return
			}
			c.player.Attributes.SetVar("lastIP", c.player.CurrentIP())
			world.RemovePlayer(c.player)
			db.DefaultPlayerService.PlayerSave(c.player)
			log.Info.Printf("Unregistered: %v\n", c.player.String())
			c.player.UpdateWG.RUnlock()
		}()
	})
}

//handlePacket Finds the mapped handler function for the specified net, and calls it with the specified parameters.
func (c *client) handlePacket(p *net.Packet) {
	handler := handlers.Handler(p.Opcode)
	if handler == nil {
		log.Info.Printf("Unhandled Packet: {opcode:%d; length:%d; payload:%v};\n", p.Opcode, len(p.FrameBuffer), p.FrameBuffer)
		return
	}

	handler(c.player, p)
}

//newClient Creates a new instance of a client, launches goroutines to handle I/O for it, and returns a reference to it.
func newClient(socket stdnet.Conn, ws2 bool) *client {
	c := &client{socket: socket}
	c.player = world.NewPlayer(-1, strings.Split(socket.RemoteAddr().String(), ":")[0])
	c.websocket = ws2
	c.readWriter = bufio.NewReadWriter(bufio.NewReader(socket), bufio.NewWriter(socket))
	c.startNetworking()
	return c
}

//Write Writes data to the client's socket from `b`.  Returns the length of the written bytes.
func (c *client) Write(src []byte) int {
	var err error
	var dataLen int
	if c.websocket {
		err = wsutil.WriteServerBinary(c.readWriter, src)
		c.readWriter.Flush()
		dataLen = len(src)
	} else {
		dataLen, err = c.readWriter.Write(src)
		c.readWriter.Flush()
	}
	if err != nil {
		log.Error.Println("Problem occurred writing data to a client:", err)
		c.player.Destroy()
		return -1
	}
	return dataLen
}
//writePacket This is a method to send a net to the client.  If this is a bare net, the net payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this net.
func (c *client) writePacket(p net.Packet) {
	if p.Opcode == 0 {
		c.Write(p.FrameBuffer)
		return
	}
	frameLength := len(p.FrameBuffer)
	header := []byte{0, 0}
	if frameLength >= 160 {
		header[0] = byte(frameLength>>8 + 160)
		header[1] = byte(frameLength)
	} else {
		header[0] = byte(frameLength)
		frameLength--
		header[1] = p.FrameBuffer[frameLength]
	}
	c.Write(append(header, p.FrameBuffer[:frameLength]...))
	return
}