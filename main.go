package main

import (
	"net/http"
	"log"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	database "github.com/JILeXanDR/golang/db"
	handlers "github.com/JILeXanDR/golang/http_handlers"
)

// TODO Данные необходимо хранить в postgresql. Реализовать валидацию. В случае любой ошибки валидации отдавать 422 ошибку.

func main() {

	database.Connect()
	defer database.Connection.Close()

	http.HandleFunc("/balance", handlers.GetBalanceHandler)
	http.HandleFunc("/deposit", handlers.DepositMoneyHandler)
	http.HandleFunc("/withdraw", handlers.WithdrawMoneyHandler)
	http.HandleFunc("/transfer", handlers.TransferMoneyHandler)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
