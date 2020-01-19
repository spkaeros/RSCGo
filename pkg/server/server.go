package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/gobwas/ws"
	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

var (
	Kill = make(chan struct{})
)

func Bind(port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Error.Printf("Can't bind to specified port: %d\n", port)
			log.Error.Println(err)
			os.Exit(1)
		}
	go func() {
		var wsUpgrader = ws.Upgrader{
			Protocol: func(protocol []byte) bool {
				// Chrome is picky, won't work without explicit protocol acceptance
				return true
			},
			ReadBufferSize:  5000,
			WriteBufferSize: 5000,
		}

		defer func() {
			err := listener.Close()
			if err != nil {
				log.Error.Println("Could not close server socket listener:", err)
				return
			}
		}()

		for {
			socket, err := listener.Accept()
			if err != nil {
				if config.Verbosity > 0 {
					log.Error.Println("Error occurred attempting to accept a client:", err)
				}
				continue
			}
			if port == config.WSPort() {
				certChain, err := tls.LoadX509KeyPair("./data/ssl/fullchain.pem", "./data/ssl/privkey.pem")
				if err != nil {
					log.Warning.Println("Problem reading cert chain for WSS protocol:", err)
					return
				}
				socket = tls.Server(socket, &tls.Config{Certificates:[]tls.Certificate{certChain}, InsecureSkipVerify:true})
				if _, err := wsUpgrader.Upgrade(socket); err != nil {
					log.Info.Println("Error upgrading websocket connection:", err)
					continue
				}
			}
			if world.Players.Size() >= config.MaxPlayers() {
				if n, err := socket.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 14}); err != nil || n != 9 {
					if config.Verbosity > 0 {
						log.Error.Println("Could not send world is full response to rejected client:", err)
					}
				}
				continue
			}

			newClient(socket, port == config.WSPort())
		}
	}()
}

func StartConnectionService() {
	Bind(config.Port())   // UNIX sockets
	Bind(config.WSPort()) // websockets
}

func runTickables(p *world.Player) {
	var toRemove []int
	for i, fn := range p.Tickables {
		if realFn, ok := fn.(func(context.Context) (reflect.Value, reflect.Value)); ok {
			_, err := realFn(context.Background())
			if !err.IsNil() {
				toRemove = append(toRemove, i)
				log.Warning.Println("Error in tickable:", err)
				continue
			}
		}
		if realFn, ok := fn.(func(context.Context, reflect.Value) (reflect.Value, reflect.Value)); ok {
			_, err := realFn(context.Background(), reflect.ValueOf(p))
			if !err.IsNil() {
				toRemove = append(toRemove, i)
				log.Warning.Println("Error in tickable:", err)
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
		if idx == len(toRemove)-1 {
			p.Tickables = p.Tickables[:idx]
			continue
		}
		p.Tickables = append(p.Tickables[:idx], p.Tickables[idx+1:])
	}
}

//Tick One game engine 'tick'.  This is to handle movement, to synchronize client, to update movement-related state variables... Runs once per 600ms.
func Tick() {
	world.Tickables.ExecuteSequentially()

	world.Players.Range(func(p *world.Player) {
		runTickables(p)
		if fn := p.DistancedAction; fn != nil {
			if fn() {
				p.ResetDistancedAction()
			}
		}
		p.TraversePath()
	})
	world.UpdateNPCPositions()
	world.Players.Range(func(p *world.Player) {
		// Everything is updated relative to our player's position, so player position packet comes first
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
	})
//	world.Players.Range(func(p *world.Player) {
//		for _, fn := range p.ResetTickables {
//			fn()
//		}
//		p.ResetTickables = p.ResetTickables[:0]
//	})
	world.Players.Range(func(p *world.Player) {
		p.ResetRegionRemoved()
		p.ResetRegionMoved()
		p.ResetSpriteUpdated()
		p.ResetAppearanceChanged()
	})
	world.ResetNpcUpdateFlags()
}

//StartGameEngine Launches a goroutine to handle updating the state of the server every 600ms in a synchronized fashion.  This is known as a single game engine 'pulse'.
func StartGameEngine() {
	go func() {
		ticker := time.NewTicker(640 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			Tick()
		}
	}()
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	log.Info.Println("Stopping server...")
	Kill <- struct{}{}
}
