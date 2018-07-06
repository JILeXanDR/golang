package db

import (
	"github.com/jinzhu/gorm"
)

var Connection *gorm.DB

func Connect() {

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=golang dbname=golang password=golang")
	if err != nil {
		panic(err)
	}

	Connection = db

	//Connection.LogMode(true)

	Connection.DropTableIfExists(&User{})
	Connection.AutoMigrate(&User{})

	Connection.Create(&User{Identifier: 1, Name: "Alexandr", Balance: 1000})
	Connection.Create(&User{Identifier: 2, Name: "Bob", Balance: 1000})
	Connection.Create(&User{Identifier: 3, Name: "Test", Balance: 1000})
}
