package http_handlers

import (
	"net/http"
	"strconv"
	"github.com/JILeXanDR/golang/common"
	"github.com/JILeXanDR/golang/db"
	"fmt"
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
		common.JsonMessageResponse(w, "Bad user ids", 424)
		return
	}

	if value <= 0 {
		common.JsonMessageResponse(w, "Amount should be greater than 0 and and integer", 424)
		return
	}

	internalErr, logicErr := common.TransferMoney(int(fromId), int(toId), value)
	if internalErr != nil {
		common.InternalServerError(w)
		return
	} else if logicErr != nil {
		common.HandleError(w, logicErr)
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

	common.JsonMessageResponse(w, message, 200)
}
