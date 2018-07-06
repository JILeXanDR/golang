package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var Connection *gorm.DB

func GetUserBalance(userId int) (res float64, err error) {
	var user = &User{}
	err = Connection.Where(&User{Identifier: userId}).First(user).Error
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}

func sender(channel chan string) {
	channel <- "ping"
	channel <- "ping"
	channel <- "ping"
}

func receiver(channel chan string) {
	for {
		log.Println(<-channel)
	}
}

func Connect() {

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=golang dbname=golang password=golang")
	if err != nil {
		panic(err)
	}

	Connection = db

	var channel = make(chan string)

	log.Println("before")

	go sender(channel)
	go receiver(channel)

	log.Println("after")

	//Connection.LogMode(true)

	Connection.DropTableIfExists(&User{})
	Connection.AutoMigrate(&User{})

	Connection.Create(&User{Identifier: 1, Name: "Alexandr", Balance: 1000})
	Connection.Create(&User{Identifier: 2, Name: "Bob", Balance: 1000})
	Connection.Create(&User{Identifier: 3, Name: "Test", Balance: 1000})
}
