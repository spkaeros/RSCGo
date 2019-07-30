package main

import (
	"bitbucket.org/zlacki/rscgo/pkg/server"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	var opts struct {
		Port int `short:"p" long:"port" description:"The port for the server to listen on," default:"43591"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	server.Start(opts.Port)
}
