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
	stdnet "net"
	"time"
	"io"
	"strings"
	
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/game/net/handlers"
	"github.com/spkaeros/rscgo/pkg/game/net/handshake"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	//	"go.uber.org/atomic"
)

//client Represents a single connecting client.
type client struct {
	world.Player
	socket    stdnet.Conn
	readSize  int
	readLimit int
	frameFin  bool
}
func (c *client) Read(data []byte) (int, error) {
	written := 0
	for written < len(data) {
		err := c.socket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(15)))
		if err != nil {
			return -1, errors.NewNetworkError("Deadline reached", true)
		}
		if c.Websocket && (c.readSize >= c.readLimit) {
			// reset buffer read index and create the next reader
			header, reader, err := wsutil.NextReader(c.socket, ws.StateServerSide)
			c.readLimit = int(header.Length)
			c.readSize = 0
			c.frameFin = header.Fin
			if err != nil {
				if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
					return -1, errors.NewNetworkError("closed conn", false)
				} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
					return -1, errors.NewNetworkError("timed out", false)
				}
				log.Warn("Problem creating reader for next websocket frame:", err)
			}
			c.ReadWriter.Reader.Reset(reader)
		}
		n, err := c.ReadWriter.Read(data[written:])
		c.readSize += n
		if err != nil {
			if err == io.EOF {
				if !c.frameFin {
					continue
				}
			}
			if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
				return -1, errors.NewNetworkError("closed conn", false)
			} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
				return -1, errors.NewNetworkError("timed out", false)
			}
			continue
		}
		written += n
	}
	return written, nil
}

//newClient Creates a new instance of a client, launches goroutines to handle I/O for it, and returns a reference to it.
func newClient(socket stdnet.Conn, webclient bool) *client {
	c := &client{socket: socket, Player: *world.NewPlayer(-1, strings.Split(socket.RemoteAddr().String(), ":")[0])}
	//destructor := sync.Once{}
	c.ReadWriter = bufio.NewReadWriter(bufio.NewReader(socket), bufio.NewWriterSize(socket, 5000))
	c.Websocket = webclient
	c.Socket = socket
	c.socket = socket
	c.Player.SetVar("client", c)
	world.AddPlayer(&c.Player)
	/*	teardownGroup := sync.WaitGroup{}
		//	go func() {
			defer func() {
				teardownGroup.Wait()
				if player, ok := world.Players.FromIndex(c.Index); ok && player.UsernameHash() != c.UsernameHash() || !ok || !c.Connected() {
					if ok {
						log.Cheatf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", c.Index, player)
					}
					return
				}
				world.RemovePlayer(&c.Player)
				c.SetConnected(false)
				log.Debug("Unregistered:'" + c.Username() + "'@'" + c.CurrentIP() + "'")
				go db.DefaultPlayerService.PlayerSave(&c.Player)
				close(c.SigDisconnect)
			}()*/
	go func() {
		for {
	
			header := make([]byte, 2)
			c := c.VarChecked("client").(*client)
			n, err := c.Read(header)
			if n < 2 || err != nil {
				if err != nil {
					log.Debug("Read error:", err)
				}
				log.Debug("Not enough bytes in header(need 2 bytes, got:", n)
				return
			}
			frameSize := int(header[0] & 0xFF)
			if frameSize >= 160 {
				frameSize = ((frameSize - 160) << 8) | int(header[1]&0xFF)
			} else {
				frameSize -= 1
			}
			
			// Upper bound is an approximation of the max size of the clientside outgoing data buffer
			if frameSize >= 23768 || frameSize < 0 {
				log.Cheatf("Invalid packet length from [%v]: %d\n", c.Player, frameSize)
				return
			}
			localData := make([]byte, frameSize)
			if frameSize > 0 {
			
				n, err := c.Read(localData)
				if n < frameSize || err != nil {
					if err != nil {
						log.Debug("Read error:", err)
					}
					log.Debug("Not enough bytes in header(need", frameSize, "bytes, got:", n)
					return
				}
			}
			if frameSize < 160 {
				localData = append(localData, header[1])
			}
			if !c.Connected() && !handshake.EarlyOperation(int(localData[0])) {
				log.Warnf("Unauthorized packet[opcode:%v,size:%v (expected:%v)] rejected from: %v\n", localData[0], len(localData), frameSize, c)
				return
			}
			//				p.incomingPackets <- net.NewPacket(localData[0], localData[1:])
			log.Debug(localData)
			p1 := net.NewPacket(localData[0], localData[1:])
			handler := handlers.Handler(p1.Opcode)
			if handler == nil {
				log.Debugf("Packet{\n\topcode:%d;\n\tlength:%d;\n\tpayload:%v\n};\n", p1.Opcode, len(p1.FrameBuffer), p1.FrameBuffer)
				return
			}
			
			handler(&c.Player, p1)
		}
	}()
/*			go func() {
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
	//				c.incomingPackets <- net.NewPacket(localData[0], localData[1:])
					p := net.NewPacket(localData[0], localData[1:])
					handler := handlers.Handler(p.Opcode)
					if handler == nil {
						log.Debugf("Packet{\n\topcode:%d;\n\tlength:%d;\n\tpayload:%v\n};\n", p.Opcode, len(p.FrameBuffer), p.FrameBuffer)
						return
					}

					handler(&c.Player, p)
				case <-c.SigKill:
					return
				}
			}
		}()*/
	/*	go func() {
				defer func() {
					c.Attributes.SetVar("lastIP", c.CurrentIP())
					close(c.OutgoingPackets)
		//			if err := c.socket.Close(); err != nil {
		//				log.Warn("Couldn't close socket:", err)
		//			}
					if player, ok := world.Players.FromIndex(c.Index); ok && player.UsernameHash() != c.UsernameHash() || !ok || !c.Connected() {
						if ok {
							log.Cheatf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", c.Index, player)
						}
						return
					}
					world.RemovePlayer(&c.Player)
					c.SetConnected(false)
					log.Debug("Unregistered:'" + c.Username() + "'@'" + c.CurrentIP() + "'")
					go db.DefaultPlayerService.PlayerSave(&c.Player)
				}()
	*/
	//		defer teardownGroup.Done()
	//		teardownGroup.Add(1)
	//		defer c.Destroy()
	/*		defer func() {
	//			c.Write([]byte{1,4})
				err := socket.Close()
				if err != nil {
					log.Debug(err)
				}
			}()*/
	//		defer c.Write([]byte{0,4})
	//		teardownGroup.Add(1)
	//		defer close(c.KillC)
	//		defer c.Destroy()
	/*		for {
				select {
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
					if p.Opcode == 4 {
						return
					}
	//			case <-c.SigKill:
				case <-c.SigDisconnect:
					return
				}
			}
		}()
	*/
	return c
}
