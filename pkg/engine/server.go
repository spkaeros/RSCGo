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
	"crypto/tls"
	stdnet "net"
	"os"
	"strconv"
	"time"

	"github.com/gobwas/ws"

	"github.com/spkaeros/rscgo/pkg/config"
	//"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/engine/tasks"
	rscerrors "github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/net/handlers"
	"github.com/spkaeros/rscgo/pkg/game/net/handshake"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

const TickMillis = time.Millisecond*640

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
					socket = tls.Server(socket, &tls.Config{Certificates: []tls.Certificate{certChain}, ServerName: "rscturmoil.com", InsecureSkipVerify: false, SessionTicketsDisabled: true})
				}
				if _, err := wsUpgrader.Upgrade(socket); err != nil {
					if config.Verbosity > 0 {
						log.Warn("Encountered a problem attempting to upgrade the HTTP(S) connection to use websockets:", err)
					}
					continue
				}
			}
			p := world.NewPlayer(socket)
			go func() {
				defer close(p.OutQueue)
				defer close(p.InQueue)
				defer p.Destroy()
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
							switch err.(type) {
							case rscerrors.NetError:
								if err.(rscerrors.NetError).Fatal {
									p.Destroy()
									return
								}
							}
							log.Debug(err)
							continue
						}
						frameSize := int(header[0] & 0xFF)
						if frameSize >= 160 {
							frameSize = ((frameSize - 160) << 8) | int(header[1]&0xFF)
						} else {
							frameSize -= 1
						}

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
						if handshake.EarlyOperation(int(localData[0])) {
							handlers.Handler(localData[0])(p, net.NewPacket(localData[0], localData[1:]))
							continue
						}
						if !p.Connected() {
							log.Warnf("Unauthorized packet[opcode:%v,size:%v (expected:%v)] rejected from: %v\n", localData[0], len(localData), frameSize, p)
							return
						}
						p.InQueue <- *net.NewPacket(localData[0], localData[1:])
					}
				}
			}()
		}
	}()
}

func handlePackets(p *world.Player) {
	//	tasks.TickList.Add(func() bool {
	//	go func() {
	go func() {
		for {
			select {
			case p1, ok := <-p.InQueue:
				if !ok || p1.Opcode == 0xFF {
					return
				}
				handler := handlers.Handler(p1.Opcode)
				if handler == nil {
					log.Debugf("Packet{\n\topcode:%d;\n\tlength:%d;\n\tpayload:%v\n};\n", p1.Opcode, len(p1.FrameBuffer), p1.FrameBuffer)
					return
				}

				// log.Debug(p1)
				handler(p, &p1)
				continue
			default:
				return
			}
		}
	}()
}

//StartGameEngine Launches a goroutine to handle updating the state of the game every 640ms in a synchronized fashion.  This is known as a single game engine 'pulse'.
func StartGameEngine() {
	//	ticker := time.NewTicker(640 * time.Millisecond)
	//	defer ticker.Stop()
	var start time.Time
	//	for range ticker.C {
	for {
		start = time.Now()

		tasks.TickList.RunAsynchronous()

		world.Players.Range(func(p *world.Player) {
			handlePackets(p)
			if p.TickAction() != nil && p.TickAction()() {
				p.ResetTickAction()
			}
			p.Tickables.Tick(interface{}(p))
			p.TraversePath()
		})

		world.UpdateNPCPositions()
		world.Players.Range(func(p *world.Player) {
			if positions := world.PlayerPositions(p); positions != nil {
				// log.Debug("position!")
				p.SendPacket(positions)
			}
			if appearances := world.PlayerAppearances(p); appearances != nil {
				// log.Debug("event!")
				p.SendPacket(appearances)
			}
			if npcUpdates := world.NPCPositions(p); npcUpdates != nil {
				// log.Debug("npcPosition!")
				p.SendPacket(npcUpdates)
			}
			if npcUpdates := world.NpcEvents(p); npcUpdates != nil {
				// log.Debug("npcEvents!")
				p.SendPacket(npcUpdates)
			}
			if objectUpdates := world.ObjectLocations(p); objectUpdates != nil {
				// log.Debug("sceneEvent!")
				p.SendPacket(objectUpdates)
			}
			if boundaryUpdates := world.BoundaryLocations(p); boundaryUpdates != nil {
				// log.Debug("boundarys!")
				p.SendPacket(boundaryUpdates)
			}
			if itemUpdates := world.ItemLocations(p); itemUpdates != nil {
				// log.Debug("lootables!")
				p.SendPacket(itemUpdates)
			}
			if clearDistantChunks := world.ClearDistantChunks(p); clearDistantChunks != nil {
				// log.Debug("chunkClear!")
				p.SendPacket(clearDistantChunks)
			}
		})

		world.Players.Range(func(p *world.Player) {
			p.FlushOutgoing()
			p.PostTickables.Tick(interface{}(p))
			p.ResetRegionRemoved()
			p.ResetRegionMoved()
			p.ResetSpriteUpdated()
			p.ResetAppearanceChanged()
		})
		world.ResetNpcUpdateFlags()

		sleepTime := TickMillis-time.Now().Sub(start)
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
}

//Stop This will stop the game instance, if it is running.
func Stop() {
	log.Debug("Stopping...")
	os.Exit(0)
}
