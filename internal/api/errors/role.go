package errors

import "net/http"

var (

	// ErrRoleExists - User already exists in the database
	ErrRoleExists = APIError{Code: http.StatusConflict, Err: "Role already exists"}

	// ErrRoleIDExists - Role with the ID already exists
	ErrRoleIDExists = APIError{Code: http.StatusConflict, Err: "Role with ID already exists"}

	// ErrRoleAllreadyGranted - tells the user that the role has already been granted to the user
	ErrRoleAllreadyGranted = APIError{Code: http.StatusConflict, Err: "The User has already been granted the role"}

	// ErrRoleNotExist - tells the user that
	ErrRoleNotExist = APIError{Code: http.StatusBadRequest, Err: "Role does not exist"}
)
