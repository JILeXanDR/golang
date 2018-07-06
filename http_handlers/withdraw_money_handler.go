package http_handlers

import (
	"net/http"
	"strconv"
	"github.com/JILeXanDR/golang/db"
	"github.com/JILeXanDR/golang/common"
	"github.com/jinzhu/gorm"
)

// снимать деньги со счетов
// Request: POST /withdraw {“user”: 101, “amount”: 50}
// Response: 200 OK
func WithdrawMoneyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		userId = r.URL.Query().Get("user")
		amount = r.URL.Query().Get("amount")
	)

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &db.User{Identifier: int(id)}
	err := db.Connection.Where(user).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			common.JsonMessageResponse(w, "User does not exist", 400)
		} else {
			common.HandleError(w, err)
		}
		return
	}

	value, _ := strconv.ParseFloat(amount, 64)

	if user.Balance < value {
		common.JsonMessageResponse(w, "Not enough money", 400)
		return
	}

	user.Balance -= value
	err = db.Connection.Save(user).Error
	if err != nil {
		common.HandleError(w, err)
		return
	}

	common.JsonResponse(w, user, 200)
}
