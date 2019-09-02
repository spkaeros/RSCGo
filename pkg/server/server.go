package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/list"
)

var (
	//LogWarning Log interface for warnings.
	LogWarning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)
	//LogInfo Log interface for debug information.
	LogInfo = log.New(os.Stdout, "[INFO] ", log.Ltime|log.Lshortfile)
	//LogError Log interface for errors.
	LogError = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
	kill     = make(chan struct{})
	//ClientList List of active clients.
	ClientList = list.New(2048)
	//Clients A map of base37 encoded username hashes to client references.  This is a common lookup and I consider this an optimization.
	Clients = make(map[uint64]*Client)
)

//TomlConfig A data structure representing the RSCGo TOML configuration file.
var TomlConfig struct {
	DataDir           string `toml:"data_directory"`
	Version           int    `toml:"version"`
	Port              int    `toml:"port"`
	MaxPlayers        int    `toml:"max_players"`
	PacketHandlerFile string `toml:"packet_handler_table"`
	Database          struct {
		PlayerDB string `toml:"player_db"`
		WorldDB  string `toml:"world_db"`
	} `toml:"database"`
	Crypto struct {
		RsaKeyFile string `toml:"rsa_key"`
		HashSalt   string `toml:"hash_salt"`
	} `toml:"crypto"`
}

//Flags This is used to interface with the go-flags package from some guy on github.
var Flags struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port      int    `short:"p" long:"port" description:"The port for the server to listen on,"`
	Config    string `short:"c" long:"config" description:"Specify the configuration file to load server settings from" default:"config.toml"`
	UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using ISAAC to encrypt packet opcodes."`
}

func startConnectionService() {
	port := TomlConfig.Port
	if Flags.Port > 0 {
		port = Flags.Port
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		LogError.Printf("Can't bind to specified port: %d\n", port)
		LogError.Println(err)
		os.Exit(1)
	}

	go func() {
		// One client every 10ms max, stops childish flooding.
		// TODO: Implement a packet filter of sorts to stop flooding behavior
		defer listener.Close()
		for range time.Tick(10 * time.Millisecond) {
			socket, err := listener.Accept()
			if err != nil {
				if len(Flags.Verbose) > 0 {
					LogError.Println("Error occurred attempting to accept a client:", err)
				}
				continue
			}
			if len(Clients) >= TomlConfig.MaxPlayers {
				if n, err := socket.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 14}); err != nil || n != 9 {
					if len(Flags.Verbose) > 0 {
						LogError.Println("Could not send world is full response to rejected client:", err)
					}
				}
				continue
			}

			client := NewClient(socket)
			client.Index = ClientList.Add(client)
		}
	}()

}

//Start Listens for and processes new clients connecting to the server.
// This method blocks while the server is running.
func Start() {
	LogInfo.Println("RSCGo starting up...")
	startConnectionService()
	if len(Flags.Verbose) > 0 {
		LogInfo.Println()
		LogInfo.Println("Launched connection service.")
	}
	startGameEngine()
	if len(Flags.Verbose) > 0 {
		LogInfo.Println("Launched game engine.")
	}
	if ok := LoadObjects(); len(Flags.Verbose) > 0 && ok {
		LogInfo.Printf("Loaded %d game objects.\n", Objects.Size())
	}
	LogInfo.Println()
	LogInfo.Println("RSCGo is now running.")
	port := TomlConfig.Port
	if Flags.Port > 0 {
		port = Flags.Port
	}
	LogInfo.Printf("Listening on port %d...\n", port)
	select {
	// TODO: Probably need to handle certain signals
	// TODO: Any other tasks I should handle in the main goroutine??
	case <-kill:
		os.Exit(0)
	}
}

//UpdateMobileEntities Updates all mobile scene entities that are traversing a path
func UpdateMobileEntities() {
	var wg sync.WaitGroup
	wg.Add(len(Clients))
	for _, c := range Clients {
		if c != nil {
			go func() {
				defer wg.Done()
				c.player.TraversePath()
			}()
		}
	}
	wg.Wait()
}

//UpdateClientState Sends the new positions to the clients
func UpdateClientState() {
	var wg sync.WaitGroup
	wg.Add(len(Clients))
	for _, c := range Clients {
		if c != nil {
			go func() {
				defer wg.Done()
				c.UpdatePositions()
			}()
		}
	}
	wg.Wait()
}

//ResetUpdateFlags Resets the variables used for client updating synchronization.
func ResetUpdateFlags() {
	var wg sync.WaitGroup
	wg.Add(len(Clients))
	for _, c := range Clients {
		if c != nil {
			go func() {
				defer wg.Done()
				c.ResetUpdateFlags()
			}()
		}
	}
	wg.Wait()
}

//Tick One game engine 'tick'.  This is to handle movement, to synchronize clients, to update movement-related state variables...
// Runs every 640ms.
func Tick() {
	UpdateMobileEntities()
	UpdateClientState()
	ResetUpdateFlags()
}

//startGameEngine Launches a goroutine to handle updating the state of the server every 640ms in a
// synchronized fashion.  This is known as a single game engine 'pulse'.  All mobile entities must have their position
// updated during this pulse to be compatible with Jagex RSClassic Client software.
// TODO: Can movement be handled concurrently per-player safely on the Jagex Client? Mob movement might not look right.
func startGameEngine() {
	go func() {
		for range time.Tick(640 * time.Millisecond) {
			Tick()
		}
	}()
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	LogInfo.Printf("Clearing client list...")
	ClientList.Clear()
	Clients = make(map[uint64]*Client)
	LogInfo.Println("done")
	LogInfo.Println("Stopping server...")
	kill <- struct{}{}
}
