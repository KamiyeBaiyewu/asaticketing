package errors

// APIError -  helps send the ritght errors paird with the right error code
type APIError struct {
	Code int
	Err  string
}

func (a APIError) Error() string {
	return a.Err
}
