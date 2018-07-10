package main

import (
	"net/http"
	"log"
	"os"
	"github.com/JILeXanDR/golang/app"
	"fmt"
)

func main() {

	app.Create("./.env")

	port := ":" + os.Getenv("PORT")

	var counter = 0

	go func() {
		counter += 1
	}()

	go func() {
		counter += 1
	}()

	fmt.Println(counter)

	log.Println("Start server at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
