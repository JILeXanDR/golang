package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"github.com/JILeXanDR/golang/db"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]db.User, 0)
	err := db.Connection.Find(&users).Error
	if err != nil {
		common.InternalServerError(w, err)
		return
	}
	common.JsonResponse(w, users, 200)
}
