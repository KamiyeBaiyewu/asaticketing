package errors

import "net/http"

var (
	// ErrUserTypeRequired - tells the user that user type has to be included in the qury
	ErrUserTypeRequired = APIError{Code: http.StatusBadRequest, Err: "UserType is required"}
	// ErrEmailAlreadyExists - tells the user that
	ErrEmailAlreadyExists = APIError{Code: http.StatusBadRequest, Err: "Email already exists"}
	// ErrCreatingUser - tells the user that
	ErrCreatingUser = APIError{Code: http.StatusBadRequest, Err: "Error creating user"}
)
