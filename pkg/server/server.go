package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
)

var (
	//LogWarning Log interface for warnings.
	LogWarning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)
	//LogInfo Log interface for debug information.
	LogInfo = log.New(os.Stdout, "[INFO] ", log.Ltime|log.Lshortfile)
	//LogError Log interface for errors.
	LogError = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
	kill     = make(chan struct{})
	//Clients A map of base37 encoded username hashes to client references.  This is a common lookup and I consider this an optimization.
	Clients = make(map[uint64]*Client)
	//ClientsIdx A map of server indexes to client references.  Also a common lookup.
	ClientsIdx = make(map[int]*Client)
	// TODO: Combine the two collections above to one custom collection type.
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
	if Flags.Port > 0 {
		TomlConfig.Port = Flags.Port
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", TomlConfig.Port))
	if err != nil {
		LogError.Printf("Can't bind to specified port: %d\n", TomlConfig.Port)
		LogError.Println(err)
		os.Exit(1)
	}

	go func() {
		// TODO: Implement a packet filter of sorts to stop flooding behavior
		defer listener.Close()
		for range time.Tick(50 * time.Millisecond) {
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

			NewClient(socket)
		}
	}()

}

//Start Listens for and processes new clients connecting to the server.
// This method blocks while the server is running.
func Start() {
	LogInfo.Println("RSCGo starting up...")
	if _, err := flags.Parse(&Flags); err != nil {
		LogError.Println(err)
		return
	}
	if Flags.Port > 65535 || Flags.Port < 0 {
		LogWarning.Println("Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		return
	}
	if !strings.HasSuffix(Flags.Config, ".toml") {
		LogWarning.Println("You entered an invalid configuration file extension.")
		LogWarning.Println("TOML is currently the only supported format for server properties.")
		LogWarning.Println()
		LogInfo.Println("Setting back to default: `config.toml`")
		Flags.Config = "config.toml"
	}
	if Flags.UseCipher {
		LogInfo.Println("TODO: Figure out why ISAAC cipher sometimes works, yet eventually desynchronizes from client.")
		LogInfo.Println("Cipher will remain disabled until such time as this issue gets resolved.")
		Flags.UseCipher = false
	}

	if _, err := toml.DecodeFile("."+string(os.PathSeparator)+Flags.Config, &TomlConfig); err != nil {
		LogWarning.Println("Error decoding TOML RSCGo general configuration file:", err)
		return
	}
	if len(Flags.Verbose) > 0 {
		LogInfo.Println()
		LogInfo.Println("Loaded TOML configuration file.")
	}

	var awaitLaunchJobs sync.WaitGroup
	awaitLaunchJobs.Add(6)
	asyncExecute(&awaitLaunchJobs, func() {
		LoadObjectDefinitions()
		if count := len(ObjectDefinitions); len(Flags.Verbose) > 0 && count > 0 {
			LogInfo.Printf("Loaded %d game object definitions.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		if count := LoadObjects(); len(Flags.Verbose) > 0 && count > 0 {
			LogInfo.Printf("Loaded %d game objects.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		initPacketHandlerTable()
		if len(Flags.Verbose) > 0 {
			LogInfo.Printf("Initialized %d packet handlers.\n", len(table.Handlers))
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		initCrypto()
		if len(Flags.Verbose) > 0 {
			LogInfo.Println("Launched cryptographic subsystem.")
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		startConnectionService()
		if len(Flags.Verbose) > 0 {
			LogInfo.Println("Launched connection service.")
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		startGameEngine()
		if len(Flags.Verbose) > 0 {
			LogInfo.Println("Launched game engine.")
		}
	})
	awaitLaunchJobs.Wait()
	LogInfo.Println()
	LogInfo.Println("RSCGo is now running.")
	LogInfo.Printf("Listening on port %d...\n", TomlConfig.Port)
	select {
	// TODO: Probably need to handle certain signals
	// TODO: Any other tasks I should handle in the main goroutine??
	case <-kill:
		os.Exit(0)
	}
}

func asyncExecute(wg *sync.WaitGroup, fn func()) {
	go func() {
		defer (*wg).Done()
		fn()
	}()
}

//Broadcast Call action passing in every active client to perform a task on everyone playing.
func Broadcast(action func(c *Client)) {
	for _, c := range Clients {
		if c != nil && c.player.Connected {
			action(c)
		}
	}
}

//Tick One game engine 'tick'.  This is to handle movement, to synchronize clients, to update movement-related state variables...
// Runs every 600ms.
func Tick() {
	//	UpdateMobileEntities()
	Broadcast(func(c *Client) {
		if c.player.IsFollowing() {
			followingClient := ClientFromIndex(c.player.FollowIndex())
			if followingClient == nil || !c.player.Location.WithinRange(followingClient.player.Location, 15) {
				c.player.ResetFollowing()
			} else if !c.player.FinishedPath() && c.player.WithinRange(followingClient.player.Location, 2) {
				c.player.ResetPath()
			} else if c.player.FinishedPath() && !c.player.WithinRange(followingClient.player.Location, 2) {
				c.player.SetPath(entity.NewPathway(followingClient.player.X, followingClient.player.Y))
			}
		}
		c.player.TraversePath()
	})
	//	UpdateClientState()
	Broadcast(func(c *Client) {
		c.UpdatePositions()
	})
	//	ResetUpdateFlags()
	Broadcast(func(c *Client) {
		c.ResetUpdateFlags()
	})
}

//startGameEngine Launches a goroutine to handle updating the state of the server every 600ms in a
// synchronized fashion.  This is known as a single game engine 'pulse'.  All mobile entities must have their position
// updated during this pulse to be compatible with Jagex RSClassic Client software.
// TODO: Can movement be handled concurrently per-player safely on the Jagex Client? Mob movement might not look right.
func startGameEngine() {
	go func() {
		for range time.Tick(600 * time.Millisecond) {
			Tick()
		}
	}()
}

//BroadcastLogin Broadcasts the login status of the user with hash as their base37 username
func BroadcastLogin(player *entity.Player, online bool) {
	Broadcast(func(c *Client) {
		if c.player.Friends(player.UserBase37) {
			if !player.FriendBlocked() || player.Friends(c.player.UserBase37) {
				c.outgoingPackets <- packets.FriendUpdate(player.UserBase37, online)
			}
		}
	})
}

//ClientFromIndex Helper function to find a specific client reference from its assigned server index.  If there is no player with the index, returns nil.
func ClientFromIndex(index int) *Client {
	if c, ok := ClientsIdx[index]; c != nil && c.player.Connected && ok {
		return c
	}

	return nil
}

//ClientFromHash Helper function to find a specific client reference from its base37 encoded username.  If there is no player with the username, returns nil.
func ClientFromHash(userHash uint64) *Client {
	if c, ok := Clients[userHash]; c != nil && c.player.Connected && ok {
		return c
	}

	return nil
}

//ClientFromUsername Helper function to find a specific client reference from its username.  If there is no player with the username, returns nil.
func ClientFromUsername(username string) *Client {
	return ClientFromHash(strutil.Base37(username))
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	LogInfo.Printf("Clearing client list...")
	Clients = make(map[uint64]*Client)
	ClientsIdx = make(map[int]*Client)
	LogInfo.Println("done")
	LogInfo.Println("Stopping server...")
	kill <- struct{}{}
}
