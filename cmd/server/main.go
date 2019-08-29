package main

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/zlacki/rscgo/pkg/server"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
)

func init() {
	if _, err := flags.Parse(&server.Flags); err != nil {
		os.Exit(100)
	}
	if server.Flags.Port > 65535 || server.Flags.Port < 0 {
		fmt.Println("WARNING: Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		os.Exit(101)
	}
	if !strings.HasSuffix(server.Flags.Config, ".toml") {
		fmt.Println("WARNING: You entered an invalid configuration file extension.")
		fmt.Println("TOML is currently the only supported format for server properties.")
		fmt.Println("Setting back to default: `config.toml`")
		server.Flags.Config = "config.toml"
	}
	if _, err := toml.DecodeFile("."+string(os.PathSeparator)+server.Flags.Config, &server.TomlConfig); err != nil {
		fmt.Println("Error decoding TOML RSCGo general configuration file:", err)
		os.Exit(137)
	}
	server.InitializeCrypto()
	server.InitPacketHandlerTable()
}

func main() {
	server.Start()
}
