package server

type NetError struct {
	msg    string
	ping   bool
	closed bool
}

func (e *NetError) Error() string {
	return e.msg
}

func Closed() *NetError {
	return &NetError{msg: "Connection reset by peer.", closed: true}
}

func Timeout() *NetError {
	return &NetError{msg: "Connection timed out.", ping: true}
}

func Deadline() *NetError {
	return &NetError{msg: "Could not set read deadline for Client listener.", closed: true}
}
