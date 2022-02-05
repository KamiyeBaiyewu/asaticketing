package errors

import (
	"fmt"
	"net/http"
)

var (
	/* Generic */
	// ErrServiceDown - Signifies that the User service is temporarily unavaialble
	ErrServiceDown = APIError{Code: http.StatusInternalServerError, Err: "Service is down try again later"}

	// ErrInvalidValues - Signifies that values submitted have errors
	ErrInvalidValues = APIError{Code: http.StatusBadRequest, Err: "Invalid values submitted"}

	// ErrNotFound - signifies that the resource was not found
	ErrNotFound = APIError{Code: http.StatusNotFound, Err: "Not found"}


	// ErrInternal - signifies that the resource was not found
	ErrInternal = APIError{Code: http.StatusInternalServerError, Err: "Try again"}

	// ErrNotExist - signifies that the resource was not found
	ErrNotExist = func(resource string) APIError {
		err := fmt.Sprintf("%s does not exist", resource)
		return APIError{Code: http.StatusNotFound, Err: err}
	}

/* 	ErrServiceDown = func(service string) error {
	message := fmt.Sprintf("%s Service not currently avalable", service)
	return errors.New(message)
} */
)
