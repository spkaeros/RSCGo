package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	listener net.Listener
)

var (
	syncTicker = time.NewTicker(time.Millisecond * 600)
	kill       = make(chan struct{})
)

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
		fmt.Println("WARNING: Attempted to start connection service without a listener!  This shouldn't happen.")
		fmt.Println("Starting listener on default port...")
		bind(43591)
	}

	go func() {
		defer listener.Close()
		// TODO: Can this ticker be made smaller safely?
		connTicker := time.NewTicker(time.Millisecond * 50)
		for range connTicker.C {
			socket, err := listener.Accept()
			if err != nil {
				fmt.Println("ERROR: Could not accept client from server listener.")
				fmt.Println(err)
				return
			}

			client := newClient(socket)
			ActiveClients.Add(client)
			fmt.Println("Registered client" + client.String())
		}
	}()

}

//Start Listens for and processes new clients connecting to the server.
// This method blocks while the server is running.
func Start(port int) {
	fmt.Printf("RSCGo starting up...\n\n")
	fmt.Print("Attempting to bind to network...")
	bind(port)
	fmt.Println("done")
	fmt.Print("Attempting to start connection service...")
	startConnectionService()
	fmt.Println("done")
	fmt.Print("Attempting to start synchronized task service...")
	startSynchronizedTaskService()
	fmt.Println("done")
	fmt.Printf("\nRSCGo is now running.\nListening on port %d...\n", port)
	select {
	case <-kill:
		os.Exit(0)
	}
}

//startSynchronizedTaskService Launches a goroutine to handle updating the state of the server every 600ms in a
// synchronized fashion.  This is known as a single game engine 'pulse'.  All mobile entities must have their position
// updated during this pulse to be compatible with Jagex RSClassic client software.
// TODO: Can movement be handled concurrently per-player safely on the Jagex client? Mob movement might not look right.
func startSynchronizedTaskService() {
	go func() {
		for range syncTicker.C {
		}
	}()
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	fmt.Print("Clearing active client list...")
	ActiveClients.Clear()
	fmt.Println("done")
	fmt.Println("Stopping server...")
	kill <- struct{}{}
}
