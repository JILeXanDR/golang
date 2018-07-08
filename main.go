package main

import (
	"net/http"
	"log"
	database "github.com/JILeXanDR/golang/db"
	handlers "github.com/JILeXanDR/golang/http_handlers"
	"github.com/joho/godotenv"
	"os"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.Connect()
	if err != nil {
		panic(err)
	}

	defer database.Connection.Close()

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/", handlers.HomePageHandler).Methods("GET")
	r.HandleFunc("/balance", handlers.GetBalanceHandler).Methods("POST")
	r.HandleFunc("/deposit", handlers.DepositMoneyHandler).Methods("POST")
	r.HandleFunc("/withdraw", handlers.WithdrawMoneyHandler).Methods("POST")
	r.HandleFunc("/transfer", handlers.TransferMoneyHandler).Methods("POST")
	r.HandleFunc("/order", handlers.OrderHandler).Methods("POST")
	r.HandleFunc("/orders", handlers.GetOrdersHandler).Methods("GET")

	http.Handle("/", r)

	port := ":" + os.Getenv("PORT")

	log.Println("Start server at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
