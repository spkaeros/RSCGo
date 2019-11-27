package main

import (
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
	"github.com/spkaeros/rscgo/pkg/server"
	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packethandlers"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"os"
	"strings"
	"sync"
	"time"
)

//Flags This is used to interface with the go-flags package from some guy on github.
var Flags struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port      int    `short:"p" long:"port" description:"The TCP port for the server to listen on, (Websocket will use the port directly above it)"`
	Config    string `short:"c" long:"config" description:"Specify the TOML configuration file to load server settings from" default:"config.toml"`
	UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using a variant of ISAAC to encrypt packet opcodes."`
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
	if _, err := flags.Parse(&Flags); err != nil {
		os.Exit(100)
	}
	if Flags.Port > 65535 || Flags.Port < 0 {
		log.Warning.Println("Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		os.Exit(101)
	}
	if !strings.HasSuffix(Flags.Config, ".toml") {
		log.Warning.Println("You entered an invalid configuration file extension.")
		log.Warning.Println("TOML is currently the only supported format for server properties.")
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
}

func main() {
	log.Info.Println("RSCGo starting up...")
	log.Info.Println()

	start := time.Now()
	// Running these init functions that are I/O heavy and synchronization between them is never important within
	//  their own goroutines should save some initialization time.
	var awaitLaunchJobs sync.WaitGroup
	// Network protocol information
	asyncExecute(&awaitLaunchJobs, world.LoadMapData)
	asyncExecute(&awaitLaunchJobs, packethandlers.Initialize)

	// Entity definitions
	asyncExecute(&awaitLaunchJobs, db.LoadObjectDefinitions)
	asyncExecute(&awaitLaunchJobs, db.LoadItemDefinitions)
	asyncExecute(&awaitLaunchJobs, db.LoadEquipmentDefinitions)
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
	asyncExecute(&awaitLaunchJobs, script.Load)
	asyncExecute(&awaitLaunchJobs, db.LoadObjectLocations)
	asyncExecute(&awaitLaunchJobs, db.LoadNpcLocations)
	awaitLaunchJobs.Wait()
	if config.Verbose() {
		log.Info.Printf("Loaded %d landscape sectors.\n", len(world.Sectors))
		log.Info.Printf("Loaded %d packet handlers.\n", packethandlers.Size())
		log.Info.Printf("Loaded %d item definitions.\n", len(world.ItemDefs))
		log.Info.Printf("Loaded %d NPC definitions.\n", len(world.NpcDefs))
		log.Info.Printf("Loaded %d object definitions.\n", len(world.Objects))
		log.Info.Printf("Loaded %d boundary definitions.\n", len(world.Boundarys))
		log.Info.Printf("Loaded %d NPCs.\n", world.NpcCounter.Load())
		log.Info.Printf("Loaded %d objects and boundaries.\n", world.ObjectCounter.Load())
		log.Info.Printf("Loaded %d inventory, %d object, %d boundary, and %d NPC action triggers.\n", len(script.InvTriggers), len(script.ObjectTriggers), len(script.BoundaryTriggers), len(script.NpcTriggers))
		log.Info.Printf("Finished initializing entities in: %dms\n", time.Since(start).Milliseconds())
	}
	server.StartGameEngine()
	if config.Verbose() {
		log.Info.Println("Launched game engine.")
		log.Info.Println()
	}
	server.StartConnectionService()
	log.Info.Println("RSCGo is now running.")
	log.Info.Printf("Listening on TCP port %d, websocket port %d...\n", config.Port(), config.WSPort())
	select {
	case <-server.Kill:
		os.Exit(0)
	}
}
