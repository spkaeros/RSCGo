package server

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
)

var (
	kill = make(chan struct{})
)

//ClientMap A thread-safe concurrent collection type for storing client references.
type ClientMap struct {
	usernames map[uint64]*Client
	indices   map[int]*Client
	lock      sync.RWMutex
}

//Clients Collection containing all of the active clients, by index and username hash, guarded by a mutex
var Clients = &ClientMap{usernames: make(map[uint64]*Client), indices: make(map[int]*Client)}

//FromUserHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func (m *ClientMap) FromUserHash(hash uint64) (*Client, bool) {
	m.lock.RLock()
	result, ok := m.usernames[hash]
	m.lock.RUnlock()
	return result, ok
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func (m *ClientMap) ContainsHash(hash uint64) bool {
	_, ret := m.FromUserHash(hash)
	return ret
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func (m *ClientMap) FromIndex(index int) (*Client, bool) {
	m.lock.RLock()
	result, ok := m.indices[index]
	m.lock.RUnlock()
	return result, ok
}

//Put Puts a client into the map.
func (m *ClientMap) Put(c *Client) {
	nextIndex := m.NextIndex()
	m.lock.Lock()
	c.Index = nextIndex
	c.player.Index = nextIndex
	m.usernames[c.player.UserBase37] = c
	m.indices[nextIndex] = c
	m.lock.Unlock()
}

//Remove Removes a client from the map.
func (m *ClientMap) Remove(c *Client) {
	m.lock.Lock()
	delete(m.usernames, c.player.UserBase37)
	delete(m.indices, c.Index)
	m.lock.Unlock()
}

//Broadcast Calls action for every active client in the collection.
func (m *ClientMap) Broadcast(action func(*Client)) {
	m.lock.RLock()
	for _, c := range m.indices {
		if c != nil && c.player.TransAttrs.VarBool("connected", false) {
			action(c)
		}
	}
	m.lock.RUnlock()
}

//Size Returns the size of the active client collection.
func (m *ClientMap) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.usernames)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func (m *ClientMap) NextIndex() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for i := 0; i < config.MaxPlayers(); i++ {
		if _, ok := m.indices[i]; !ok {
			return i
		}
	}
	return -1
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
		config.TomlConfig.Port = Flags.Port
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port()))
	if err != nil {
		log.Error.Printf("Can't bind to specified port: %d\n", config.Port())
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
			if Clients.Size() >= config.MaxPlayers() {
				if n, err := socket.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 14}); err != nil || n != 9 {
					if len(Flags.Verbose) > 0 {
						log.Error.Println("Could not send world is full response to rejected client:", err)
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
		initPacketHandlerTable()
		if len(Flags.Verbose) > 0 {
			log.Info.Printf("Initialized %d packet handlers.\n", len(table.Handlers))
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
		startConnectionService()
		if len(Flags.Verbose) > 0 {
			log.Info.Println("Launched connection service.")
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		startGameEngine()
		if len(Flags.Verbose) > 0 {
			log.Info.Println("Launched game engine.")
		}
	})
	awaitLaunchJobs.Wait()
	log.Info.Println()
	log.Info.Println("RSCGo is now running.")
	log.Info.Printf("Listening on port %d...\n", config.Port())
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

//Tick One game engine 'tick'.  This is to handle movement, to synchronize clients, to update movement-related state variables... Runs once per 600ms.
func Tick() {
	Clients.Broadcast(func(c *Client) {
		//TODO: Handle this in a less hacky way.  Sticks out like a sore thumb.
		if c.player.IsFollowing() {
			followingClient, ok := Clients.FromIndex(c.player.FollowIndex())
			if followingClient == nil || !ok || !c.player.Location.WithinRange(followingClient.player.Location, 15) {
				c.player.ResetFollowing()
			} else if !c.player.FinishedPath() && c.player.WithinRange(followingClient.player.Location, 2) {
				c.player.ResetPath()
			} else if c.player.FinishedPath() && !c.player.WithinRange(followingClient.player.Location, 2) {
				c.player.SetPath(world.NewPathwayFromLocation(&followingClient.player.Location))
			}
		}
		var tmpFns []func() bool
		c.player.ActionLock.Lock()
		for _, fn := range c.player.DistancedActions {
			if !fn() {
				tmpFns = append(tmpFns, fn)
			}
		}
		c.player.DistancedActions = tmpFns
		c.player.ActionLock.Unlock()
		c.player.TraversePath()
	})
	Clients.Broadcast(func(c *Client) {
		c.UpdatePositions()
	})
	Clients.Broadcast(func(c *Client) {
		c.ResetUpdateFlags()
		world.ResetNpcUpdateFlags()
	})
}

//startGameEngine Launches a goroutine to handle updating the state of the server every 600ms in a synchronized fashion.  This is known as a single game engine 'pulse'.
func startGameEngine() {
	go func() {
		ticker := time.NewTicker(600*time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			Tick()
		}
	}()
}

//BroadcastLogin Broadcasts the login status of player to the whole server.
func BroadcastLogin(player *world.Player, online bool) {
	Clients.Broadcast(func(c *Client) {
		if c.player.Friends(player.UserBase37) {
			if !player.FriendBlocked() || player.Friends(c.player.UserBase37) {
				c.player.FriendList[player.UserBase37] = online
				c.outgoingPackets <- packets.FriendUpdate(player.UserBase37, online)
			}
		}
	})
}

//Stop This will stop the server instance, if it is running.
func Stop() {
	log.Info.Println("Stopping server...")
	kill <- struct{}{}
}
