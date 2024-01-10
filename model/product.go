package model

import "time"

type Product struct {
	ID        *uint64    `gorm:"primary_key:auto_increment" json:"id"`
	StoreID   *uint64 
	Store     *Store    `gorm:"foreignKey:StoreID"` 
	Name      *string    `gorm:"type:varchar(255)" json:"name"`
	Price     *int64     `gorm:"type:integer" json:"price"`
	IsPurchased *bool `gorm:"type:boolean" json:"is_purchased"`
	ParentID   *uint64 
	Parent     *Product    `gorm:"foreignKey:ParentID"` 
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}