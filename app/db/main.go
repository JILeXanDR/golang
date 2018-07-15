package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"fmt"
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

var Connection *gorm.DB

func migrations() {

	var models = []interface{}{&User{}, &Order{}}

	Connection.DropTableIfExists(models...)
	Connection.AutoMigrate(models...)

	var users = []User{
		{Identifier: 1, Name: "Alexandr", Balance: 1000, Phone: "380939411685"},
		{Identifier: 2, Name: "Bob", Balance: 1000, Phone: "380980808421"},
		{Identifier: 3, Name: "Test", Balance: 1000},
	}

	for _, user := range users {
		log.Println(Connection.Create(&user).Error)
	}

	var geList = func(list []string) postgres.Jsonb {
		metadata, err := json.Marshal(list)
		if err != nil {
			panic(err)
		}
		return postgres.Jsonb{metadata}
	}

	fiveDaysAgoDuration, _ := time.ParseDuration("-120h")
	fiveDaysAgo := time.Now().Add(fiveDaysAgoDuration)

	var orders = []Order{
		{
			Phone:             "380939411685",
			DeliveryAddressId: "EnHQstGD0LvQuNGG0Y8g0JTQvtCx0YDQvtCy0L7Qu9GM0YHRjNC60L7Qs9C",
			DeliveryAddress:   "вулиця Добровольського, 6, Черкаси, Черкаська область, Україна",
			Comment:           "4 подъезд квартира 117",
			Name:              "Саша",
			Status:            STATUS_CREATED,
			List:              geList([]string{"Латте апельсиновый (фабрика кофе)", "мак-меню"}),
			CreatedAt:         fiveDaysAgo,
			//User:              users[0],
			UserId: users[0].ID,
		},
		{
			Phone:             "380939411685",
			DeliveryAddressId: "EnTQstGD0LvQuNGG0Y8g0J7RgdGC0LDRhNGW0Y8g0JTQsNGI0LrQvtCy0LjRh9CwLCAzLCDQp9C10YDQutCw0YHQuCwg0KfQtdGA0LrQsNGB0YzQutCwINC",
			DeliveryAddress:   "вулиця Остафія Дашковича, 3, Черкаси, Черкаська область, Україна",
			Comment:           "",
			Name:              "Саша",
			Status:            STATUS_CONFIRMED,
			List:              geList([]string{"Латте апельсиновый (фабрика кофе)", "мак-меню"}),
			//User:              users[1],
			UserId: users[1].ID,
		},
		{
			Phone:             "380939411685",
			DeliveryAddressId: "EnTQstGD0LvQuNGG0Y8g0J7RgdGC0LDRhNGW0Y8g0JTQsNGI0LrQvtCy0LjRh9CwLCAzLCDQp9C10YDQutCw0YHQuCwg0KfQtdGA0LrQsNGB0YzQutCwINC",
			DeliveryAddress:   "вулиця Остафія Дашковича, 3, Черкаси, Черкаська область, Україна",
			Comment:           "",
			Name:              "Саша",
			Status:            STATUS_CANCELED,
			List:              geList([]string{"Латте апельсиновый (фабрика кофе)", "мак-меню"}),
			//User:              users[1],
			UserId: users[1].ID,
		},
		{
			Phone:             "380939411685",
			DeliveryAddressId: "EnTQstGD0LvQuNGG0Y8g0J7RgdGC0LDRhNGW0Y8g0JTQsNGI0LrQvtCy0LjRh9CwLCAzLCDQp9C10YDQutCw0YHQuCwg0KfQtdGA0LrQsNGB0YzQutCwINC",
			DeliveryAddress:   "вулиця Остафія Дашковича, 3, Черкаси, Черкаська область, Україна",
			Comment:           "",
			Name:              "Саша",
			Status:            STATUS_PROCESSING,
			List:              geList([]string{"Круасан с малиной (Львовский круасан)", "Латте большой"}),
			//User:              users[1],
			UserId: users[1].ID,
		},
		{
			Phone:             "380939411685",
			DeliveryAddressId: "EnTQstGD0LvQuNGG0Y8g0J7RgdGC0LDRhNGW0Y8g0JTQsNGI0LrQvtCy0LjRh9CwLCAzLCDQp9C10YDQutCw0YHQuCwg0KfQtdGA0LrQsNGB0YzQutCwINC",
			DeliveryAddress:   "вулиця Остафія Дашковича, 3, Черкаси, Черкаська область, Україна",
			Comment:           "ничего",
			Name:              "Саша",
			Status:            STATUS_DELIVERED,
			List:              geList([]string{"Латте апельсиновый (фабрика кофе)"}),
			CreatedAt:         fiveDaysAgo,
			//User:              users[1],
			UserId: users[1].ID,
		},
	}

	for _, order := range orders {
		log.Println(Connection.Create(&order).Error)
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
