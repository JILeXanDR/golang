package http_handlers

import (
	"net/http"
	"os"
	"github.com/JILeXanDR/golang/app/session"
	"github.com/JILeXanDR/golang/app/response_writer"
)

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	sess.Values["user"] = "anon."
	sess.Save(r, w)
	http.ServeFile(w, r, os.Getenv("ROOT_DIR")+"/public/index.html")
}
