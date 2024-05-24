package error

import (
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func Test_ErrToErrorObject(t *testing.T) {
	id := "test_id"

	// Test case for ServiceError
	serviceErr := ServiceError{
		Err: errors.New("service error"),
		CommonError: CommonError{Code: "SERVICE_ERROR",
			Message: "Service error occurred"},
	}
	serviceObj := errToErrorObject(id, serviceErr)
	if serviceObj.ID != id || serviceObj.Status != http.StatusInternalServerError || serviceObj.Code != "SERVICE_ERROR" || serviceObj.Detail != "Service error occurred" {
		t.Errorf("errToErrorObject() for ServiceError: got %v, want %v", serviceObj, &ErrorObject{ID: id, Status: http.StatusInternalServerError, Code: "SERVICE_ERROR", Detail: "Service error occurred"})
	}

	// Test case for TransportError
	transportBadRequestErr := TransportError{
		StatusCode:  http.StatusBadRequest,
		Field:       "test_field",
		CommonError: CommonError{Code: "TRANSPORT_ERROR", Message: "Transport error occurred"},
	}
	transportBadRequestObj := errToErrorObject(id, transportBadRequestErr)
	expectedTransportBadRequestObj := &ErrorObject{
		ID:     id,
		Status: http.StatusBadRequest,
		Code:   "TRANSPORT_ERROR",
		Source: &Source{
			Field:   "field",
			Message: "Transport error occurred",
		},
	}
	if transportBadRequestObj.ID != id || transportBadRequestObj.Status != http.StatusBadRequest || transportBadRequestObj.Code != "TRANSPORT_ERROR" || transportBadRequestObj.Source.Field != "test_field" || transportBadRequestObj.Source.Message != "Transport error occurred" {
		t.Errorf("errToErrorObject() for TransportError: got %v, want %v", transportBadRequestObj, expectedTransportBadRequestObj)
	}
}
