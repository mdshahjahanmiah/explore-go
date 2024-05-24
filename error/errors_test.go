package error

import (
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func Test_ErrToErrorObject(t *testing.T) {
	id := "test_id"

	// Test case for TransportError
	transportRequestErr := TransportError{
		Err: errors.New("transport error"),
		CommonError: CommonError{Code: "TRANSPORT_ERROR",
			Message: "Transport error occurred"},
	}
	serviceObj := errToErrorObject(id, transportRequestErr)
	if serviceObj.ID != id || serviceObj.Status != http.StatusInternalServerError || serviceObj.Code != "TRANSPORT_ERROR" || serviceObj.Detail != "Transport error occurred" {
		t.Errorf("errToErrorObject() for TransportError: got %v, want %v", serviceObj, &ErrorObject{ID: id, Status: http.StatusInternalServerError, Code: "TRANSPORT_ERROR", Detail: "Transport error occurred"})
	}

	// Test case for ServiceError
	serviceBadRequestErr := ServiceError{
		StatusCode:  http.StatusBadRequest,
		Field:       "test_field",
		CommonError: CommonError{Code: "SERVICE_ERROR", Message: "Service error occurred"},
	}
	serviceBadRequestObj := errToErrorObject(id, serviceBadRequestErr)
	expectedServiceBadRequestObj := &ErrorObject{
		ID:     id,
		Status: http.StatusBadRequest,
		Code:   "SERVICE_ERROR",
		Source: &Source{
			Field:   "field",
			Message: "Service error occurred",
		},
	}
	if serviceBadRequestObj.ID != id || serviceBadRequestObj.Status != http.StatusBadRequest || serviceBadRequestObj.Code != "SERVICE_ERROR" || serviceBadRequestObj.Source.Field != "test_field" || serviceBadRequestObj.Source.Message != "Service error occurred" {
		t.Errorf("errToErrorObject() for ServiceError: got %v, want %v", serviceBadRequestObj, expectedServiceBadRequestObj)
	}
}
