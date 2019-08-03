package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	isaacSeed              [4]uint32
	isaacStream            *IsaacSeed
	uID                    uint8
	ip                     string
	index                  int
	kill                   chan struct{}
	player                 *entity.Player
	socket                 net.Conn
	packetQueue            chan *Packet
}

//StartReader Creates a new goroutine to handle all incoming network events for the receiver Client.
// This goroutine will also automatically handle cleanup for Client disconnections, and handle incoming I/O errors
// and disconnect the related Client appropriately.
func (c *Client) StartReader() {
	go func() {
		defer close(c.packetQueue)
		for {
			select {
			default:
				p, err := c.ReadPacket()
				if err != nil {
					if err, ok := err.(*NetError); ok {
						if err.closed || err.ping {
							return
						}
						fmt.Printf("Rejected Packet from: '%s'\n", getIPFromConn(c.socket))
						fmt.Println(err)
					}
					continue
				}
				c.packetQueue <- p
			case <-c.kill:
				return
			}
		}
	}()
	go func() {
		defer func() {
			if err := c.socket.Close(); err != nil {
				// This shouldn't reasonably happen.
				fmt.Println("WARNING: Error closing socket!", err)
			}
			ActiveClients.Remove(c.index)
		}()
		defer close(c.kill)
		for {
			select {
			case p := <-c.packetQueue:
				if p == nil {
					return
				}
				c.HandlePacket(p)
			case <-c.kill:
				return
			}
		}
	}()
}

//NewClient Creates a new instance of a Client, registers it with the global ClientList, and returns it.
func NewClient(socket net.Conn) *Client {
	c := &Client{socket: socket, packetQueue: make(chan *Packet, 1), ip: getIPFromConn(socket), index: -1, kill: make(chan struct{}, 1), player: entity.NewPlayer()}
	c.StartReader()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return "{idx:'" + strconv.Itoa(c.index) + "', ip:'" + c.ip + "'};"
}
