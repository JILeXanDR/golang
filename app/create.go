package app

import (
	"github.com/joho/godotenv"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/JILeXanDR/golang/db"
	"github.com/JILeXanDR/golang/http_handlers"
)

var loaded = false

// it loads env vars, creates db connection and adds route handlers
func Create(env string) {

	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = db.Connect()
	if err != nil {
		panic(err)
	}

	//defer db.Connection.Close()

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/", http_handlers.HomePageHandler).Methods("GET")
	r.HandleFunc("/balance", http_handlers.GetBalanceHandler).Methods("POST")
	r.HandleFunc("/deposit", http_handlers.DepositMoneyHandler).Methods("POST")
	r.HandleFunc("/withdraw", http_handlers.WithdrawMoneyHandler).Methods("POST")
	r.HandleFunc("/transfer", http_handlers.TransferMoneyHandler).Methods("POST")
	r.HandleFunc("/order", http_handlers.OrderHandler).Methods("POST")
	r.HandleFunc("/orders", http_handlers.GetOrdersHandler).Methods("GET")

	http.Handle("/", r)
}

func CreateTest() {
	if !loaded {
		Create("./../.env")
		loaded = true
	}
}
