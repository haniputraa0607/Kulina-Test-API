package model

import "time"

type Product struct {
	ID        uint64    `gorm:"primary_key:auto_increment" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Price     int64     `gorm:"type:integer" json:"price"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}