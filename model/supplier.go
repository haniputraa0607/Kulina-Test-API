package model

import "time"

type Supplier struct {
	ID        *uint64    `gorm:"primary_key:auto_increment" json:"id"`
	Name      *string    `gorm:"type:varchar(255)" json:"name"`
	Email     *string    `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password  *string    `gorm:"->;<-;not null" json:"-"`
	Store     *Store     `json:"store,omitempty"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

type Store struct {
	ID               *uint64 `gorm:"primary_key:auto_increment" json:"id"`
	SupplierID       *uint64
	Supplier         *Supplier `gorm:"foreignKey:SupplierID"`
	PostalCode       *string   `gorm:"type:varchar(255)" json:"postal_code"`
	City             *string   `gorm:"type:varchar(255)" json:"city"`
	Province         *string   `gorm:"type:varchar(255)" json:"province"`
	Detail           *string   `gorm:"type:text" json:"detail"`
	Longitude        *float64  `gorm:"type:double" json:"longitude"`
	Latitude         *float64  `gorm:"type:double" json:"latitude"`
	StoreSellingArea *[]StoreSellingArea
	CreatedAt        *time.Time `json:"-"`
	UpdatedAt        *time.Time `json:"-"`
}

type StoreSellingArea struct {
	ID         *uint64 `gorm:"primary_key:auto_increment" json:"id"`
	StoreID    *uint64
	Store      *Store     `gorm:"foreignKey:StoreID"`
	PostalCode *string    `gorm:"type:varchar(255)" json:"postal_code"`
	City       *string    `gorm:"type:varchar(255)" json:"city"`
	Province   *string    `gorm:"type:varchar(255)" json:"province"`
	Detail     *string    `gorm:"type:text" json:"detail"`
	Longitude  *float64   `gorm:"type:double" json:"longitude"`
	Latitude   *float64   `gorm:"type:double" json:"latitude"`
	CreatedAt  *time.Time `json:"-"`
	UpdatedAt  *time.Time `json:"-"`
}
