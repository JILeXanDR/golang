package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"github.com/JILeXanDR/golang/db"
)

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orders = &[]db.Order{}
	err := db.Connection.Find(orders).Error
	if err != nil {
		common.InternalServerError(w)
		return
	}

	common.JsonResponse(w, orders, 200)
}
