package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/packethandlers"
	"bitbucket.org/zlacki/rscgo/pkg/server/script"
	"fmt"
	"github.com/gobwas/ws"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
)

var (
	kill = make(chan struct{})
)

//Flags This is used to interface with the go-flags package from some guy on github.
var Flags struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port      int    `short:"p" long:"port" description:"The port for the server to listen on,"`
	Config    string `short:"c" long:"config" description:"Specify the configuration file to load server settings from" default:"config.toml"`
	UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using ISAAC to encrypt packet opcodes."`
}

var wsUpgrader = ws.Upgrader{
	Protocol: func(protocol []byte) bool {
		// Chrome is picky, won't work without explicit protocol acceptance
		return true
	},
}

func startConnectionService() {
	if Flags.Port > 0 {
		config.TomlConfig.Port = Flags.Port
	}
	bind := func(offset int) {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port()+offset))
		if err != nil {
			log.Error.Printf("Can't bind to specified port: %d\n", config.Port()+offset)
			log.Error.Println(err)
			os.Exit(1)
		}

		go func() {
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
					if len(Flags.Verbose) > 0 {
						log.Error.Println("Error occurred attempting to accept a client:", err)
					}
					continue
				}
				if offset != 0 {
					if _, err := wsUpgrader.Upgrade(socket); err != nil {
						log.Info.Println("Error upgrading websocket connection:", err)
						continue
					}
				}
				if clients.Size() >= config.MaxPlayers() {
					if n, err := socket.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 14}); err != nil || n != 9 {
						if len(Flags.Verbose) > 0 {
							log.Error.Println("Could not send world is full response to rejected client:", err)
						}
					}
					continue
				}

				c := NewClient(socket)
				if offset != 0 {
					c.Player().Websocket = true
				}
			}
		}()
	}

	bind(0) // UNIX sockets
	bind(1) // websockets
}

//Start Listens for and processes new client connecting to the server.
// This method blocks while the server is running.
func Start() {
	log.Info.Println("RSCGo starting up...")
	if _, err := flags.Parse(&Flags); err != nil {
		log.Error.Println(err)
		return
	}
	if Flags.Port > 65535 || Flags.Port < 0 {
		log.Warning.Println("Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		return
	}
	if !strings.HasSuffix(Flags.Config, ".toml") {
		log.Warning.Println("You entered an invalid configuration file extension.")
		log.Warning.Println("TOML is currently the only supported format for server properties.")
		log.Warning.Println()
		log.Info.Println("Setting back to default: `config.toml`")
		Flags.Config = "config.toml"
	}
	if Flags.UseCipher {
		log.Info.Println("TODO: Figure out why ISAAC cipher sometimes works, yet eventually desynchronizes from client.")
		log.Info.Println("Cipher will remain disabled until such time as this issue gets resolved.")
		Flags.UseCipher = false
	}

	if _, err := toml.DecodeFile("."+string(os.PathSeparator)+Flags.Config, &config.TomlConfig); err != nil {
		log.Warning.Println("Error decoding TOML RSCGo general configuration file:", err)
		return
	}
	if len(Flags.Verbose) > 0 {
		log.Info.Println()
		log.Info.Println("Loaded TOML configuration file.")
	}

	var awaitLaunchJobs sync.WaitGroup
	asyncExecute(&awaitLaunchJobs, func() {
		db.LoadItemDefinitions()
		if count := len(db.Items); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d item definitions.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		db.LoadNpcDefinitions()
		if count := len(db.Npcs); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d NPC definitions.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		err := db.LoadObjectDefinitions()
		if err != nil {
			log.Warning.Println("Could not load game object definitions:", err)
			return
		}
		if count := len(db.Objects); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d game object definitions.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		db.LoadBoundaryDefinitions()
		if count := len(db.Boundarys); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d boundary definitions.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		if count := db.LoadObjectLocations(); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d game objects.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		if count := db.LoadNpcLocations(); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d NPCs.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		script.LoadObjectTriggers()
		log.Info.Printf("Loaded %d object action triggers.\n", len(script.ObjectTriggers))
	})
	asyncExecute(&awaitLaunchJobs, func() {
		script.LoadBoundaryTriggers()
		log.Info.Printf("Loaded %d boundary action triggers.\n", len(script.BoundaryTriggers))
	})
	asyncExecute(&awaitLaunchJobs, func() {
		packethandlers.Initialize()
		if len(Flags.Verbose) > 0 {
			log.Info.Printf("Initialized %d packet handlers.\n", packethandlers.CountPacketHandlers())
		}
	})
	// TODO: Re-enable RSA
	/*
		asyncExecute(&awaitLaunchJobs, func() {
			loadRsaKey()
			if len(Flags.Verbose) > 0 {
				log.Info.Println("Loaded RSA key data.")
			}
		})
	*/
	asyncExecute(&awaitLaunchJobs, func() {
		startGameEngine()
		if len(Flags.Verbose) > 0 {
			log.Info.Println("Launched game engine.")
		}
	})
	awaitLaunchJobs.Wait()
	log.Info.Println()
	log.Info.Println("RSCGo is now running.")
	startConnectionService()
	log.Info.Printf("Listening on TCP port %d, websocket port %d...\n", config.Port(), config.Port()+1)
	select {
	case <-kill:
		os.Exit(0)
	}
}

//asyncExecute First this will add 1 task to the specified waitgroup, then it will execute the function fn in its own goroutine, and upon exiting this goroutine will indicate to wg that the task we added is finished.
func asyncExecute(wg *sync.WaitGroup, fn func()) {
	(*wg).Add(1)
	go func() {
		defer (*wg).Done()
		fn()
	}()
}

//Tick One game engine 'tick'.  This is to handle movement, to synchronize client, to update movement-related state variables... Runs once per 600ms.
func Tick() {
	clients.Range(func(c clients.Client) {
		if fn := c.Player().DistancedAction; fn != nil {
			if fn() {
				c.Player().ResetDistancedAction()
			}
		}
		c.Player().TraversePath()
	})
	world.UpdateNPCPositions()
	clients.Range(func(c clients.Client) {
		c.UpdatePositions()
	})
	clients.Range(func(c clients.Client) {
		c.ResetUpdateFlags()
	})
	world.ResetNpcUpdateFlags()
	fns := script.ActiveTriggers
	script.ActiveTriggers = script.ActiveTriggers[:0]
	for _, fn := range fns {
		go fn()
	}

}

//startGameEngine Launches a goroutine to handle updating the state of the server every 600ms in a synchronized fashion.  This is known as a single game engine 'pulse'.
func startGameEngine() {
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
	kill <- struct{}{}
}
