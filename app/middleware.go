package app

import (
	"net/http"
	"log"
	"github.com/JILeXanDR/golang/app/session"
	"github.com/JILeXanDR/golang/app/services"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := session.GetSession(r)

		if auth := sess.Values["auth"]; auth != nil {
			if phoneSession := auth.(services.PhoneSession); phoneSession.Confirmed {
				w.Header().Set("Phone", phoneSession.Phone)
			}
		}

		next.ServeHTTP(w, r)
	})
}
