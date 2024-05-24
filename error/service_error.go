package error

import (
	"fmt"
	"github.com/pkg/errors"
)

// ServiceError represents an error encountered during service operations.
type ServiceError struct {
	Err error
	CommonError
}

// NewServiceError creates a new ServiceError with the provided error and error code.
func NewServiceError(err error, code string) ServiceError {
	return ServiceError{
		Err:         errors.WithStack(err),
		CommonError: CommonError{Code: code, Message: err.Error()},
	}
}

// Error returns the string representation of the ServiceError.
func (serviceError ServiceError) Error() string {
	return fmt.Sprintf("Error: %s", serviceError.Message)
}
