package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Channel
	cipherKey int64
	uID       uint8
	ip        string
	index     int
	kill      chan struct{}
	player    *entity.Player
}

//Unregister Clean up resources and Unregister the receiver from the global ClientList.
func (c *Client) Unregister() {
	close(c.kill)
	fmt.Println("Unregistering Client" + c.String())
	if err := c.socket.Close(); err != nil {
		fmt.Printf("WARNING: Error closing listener for Client%s\n", c.String())
		fmt.Println(err)
	}
	ActiveClients.Remove(c.index)
}

//StartReader Creates a new goroutine to handle all incoming network events for the receiver Client.
// This goroutine will also automatically handle cleanup for Client disconnections, and handle incoming I/O errors
// and disconnect the related Client appropriately.
func (c *Client) StartReader() {
	go func() {
		defer c.Unregister()
		for {
			select {
			case <-c.kill:
				return
			case <-time.After(time.Millisecond * 5):
				p, err := c.NextPacket()
				if err != nil {
					fmt.Println(err.Error())
					if err.ping || err.closed {
						return
					}
				}
				c.HandlePacket(p)
				continue
			}
		}
	}()
}

//NewClient Creates a new instance of a Client, registers it with the global ClientList, and returns it.
func NewClient(socket net.Conn) *Client {
	c := &Client{Channel: Channel{socket: socket}, cipherKey: -1, ip: getIPFromConn(socket), index: -1, kill: make(chan struct{}, 1), player: entity.NewPlayer()}
	c.StartReader()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return "{idx:'" + strconv.Itoa(c.index) + "', ip:'" + c.ip + "'};"
}
