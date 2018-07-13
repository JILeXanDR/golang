package http_handlers

import (
	"net/http"
	"github.com/asaskevich/govalidator"
	"strconv"
	"github.com/JILeXanDR/golang/app/db"
	"github.com/jinzhu/gorm"
	"github.com/JILeXanDR/golang/app/response_writer"
)

// зачислять деньги на счет пользователям (создать пользователя, если еще не существует)
// Request: POST /deposit {“user”: 101, “amount”: 100}
// Response: 200 OK
func DepositMoneyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		userId = r.URL.Query().Get("user")
		amount = r.URL.Query().Get("amount")
	)

	if len(userId) == 0 || !govalidator.IsNumeric(userId) {
		response_writer.ValidationError(w, "Bad user")
		return
	}

	if len(amount) == 0 || !govalidator.IsNumeric(amount) {
		response_writer.ValidationError(w, "Bad amount")
		return
	}

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &db.User{Identifier: int(id)}
	err := db.Connection.Where(user).FirstOrCreate(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// create new user if does not exist
			user.Name = "new"
			user.Balance = 0
			err := db.Connection.Create(user).Error
			if err != nil {
				response_writer.HandleError(w, err)
				return
			}
		} else {
			response_writer.HandleError(w, err)
			return
		}
	}

	value, _ := strconv.ParseFloat(amount, 64)

	user.Balance += value
	err = db.Connection.Save(user).Error
	if err != nil {
		response_writer.HandleError(w, err)
		return
	}

	response_writer.JsonResponse(w, user, 200)
}
