package error

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	requestId = "request_id"
)

// CommonError struct represents an error with a code and a corresponding message.
type CommonError struct {
	Code    string
	Message string
}

// Source represents a field and its corresponding message in a JSON object.
type Source struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorObject represents an object with an ID, status, code, optional detail, and optional source.
type ErrorObject struct {
	ID     string  `json:"id"`
	Status int     `json:"status"`
	Code   string  `json:"code"`
	Detail string  `json:"detail,omitempty"`
	Source *Source `json:"source,omitempty"`
}

// Payload represents a payload containing a list of errors.
type Payload struct {
	Errors []*ErrorObject `json:"errors"`
}

// Error returns the string representation of the Object error.
func (e *ErrorObject) Error() string {
	return fmt.Sprintf("Error: %s %s\n", e.Source.Field, e.Source.Message)
}

// EncodeError encodes the given error to JSON format and writes it to the HTTP response.
func EncodeError(ctx context.Context, e error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	reqID, _ := getRequestId(ctx)
	errObj := errToErrorObject(reqID, e)
	w.WriteHeader(errObj.Status)

	// Convert an error object to JSON and write to response.
	json.NewEncoder(w).Encode(&Payload{Errors: []*ErrorObject{errObj}})
}

// getRequestId extracts the request ID from the provided context.
func getRequestId(ctx context.Context) (string, error) {
	if reqID, ok := ctx.Value(requestId).(string); ok && reqID != "" {
		return reqID, nil
	}
	return "", errors.New("request_id not found in request context")
}

// errToErrorObject transforms the given error into an ErrorObject.
func errToErrorObject(id string, err error) *ErrorObject {
	errObj := &ErrorObject{ID: id}

	switch v := err.(type) {
	case TransportError:
		errObj.Status = http.StatusInternalServerError
		errObj.Code = v.Code
		errObj.Detail = v.Message
	case ServiceError:
		errObj.Status = v.StatusCode
		errObj.Code = v.Code
		errObj.Source = &Source{Field: v.Field, Message: v.Message}
	case error:
		errObj.Status = http.StatusInternalServerError
		errObj.Code = err.Error()
	default:
		errObj.Status = http.StatusInternalServerError
		errObj.Code = err.Error()
	}

	return errObj
}
