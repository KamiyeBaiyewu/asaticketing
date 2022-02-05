package utils

import (
	"encoding/json"
	"net/http"

	apiError "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/sirupsen/logrus"
)

// GenericError - represents error structure for generic error (All error responses are the same name)
type GenericError struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

// WriteError returns a JSON error message and HTTP Status Code
func WriteError(w http.ResponseWriter, code int, err interface{}, data interface{}) {

	var response GenericError
	switch err.(type) {
	// apiError already extedns the golang generic erros
	case apiError.APIError: // if the errors is an API error then convert it into a generic err
	apiErr := err.(apiError.APIError)
		response = GenericError{
			Code:  apiErr.Code,
			Error: apiErr.Error(),
			Data:  data,
		}
		break
	case error: //incase a golang error was sent
		errErr := err.(error)
		response = GenericError{
			Code:  code,
			Error: errErr.Error(),
			Data:  data,
		}
		break
	default: // any error at this point has to be a string
		strErr := err.(string)
		response = GenericError{
			Code:  code,
			Error: strErr,
			Data:  data,
		}
		break
	}
	WriteJSON(w, code, response)
}

// WriteJSON returns a JSON data and HTTP Status Code
func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Warn("Error writing response")
	}
}
