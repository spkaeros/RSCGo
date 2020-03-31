package errors

import "errors"

//RscError An RSCGo generic error type.
type RscError = error

//NewRscError Returns tag surrounded by brackets and prepended to msg as a new RscError.
func NewRscError(tag string, msg string) error {
	return errors.New("[" + tag + "] " + msg)
}

//DatabaseError An RSCGo error type for database-related errors.
type DatabaseError struct {
	RscError
}

//NewDatabaseError Returns a new database-related error.
func NewDatabaseError(s string) error {
	return NewRscError("DatabaseError", s)
}

//ArgumentsError An RSCGo error type for function parameter-related errors.
type ArgumentsError struct {
	RscError
}

//NewArgsError Returns a new ArgumentError with the specified message appended to the error tag.
func NewArgsError(s string) error {
	return NewRscError("InvalidArgument", s)
}

//NetError A RSCGo network-related error.
type NetError struct {
	RscError
	Fatal bool
}

//NewNetworkError Returns a new database-related error.
func NewNetworkError(s string, fatal bool) NetError {
	return NetError{NewRscError("NetworkError", s), fatal}
}
