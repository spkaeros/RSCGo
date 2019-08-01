package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func (c *Client) Write(b []byte) {
	l, err := c.socket.Write(b)
	if err != nil {
		fmt.Println("ERROR: Could not Write to Client socket.")
		fmt.Println(err)
	}
	if l != len(b) {
		fmt.Printf("WARNING: Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
}

func (c *Client) Read(len int) ([]byte, error) {
	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		// This shouldn't happen
		return nil, Deadline()
	}
	buf := make([]byte, len)
	length, err := c.socket.Read(buf)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil, Timeout()
		}
		if strings.Contains(err.Error(), "use of closed") {
			return nil, &NetError{msg: "Trying to read a closed socket.", closed: true}
		}
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") {
			return nil, Closed()
		}
	}
	if length != len {
		return nil, &NetError{msg: "Client.Read: unexpected length.  Expected " + strconv.Itoa(len) + ", got " + strconv.Itoa(length) + "."}
	}

	return buf, nil
}
