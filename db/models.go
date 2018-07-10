package db

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type User struct {
	gorm.Model
	Identifier int     `gorm:"not null;unique"`
	Name       string  `gorm:"type:varchar(64);not null"`
	Balance    float64 `gorm:"not null"`
}

const (
	STATUS_CREATED    = "created"
	STATUS_CONFIRMED  = "confirmed"
	STATUS_PROCESSING = "processing"
	STATUS_DELIVERED  = "delivered"
)

type Order struct {
	gorm.Model
	List            postgres.Jsonb `gorm:"type:json;not null";json:"list"`
	Phone           string         `gorm:"type:varchar(20);not null";json:"phone"`
	DeliveryAddress string         `gorm:"type:varchar(200);not null";json:"delivery_address"`
	Name            string         `gorm:"type:varchar(100);not null";json:"name"`
	Comment         string         `gorm:"type:varchar(255);not null";json:"comment"`
	Status          string         `gorm:"type:varchar(20);not null";json:"status"`
	DeliveredAt     time.Time      `gorm:"default null";json:"delivered_at"`
}
