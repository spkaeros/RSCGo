package main

import (
	"os"
	"time"
	"math"
	"strconv"

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

}
func startAsync(fns ...func()) {
//	wg := sync.WaitGroup{}
//	wg.Add(len(fns))
//	(*wg1).Add(1)
//	log.Debug(fns)
	for _, v := range fns {
//		go func() {
//			defer wg.Done()	
			v()
//		}()
	}
//	(*wg1).Done()
}

func main() {
	if _, err := flags.Parse(&Flags); err != nil {
		os.Exit(100)
	}
	if Flags.Port > 65535 || Flags.Port < 0 {
		log.Fatal("Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		os.Exit(101)
	}
	if len(Flags.Config) == 0 {
		Flags.Config = "config.toml"
	}
//	if Flags.UseCipher {
//		log.Debug("TODO: Figure out why ISAAC cipher sometimes works, yet eventually desynchronizes from client's ISAAC stream.")
//		log.Debug("This cipher will remain disabled until such time as this issue gets resolved.  Possibly to be replaced by full stream encryption anyways")
//		Flags.UseCipher = false
//	}
	if _, err := toml.DecodeFile("./data/dbio.conf", &config.TomlConfig.Database); err != nil {
		log.Warn("Error reading database config file:", err)
		os.Exit(110)
	}
	db.DefaultPlayerService = db.NewPlayerServiceSql()
	db.ConnectEntityService()
	if _, err := toml.DecodeFile(Flags.Config, &config.TomlConfig); err != nil {
		log.Warn("Error reading general config file:", err)
		os.Exit(102)
	}
	config.Verbosity = int(math.Min(math.Max(float64(len(Flags.Verbose)), 0), 4))
	if Flags.Port > 0 {
		config.TomlConfig.Port = Flags.Port
	}
	log.Debugln("RSCGo starting up...")
	log.Debugln()

	start := time.Now()

//	var wg sync.WaitGroup
//	wg.Add(1)
//	go func() {
		// Running these init functions that are I/O heavy and synchronization between them is never important within
		db.LoadTileDefinitions()
		//  their own goroutines should save some initialization time.
		db.LoadObjectDefinitions()
		db.LoadBoundaryDefinitions()
		db.LoadNpcDefinitions()
		db.LoadItemDefinitions()
//		wg.Done()
//	}()
//	wg.Wait()
//	phaseEverything := sync.WaitGroup{}
	/*phaseInstances := */
//	(*startAsync(db.LoadObjectLocations, db.LoadItemLocations, handlers.UnmarshalPackets)).Wait()
	db.LoadObjectLocations()
	db.LoadNpcLocations()
	db.LoadItemLocations()
	handlers.UnmarshalPackets()
	world.LoadCollisionData()
//	phaseInstances.Wait()
//	phaseEverything.Wait()
	world.RunScripts()
	// Entity definitions

	// Entity locations

	// Network protocol information
	if config.Verbose() {
		log.Debugln("Loaded " + strconv.Itoa(len(world.Sectors)-1) + " map sectors")
		log.Debugln("Loaded " + strconv.Itoa(handlers.PacketCount()) + " (" + strconv.Itoa(handlers.HandlerCount()) + ") packets (handlers)")
		log.Debugln("Loaded " + strconv.Itoa(int(world.ItemIndexer.Load())) + " items and " + strconv.Itoa(len(world.ItemDefs)-1) + " item definitions")
		log.Debugln("Loaded " + strconv.Itoa(world.Npcs.Size()) + " NPCs and " + strconv.Itoa(len(world.NpcDefs)-1) + " NPC definitions")
		log.Debugln("Loaded " + strconv.Itoa(len(world.ObjectDefs)-1) + " scenary definitions, and " + strconv.Itoa(len(world.BoundaryDefs)-1) + " boundary definitions")
		log.Debugln("Loaded " + strconv.Itoa(int(world.ObjectCounter.Load())) + " scenary / boundary objects")
		log.Debugf("Triggers[\n\t%d item actions,\n\t%d scenary actions,\n\t%d boundary actions,\n\t%d npc actions,\n\t%d item->boundary actions,\n\t%d item->scenary actions,\n\t%d attacking NPC actions,\n\t%d killing NPC actions\n];\n", len(world.ItemTriggers), len(world.ObjectTriggers), len(world.BoundaryTriggers), len(world.NpcTriggers), len(world.InvOnBoundaryTriggers), len(world.InvOnObjectTriggers), len(world.NpcAtkTriggers), len(world.NpcDeathTriggers))
		log.Debugln("Finished loading entitys; took " + strconv.Itoa(int(time.Since(start).Milliseconds())) + "ms")
	}
	engine.Bind(config.Port())
	engine.Bind(config.Port()+1)
	log.Debugf("Listening at TCP port %d, and TCP websockets port %d on all addresses.\n", config.Port(), config.Port()+1)
	log.Debugln()
	log.Debugln("RSCGo has finished initializing world; we hope you enjoy it")
	engine.StartGameEngine()
}