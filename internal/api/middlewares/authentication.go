package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/auth"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/utils"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
)

type principalContextKeyType struct{}

var principalContextKey principalContextKeyType

// Authentication middleware for checking if token submitted is valid before next handler
func Authentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := CheckToken(r)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// CheckToken - Gets the token and places the pricipal inside the request context
func CheckToken(r *http.Request) (*http.Request, error) {
	// extract the token from the request
	token, err := getToken(r)
	// println(token)
	if err != nil {
		return r, nil
	}

	if token == "" {
		return r, nil
	}

	principal, err := auth.VerifyToken(token)
	if err != nil {

		return r, err
	}
	// log.Printf("Principlal => %+v\n", principal)

	return r.WithContext(WithPricipalContext(r.Context(), *principal)), nil
}

// WithPricipalContext inserts pricipal information inside request context
func WithPricipalContext(ctx context.Context, principal model.Principal) context.Context {

	return context.WithValue(ctx, principalContextKey, principal)
}

// getToken - retrievs the token from authorazation header
func  getToken(r *http.Request) (string, error) {

	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("Invalid Token")
	}

	tokenParts := strings.SplitN(token, " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" || len(tokenParts[1]) == 0 {
		return "", errors.New("Authorization headaer format must be Bearer {token}")
	}

	return tokenParts[1], nil
}


// GetPrincipal - gets the principal information from the request profile
func GetPrincipal(r *http.Request) model.Principal {

	if principal, ok := r.Context().Value(principalContextKey).(model.Principal); ok {

		return principal
	}
	return model.NilPricipal
}
