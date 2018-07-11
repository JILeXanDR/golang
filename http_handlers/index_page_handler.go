package http_handlers

import (
	"net/http"
	"time"
	"os"
	"github.com/nu7hatch/gouuid"
	"log"
)

var reloads = make(map[string]uint64)

const cookieName = "uuid"
const cookieExpiration = 365 * 24 * time.Hour // 1 year

func setCookies(w http.ResponseWriter, r *http.Request) {

	var value string

	cookieUuid, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		// generate value for cookie
		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		value = id.String()
	} else {
		value = cookieUuid.Value
	}

	reloads[value]++
	reloads[value] = reloads[value] * reloads[value]

	log.Println(reloads)

	cookie := http.Cookie{
		Name:    cookieName,
		Value:   value,
		Expires: time.Now().Add(cookieExpiration),
	}

	http.SetCookie(w, &cookie)
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	setCookies(w, r)
	http.ServeFile(w, r, os.Getenv("ROOT_DIR")+"/public/index.html")
}
