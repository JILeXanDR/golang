package app

import (
	"github.com/joho/godotenv"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/JILeXanDR/golang/db"
	"github.com/JILeXanDR/golang/http_handlers"
	"os"
	"golang.org/x/net/websocket"
	"github.com/JILeXanDR/golang/app/ws"
)

var loaded = false

// it loads env vars, creates db connection and adds route handlers
func Create(envFile string) {

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = db.Connect()
	if err != nil {
		panic(err)
	}

	// TODO FIX
	//defer db.Connection.Close()

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/", http_handlers.IndexPageHandler)

	r.HandleFunc("/balance", http_handlers.GetBalanceHandler).Methods("POST")
	r.HandleFunc("/deposit", http_handlers.DepositMoneyHandler).Methods("POST")
	r.HandleFunc("/withdraw", http_handlers.WithdrawMoneyHandler).Methods("POST")
	r.HandleFunc("/transfer", http_handlers.TransferMoneyHandler).Methods("POST")

	r.HandleFunc("/orders", http_handlers.GetOrdersHandler).Methods("GET")
	r.HandleFunc("/orders", http_handlers.CreateOrderHandler).Methods("POST")
	r.HandleFunc("/orders/{id:[0-9]+}/cancel", http_handlers.CancelOrderHandler).Methods("GET")   // TODO change to POST
	r.HandleFunc("/orders/{id:[0-9]+}/confirm", http_handlers.ConfirmOrderHandler).Methods("GET") // TODO change to POST
	r.HandleFunc("/orders/{id:[0-9]+}/deliver", http_handlers.DeliverOrderHandler).Methods("GET") // TODO change to POST
	r.HandleFunc("/orders/{id:[0-9]+}/process", http_handlers.ProcessOrderHandler).Methods("GET") // TODO change to POST
	r.HandleFunc("/find-address", http_handlers.FindAddressHandler).Methods("GET")

	r.Handle("/io", websocket.Handler(ws.EchoHandler))

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(os.Getenv("ROOT_DIR")+"/public"))))
	http.Handle("/", r)
}

func CreateTest() {
	if !loaded {
		Create("./../.env")
		loaded = true
	}
}
