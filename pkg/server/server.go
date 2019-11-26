package server

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"net"
	"os"
	"time"

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
				if _, err := wsUpgrader.Upgrade(socket); err != nil {
					log.Info.Println("Error upgrading websocket connection:", err)
					continue
				}
			}
			if clients.Size() >= config.MaxPlayers() {
				if n, err := socket.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 14}); err != nil || n != 9 {
					if config.Verbosity > 0 {
						log.Error.Println("Could not send world is full response to rejected client:", err)
					}
				}
				continue
			}

			NewClient(socket, port == config.WSPort())
		}
	}()
}

func StartConnectionService() {
	Bind(config.Port())   // UNIX sockets
	Bind(config.WSPort()) // websockets
}

//Tick One game engine 'tick'.  This is to handle movement, to synchronize client, to update movement-related state variables... Runs once per 600ms.
func Tick() {
	select {
	case fn := <-script.EngineChannel:
		fn()
	default:
		break
	}
	clients.Range(func(c clients.Client) {
		if fn := c.Player().DistancedAction; fn != nil {
			if fn() {
				c.Player().ResetDistancedAction()
			}
		}
		nextTile := c.Player().TraversePath()
		if nextTile.LongestDelta(c.Player().Location) > 0 {
			c.Player().SetLocation(nextTile)
			c.Player().Move()
		}
	})
	go world.UpdateNPCPaths()
	world.UpdateNPCPositions()
	clients.Range(func(c clients.Client) {
		c.UpdatePositions()
	})
	clients.Range(func(c clients.Client) {
		if time.Since(c.Player().Transients().VarTime("deathTime")) < 5*time.Second {
			// Ugly hack to work around a client bug with region loading.
			return
		}
		c.ResetUpdateFlags()
	})
	world.ResetNpcUpdateFlags()
}

//StartGameEngine Launches a goroutine to handle updating the state of the server every 600ms in a synchronized fashion.  This is known as a single game engine 'pulse'.
func StartGameEngine() {
	go func() {
		ticker := time.NewTicker(600 * time.Millisecond)
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
