package http_handlers

import (
	"net/http"
	"strconv"
	"github.com/JILeXanDR/golang/app/db"
	"fmt"
	"github.com/JILeXanDR/golang/app/response_writer"
	"github.com/JILeXanDR/golang/app/services"
)

// переводить деньги от одного пользователя другому.
// Request: POST /transfer {“from”: 101, “to”: 205, amount: 25}
// Response: 200 OK
func TransferMoneyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		from   = r.URL.Query().Get("from")
		to     = r.URL.Query().Get("to")
		amount = r.URL.Query().Get("amount")
	)

	fromId, _ := strconv.ParseInt(from, 10, 0)
	toId, _ := strconv.ParseInt(to, 10, 0)
	value, _ := strconv.ParseFloat(amount, 64)

	if fromId == 0 || toId == 0 {
		response_writer.JsonMessageResponse(w, "Bad user ids", 424)
		return
	}

	internalErr, logicErr := services.TransferMoney(int(fromId), int(toId), value)
	if internalErr != nil {
		response_writer.InternalServerError(w, internalErr)
		return
	} else if logicErr != nil {
		response_writer.HandleError(w, logicErr)
		return
	}

	user := &db.User{}
	db.Connection.Where(&db.User{Identifier: int(fromId)}).First(user)

	recipient := &db.User{}
	db.Connection.Where(&db.User{Identifier: int(toId)}).First(recipient)

	message := fmt.Sprintf(
		"Operation was succesful. You have transfered '%v' money to '%v (Identifier=%v)'. Your current balance is '%v'",
		value,
		recipient.Name,
		recipient.Identifier,
		user.Balance,
	)

	response_writer.JsonMessageResponse(w, message, 200)
}
