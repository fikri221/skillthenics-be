package json

import "net/http"

// StatusError is an interface for errors that carry an HTTP status code.
type StatusError interface {
	error
	Status() int
}

// statusError is a concrete implementation of StatusError.
type statusError struct {
	status  int
	message string
}

func (e *statusError) Error() string {
	return e.message
}

func (e *statusError) Status() int {
	return e.status
}

// NewError creates a new error with a specific HTTP status code.
func NewError(status int, message string) error {
	return &statusError{
		status:  status,
		message: message,
	}
}

// Common errors
var (
	ErrInternal     = NewError(http.StatusInternalServerError, "Internal Server Error")
	ErrNotFound     = NewError(http.StatusNotFound, "Resource Not Found")
	ErrUnauthorized = NewError(http.StatusUnauthorized, "Unauthorized")
	ErrBadRequest   = NewError(http.StatusBadRequest, "Bad Request")
)
