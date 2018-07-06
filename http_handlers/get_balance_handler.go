package http_handlers

import (
	"net/http"
	"github.com/asaskevich/govalidator"
	"github.com/JILeXanDR/golang/common"
	"strconv"
	"github.com/JILeXanDR/golang/db"
	"github.com/jinzhu/gorm"
)

type balanceResponse struct {
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
		common.ValidationError(w, "Query parameter 'user' is invalid")
		return
	}

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &db.User{}
	err := db.Connection.Where(&db.User{Identifier: int(id)}).First(user).Error
	if err == gorm.ErrRecordNotFound {
		common.JsonMessageResponse(w, "User does not exist", 404)
		return
	} else if err != nil {
		common.InternalServerError(w)
		return
	}

	common.JsonResponse(w, balanceResponse{Balance: user.Balance}, 200)
}
