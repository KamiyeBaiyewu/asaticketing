package errors

import "net/http"

var (

	// ErrPolicyExists - User already exists in the database
	ErrPolicyExists = APIError{Code: http.StatusConflict, Err: "Policy already exists"}

	// ErrPolicyAllreadyGranted - tells the user that the policy has already been granted to the user
	ErrPolicyAllreadyGranted = APIError{Code: http.StatusConflict, Err: "The Role has already been granted the policy"}
)
