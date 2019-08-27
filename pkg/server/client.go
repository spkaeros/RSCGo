/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-22-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-27-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

//Client Represents a single connecting client.
type Client struct {
	isaacSeed       []uint64
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
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
			select {
			default:
				p, err := c.ReadPacket()
				if err != nil {
					if err, ok := err.(errors.NetError); ok {
						if err.Closed || err.Ping {
							// TODO: I need to make sure this doesn't cause a panic due to kill being closed already
							c.kill <- struct{}{}
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
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
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
			c.player.Removing = true
			close(c.kill)
			close(c.outgoingPackets)
			close(c.packetQueue)
			if err := c.socket.Close(); err != nil {
				LogError.Println("Couldn't close socket:", err)
			}
			hash := strutil.Base37(c.player.Username)
			if c1, ok := Clients[hash]; c1 == c && ok {
				delete(Clients, hash)
			}
			if ok := ClientList.Remove(c.index); ok {
				LogInfo.Printf("Unregistered: %v\n", c)
			}
		}()
		ticker := time.NewTicker(25 * time.Millisecond)
		for range ticker.C {
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
		LogInfo.Printf("Denied Player[%v]: {ip:'%v', username:'%v', Response='%v'}\n", c.index, c.ip, c.player.Username, i)
		select {
		case <-time.After(100 * time.Millisecond):
			c.kill <- struct{}{}
		}
	} else {
		LogInfo.Printf("Registered Player[%v]: {ip:'%v', username:'%v'}\n", c.index, c.ip, c.player.Username)
		c.player.AppearanceChanged = true
		c.player.SetCoords(220, 445)
		for i := 0; i < 18; i++ {
			level := 1
			exp := 0
			if i == 3 {
				level = 10
				exp = 1154
			}
			c.player.Skillset.Current[i] = level
			c.player.Skillset.Maximum[i] = level
			c.player.Skillset.Experience[i] = exp
		}
		c.outgoingPackets <- packets.PlayerInfo(c.player)
		c.outgoingPackets <- packets.PlayerStats(c.player)
		c.outgoingPackets <- packets.EquipmentStats(c.player)
		c.outgoingPackets <- packets.FightMode(c.player)
		c.outgoingPackets <- packets.FriendList(c.player)
		c.outgoingPackets <- packets.ClientSettings(c.player)
		c.outgoingPackets <- packets.Fatigue(c.player)
		c.outgoingPackets <- packets.WelcomeMessage
		c.outgoingPackets <- packets.ServerInfo(ClientList.Size())
		c.outgoingPackets <- packets.LoginBox(0, c.ip)
	}
}

//NewClient Creates a new instance of a Client, launches goroutines to handle I/O for it, and returns a reference to it.
func NewClient(socket net.Conn) *Client {
	c := &Client{socket: socket, isaacSeed: make([]uint64, 2), packetQueue: make(chan *packets.Packet, 25), ip: strings.Split(socket.RemoteAddr().String(), ":")[0], index: -1, kill: make(chan struct{}, 1), player: entity.NewPlayer(), buffer: make([]byte, 5000), outgoingPackets: make(chan *packets.Packet, 25)}
	c.StartNetworking()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return fmt.Sprintf("Client[%v] {username:'%v', ip:'%v'}", c.index, c.player.Username, c.ip)
}
