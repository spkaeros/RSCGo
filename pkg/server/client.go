package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

type Client struct {
	isaacSeed   []uint32
	isaacStream *IsaacSeed
	uID         uint8
	ip          string
	index       int
	kill        chan struct{}
	player      *entity.Player
	socket      net.Conn
	packetQueue chan *packets.Packet
	buffer      []byte
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
					if err, ok := err.(errors.NetError); ok {
						if err.Closed || err.Ping {
							return
						}
						LogDebug(0, "Rejected Packet from: '%s'\n", c.ip)
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

func (c *Client) sendLoginResponse(i byte) {
	c.WritePacket(packets.NewBarePacket([]byte{i}))
	if i != 0 {
		c.kill <- struct{}{}
	} else {
		c.WritePacket(packets.PlayerInfo(c.index, 0))
	}
}

//NewClient Creates a new instance of a Client, registers it with the global ClientList, and returns it.
func NewClient(socket net.Conn) *Client {
	c := &Client{socket: socket, isaacSeed: make([]uint32, 4), packetQueue: make(chan *packets.Packet, 1), ip: strings.Split(socket.RemoteAddr().String(), ":")[0], index: -1, kill: make(chan struct{}, 1), player: entity.NewPlayer(), buffer: make([]byte, 5000)}
	c.StartReader()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return "{idx:'" + strconv.Itoa(c.index) + "', ip:'" + c.ip + "'};"
}
