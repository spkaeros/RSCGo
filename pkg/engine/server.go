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
	"context"
	"crypto/tls"
	stdnet "net"
	"os"
	"reflect"
	"strconv"
	"strings"
	`sync`
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/engine/tasks"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/net/handlers"
	"github.com/spkaeros/rscgo/pkg/game/net/handshake"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

//Bind binds to the TCP port at port, and the websocket port at port+1.
func Bind(port int) {
	listener, err := stdnet.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Can't bind to specified port: %d\n", port)
		log.Fatal(err)
		os.Exit(1)
	}
	go func() {
		var wsUpgrader = ws.Upgrader{
			Protocol: func(protocol []byte) bool {
				// Chrome is picky, won't work without explicit protocol acceptance
				return true
			},
			ReadBufferSize:  23768,
			WriteBufferSize: 23768,
		}

		defer func() {
			err := listener.Close()
			if err != nil {
				log.Error.Println("Could not close game socket listener:", err)
			}
		}()

		certChain, certErr := tls.LoadX509KeyPair("./data/ssl/fullchain.pem", "./data/ssl/privkey.pem")

		for {
			socket, err := listener.Accept()
			if err != nil {
				if config.Verbosity > 0 {
					log.Error.Println("Error occurred attempting to accept a client:", err)
				}
				continue
			}
			if port == config.WSPort() {
				if certErr == nil {
					// set up socket to use TLS if we have certs that we can load
					socket = tls.Server(socket, &tls.Config{Certificates: []tls.Certificate{certChain}, ServerName: "rscturmoil.com", InsecureSkipVerify: true, SessionTicketsDisabled: true})
				}
				if _, err := wsUpgrader.Upgrade(socket); err != nil {
					if config.Verbosity > 0 {
						log.Warn("Encountered a problem attempting to upgrade the HTTP(S) connection to use websockets:", err)
					}
					continue
				}
			}
			p := world.NewPlayer(world.Players.NextIndex(), strings.Split(socket.RemoteAddr().String(), ":")[0])
			p.Socket = socket
			p.Reader = bufio.NewReader(socket)
			go func() {
				defer func() {
					p.Attributes.SetVar("lastIP", p.CurrentIP())
					if err := p.Socket.Close(); err != nil {
						log.Warn("Couldn't close socket:", err)
					}
					if player, ok := world.Players.FromIndex(p.Index); ok && player.UsernameHash() != p.UsernameHash() || !ok || !p.Connected() {
						if ok {
							log.Cheatf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", p.Index, player)
						}
						return
					}
					if p.Connected() {
						go db.DefaultPlayerService.PlayerSave(p)
					}
					world.RemovePlayer(p)
					p.SetConnected(false)
					log.Debug("Unregistered:{'" + p.Username() + "'@'" + p.CurrentIP() + "'}")
				}()
				for {
					select {
					case <-p.SigKill:
						return
					case packet, ok := <-p.OutgoingPackets:
						if ok && packet != nil {
							if strings.HasSuffix(p.Socket.LocalAddr().String(), "43595") {
								writer := wsutil.NewWriter(socket, ws.StateServerSide, ws.OpBinary)
								if packet.Opcode == 0 {
									writer.Write(packet.FrameBuffer)
									writer.Flush()
									continue
								}
								header := []byte{0, 0}
								frameLength := len(packet.FrameBuffer)
								if frameLength >= 160 {
									header[0] = byte(frameLength>>8 + 160)
									header[1] = byte(frameLength)
								} else {
									header[0] = byte(frameLength)
									if frameLength > 0 {
										frameLength--
										header[1] = packet.FrameBuffer[frameLength]
									}
								}
								writer.Write(append(header, packet.FrameBuffer[:frameLength]...))
								writer.Flush()
							} else {
								writer := bufio.NewWriter(socket)
								if packet.Opcode == 0 {
									writer.Write(packet.FrameBuffer)
									writer.Flush()
									continue
								}
								header := []byte{0, 0}
								frameLength := len(packet.FrameBuffer)
								if frameLength >= 160 {
									header[0] = byte(frameLength>>8 + 160)
									header[1] = byte(frameLength)
								} else {
									header[0] = byte(frameLength)
									if frameLength > 0 {
										frameLength--
										header[1] = packet.FrameBuffer[frameLength]
									}
								}
								writer.Write(append(header, packet.FrameBuffer[:frameLength]...))
								writer.Flush()
							}
						}
					}
				}
			}()
			go func() {
				for {
					select {
					case <-p.SigKill:
						return
					default:
						header := make([]byte, 2)
						n, err := p.Read(header)
						if n < 2 {
							continue
						} else if err != nil {
							log.Debug(err)
							continue
						}
						frameSize := int(header[0] & 0xFF)
						if frameSize >= 160 {
							frameSize = ((frameSize - 160) << 8) | int(header[1]&0xFF)
						} else {
							frameSize -= 1
						}

						// Upper bound is an approximation of the max size of the clientside outgoing data buffer
						if frameSize >= 24573 || frameSize < 0 {
							log.Cheatf("Invalid packet length from [%v]: %d\n", p, frameSize)
							return
						}
						localData := make([]byte, frameSize)
						if frameSize > 0 {
							n, err := p.Read(localData)
							if n < frameSize {
								continue
							} else if err != nil {
								log.Debug(err)
								continue
							}
						}
						if frameSize < 160 {
							localData = append(localData, header[1])
						}
						if !p.Connected() && !handshake.EarlyOperation(int(localData[0])) {
							log.Warnf("Unauthorized packet[opcode:%v,size:%v (expected:%v)] rejected from: %v\n", localData[0], len(localData), frameSize, p)
							return
						}
						p1 := net.NewPacket(localData[0], localData[1:])
						handler := handlers.Handler(p1.Opcode)
						if handler == nil {
							log.Debugf("Packet{\n\topcode:%d;\n\tlength:%d;\n\tpayload:%v\n};\n", p1.Opcode, len(p1.FrameBuffer), p1.FrameBuffer)
							continue
						}

						handler(p, p1)
					}
				}
			}()
		}
	}()
}

func runTickables(p *world.Player) {
	var toRemove []int
	for i, fn := range p.Tickables {
		if realFn, ok := fn.(func(context.Context) (reflect.Value, reflect.Value)); ok {
			_, err := realFn(context.Background())
			if !err.IsNil() {
				toRemove = append(toRemove, i)
				log.Warn("Error in tickable:", err)
				continue
			}
		}
		if realFn, ok := fn.(func(context.Context, reflect.Value) (reflect.Value, reflect.Value)); ok {
			_, err := realFn(context.Background(), reflect.ValueOf(p))
			if !err.IsNil() {
				toRemove = append(toRemove, i)
				log.Warn("Error in tickable:", err)
				continue
			}
		}
		if realFn, ok := fn.(func()); ok {
			realFn()
		}
		if realFn, ok := fn.(func() bool); ok {
			if realFn() {
				toRemove = append(toRemove, i)
			}
		}
		if realFn, ok := fn.(func(*world.Player)); ok {
			realFn(p)
		}
		if realFn, ok := fn.(func(*world.Player) bool); ok {
			if realFn(p) {
				toRemove = append(toRemove, i)
			}
		}
	}
	for _, idx := range toRemove {
		p.Tickables[idx] = nil
		p.Tickables = p.Tickables[:idx]
		if idx < len(p.Tickables)-1 {
			p.Tickables = append(p.Tickables[idx+1:])
		}
	}
}
	/*	world.Players.Range(func(p *world.Player) {
			for _, fn := range p.ResetTickables {
				fn()
			}
			p.ResetTickables = p.ResetTickables[:0]
		})
	*/
/*
		world.Players.Range(func(p *world.Player) {

			if p.Unregistering.Load() && p.Unregistering.CAS(true, false) {
				p.Attributes.SetVar("lastIP", p.CurrentIP())
				if err := p.Socket.Close(); err != nil {
					log.Warn("Couldn't close socket:", err)
				}
				if player, ok := world.Players.FromIndex(p.Index); ok && player.UsernameHash() != p.UsernameHash() || !ok || !p.Connected() {
					if ok {
						log.Cheatf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", p.Index, player)
					}
					return
				}
				world.RemovePlayer(p)
				p.SetConnected(false)
				log.Debug("Unregistered:{'" + p.Username() + "'@'" + p.CurrentIP() + "'}")
				go db.DefaultPlayerService.PlayerSave(p)
			}
		})
*/

//StartGameEngine Launches a goroutine to handle updating the state of the game every 640ms in a synchronized fashion.  This is known as a single game engine 'pulse'.
func StartGameEngine() {
	ticker := time.NewTicker(640 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		tasks.TickList.RunAsynchronous()
//		tasks.TickList.RunSynchronous()
		wait := sync.WaitGroup{}
		world.Players.Range(func(p *world.Player) {
			wait.Add(1)
			go func() {
				defer wait.Done()
				runTickables(p)
				if fn := p.DistancedAction; fn != nil {
					if fn() {
						p.ResetDistancedAction()
					}
				}
				p.TraversePath()
			}()
		})
		wait.Wait()
		world.UpdateNPCPositions()
		
		world.Players.Range(func(p *world.Player) {
			wait.Add(1)
			go func() {
				defer wait.Done()
				// Everything is updated relative to our player's position, so player position net comes first
				if positions := world.PlayerPositions(p); positions != nil {
					p.SendPacket(positions)
				}
				if appearances := world.PlayerAppearances(p); appearances != nil {
					p.SendPacket(appearances)
				}
				if npcUpdates := world.NPCPositions(p); npcUpdates != nil {
					p.SendPacket(npcUpdates)
				}
				if objectUpdates := world.ObjectLocations(p); objectUpdates != nil {
					p.SendPacket(objectUpdates)
				}
				if boundaryUpdates := world.BoundaryLocations(p); boundaryUpdates != nil {
					p.SendPacket(boundaryUpdates)
				}
				if itemUpdates := world.ItemLocations(p); itemUpdates != nil {
					p.SendPacket(itemUpdates)
				}
				if clearDistantChunks := world.ClearDistantChunks(p); clearDistantChunks != nil {
					p.SendPacket(clearDistantChunks)
				}
			}()
		})
		wait.Wait()
		
		world.Players.Range(func(p *world.Player) {
			wait.Add(1)
			go func() {
				defer wait.Done()
				p.ResetRegionRemoved()
				p.ResetRegionMoved()
				p.ResetSpriteUpdated()
				p.ResetAppearanceChanged()
			}()
		})
		wait.Wait()
		
		world.ResetNpcUpdateFlags()
	}
}

//Stop This will stop the game instance, if it is running.
func Stop() {
	log.Debug("Stopping...")
	os.Exit(0)
}
