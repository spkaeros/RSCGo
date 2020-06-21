/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package main

import (
	"crypto/tls"
	stdnet "net"
	"os"
	"sync"
	"strconv"
	"time"
	"strings"
	"math"
	"bufio"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/jessevdk/go-flags"
	"github.com/BurntSushi/toml"

	"github.com/spkaeros/rscgo/pkg/crypto"
	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/tasks"
	rscerrors "github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/isaac"
	"github.com/spkaeros/rscgo/pkg/strutil"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game"
	"github.com/spkaeros/rscgo/pkg/game/net/handshake"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	
	_ "github.com/spkaeros/rscgo/pkg/game/net/handlers"
)

const (
	TickMillis = time.Millisecond*640
)

type (
	Flags struct {
		Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
		Port      int    `short:"p" long:"port" description:"The TCP port for the game to listen on, (Websocket will use the port directly above it)"`
		Config    string `short:"c" long:"config" description:"Specify the TOML configuration file to load game settings from" default:"config.toml"`
		UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using a variant of ISAAC to encrypt net opcodes."`
	}
	Server struct {
		port int
		listener stdnet.Listener
		*time.Ticker
	}
)

var (
	cliFlags = &Flags{}
	start = time.Now()
	newPlayers chan *world.Player
	tlsCerts, tlsError = tls.LoadX509KeyPair("./data/ssl/fullchain.pem", "./data/ssl/privkey.pem")
	wsUpgrader = ws.Upgrader{
		Protocol: func(protocol []byte) bool {
			// Chrome is picky, won't work without explicit protocol acceptance
			return true
		},
		ReadBufferSize:  23768,
		WriteBufferSize: 23768,
	}
)

func main() {
	// Initialize sane defaults as fallback configuration options, if the config.toml file is not found or if some values are left out of it
	config.TomlConfig.MaxPlayers = 1250
	config.TomlConfig.DataDir = "./data/"
	config.TomlConfig.DbioDefs = "./data/dbio.conf"
	config.TomlConfig.PacketHandlerFile = "./data/packets.toml"
	config.TomlConfig.Crypto.HashComplexity = 15
	config.TomlConfig.Crypto.HashLength = 32
	config.TomlConfig.Crypto.HashMemory = 8
	config.TomlConfig.Crypto.HashSalt = "rscgo./GOLANG!RULES/.1994"
	config.TomlConfig.Version = 204
	config.TomlConfig.Port = 43594 // +1 for websockets

	if _, err := flags.Parse(cliFlags); err != nil {
		log.Warn("Error parsing command arguments:", cliFlags)
		return
	}
	// Default to config.toml for config file
	if len(cliFlags.Config) == 0 {
		cliFlags.Config = "config.toml"
	}
	if _, err := toml.DecodeFile(cliFlags.Config, &config.TomlConfig); err != nil {
		log.Warn("Error decoding server TOML configuration file `" + cliFlags.Config + "`:", err)
		log.Fatal("Error decoding server TOML configuration file:", "`" + cliFlags.Config + "`")
		log.Fatal(err)
		os.Exit(101)
		return 
	}

	// TODO: data backend default to JSON or BSON maybe?
	config.TomlConfig.Database.PlayerDriver = "sqlite3"
	config.TomlConfig.Database.PlayerDB = "file:./data/players.db"
	config.TomlConfig.Database.WorldDriver = "sqlite3"
	config.TomlConfig.Database.WorldDB = "file:./data/world.db"
	if _, err := toml.DecodeFile(config.TomlConfig.DbioDefs, &config.TomlConfig.Database); err != nil {
		log.Warn("Error reading database config file:", err)
		return
	}

	run(game.UnmarshalPackets)

	if cliFlags.Port > 0 {
		config.TomlConfig.Port = cliFlags.Port
	}
	if config.Port() >= 65534 || config.Port() < 0 {
		log.Warn("Error: Invalid port number specified.")
		log.Warn("Valid port numbers are 1-65533 (needs the port 1 above it open to bind a websockets listener).")
		return 
	}

	config.Verbosity = len(cliFlags.Verbose)
	run(db.ConnectEntityService, func() {
		db.DefaultPlayerService, world.DefaultPlayerService = db.NewPlayerServiceSql(), db.NewPlayerServiceSql()
	})
	// Three init phases after data backend is connected--Entity definitions, then tile collision bitmask loading, followed by entity spawn locations
	// So, the order here of these three phases is important.  If you attempt to load object spawn locations during the same phase as the collision
	// data, it will result in a world filled with objects that are not solid.  Many similar bugs possible.  Best just to leave this be.
	run(db.LoadTileDefinitions, db.LoadObjectDefinitions, db.LoadBoundaryDefinitions, db.LoadItemDefinitions, db.LoadNpcDefinitions)
	run(world.LoadCollisionData, world.RunScripts)
	run(db.LoadObjectLocations, db.LoadNpcLocations, db.LoadItemLocations)

	if config.Verbose() {
		log.Debug("Loaded", len(world.Sectors), "map sectors")
		log.Debug("Loaded", game.PacketCount(), "packets (with", game.HandlerCount(), "handlers)")
		log.Debug("Loaded", world.ItemIndexer.Load(), "items and", len(definitions.Items), "item definitions")
		log.Debug("Loaded", world.Npcs.Size(), "NPCs and", len(definitions.Npcs), "NPC definitions")
		log.Debug("Loaded", len(definitions.ScenaryObjects), "scenary definitions, and", len(definitions.BoundaryObjects), "boundary definitions")
		log.Debug("Loaded", world.ObjectCounter.Load(), "scenary / boundary objects")
		log.Debug("Loading all game entitys took:", time.Since(start).Seconds(), "seconds")
		if config.Verbosity >= 2 {
			log.Debugf("Triggers[\n\t%d item actions,\n\t%d scenary actions,\n\t%d boundary actions,\n\t%d npc actions,\n\t%d item->boundary actions,\n\t%d item->scenary actions,\n\t%d attacking NPC actions,\n\t%d killing NPC actions\n];\n", len(world.ItemTriggers), len(world.ObjectTriggers), len(world.BoundaryTriggers), len(world.NpcTriggers), len(world.InvOnBoundaryTriggers), len(world.InvOnObjectTriggers), len(world.NpcAtkTriggers), len(world.NpcDeathTriggers))
		}
	}
	log.Debug("Listening at TCP port " + strconv.Itoa(config.Port()))// + " (TCP), " + strconv.Itoa(config.WSPort()) + " (websockets)")
	log.Debug()
	log.Debug("RSCGo has finished initializing world; we hope you enjoy it")
	Instance.Start()
}


var Instance = &Server{Ticker: time.NewTicker(TickMillis)}
func readPacket(player *world.Player) (*net.Packet, error) {
	header := make([]byte, 2)
	n, err := player.Read(header)
	if err != nil {
		switch err.(type) {
		case rscerrors.NetError:
			if err.(rscerrors.NetError).Fatal {
				player.Destroy()
			}
		}
		log.Warn("Error reading packet header:", err)
		return nil, rscerrors.NewNetworkError("Error reading header for packet:" + err.Error(), true)
	}
	if n < 2 {
		return nil, rscerrors.NewNetworkError("Invalid packet-frame length recv; got " + strconv.Itoa(n), false)
	}
	length := int(header[0] & 0xFF)
	if length >= 160 {
		length = (length-160) << 8|int(header[1] & 0xFF)
	} else {
		length -= 1
	}

	frame := make([]byte, length)
	if length > 0 {
		// for written := 0; written < length; {
			_, err := player.Read(frame)
			if err != nil {
				log.Warn("Error reading packet frame:", err)
				return nil, err
			}
			// written += n
		// }
	}

	if length < 160 {
		frame = append(frame, header[1])
	}
	return net.NewPacket(frame[0], frame[1:]), nil
}

//run Helper function for concurrently running a bunch of functions and waiting for them to complete
func run(fns ...func()) {
	w := &sync.WaitGroup{}
	do := func(fn func()) {
		w.Add(1)
		go func(fn func()) {
			defer w.Done()
			fn()
		}(fn)
	}

	for _, fn := range fns {
		do(fn)
	}
	w.Wait()
}


func (s *Server) tlsAccept(l stdnet.Listener) *world.Player {
	socket, err := l.Accept()
	if err != nil {
		log.Errorf("Error: Could not accept new player websocket (%v):%v\n", socket,  err.Error())
		return nil
	}
	if tlsError == nil {
		// This block only runs if the certificate chain was initialized right
		if tmpSock := tls.Server(socket, &tls.Config{Certificates: []tls.Certificate{tlsCerts}, ServerName: "rscturmoil.com", InsecureSkipVerify: false, SessionTicketsDisabled: true}); tmpSock != nil {
			// If we encountered some problem setting up TLS, this should prevent us from losing our original non-encrypted socket hopefully
			socket = tmpSock
		}
	}

	p := world.NewPlayer(socket)
	// TODO: See if we can get TLS working on one port for either TCP sockets or websockets
	_, err = wsUpgrader.Upgrade(p.Socket)
	p.Websocket = err == nil
	if p.IsWebsocket() {
		p.Writer = wsutil.NewWriter(p.Socket, ws.StateServerSide, ws.OpBinary)
	} else {
		p.Writer = bufio.NewWriter(p.Socket)
	}

	return p
}


//Bind binds to the TCP port at port, and the websocket port at port+1.
func (s *Server) Bind(port int) bool {
	listener, err := stdnet.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Can't bind to specified port: %d\n", port)
		log.Fatal(err)
		os.Exit(1)
	}
	s.listener = listener

	go func() {
		defer func() {
			if err := s.listener.Close(); err != nil {
				log.Fatal("closing listener failed:", err)
				os.Exit(1)
			}
		}()
		for {
			player := s.tlsAccept(s.listener)
			login, err := readPacket(player)
			if err != nil {
				if err, ok := err.(rscerrors.NetError); ok {
					if err.Fatal {
						player.Socket.Close()
						continue
					}
				}
				player.Socket.Close()
				continue
			}
			if login == nil {
				player.Socket.Close()
				continue
			}
			if login.Opcode == 32 {
				second, err := readPacket(player)
				if err != nil {
					if err, ok := err.(rscerrors.NetError); ok {
						if err.Fatal {
							player.Socket.Close()
							continue
						}
					}
					player.Socket.Close()
					continue
				}
				login = second
			}
			if login.Opcode == 2 {
				defer func() {
					close(player.InQueue)
					err := player.Socket.Close()
					if err != nil {
						log.Debug("Error closing socket:", err)
						return
					}
					player.Inventory.Owner = nil
				}()
				if version := login.ReadUint16(); version != config.Version() {
					player.WritePacket(world.HandshakeResponse(int(handshake.ResponseUpdated)))
					player.Writer.Flush()
					return
				}
				username := strutil.Base37.Decode(login.ReadUint64())
				password := strings.TrimSpace(login.ReadString())
				reply := func(i handshake.ResponseCode, reason string) {
					player.WritePacket(world.HandshakeResponse(int(i)))
					player.Writer.Flush()
					if reason == "" {
						log.Debug("[REGISTER] Player", "'" + username + "'", "created successfully for:", player.CurrentIP())
						return
					}
					log.Debug("[REGISTER] Player creation failed for:", "'" + username + "'@'" + player.CurrentIP() + "'")
					return
				}
				go func() {
					if userLen, passLen := len(username), len(password); userLen < 2 || userLen > 12 || passLen < 5 || passLen > 20 {
						reply(handshake.ResponseBadInputLength, "Password and/or username too long and/or too short.")
						return
					}
					dataService := db.DefaultPlayerService
					if dataService.PlayerNameExists(username) {
						reply(handshake.ResponseUsernameTaken, "Username is taken by another player already.")
						return
					}

					if !dataService.PlayerCreate(username, crypto.Hash(password), player.CurrentIP()) {
						reply(8, "Data backend seems to have failed creating a player")
						return
					}
					reply(handshake.ResponseRegisterSuccess, "")
				}()
			}
			if login.Opcode == 0 {
				sendReply := func(i handshake.ResponseCode, reason string) {
					player.WritePacket(world.HandshakeResponse(int(i)))
					player.Writer.Flush()
					if i.IsValid() {
						go func() {
							defer close(player.InQueue)
							defer player.Destroy()
							defer player.WritePacket(world.Logout)
							player.Initialize()
							for {
								select {
								default:
									if p, err := readPacket(player); err != nil {
										if err, ok := err.(rscerrors.NetError); ok {
											if err.Fatal {
												return
											}
										}
										return
									} else if p == nil {
										continue
									} else {
										player.InQueue <- p
									}
								}
							}
						}()
						log.Debug("[LOGIN]", player.Username() + "@" + player.CurrentIP(), "successfully logged in")
					} else {
						log.Debug("[LOGIN]", player.Username() + "@" + player.CurrentIP(), "failed to login (" + reason + ")")
						close(player.InQueue)
						player.Destroy()
					}
				}
				if !world.UpdateTime.IsZero() {
					sendReply(handshake.ResponseLoginServerRejection, "System update in progress")
					continue
				}
				if world.Players.Size() >= config.MaxPlayers() {
					sendReply(handshake.ResponseWorldFull, "Out of usable player slots")
					continue
				}
				if handshake.LoginThrottle.Recent(player.CurrentIP(), time.Minute*5) >= 5 {
					sendReply(handshake.ResponseSpamTimeout, "Too many recent invalid login attempts (5 in 5 minutes)")
					continue
				}

				player.SetReconnecting(login.ReadBoolean())
				if ver := login.ReadUint16(); ver != config.Version() {
					sendReply(handshake.ResponseUpdated, "Invalid client version (" + strconv.Itoa(ver) + ")")
					continue
				}

				rsaSize := login.ReadUint16()
				data := make([]byte, rsaSize)
				rsaRead := login.Read(data)
				if rsaRead < rsaSize {
					sendReply(handshake.ResponseLoginServerRejection, "Invalid RSA block")
					continue
				}
				packetDec := net.NewPacket(0, crypto.DecryptRSA(data))
				player.SetVar("ourRng", isaac.New(packetDec.ReadUint64()))
				player.SetVar("theirRng", isaac.New(packetDec.ReadUint64()))
				player.SetVar("username", strutil.Base37.Encode(strings.TrimSpace(packetDec.ReadString())))
				password := strings.TrimSpace(packetDec.ReadString())
				if world.Players.ContainsHash(player.UsernameHash()) {
					sendReply(handshake.ResponseLoggedIn, "Player with same username is already logged in")
					continue
				}
				var dataService = db.DefaultPlayerService
				if !dataService.PlayerNameExists(player.Username()) || !dataService.PlayerValidLogin(player.UsernameHash(), crypto.Hash(password)) {
					handshake.LoginThrottle.Add(player.CurrentIP())
					sendReply(handshake.ResponseBadPassword, "Invalid credentials")
					continue
				}
				if !dataService.PlayerLoad(player) {
					sendReply(handshake.ResponseDecodeFailure, "Could not load player profile; is the dataService setup properly?")
					continue
				}

				if player.Reconnecting() {
					sendReply(handshake.ResponseReconnected, "")
					continue
				}
				switch player.Rank() {
				case 2:
					sendReply(handshake.ResponseAdministrator, "")
				case 1:
					sendReply(handshake.ResponseModerator, "")
				default:
					sendReply(handshake.ResponseLoginSuccess, "")
				}

				continue
			}
		}
	}()
	config.Verbosity = int(math.Min(math.Max(float64(len(cliFlags.Verbose)), 0), 4))
	return false
}

func (s *Server) handlePackets(p *world.Player) {
	go func() {
		for {
			select {
			default:
				return
			case p1, ok := <-p.InQueue:
				if !ok || p1 == nil {
					return
				}
				if handlePacket := game.Handler(p1.Opcode); handlePacket != nil {
					handlePacket(p, p1)
				}
			}
		}
	}()
}

func (s *Server) Start() {
	s.Bind(config.Port())
	defer s.Ticker.Stop()
	for range s.C {
		tasks.TickList.Tick()

		world.Players.Range(func(p *world.Player) {
			if p == nil {
				return
			}
			s.handlePackets(p)
			p.Tickables.Call(interface{}(p))
			
			if fn := p.TickAction(); fn != nil && !fn() {
				p.ResetTickAction()
			}
			p.TraversePath()
		})
		world.UpdateNPCPositions()

		world.Players.Range(func(p *world.Player) {
			if p == nil {
				return
			}
			if positions := world.PlayerPositions(p); positions != nil {
				p.SendPacket(positions)
			}
			if appearances := world.PlayerAppearances(p); appearances != nil {
				p.SendPacket(appearances)
			}
			if npcUpdates := world.NPCPositions(p); npcUpdates != nil {
				p.SendPacket(npcUpdates)
			}
			if npcUpdates := world.NpcEvents(p); npcUpdates != nil {
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
		})

		world.Players.Range(func(p *world.Player) {
			if p == nil {
				return
			}
			p.PostTickables.Call(interface{}(p))
			p.ResetRegionRemoved()
			p.ResetRegionMoved()
			p.ResetSpriteUpdated()
			p.ResetAppearanceChanged()
			// p.Writer.Flush()
		})
		world.ResetNpcUpdateFlags()
		world.Ticks.Inc()
	}
}

//Stop This will stop the game instance, if it is running.
func (s *Server) Stop() {
	log.Debug("Stopping...")
	os.Exit(0)
}
