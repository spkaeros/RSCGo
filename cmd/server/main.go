package main

import (
	"bitbucket.org/zlacki/rscgo/pkg/server"
	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packethandlers"
	"bitbucket.org/zlacki/rscgo/pkg/server/script"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
	"sync"
)

//Flags This is used to interface with the go-flags package from some guy on github.
var Flags struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Display more verbose output"`
	Port      int    `short:"p" long:"port" description:"The port for the server to listen on,"`
	Config    string `short:"c" long:"config" description:"Specify the configuration file to load server settings from" default:"config.toml"`
	UseCipher bool   `short:"e" long:"encryption" description:"Enable command opcode encryption using ISAAC to encrypt packet opcodes."`
}

//asyncExecute First this will add 1 task to the specified waitgroup, then it will execute the function fn in its own goroutine, and upon exiting this goroutine will indicate to wg that the task we added is finished.
func asyncExecute(wg *sync.WaitGroup, fn func()) {
	(*wg).Add(1)
	go func() {
		defer (*wg).Done()
		fn()
	}()
}

func main() {
	log.Info.Println("RSCGo starting up...")
	if _, err := flags.Parse(&Flags); err != nil {
		log.Error.Println(err)
		return
	}
	config.Verbosity = len(Flags.Verbose)
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
	if Flags.Port > 0 {
		config.TomlConfig.Port = Flags.Port
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
		if count := db.LoadNpcLocations(); len(Flags.Verbose) > 0 && count > 0 {
			log.Info.Printf("Loaded %d NPCs.\n", count)
		}
	})
	asyncExecute(&awaitLaunchJobs, func() {
		script.LoadObjectTriggers()
		log.Info.Printf("Loaded %d object action triggers.\n", len(script.ObjectTriggers))
	})
	asyncExecute(&awaitLaunchJobs, func() {
		script.LoadItemTriggers()
		log.Info.Printf("Loaded %d item action triggers.\n", len(script.ItemTriggers))
	})
	asyncExecute(&awaitLaunchJobs, func() {
		script.LoadBoundaryTriggers()
		log.Info.Printf("Loaded %d boundary action triggers.\n", len(script.BoundaryTriggers))
	})
	asyncExecute(&awaitLaunchJobs, func() {
		packethandlers.Initialize()
		if len(Flags.Verbose) > 0 {
			log.Info.Printf("Loaded %d packet handlers.\n", packethandlers.CountPacketHandlers())
		}
	})
	// TODO: Re-enable RSA
	/*
		asyncExecute(&awaitLaunchJobs, func() {
			loadRsaKey()
			if len(Flags.Verbosity) > 0 {
				log.Info.Println("Loaded RSA key data.")
			}
		})
	*/
	awaitLaunchJobs.Wait()
	server.StartGameEngine()
	if len(Flags.Verbose) > 0 {
		log.Info.Println("Launched game engine.")
	}
	log.Info.Println()
	log.Info.Println("RSCGo is now running.")
	server.StartConnectionService()
	log.Info.Printf("Listening on TCP port %d, websocket port %d...\n", config.Port(), config.Port()+1)
	select {
	case <-server.Kill:
		os.Exit(0)
	}
}
