package session

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "session")
}
