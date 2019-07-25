package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type client struct {
	channel
	cipherKey int64
	uID       uint8
	ip        string
	index     int
	kill      chan struct{}
}

//unregister Clean up resources and unregister the receiver from the global clientList.
func (c *client) unregister() {
	c.kill <- struct{}{}
	close(c.kill)
	close(c.send)
	fmt.Println("Unregistering client" + c.String())
	if err := c.socket.Close(); err != nil {
		fmt.Printf("WARNING: Error closing listener for client%s\n", c.String())
		fmt.Println(err)
	}
	activeClients.remove(c.index)
}

//startWriter todo: do I need this?
func (c *client) startWriter() {
	go func() {
		for {
			select {
			case p := <-c.send:
				c.sendPacket(p)
				continue
			case <-c.kill:
				return;
			}
		}
	}()
}

//startReader Creates a new goroutine to handle all incoming network events for the receiver client.
// This goroutine will also automatically handle cleanup for client disconnections, and handle incoming I/O errors
// and disconnect the related client appropriately.
func (c *client) startReader() {
	go func() {
		defer func() {
			c.unregister()
		}()

		for {
			select {
			case <-c.kill:
				return
			default:
				headerBuffer := make([]byte, 3)

				if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
					fmt.Println("Rejecting client" + c.String())
					fmt.Println("ERROR: Could not set read timeout for client listenerRef.")
					fmt.Println(err)
					return
				}
				headerLength, err := c.socket.Read(headerBuffer)

				if err == io.EOF {
					fmt.Println("Connection reset by peer while attempting to read a packet header (io.EOF)")
					return
				} else if err, ok := err.(net.Error); ok && err.Timeout() {
					fmt.Println("Connection timed out while attempting to read a packet header. (net.Error)")
					return
				} else if err != nil {
					fmt.Println("ERROR: Unexpected I/O error encountered while reading packet header for client"+c.String(), err)
					continue
				} else if headerLength != 3 {
					fmt.Printf("DEBUG: Packet header unexpected length.  Expected 3 bytes, got %d bytes\n", headerLength)
					continue
				}

				length := int(headerBuffer[0] & 0xFF)
				if length >= 160 {
					length = (length-160)*256 + int(headerBuffer[1]&0xFF)
				} else {
					length--
				}

				opcode := headerBuffer[2] & 0xFF
				length--

				payloadBuffer := make([]byte, length)

				if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
					fmt.Println("Rejecting client" + c.String())
					fmt.Println("ERROR: Could not set read timeout for client listenerRef.")
					fmt.Println(err)
					return
				}
				payloadLength, err := c.socket.Read(payloadBuffer)

				if err == io.EOF {
					fmt.Println("Connection reset by peer while attempting to read a packet frame (io.EOF)")
					return
				} else if err, ok := err.(net.Error); ok && err.Timeout() {
					fmt.Println("Connection timed out while attempting to read a packet frame. (net.Error)")
					return
				} else if err != nil {
					fmt.Println("ERROR: Unexpected I/O error encountered while reading packet header for client"+c.String(), err)
					continue
				} else if payloadLength != length {
					fmt.Printf("DEBUG: Packet frame unexpected length.  Expected %d bytes, got %d bytes\n", length, payloadLength)
					continue
				}

				if length < 160 {
					payloadBuffer = append(payloadBuffer, headerBuffer[1])
					length++
				}

				c.handlePacket(newPacket(opcode, payloadBuffer, length))
			}
		}
	}()
}

//newClient Creates a new instance of a client, registers it with the global clientList, and returns it.
func newClient(socket net.Conn) *client {
	return &client{channel: channel{socket: socket, send: make(chan *packet)}, cipherKey: -1, ip: getIPFromConn(socket), index: -1, kill: make(chan struct{}, 1)}
}

//String Returns a string populated with some of the more identifying fields from the receiver client.
func (c *client) String() string {
	return "{idx:'" + strconv.Itoa(c.index) + "', ip:'" + c.ip + "'};"
}
