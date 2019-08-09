package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func (c *Client) Write(b []byte) int {
	l, err := c.socket.Write(b)
	if err != nil {
		// TODO: Severe enough to kill the client?  More than likely, yes.
		LogDebug(0, "ERROR: Could not Write to Client socket.\n")
		fmt.Println(err)
	}
	if l != len(b) {
		LogDebug(1, "WARNING: Wrong number of bytes written to Client socket.  Expected %d, got %d.\n", len(b), l)
	}
	return l
}

func (c *Client) Read(dst []byte) error {
	if err := c.socket.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		// This shouldn't happen
		return errors.ConnDeadline
	}
	length, err := c.socket.Read(dst)
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return errors.ConnTimedOut
		}
	} else if length != len(dst) {
		return errors.NewNetworkError("Client.Read: unexpected length.  Expected " + strconv.Itoa(len(dst)) + ", got " + strconv.Itoa(length) + ".")
	}

	return nil
}