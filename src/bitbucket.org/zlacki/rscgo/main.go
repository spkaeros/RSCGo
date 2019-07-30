package main

import (
	"bitbucket.org/zlacki/rscgo/server"
	"flag"
)

var port = 43591

func main() {
	flag.IntVar(&port,"p", 43591, "The TCP port that the server should bind to")
	flag.Parse()
	server.Start(port)
}
