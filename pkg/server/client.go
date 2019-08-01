package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"fmt"
	"net"
	"strconv"
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

//StartReader Creates a new goroutine to handle all incoming network events for the receiver Client.
// This goroutine will also automatically handle cleanup for Client disconnections, and handle incoming I/O errors
// and disconnect the related Client appropriately.
func (c *Client) StartReader() {
	go func() {
		defer close(c.nextPacket)
		for {
			select {
			default:
				headerBuf, err := c.Read(3)
				if err != nil {
					if err, ok := err.(*NetError); ok {
						if err.closed || err.ping {
							return
						}
						fmt.Printf("Rejected Packet from: '%s'\n", getIPFromConn(c.socket))
					}
					continue
				}

				length := int(headerBuf[0] & 0xFF)
				if length >= 160 {
					length = (length-160)*256 + int(headerBuf[1]&0xFF)
				} else {
					// TODO: Should it be <= 160, and should it be >= 1?
					// If the payload length is less than 160 bytes, the 2nd byte in the header is used to store the last byte
					//  of payload data.  Subtract one from length so that we don't try to read it from the end of the payload.
					length--
				}

				// Opcode byte is included in the length variable, but we read it into the header buffer since it should be there.
				opcode := headerBuf[2] & 0xFF
				length--

				payloadBuf, err := c.Read(length)
				if err != nil {
					fmt.Println("Problem reading next packet payload:", err)
					if err, ok := err.(*NetError); ok {
						if err.closed || err.ping {
							return
						}
					}
					continue
				}

				if length < 160 {
					// 1-byte length in header causes the client to put the last byte of payload data in the header
					payloadBuf = append(payloadBuf, headerBuf[1])
					length++
				}

				c.nextPacket <- NewPacket(opcode, payloadBuf, length)
			case <-c.kill:
				return
			}
		}
	}()
	go func() {
		defer func() {
			fmt.Println("Unregistering client" + c.String())
			if err := c.socket.Close(); err != nil {
				fmt.Printf("WARNING: Error closing listener for client%s\n", c.String())
				fmt.Println(err)
			}
			ActiveClients.Remove(c.index)
		}()
		defer close(c.kill)
		for {
			select {
			case p := <-c.nextPacket:
				/*				if err != nil {
								fmt.Println(err.Error())
								if err.(*NetError).ping || err.(*NetError).closed {
									return
								}
							}*/
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
	c := &Client{Channel: Channel{socket: socket, nextPacket: make(chan *Packet, 1)}, cipherKey: -1, ip: getIPFromConn(socket), index: -1, kill: make(chan struct{}, 1), player: entity.NewPlayer()}
	c.StartReader()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return "{idx:'" + strconv.Itoa(c.index) + "', ip:'" + c.ip + "'};"
}
