package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	listener      net.Listener
	syncTicker    = time.NewTicker(time.Millisecond * 600)
	kill          = make(chan struct{})
	Version       = -1
	DataDirectory = "."
)

var Flags struct {
	Verbose []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port    int    `short:"p" long:"port" description:"The port for the server to listen on," default:"43591"`
	Config  string `short:"c" long:"config" description:"Specify the configuration file to load server settings from" default:"config.ini"`
	UseCipher bool `short:"e" long:"encryption" description:"Enable command opcode encryption using ISAAC to encrypt packet opcodes."`
}

func bind(port int) {
	var err error
	portS := strconv.Itoa(port)
	listener, err = net.Listen("tcp", ":"+portS)
	if err != nil {
		fmt.Println("ERROR: Can't bind to specified port: " + portS)
		fmt.Println(err)
		os.Exit(1)
	}
}

func startConnectionService() {
	if listener == nil {
		if len(Flags.Verbose) > 0 {
			fmt.Println("WARNING: Attempted to start connection service without a listener!  This shouldn't happen.")
			fmt.Println("Starting listener on default port...")
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
					fmt.Println("ERROR: Could not accept Client from server listener.")
					fmt.Println(err)
				}
				return
			}

			client := NewClient(socket)
			ActiveClients.Add(client)
			fmt.Println("Registered Client" + client.String())
		}
	}()

}

func LogDebug(lvl int, s string, args ...interface{}) {
	if len(Flags.Verbose) > lvl {
		fmt.Printf(s, args...)
	}
}

//Start Listens for and processes new clients connecting to the server.
// This method blocks while the server is running.
func Start() {
	fmt.Printf("RSCGo starting up...\n\n")
	LogDebug(0, "Attempting to bind to network...")
	bind(Flags.Port)
	LogDebug(0, "done\nAttempting to start connection service...")
	startConnectionService()
	LogDebug(0, "done\nAttempting to start synchronized task service...")
	startSynchronizedTaskService()
	LogDebug(0, "done\n\n")
	fmt.Printf("RSCGo is now running.\nListening on port %d...\n", Flags.Port)
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
				if c != nil {
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
