package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"fmt"
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

func Connect() (err error) {

	var conn = fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		return err
	}

	Connection = db

	if os.Getenv("DB_LOG_MODE_ENABLED") == "true" {
		Connection.LogMode(true)
	}

	var models = []interface{}{&User{}, &Order{}}

	Connection.DropTableIfExists(models...)
	Connection.AutoMigrate(models...)

	Connection.Create(&User{Identifier: 1, Name: "Alexandr", Balance: 1000})
	Connection.Create(&User{Identifier: 2, Name: "Bob", Balance: 1000})
	Connection.Create(&User{Identifier: 3, Name: "Test", Balance: 1000})

	return nil
}
