package errors

//NetError A RSCGo network-related error.
type NetError struct {
	msg    string
	Ping   bool
	Closed bool
}

func (e NetError) Error() string {
	return e.msg
}

//NewNetworkError Returns a new NetError struct with the specified message.
func NewNetworkError(s string) error {
	return NetError{msg: s, Closed: true}
}

//ConnClosed Error to return when the connection closes normally.
var ConnClosed = NewNetworkError("Connection closed.")

//ConnTimedOut Error to return when the connection is inactive for 10 seconds.
var ConnTimedOut = NewNetworkError("Connection timed out.")

//ConnDeadline Error to return when the connection's deadline for reading data can not be properly set.
var ConnDeadline = NewNetworkError("Problem setting connection deadline.")

//BufferOverflow Error to return when we accidentally try to read from an empty packet.
var BufferOverflow = NewNetworkError("Attempted to read too much data from packet; would have caused buffer overflow.")
