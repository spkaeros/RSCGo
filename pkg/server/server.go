package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

var (
	listener      net.Listener
	LogWarning    = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)
	LogInfo       = log.New(os.Stdout, "[INFO] ", log.Ltime|log.Lshortfile)
	LogError      = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
	syncTicker    = time.NewTicker(time.Millisecond * 600)
	kill          = make(chan struct{})
	Version       = -1
	DataDirectory = "."
)

func init() {

}

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

func LogDebug(lvl int, s string, args ...interface{}) {
	if len(Flags.Verbose) > lvl {
		LogInfo.Printf(s, args...)
	}
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

//startSynchronizedTaskService Launches a goroutine to handle updating the state of the server every 600ms in a
// synchronized fashion.  This is known as a single game engine 'pulse'.  All mobile entities must have their position
// updated during this pulse to be compatible with Jagex RSClassic Client software.
// TODO: Can movement be handled concurrently per-player safely on the Jagex Client? Mob movement might not look right.
func startSynchronizedTaskService() {
	go func() {
		for range syncTicker.C {
			for _, c := range ActiveClients.clients {
				if c, ok := c.(*Client); ok {
					p := packets.NewOutgoingPacket(145)
					p.AddBits(220, 11)
					p.AddBits(445, 13)
					p.AddBits(0, 4)
					p.AddBits(0, 8)
					c.WritePacket(p)
					// TODO: Update movement, update client-side collections
				}
			}
		}
	}()
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	LogDebug(0, "Clearing active Client list...")
	ActiveClients.Clear()
	LogDebug(0, "done\n")
	fmt.Println("Stopping server...")
	kill <- struct{}{}
}
