package error

import "fmt"

// ServiceError represents an error encountered during service operations.
type ServiceError struct {
	StatusCode int
	Field      string
	CommonError
}

// NewServiceError creates a new ServiceError with the provided error and error code.
func NewServiceError(err error, code, field string, statusCode int) ServiceError {
	return ServiceError{
		StatusCode:  statusCode,
		Field:       field,
		CommonError: CommonError{Code: code, Message: err.Error()},
	}
}

// Error returns the string representation of the ServiceError.
func (serviceError ServiceError) Error() string {
	return fmt.Sprintf("error: %s %s", serviceError.Field, serviceError.Message)
}
