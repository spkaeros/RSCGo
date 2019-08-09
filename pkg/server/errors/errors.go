package errors

type NetError struct {
	msg    string
	Ping   bool
	Closed bool
}

func (e NetError) Error() string {
	return e.msg
}

func NewNetworkError(s string) error {
	return NetError{ msg: s, Closed:true }
}

var ConnClosed = NewNetworkError("Connection closed.")
var ConnTimedOut = NewNetworkError("Connection timed out.")
var ConnDeadline = NewNetworkError("Problem setting connection deadline.")