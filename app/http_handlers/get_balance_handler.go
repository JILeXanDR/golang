package http_handlers

import (
	"net/http"
	"github.com/asaskevich/govalidator"
	"strconv"
	"github.com/JILeXanDR/golang/app/db"
	"github.com/jinzhu/gorm"
	"github.com/JILeXanDR/golang/app/response_writer"
)

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

// Просмотреть баланс
// Request: GET /balance?user=101
// Response: 200 OK {“balance”: 1000}
func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	var (
		userId = r.URL.Query().Get("user")
	)

	if !govalidator.IsInt(userId) {
		response_writer.ValidationError(w, "Query parameter 'user' is invalid")
		return
	}

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &db.User{}
	err := db.Connection.Where(&db.User{Identifier: int(id)}).First(user).Error
	if err == gorm.ErrRecordNotFound {
		response_writer.JsonMessageResponse(w, "User does not exist", 404)
		return
	} else if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	response_writer.JsonResponse(w, BalanceResponse{Balance: user.Balance}, 200)
}
