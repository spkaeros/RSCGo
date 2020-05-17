package main

import (
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/db"
	"github.com/spkaeros/rscgo/pkg/definitions"
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
	// Parse CLI arguments
	if _, err := flags.Parse(&Flags); err != nil {
		os.Exit(100)
	}
	// Bounds check for TCP port
	if Flags.Port >= 65534 || Flags.Port < 0 {
		log.Fatal("Invalid port number specified.  Valid port numbers are 1-65533 inclusive.  Got:", Flags.Port)
		os.Exit(100)
	}
	if len(Flags.Config) == 0 {
		Flags.Config = "config.toml"
	}

	// Initialize sane defaults as fallback configuration options, if the config.toml file is not found or if some values are left out of it
	config.TomlConfig.MaxPlayers = 1250
	config.TomlConfig.DataDir = "./data/"
	config.TomlConfig.PacketHandlerFile = "packets.toml"
	config.TomlConfig.Crypto.HashComplexity = 15
	config.TomlConfig.Crypto.HashLength = 32
	config.TomlConfig.Crypto.HashMemory = 8
	config.TomlConfig.Crypto.HashSalt = "rscgo./GOLANG!RULES/.1994"
	config.TomlConfig.Version = 204
	config.TomlConfig.Port = 43594 // = 43595 for websocket connections

	if _, err := toml.DecodeFile(Flags.Config, &config.TomlConfig); err != nil {
		log.Warn("Error reading general config file:", err)
		os.Exit(101)
	}

	// TODO: Default serialization to JSON or BSON maybe?
	config.TomlConfig.Database.PlayerDriver = "sqlite3"
	config.TomlConfig.Database.PlayerDB = "file:./data/players.db"
	// TODO: Default serialization to JSON or BSON maybe?
	config.TomlConfig.Database.WorldDriver = "sqlite3"
	config.TomlConfig.Database.WorldDB = "file:./data/world.db"

	if _, err := toml.DecodeFile("data/dbio.conf", &config.TomlConfig.Database); err != nil {
		log.Warn("Error reading database config file:", err)
		os.Exit(102)
	}
	db.ConnectEntityService()
	db.DefaultPlayerService = db.NewPlayerServiceSql()
}

type asInit struct {
	sync.WaitGroup
}

func newRunner() asInit {
	return asInit{sync.WaitGroup{}}
}

func (w *asInit) runAll(fn ...func()) {
	for _, v := range fn {
		w.run(v)
	}
	w.Wait()
}

func (w *asInit) run(fn func()) {
	w.Add(1)
	go func() {
		defer w.Done()
		fn()
	}()
}

func (w *asInit) executeAll() {
	w.WaitGroup.Wait()
	return
}

func main() {
	//	if Flags.UseCipher {
	//		log.Debug("TODO: Figure out why ISAAC cipher sometimes works, yet eventually desynchronizes from client's ISAAC stream.")
	//		log.Debug("This cipher will remain disabled until such time as this issue gets resolved.  Possibly to be replaced by full stream encryption anyways")
	//		Flags.UseCipher = false
	//	}

	config.Verbosity = int(math.Min(math.Max(float64(len(Flags.Verbose)), 0), 4))
	if Flags.Port > 0 {
		config.TomlConfig.Port = Flags.Port
	}
	log.Debugln("RSCGo starting up...")
	log.Debugln()

	start := time.Now()
	runner := newRunner()
	// Entity definitions
	runner.runAll(handlers.UnmarshalPackets, db.LoadTileDefinitions, db.LoadObjectDefinitions, db.LoadBoundaryDefinitions, db.LoadItemDefinitions, db.LoadNpcDefinitions)
	// Entity locations
	runner.runAll(db.LoadObjectLocations, db.LoadNpcLocations, db.LoadItemLocations, world.RunScripts, world.LoadCollisionData)


	// Network protocol information
	if config.Verbose() {
		log.Debugln("Loaded " + strconv.Itoa(len(world.Sectors)-1) + " map sectors")
		log.Debugln("Loaded " + strconv.Itoa(handlers.PacketCount()) + " (" + strconv.Itoa(handlers.HandlerCount()) + ") packets (handlers)")
		log.Debugln("Loaded " + strconv.Itoa(int(world.ItemIndexer.Load())) + " items and " + strconv.Itoa(len(definitions.Items)-1) + " item definitions")
		log.Debugln("Loaded " + strconv.Itoa(world.Npcs.Size()) + " NPCs and " + strconv.Itoa(len(definitions.Npcs)-1) + " NPC definitions")
		log.Debugln("Loaded " + strconv.Itoa(len(definitions.ScenaryObjects)-1) + " scenary definitions, and " + strconv.Itoa(len(definitions.BoundaryObjects)-1) + " boundary definitions")
		log.Debugln("Loaded " + strconv.Itoa(int(world.ObjectCounter.Load())) + " scenary / boundary objects")
		log.Debugln("Finished loading game data; took " + strconv.Itoa(int(time.Since(start).Milliseconds())) + "ms")
		if config.Verbosity >= 2 {
			log.Debugf("Triggers[\n\t%d item actions,\n\t%d scenary actions,\n\t%d boundary actions,\n\t%d npc actions,\n\t%d item->boundary actions,\n\t%d item->scenary actions,\n\t%d attacking NPC actions,\n\t%d killing NPC actions\n];\n", len(world.ItemTriggers), len(world.ObjectTriggers), len(world.BoundaryTriggers), len(world.NpcTriggers), len(world.InvOnBoundaryTriggers), len(world.InvOnObjectTriggers), len(world.NpcAtkTriggers), len(world.NpcDeathTriggers))
		}
	}
	engine.Bind(config.Port())
	engine.Bind(config.Port() + 1)
	log.Debugf("Listening at TCP port %d, and TCP websockets port %d on all addresses.\n", config.Port(), config.Port()+1)
	log.Debugln()
	log.Debugln("RSCGo has finished initializing world; we hope you enjoy it")
	engine.StartGameEngine()
}
