package server

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

//Client Represents a single connecting client.
type Client struct {
	isaacSeed       []uint32
	isaacStream     *IsaacSeed
	uID             uint8
	ip              string
	index           int
	kill            chan struct{}
	player          *entity.Player
	socket          net.Conn
	packetQueue     chan *packets.Packet
	outgoingPackets chan *packets.Packet
	buffer          []byte
}

//StartNetworking Starts up 3 new goroutines; one for reading incoming data from the socket, one for writing outgoing data to the socket, and one for client state updates and parsing plus handling incoming packets.  When the clients kill signal is sent through the kill channel, the state update and packet handling goroutine will wait for both the reader and writer goroutines to complete their operations before unregistering the client.
func (c *Client) StartNetworking() {
	var waitForTermination sync.WaitGroup
	waitForTermination.Add(2)
	go func() {
		defer waitForTermination.Done()
		for {
			select {
			default:
				p, err := c.ReadPacket()
				if err != nil {
					if err, ok := err.(errors.NetError); ok {
						if err.Closed || err.Ping {
							return
						}
						LogError.Printf("Rejected Packet from: '%s'\n", c.ip)
						LogError.Println(err)
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
		defer waitForTermination.Done()
		for {
			select {
			case p := <-c.outgoingPackets:
				if p == nil {
					return
				}
				c.WritePacket(p)
			case <-c.kill:
				return
			}
		}
	}()
	go func() {
		defer func() {
			waitForTermination.Wait()
			entity.GetRegion(c.player.X(), c.player.Y()).RemovePlayer(c.player)
			close(c.outgoingPackets)
			close(c.packetQueue)
			if err := c.socket.Close(); err != nil {
				LogError.Println("Couldn't close socket:", err)
			}
			if ok := ClientList.Remove(c.index); ok {
				LogInfo.Printf("Unregistered: %v\n", c)
			}
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
	c.outgoingPackets <- packets.LoginResponse(int(i))
	if i != 0 {
		c.kill <- struct{}{}
	} else {
		c.outgoingPackets <- packets.PlayerInfo(c.index, (c.player.Y()+100)/1000)
		c.outgoingPackets <- packets.ServerMessage("Welcome to RuneScape")
		c.outgoingPackets <- packets.ServerInfo(ClientList.Size())
		c.outgoingPackets <- packets.LoginBox(0, c.ip)
	}
}

//NewClient Creates a new instance of a Client, launches goroutines to handle I/O for it, and returns a reference to it.
func NewClient(socket net.Conn) *Client {
	c := &Client{socket: socket, isaacSeed: make([]uint32, 4), packetQueue: make(chan *packets.Packet, 25), ip: strings.Split(socket.RemoteAddr().String(), ":")[0], index: -1, kill: make(chan struct{}, 1), player: entity.NewPlayer(), buffer: make([]byte, 5000), outgoingPackets: make(chan *packets.Packet, 25)}
	c.StartNetworking()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return fmt.Sprintf("Client[%v] {username:'%v', ip:'%v'}", c.index, c.player.Username, c.ip)
}
