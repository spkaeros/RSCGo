package main

import (
	"bitbucket.org/zlacki/rscgo/pkg/server"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
)

func init() {
	if _, err := flags.Parse(&server.Flags); err != nil {
		os.Exit(100)
	}
	if server.Flags.Port >= 65535 || server.Flags.Port <= 0 {
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
}

func main() {
	server.Start(server.Flags.Port)
}
