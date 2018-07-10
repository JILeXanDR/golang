package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"fmt"
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
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

func migrations() {

	var models = []interface{}{&User{}, &Order{}}

	Connection.DropTableIfExists(models...)
	Connection.AutoMigrate(models...)

	// test users
	Connection.Create(&User{Identifier: 1, Name: "Alexandr", Balance: 1000})
	Connection.Create(&User{Identifier: 2, Name: "Bob", Balance: 1000})
	Connection.Create(&User{Identifier: 3, Name: "Test", Balance: 1000})

	var list = []string{
		"Латте апельсиновый (фабрика кофе)",
		"мак-меню",
	}

	metadata, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	// test orders
	Connection.Create(&Order{
		Phone:           "0939411685",
		DeliveryAddress: "Добровольського 6",
		Comment:         "4 подъезд квартира 117",
		Name:            "Саша",
		Status:          STATUS_CREATED,
		List:            postgres.Jsonb{metadata},
	})
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

	migrations()

	return nil
}
