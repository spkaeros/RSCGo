package main

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/engine"
	"github.com/spkaeros/rscgo/pkg/game/net/handlers"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

//Flags This is used to interface with the go-flags package from some guy on github.
var Flags struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port      int    `short:"p" long:"port" description:"The TCP port for the game to listen on, (Websocket will use the port directly above it)"`
	Config    string `short:"c" long:"config" description:"Specify the TOML configuration file to load game settings from" default:"config.toml"`
	UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using a variant of ISAAC to encrypt net opcodes."`
}

//asyncExecute First this will add 1 task to the specified waitgroup, then it will execute the function fn in its own goroutine, and upon exiting this goroutine will indicate to wg that the task we added is finished.
func asyncExecute(wg *sync.WaitGroup, fn func()) {
	(*wg).Add(1)
	go func() {
		defer (*wg).Done()
		fn()
	}()
}

func init() {
	config.TomlConfig.MaxPlayers = 1250
	config.TomlConfig.DataDir = "./data/"
	config.TomlConfig.Database.PlayerDriver = "sqlite3"
	config.TomlConfig.Database.PlayerDB = "file:./data/players.db"
	config.TomlConfig.Database.WorldDriver = "sqlite3"
	config.TomlConfig.Database.WorldDB = "file:./data/world.db"
	config.TomlConfig.PacketHandlerFile = "packets.toml"
	config.TomlConfig.Crypto.HashComplexity = 15
	config.TomlConfig.Crypto.HashLength = 32
	config.TomlConfig.Crypto.HashMemory = 8
	config.TomlConfig.Crypto.HashSalt = "rscgo./GOLANG!RULES/.1994"
	config.TomlConfig.Version = 204
	config.TomlConfig.Port = 43594 // = 43595 for websocket connections
	//TomlConfig.Crypto.RsaKeyFile = "rsa.der"

}

func main() {
	if _, err := flags.Parse(&Flags); err != nil {
		os.Exit(100)
	}
	if Flags.Port > 65535 || Flags.Port < 0 {
		log.Warning.Println("Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		os.Exit(101)
	}
	if !strings.HasSuffix(Flags.Config, ".toml") {
		log.Warning.Println("You entered an invalid configuration file extension.")
		log.Warning.Println("TOML is currently the only supported format for game properties.")
		log.Warning.Println()
		log.Info.Println("Setting back to default: `config.toml`")
		Flags.Config = "config.toml"
	}
	if Flags.UseCipher {
		log.Info.Println("TODO: Figure out why ISAAC cipher sometimes works, yet eventually desynchronizes from client's ISAAC stream.")
		log.Info.Println("Cipher will remain disabled until such time as this issue gets resolved.  Possibly to be replaced by full stream encryption anyways")
		Flags.UseCipher = false
	}
	if _, err := toml.DecodeFile("."+string(os.PathSeparator)+Flags.Config, &config.TomlConfig); err != nil {
		log.Warning.Println("Error decoding TOML RSCGo general configuration file:", err)
		os.Exit(102)
	}
	config.Verbosity = len(Flags.Verbose)
	if Flags.Port > 0 {
		config.TomlConfig.Port = Flags.Port
	}
	db.DefaultPlayerService = db.NewPlayerServiceSql()
	db.ConnectEntityService()
	log.Info.Println("RSCGo starting up...")
	log.Info.Println()

	start := time.Now()
	// Running these init functions that are I/O heavy and synchronization between them is never important within
	//  their own goroutines should save some initialization time.
	var awaitLaunchJobs sync.WaitGroup
	// Network protocol information
	asyncExecute(&awaitLaunchJobs, world.LoadCollisionData)
	asyncExecute(&awaitLaunchJobs, handlers.UnmarshalPackets)

	// Entity definitions
	asyncExecute(&awaitLaunchJobs, db.LoadObjectDefinitions)
	asyncExecute(&awaitLaunchJobs, db.LoadItemDefinitions)
	asyncExecute(&awaitLaunchJobs, db.LoadNpcDefinitions)
	asyncExecute(&awaitLaunchJobs, db.LoadBoundaryDefinitions)
	asyncExecute(&awaitLaunchJobs, db.LoadTileDefinitions)

	// Entity locations
	//	asyncExecute(&awaitLaunchJobs, db.LoadObjectLocations)
	// Entity action scripting triggers
	//	asyncExecute(&awaitLaunchJobs, script.LoadObjectTriggers)
	//	asyncExecute(&awaitLaunchJobs, script.LoadBoundaryTriggers)
	//	asyncExecute(&awaitLaunchJobs, script.LoadItemTriggers)
	awaitLaunchJobs.Wait()
	asyncExecute(&awaitLaunchJobs, db.LoadObjectLocations)
	asyncExecute(&awaitLaunchJobs, db.LoadNpcLocations)
	asyncExecute(&awaitLaunchJobs, db.LoadItemLocations)
	awaitLaunchJobs.Wait()
	world.RunScripts()

	if config.Verbose() {
		log.Info.Printf("Loaded %d landscape sectors.\n", len(world.Sectors))
		log.Info.Printf("Loaded %d packets, %d of which have handlers.\n", handlers.PacketCount(), handlers.HandlerCount())
		log.Info.Printf("Loaded %d item definitions.\n", len(world.ItemDefs))
		log.Info.Printf("Loaded %d NPC definitions.\n", len(world.NpcDefs))
		log.Info.Printf("Loaded %d object definitions.\n", len(world.ObjectDefs))
		log.Info.Printf("Loaded %d boundary definitions.\n", len(world.BoundaryDefs))
		log.Info.Printf("Loaded %d NPCs.\n", world.NpcCounter.Load())
		log.Info.Printf("Loaded %d ground items.\n", world.ItemIndexer.Load())
		log.Info.Printf("Loaded %d objects and boundaries.\n", world.ObjectCounter.Load())
		log.Info.Printf("Bind[%d item, %d obj, %d bound, %d npc, %d invBound, %d invObject, %d npcAtk, %d npcKill] loaded\n", len(world.ItemTriggers), len(world.ObjectTriggers), len(world.BoundaryTriggers), len(world.NpcTriggers), len(world.InvOnBoundaryTriggers), len(world.InvOnObjectTriggers), len(world.NpcAtkTriggers), len(world.NpcDeathTriggers))
		log.Info.Printf("Finished initializing entities in: %dms\n", time.Since(start).Milliseconds())
	}
	engine.StartGameEngine()
	if config.Verbose() {
		log.Info.Println("Launched game engine.")
		log.Info.Println()
	}
	engine.StartConnectionService()
	log.Info.Println("RSCGo is now running.")
	log.Info.Printf("Listening on TCP port %d, websocket port %d...\n", config.Port(), config.WSPort())
	select {
	case <-engine.Kill:
		os.Exit(0)
	}
}
