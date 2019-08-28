package main

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/zlacki/rscgo/pkg/server"
	"github.com/go-ini/ini"
	"github.com/jessevdk/go-flags"
)

func init() {
	if _, err := flags.Parse(&server.Flags); err != nil {
		os.Exit(100)
	}
	if server.Flags.Port > 65535 || server.Flags.Port <= 0 {
		fmt.Println("WARNING: Invalid port number specified.  Valid port numbers are between 0 and 65535.")
		fmt.Println("Setting back to default: `43591`")
		server.Flags.Port = 43591
	}
	if !strings.HasSuffix(server.Flags.Config, ".ini") {
		fmt.Println("WARNING: You entered an invalid configuration file extension.")
		fmt.Println("INI is currently the only supported format for server properties.")
		fmt.Println("Setting back to default: `config.ini`")
		server.Flags.Config = "config.ini"
	}
	cfg, err := ini.Load(server.Flags.Config)
	if err != nil {
		fmt.Printf("Failed to load server config file: %s\n%v", server.Flags.Config, err)
		os.Exit(101)
	}
	server.Version, err = cfg.Section("client").Key("version").Int()
	if err != nil {
		fmt.Println("Failed loading server version number from config file:", err)
		os.Exit(102)
	}
	server.DataDirectory += string(os.PathSeparator) + cfg.Section("server").Key("dataDir").String()
	server.ReadRSAKeyFile(cfg.Section("server").Key("rsaKey").String())
	server.InitializeHashing(cfg.Section("server").Key("salt").String())
}

func main() {
	server.Start()
}
