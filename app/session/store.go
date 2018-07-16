package session

import (
	"github.com/gorilla/sessions"
	"net/http"
	"encoding/gob"
	"github.com/JILeXanDR/golang/app/services"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))
var session *sessions.Session

func GetSession(r *http.Request) (*sessions.Session, error) {
	var err error
	if session == nil {
		gob.Register(services.PhoneSession{})
		session, err = store.Get(r, "session")
	}

	return session, err
}
