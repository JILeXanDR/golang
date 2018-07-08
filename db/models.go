package db

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Identifier int     `gorm:"not null;unique"`
	Name       string  `gorm:"type:varchar(64);not null"`
	Balance    float64 `gorm:"not null"`
}

//type PostgresJson struct {
//	json.RawMessage
//}
//
//func (j PostgresJson) Value() (driver.Value, error) {
//	return j.MarshalJSON()
//}
//
//func (j *PostgresJson) Scan(src interface{}) error {
//	if bytes, ok := src.([]byte); ok {
//		return json.Unmarshal(bytes, j)
//
//	}
//	return errors.New(fmt.Sprint("Failed to unmarshal JSON from DB", src))
//}

type Order struct {
	gorm.Model
	List            postgres.Jsonb `gorm:"type:json;not null";json:"list"`
	Phone           string         `gorm:"type:varchar(20);not null";json:"phone"`
	DeliveryAddress string         `gorm:"type:varchar(200);not null";json:"delivery_address"`
}
