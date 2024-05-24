package error

import "fmt"

// TransportError represents an error that occurred during transport, such as HTTP errors.
type TransportError struct {
	StatusCode  int    // holds the HTTP status code.
	Field       string // represents the field associated with the error.
	CommonError        // contains common error details.
}

// NewTransportError creates a new TransportError with the provided error, error code, field, and HTTP status code.
func NewTransportError(err error, code, field string, statusCode int) TransportError {
	return TransportError{
		StatusCode:  statusCode,
		Field:       field,
		CommonError: CommonError{Code: code, Message: err.Error()},
	}
}

// Error returns the string representation of the TransportError.
func (transportError TransportError) Error() string {
	return fmt.Sprintf("Error: %s %s", transportError.Field, transportError.Message)
}
