package db

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
	"strings"
)

//type model struct {
//	//	ID        uint      `gorm:"primary_key" json:"id"`
//	//	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
//	//	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
//	ID        uint       `gorm:"primary_key"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt *time.Time `sql:"index"`
//}
//
//func (m *model) BeforeCreate() (err error) {
//	now := time.Now()
//	m.CreatedAt = now
//	m.UpdatedAt = now
//	return
//}
//
//func (m *model) BeforeSave() (err error) {
//	m.UpdatedAt = time.Now()
//	return
//}

type User struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Identifier int       `gorm:"not null;unique" json:"identifier"`
	Name       string    `gorm:"type:varchar(64);not null" json:"name"`
	Balance    float64   `gorm:"not null" json:"balance"`
	Phone      string    `json:"phone"`
	Orders     []Order   `json:"orders"`
}

const (
	// 1. присваивается всем новым заказам
	STATUS_CREATED = "created"

	// 2. подтверждение заказа менеджером
	STATUS_CONFIRMED = "confirmed"
	// 2. отмена заказа менеджером
	STATUS_CANCELED = "canceled"

	// 3. послее взятия заказа в обработку
	STATUS_PROCESSING = "processing"

	// 4. после доставки заказа получателю
	STATUS_DELIVERED = "delivered"
)

type Order struct {
	ID                uint           `gorm:"primary_key" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	List              postgres.Jsonb `gorm:"type:json;not null" json:"list"`
	Phone             string         `gorm:"type:varchar(20);not null" json:"phone"`
	DeliveryAddressId string         `gorm:"type:varchar(1000);not null" json:"delivery_address_id"`
	DeliveryAddress   string         `gorm:"type:varchar(1000);not null" json:"delivery_address"`
	Name              string         `gorm:"type:varchar(100);not null" json:"name"`
	Comment           string         `gorm:"type:varchar(255);not null" json:"comment"`
	Status            string         `gorm:"type:varchar(20);not null" json:"status"`
	CancelReason      string         `gorm:"type:varchar(255);not null" json:"cancel_reason"`
	DeliveredAt       time.Time      `gorm:"default NULL" json:"delivered_at"`
	User              User           `gorm:"foreignkey:UserId" json:"user"`
	UserId            uint           `json:"user_id"`
}

func (order *Order) AfterFind() (err error) {
	// FIXME
	order.DeliveryAddress = strings.Replace(order.DeliveryAddress, ", Черкаси, Черкаська область, Україна", "", -1)
	return
}
