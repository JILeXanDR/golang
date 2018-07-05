package main

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/asaskevich/govalidator"
	"strconv"
)

var db *gorm.DB

// TODO Реализовать сервис по работе со счетами пользователей. В сервисе есть пользователи (id, name, balance).
// TODO Данные необходимо хранить в postgresql. Реализовать валидацию. В случае любой ошибки валидации отдавать 422 ошибку.
// TODO можно использовать любые framework’и

type User struct {
	gorm.Model
	Identifier int     `gorm:"not null;unique"`
	Name       string  `gorm:"type:varchar(64);not null"`
	Balance    float64 `gorm:"not null"`
}

type ValidationErrorResponse struct {
	Message string `json:"message"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

func TransferMoney(from *User, to *User, amount float64) (err error) {

	transaction := db.Begin()

	from.Balance -= amount
	err = db.Save(from).Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	to.Balance += amount
	err = db.Save(to).Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	transaction.Commit()

	return nil
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Просмотреть баланс
// Request: GET /balance?user=101
// Response: 200 OK {“balance”: 1000}
func getBalanceHandler(w http.ResponseWriter, r *http.Request) {

	var (
		userId = r.URL.Query().Get("user")
	)

	if !govalidator.IsInt(userId) {
		jsonResponse(w, &ValidationErrorResponse{Message: "Query parameter 'user' is invalid"}, 422)
		return
	}

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &User{Identifier: int(id)}
	db.Where(user).First(user)

	if user.ID == 0 {
		jsonResponse(w, &ValidationErrorResponse{Message: "User does not exist"}, 404)
		return
	}

	jsonResponse(w, BalanceResponse{Balance: user.Balance}, 200)
}

// зачислять деньги на счет пользователям (создать пользователя, если еще не существует)
// Request: POST /deposit {“user”: 101, “amount”: 100}
// Response: 200 OK
func depositMoneyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		userId = r.URL.Query().Get("user")
		amount = r.URL.Query().Get("amount")
	)

	if len(userId) == 0 || !govalidator.IsNumeric(userId) {
		jsonResponse(w, &ValidationErrorResponse{Message: "Bad user"}, 422)
		return
	}

	if len(amount) == 0 || !govalidator.IsNumeric(amount) {
		jsonResponse(w, &ValidationErrorResponse{Message: "Bad amount"}, 422)
		return
	}

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &User{Identifier: int(id)}
	db.Where(user).FirstOrCreate(user)

	if user.ID == 0 {
		user.Name = "new"
		user.Balance = 0
		db.Create(user)
	}

	value, _ := strconv.ParseFloat(amount, 64)

	user.Balance += value
	db.Save(user)

	jsonResponse(w, user, 200)
}

// снимать деньги со счетов
// Request: POST /withdraw {“user”: 101, “amount”: 50}
// Response: 200 OK
func withdrawMoneyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		userId = r.URL.Query().Get("user")
		amount = r.URL.Query().Get("amount")
	)

	id, _ := strconv.ParseInt(userId, 10, 0)

	user := &User{Identifier: int(id)}
	db.Where(user).First(user)

	if user.ID == 0 {
		jsonResponse(w, &ValidationErrorResponse{Message: "User does not exist"}, 404)
		return
	}

	value, _ := strconv.ParseFloat(amount, 64)

	if user.Balance < value {
		jsonResponse(w, &ValidationErrorResponse{Message: "Not enough money"}, 400)
		return
	}

	user.Balance -= value
	db.Save(user)

	jsonResponse(w, user, 200)
}

// переводить деньги от одного пользователя другому.
// Request: POST /transfer {“from”: 101, “to”: 205, amount: 25}
// Response: 200 OK
func transferMoneyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		from   = r.URL.Query().Get("from")
		to     = r.URL.Query().Get("to")
		amount = r.URL.Query().Get("amount")
	)

	fromId, _ := strconv.ParseInt(from, 10, 0)
	toId, _ := strconv.ParseInt(to, 10, 0)
	value, _ := strconv.ParseFloat(amount, 64)

	fromUser := &User{Identifier: int(fromId)}
	err := db.Where(fromUser).First(fromUser).Error
	if err != nil {
		panic(err)
	}

	if fromUser.ID == 0 {
		jsonResponse(w, &ValidationErrorResponse{Message: "Sender user does not exist"}, 404)
		return
	}

	if fromUser.Balance < value {
		jsonResponse(w, &ValidationErrorResponse{Message: "Not enough money"}, 404)
		return
	}

	toUser := &User{Identifier: int(toId)}
	db.Where(toUser).First(toUser)

	if toUser.ID == 0 {
		jsonResponse(w, &ValidationErrorResponse{Message: "Recipient user does not exist"}, 404)
		return
	}

	err = TransferMoney(fromUser, toUser, value)
	if err != nil {
		jsonResponse(w, "Could not transfer money", 500)
		return
	}

	jsonResponse(w, fromUser, 200)
}

func main() {

	db, _ = gorm.Open("postgres", "host=localhost port=5432 user=golang dbname=golang password=golang")
	defer db.Close()

	db.LogMode(true)

	db.DropTableIfExists(&User{})
	db.AutoMigrate(&User{})

	user1 := &User{Identifier: 1, Name: "Alexandr", Balance: 1000}
	user2 := &User{Identifier: 2, Name: "Bob", Balance: 1000}
	user3 := &User{Identifier: 3, Name: "Test", Balance: 1000}

	db.Create(user1)
	db.Create(user2)
	db.Create(user3)

	http.HandleFunc("/balance", getBalanceHandler)
	http.HandleFunc("/deposit", depositMoneyHandler)
	http.HandleFunc("/withdraw", withdrawMoneyHandler)
	http.HandleFunc("/transfer", transferMoneyHandler)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
