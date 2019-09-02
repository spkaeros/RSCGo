package errors

//DatabaseError A database-related error.
type DatabaseError struct {
	msg string
}

func (e DatabaseError) Error() string {
	return e.msg
}

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

//NewDatabaseError Returns a new database-related error.
func NewDatabaseError(s string) error {
	return DatabaseError{msg: s}
}

//ConnClosed Error to return when the connection closes normally.
var ConnClosed = NewNetworkError("Connection closed.")

//ConnTimedOut Error to return when the connection is inactive for 10 seconds.
var ConnTimedOut = NewNetworkError("Connection timed out.")

//ConnDeadline Error to return when the connection's deadline for reading data can not be properly set.
var ConnDeadline = NewNetworkError("Connection deadline could not be set.")

//BufferOverflow Error to return when we accidentally try to read from an empty packet.
var BufferOverflow = NewNetworkError("Attempted to read too much data from packet.")
