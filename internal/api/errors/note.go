package errors

import "net/http"

var (




	// ErrNoteNotExist - the note was nott found
	ErrNoteNotExist = APIError{Code: http.StatusBadRequest, Err: "Note does not exist"}
)
