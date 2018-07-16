package app

import (
	"github.com/gorilla/mux"
	"github.com/JILeXanDR/golang/app/http_handlers"
	"golang.org/x/net/websocket"
	"github.com/JILeXanDR/golang/app/ws"
	"net/http"
	"os"
	"github.com/JILeXanDR/golang/app/response_writer"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path
	if filePath := os.Getenv("ROOT_DIR") + path; fileExists(filePath) {
		http.ServeFile(w, r, filePath)
	} else {
		http_handlers.IndexPageHandler(w, r)
	}
}

func apiNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response_writer.JsonMessageResponse(w, "Endpoint is not defined", 404)
}

func apiIndexHandler(w http.ResponseWriter, r *http.Request) {
	response_writer.JsonMessageResponse(w, "Base API endpoint", 200)
}

func GetRouter() (router *mux.Router) {

	router = mux.NewRouter()
	router.StrictSlash(true)

	router.Use(loggingMiddleware)
	router.Use(headerMiddleware)

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.HandleFunc("/", http_handlers.IndexPageHandler)

	api := router.PathPrefix("/api").Subrouter()

	api.NotFoundHandler = http.HandlerFunc(apiNotFoundHandler)

	api.HandleFunc("", apiIndexHandler)

	api.HandleFunc("/users", http_handlers.GetUsers).Methods("GET")
	api.HandleFunc("/balance", http_handlers.GetBalanceHandler).Methods("POST")
	api.HandleFunc("/deposit", http_handlers.DepositMoneyHandler).Methods("POST")
	api.HandleFunc("/withdraw", http_handlers.WithdrawMoneyHandler).Methods("POST")
	api.HandleFunc("/transfer", http_handlers.TransferMoneyHandler).Methods("POST")

	api.HandleFunc("/confirm-phone", http_handlers.ConfirmPhoneHandler).Methods("POST")
	api.HandleFunc("/check-code", http_handlers.CheckCodeHandler).Methods("POST")

	orders := api.PathPrefix("/orders").Subrouter()

	orders.HandleFunc("", http_handlers.GetOrdersHandler).Methods("GET")
	orders.HandleFunc("", http_handlers.CreateOrderHandler).Methods("POST")

	order := orders.PathPrefix("/{id:[0-9]+}").Subrouter()

	order.HandleFunc("", http_handlers.GetOrderHandler).Methods("GET")
	order.HandleFunc("/cancel", http_handlers.CancelOrderHandler).Methods("GET")   // TODO change to POST
	order.HandleFunc("/confirm", http_handlers.ConfirmOrderHandler).Methods("GET") // TODO change to POST
	order.HandleFunc("/deliver", http_handlers.DeliverOrderHandler).Methods("GET") // TODO change to POST
	order.HandleFunc("/process", http_handlers.ProcessOrderHandler).Methods("GET") // TODO change to POST

	router.HandleFunc("/find-address", http_handlers.FindAddressHandler).Methods("GET")

	router.Handle("/io", websocket.Handler(ws.EchoHandler))

	return router
}
