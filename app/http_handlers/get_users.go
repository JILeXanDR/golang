package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/app/db"
	"github.com/JILeXanDR/golang/app/response_writer"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]db.User, 0)
	err := db.Connection.Find(&users).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}
	response_writer.JsonResponse(w, users, 200)
}
