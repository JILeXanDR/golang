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

	var orders = []Order{
		Order{
			Phone:             "0939411685",
			DeliveryAddressId: "EnHQstGD0LvQuNGG0Y8g0JTQvtCx0YDQvtCy0L7Qu9GM0YHRjNC60L7Qs9C",
			DeliveryAddress:   "вулиця Добровольського, 6, Черкаси, Черкаська область, Україна",
			Comment:           "4 подъезд квартира 117",
			Name:              "Саша",
			Status:            STATUS_CREATED,
			List:              postgres.Jsonb{metadata},
		},
		Order{
			Phone:             "0939411685",
			DeliveryAddressId: "EnTQstGD0LvQuNGG0Y8g0J7RgdGC0LDRhNGW0Y8g0JTQsNGI0LrQvtCy0LjRh9CwLCAzLCDQp9C10YDQutCw0YHQuCwg0KfQtdGA0LrQsNGB0YzQutCwINC",
			DeliveryAddress:   "вулиця Остафія Дашковича, 3, Черкаси, Черкаська область, Україна",
			Comment:           "",
			Name:              "Саша",
			Status:            STATUS_CREATED,
			List:              postgres.Jsonb{metadata},
		},
	}

	for _, order := range orders {
		Connection.Create(&order)
	}
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
