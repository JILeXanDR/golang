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
	// first status
	STATUS_CREATED = "created"

	// second status
	STATUS_CONFIRMED = "confirmed"
	STATUS_CANCELED  = "canceled"

	// third status
	STATUS_PROCESSING = "processing"

	// last status
	STATUS_DELIVERED = "delivered"
)

type Order struct {
	gorm.Model
	List              postgres.Jsonb `gorm:"type:json;not null" json:"list"`
	Phone             string         `gorm:"type:varchar(20);not null" json:"phone"`
	DeliveryAddressId string         `gorm:"type:varchar(1000);not null" json:"delivery_address_id"`
	DeliveryAddress   string         `gorm:"type:varchar(1000);not null" json:"delivery_address"`
	Name              string         `gorm:"type:varchar(100);not null" json:"name"`
	Comment           string         `gorm:"type:varchar(255);not null" json:"comment"`
	Status            string         `gorm:"type:varchar(20);not null" json:"status"`
	CancelReason      string         `gorm:"type:varchar(255);not null" json:"cancel_reason"`
	DeliveredAt       time.Time      `gorm:"default null" json:"delivered_at"`
}
