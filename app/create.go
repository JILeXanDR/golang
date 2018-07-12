package app

import (
	"github.com/joho/godotenv"
	"log"
	"github.com/JILeXanDR/golang/db"
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
}

func CreateTest() {
	if !loaded {
		Create("./../.env")
		loaded = true
	}
}
