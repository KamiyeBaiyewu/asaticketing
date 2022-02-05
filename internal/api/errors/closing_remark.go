package errors


import "net/http"

var (
	// ErrClosingRemarkExists - tells the user that there can only be one ticket closing remark
	ErrClosingRemarkExists = APIError{Code: http.StatusConflict, Err: "A closing remark already exist"}
)