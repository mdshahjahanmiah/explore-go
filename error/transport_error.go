package error

import (
	"fmt"
	"github.com/pkg/errors"
)

// TransportError represents an error that occurred during transport, such as HTTP errors.
type TransportError struct {
	Err error
	CommonError
}

// NewTransportError creates a new TransportError with the provided error, error code, field, and HTTP status code.
func NewTransportError(err error, code string) TransportError {
	return TransportError{
		Err:         errors.WithStack(err),
		CommonError: CommonError{Code: code, Message: err.Error()},
	}
}

// Error returns the string representation of the TransportError.
func (transportError TransportError) Error() string {
	return fmt.Sprintf("Error: %s", transportError.Message)
}
