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
	world.Player
	socket     stdnet.Conn
	destroyer  sync.Once
	readWriter *bufio.ReadWriter
	websocket  bool
	readSize  int
	readLimit  int
	frameFin bool
	incomingPackets chan *net.Packet
}

//newClient Creates a new instance of a client, launches goroutines to handle I/O for it, and returns a reference to it.
func newClient(socket stdnet.Conn, webclient bool) *client {
	c := &client{socket: socket, Player: *world.NewPlayer(-1, strings.Split(socket.RemoteAddr().String(), ":")[0]),
			websocket: webclient, readWriter: bufio.NewReadWriter(bufio.NewReader(socket), bufio.NewWriter(socket)), incomingPackets: make(chan *net.Packet, 25)}
//	defer c.startNetworking()
//	terminate := make(chan struct{})
//	awaitDeath := sync.WaitGroup{}
	destroy := func() {
		c.Attributes.SetVar("lastIP", c.CurrentIP())
//		c.UpdateWG.Lock()
		close(c.OutgoingPackets)
		close(c.incomingPackets)
//		if err := c.socket.Close(); err != nil {
//			log.Warn("Couldn't close socket:", err)
//		}
		if player, ok := world.Players.FromIndex(c.Index); ok && player.UsernameHash() != c.UsernameHash() || !ok || !c.Connected() {
			if ok {
				log.Cheatf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", c.Index, player)
			}
			return
		}
		c.SetConnected(false)
		go db.DefaultPlayerService.PlayerSave(&c.Player)
		world.RemovePlayer(&c.Player)
		log.Debug("Unregistered:'" + c.Username() + "'@'" + c.CurrentIP() + "'")
//		c.UpdateWG.Unlock()
	}
	go func() {
//		defer c.destroyer.Do(destroy)
		defer c.Destroy()
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
						log.Warn("Problem creating reader for next websocket frame:", err)
					}
					c.readWriter.Reader.Reset(reader)
				}
				n, err := c.readWriter.Read(data[written:])
				c.readSize += n
				if err != nil {
					if err == io.EOF {
						if !c.frameFin {
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
			return written
		}
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
					log.Cheatf("Invalid packet length from [%v]: %d\n", c, frameSize)
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
				if !c.Connected() && !handshake.EarlyOperation(int(localData[0])) {
					log.Warnf("Unauthorized packet[opcode:%v,size:%v (expected:%v)] rejected from: %v\n", localData[0], len(localData), frameSize, c)
					continue
				}
				c.incomingPackets <- net.NewPacket(localData[0], localData[1:])
			case <-c.KillC:
				return
			}
		}
	}()
	go func() {
		defer c.destroyer.Do(destroy)
		for {
			select {
			case p := <-c.incomingPackets:
				if p == nil {
					log.Warn(c, "tried processing nil packet:", p)
					continue
				}
				handler := handlers.Handler(p.Opcode)
				if handler == nil {
					log.Debugf("Packet{\n\topcode:%d;\n\tlength:%d;\n\tpayload:%v\n};\n", p.Opcode, len(p.FrameBuffer), p.FrameBuffer)
					return
				}

				handler(&c.Player, p)
			case p := <-c.OutgoingPackets:
				if p == nil {
					return
				}
				if p.Opcode != 0 {
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
					p.FrameBuffer = append(header, p.FrameBuffer[:frameLength]...)
				}
				c.Write(p.FrameBuffer)
			case <-c.KillC:
				return
			}
		}
	}()

	return c
}

//Write Writes data to the client's socket from `b`.  Returns the length of the written bytes.
func (c *client) Write(src []byte) int {
	var err error
	var dataLen int
	if c.websocket {
		err = wsutil.WriteServerBinary(c.readWriter, src)
		dataLen = len(src)
	} else {
		dataLen, err = c.readWriter.Write(src)
	}
	if err != nil {
		log.Warn("Problem occurred writing data to a client:", err)
		c.Destroy()
		return -1
	}
	c.readWriter.Flush()
	return dataLen
}