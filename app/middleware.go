package app

import (
	"net/http"
	"log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
