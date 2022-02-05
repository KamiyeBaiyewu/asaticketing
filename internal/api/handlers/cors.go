package handlers

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// CORS - Cross-origin resource sharing helps with API source
func CORS(handler http.Handler)  http.Handler {
	
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS","PATCH","DELETE"}),
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
	)

	return cors(handler)
}
