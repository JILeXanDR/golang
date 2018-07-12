package main

import (
	"log"
	"github.com/JILeXanDR/golang/app"
	"os"
	"net/http"
)

func main() {

	app.Create("./.env")

	port := ":" + os.Getenv("PORT")

	http.Handle("/", app.GetRouter())

	log.Println("Start server at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
