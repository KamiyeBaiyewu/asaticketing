package v1

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/opentracing/opentracing-go"
)

/*
	Helper
*/
// Aliases for endpoints
type apiEndpoint utils.Endpoint

// type userParametersRequest

func newAPIEndpoint(method string, path string, handlerFunc http.HandlerFunc, c ...alice.Constructor) apiEndpoint {
	return apiEndpoint{method, path, alice.New(c...).ThenFunc(handlerFunc).(http.HandlerFunc)}
}

func tagSpan(span opentracing.Span, key, value string) {

	if span != nil {
		span.SetTag(key, value)
	}
}
