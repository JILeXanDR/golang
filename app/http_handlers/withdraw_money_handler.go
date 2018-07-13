package http_handlers

import (
	"net/http"
	"strconv"
	"github.com/JILeXanDR/golang/app/db"
	"github.com/jinzhu/gorm"
	"github.com/JILeXanDR/golang/app/response_writer"
	"github.com/JILeXanDR/golang/errors"
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
			response_writer.JsonMessageResponse(w, "User does not exist", 400)
		} else {
			response_writer.HandleError(w, err)
		}
		return
	}

	value, _ := strconv.ParseFloat(amount, 64)

	if user.Balance < value {
		response_writer.JsonMessageResponse(w, errors.ErrNotEnoughMoney.Error(), 400)
		return
	}

	user.Balance -= value
	err = db.Connection.Save(user).Error
	if err != nil {
		response_writer.HandleError(w, err)
		return
	}

	response_writer.JsonResponse(w, user, 200)
}
