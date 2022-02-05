package handlers

import (
	"errors"
	"log"
	"net/http"
)

// Recover - helps to recover websserver in case of a panic
func Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Printf("Error on panic => %+v\n", err)
				//sendMeMail(err)
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				http.Error(w, "Error Processing Request", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
