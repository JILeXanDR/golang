package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Connection *gorm.DB

func GetUserBalance(userId int) (float64) {
	var user = &User{}
	Connection.Where(&User{Identifier: userId}).First(user)

	return user.Balance
}

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
