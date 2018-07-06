package db

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Identifier int     `gorm:"not null;unique"`
	Name       string  `gorm:"type:varchar(64);not null"`
	Balance    float64 `gorm:"not null"`
}
