package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

var (
	listener net.Listener
	//LogWarning Log interface for warnings.
	LogWarning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)
	//LogInfo Log interface for debug information.
	LogInfo = log.New(os.Stdout, "[INFO] ", log.Ltime|log.Lshortfile)
	//LogError Log interface for errors.
	LogError   = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
	syncTicker = time.NewTicker(time.Millisecond * 650)
	kill       = make(chan struct{})
	//Version Client version.
	Version = -1
	//DataDirectory The directory for data files to be read from.  This should be expanded by the CLI flag parser.
	DataDirectory = "."
)

func init() {

}

//Flags This is used to interface with the go-flags package from some guy on github.
var Flags struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port      int    `short:"p" long:"port" description:"The port for the server to listen on," default:"43591"`
	Config    string `short:"c" long:"config" description:"Specify the configuration file to load server settings from" default:"config.ini"`
	UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using ISAAC to encrypt packet opcodes."`
}

func bind(port int) {
	var err error
	portS := strconv.Itoa(port)
	listener, err = net.Listen("tcp", ":"+portS)
	if err != nil {
		LogError.Printf("Can't bind to specified port: %d\n", port)
		LogError.Println(err)
		os.Exit(1)
	}
}

func startConnectionService() {
	if listener == nil {
		if len(Flags.Verbose) > 0 {
			LogWarning.Println("Attempted to start connection service without a listener!  This shouldn't happen.")
			LogWarning.Println("Starting listener on default port...")
		}
		bind(43591)
	}

	go func() {
		defer listener.Close()
		// TODO: Can this ticker be made smaller safely?
		connTicker := time.NewTicker(time.Millisecond * 50)
		for range connTicker.C {
			socket, err := listener.Accept()
			if err != nil {
				if len(Flags.Verbose) > 0 {
					LogError.Println("Could not accept connection via server listener.")
					LogError.Println(err)
				}
				return
			}

			client := NewClient(socket)
			ActiveClients.Add(client)
		}
	}()

}

//Start Listens for and processes new clients connecting to the server.
// This method blocks while the server is running.
func Start() {
	LogInfo.Printf("RSCGo starting up...")
	if len(Flags.Verbose) > 0 {
		LogInfo.Printf("Attempting to bind to network...")
	}
	bind(Flags.Port)
	if len(Flags.Verbose) > 0 {
		LogInfo.Printf("done")
		LogInfo.Printf("Attempting to start connection service...")
	}
	startConnectionService()
	if len(Flags.Verbose) > 0 {
		LogInfo.Printf("done")
		LogInfo.Printf("Attempting to start synchronized task service...")
	}
	startSynchronizedTaskService()
	LogInfo.Printf("done\n\n")
	LogInfo.Printf("RSCGo is now running.\n")
	LogInfo.Printf("Listening on port %d...\n", Flags.Port)
	// TODO: Probably need to handle certain signals, for usability sake.
	// TODO: Implement some form of data store for static game data, e.g entity information, seldom-changed config
	//  settings and the like.
	// TODO: Implement a data store for dynamic game data, e.g player information, and so on.
	select {
	case <-kill:
		os.Exit(0)
	}
}

//startSynchronizedTaskService Launches a goroutine to handle updating the state of the server every 650ms in a
// synchronized fashion.  This is known as a single game engine 'pulse'.  All mobile entities must have their position
// updated during this pulse to be compatible with Jagex RSClassic Client software.
// TODO: Can movement be handled concurrently per-player safely on the Jagex Client? Mob movement might not look right.
func startSynchronizedTaskService() {
	go func() {
		for range syncTicker.C {
			// Loop once to actually move the mobs
			var wg sync.WaitGroup
			wg.Add(ActiveClients.Size())
			for _, c := range ActiveClients.values {
				if c, ok := c.(*Client); ok {
					go func() {
						if c.player.Path != nil {
							c.player.TraversePath()
						}
						wg.Done()
					}()
				}
			}
			wg.Wait()
			// Loop again to update the clients about what the mobs have been up to in the prior loop.
			wg.Add(ActiveClients.Size())
			for _, c := range ActiveClients.values {
				if c, ok := c.(*Client); ok {
					go func() {
						var localPlayers []*entity.Player
						for _, r := range entity.SurroundingRegions(c.player.X(), c.player.Y()) {
							for _, p := range r.Players {
								if strutil.Base37(p.Username) != strutil.Base37(c.player.Username) {
									localPlayers = append(localPlayers, p)
								}
							}
						}
						c.WritePacket(packets.PlayerPositions(c.player, localPlayers))
						c.WritePacket(packets.PlayerAppearances(c.index, strutil.Base37(c.player.Username)))
						// TODO: Update movement, update client-side collections
						wg.Done()
					}()
				}
			}
			wg.Wait()
		}
	}()
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	LogInfo.Printf("Clearing active Client list...")
	ActiveClients.Clear()
	LogInfo.Println("done")
	LogInfo.Println("Stopping server...")
	kill <- struct{}{}
}
