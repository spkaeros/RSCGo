package isaac

type isaacError struct {
	msg string
}

func (err *isaacError) Error() string {
	return err.msg
}
