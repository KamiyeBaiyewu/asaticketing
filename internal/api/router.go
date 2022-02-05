package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/handlers"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/middlewares"
	v1 "github.com/lilkid3/ASA-Ticket/Backend/internal/api/v1"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/env"
)

// NewRouter create a new router for API Service
func NewRouter(env *env.Env) (http.Handler, error) {

	authorizer := middlewares.NewAuthorizer(env)
	router := mux.NewRouter()
	router.Use(middlewares.UUIDMiddleware)

	// Load all the V1 routes
	v1.LoadRoutes(router, env, authorizer)

	router.HandleFunc("/version", v1.VersionHandler)

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Response: ticles"))
	})

	handler := handlers.CORS(router)
	handler = handlers.Recover(handler)
	return handler, nil
}
