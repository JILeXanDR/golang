package app

import (
	"github.com/joho/godotenv"
	"log"
	"github.com/JILeXanDR/golang/app/db"
	"os"
)

var loaded = false

// it loads env vars, creates db connection and adds route handlers
func Create(envFile string) {

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Getenv("ENV")) == 0 {
		panic("ENV is not defined. Check your .env file")
	}

	log.Printf("Going to start application in '%v' env", os.Getenv("ENV"))

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
