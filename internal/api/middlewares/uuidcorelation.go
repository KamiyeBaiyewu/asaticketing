package middlewares

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
)

// UUIDMiddleware generates and inject correlationId in the context
//this will be used in register fuc in http_server.go
func UUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get request context
		ctx := r.Context()
		//get the global tracer
		tracer := opentracing.GlobalTracer()
		//this is where we start our span for this operation, this will be the parent for this method
		span := tracer.StartSpan("UUIDMiddleware")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
		u2 := uuid.NewV4()
		ctx = context.WithValue(ctx, "correlationid", u2.String())
		span.SetTag("correlationid",u2.String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// UUIDMiddleware generates and inject correlationId in the context
//this will be used in register fuc in http_server.go
/* func UUIDMiddleware(env *env.Env) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			println("Passing UUID middleware")
			//get request context
			ctx := r.Context()
			//get the global tracer
			//tracer := opentracing.GlobalTracer()
			tracer := env.Tracer
			//this is where we start our span for this operation, this will be the parent for this method
			span := tracer.StartSpan("UUIDMiddleware")
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
			u2 := uuid.NewV4()
			ctx = context.WithValue(ctx, "correlationid", u2.String())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})

	}
} */
