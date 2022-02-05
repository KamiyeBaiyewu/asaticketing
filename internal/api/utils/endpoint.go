package utils

import (
	"net/http"

	"github.com/justinas/alice"
)

// Endpoint is a structure of API endpoints
type Endpoint struct {
	Method string
	Path   string
	Func   http.HandlerFunc

	// permissionTypes []auth.PolicyType
}

// NewEndpoint creates a new endpoint
func NewEndpoint(method string, path string, handlerFunc http.HandlerFunc, c ...alice.Constructor) Endpoint {
	return Endpoint{method, path, alice.New(c...).ThenFunc(handlerFunc).(http.HandlerFunc)}
}
