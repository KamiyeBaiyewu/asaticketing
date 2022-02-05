package errors

import "net/http"

var (

	// ErrObjectExists - Object already exists in the database
	ErrObjectExists = APIError{Code: http.StatusConflict, Err: "Object already exists"}
)
