package model

import "time"

type User struct {
	ID        *uint64    `gorm:"primary_key:auto_increment" json:"id"`
	Name      *string    `gorm:"type:varchar(255)" json:"name"`
	Email     *string    `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password  *string    `gorm:"->;<-;not null" json:"-"`
	Address   *[]Address `json:"address,omitempty"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

type Address struct {
	ID         *uint64 `gorm:"primary_key:auto_increment" json:"id"`
	UserID     *uint64
	User       *User      `gorm:"foreignKey:UserID"`
	PostalCode *string    `gorm:"type:varchar(255)" json:"postal_code"`
	City       *string    `gorm:"type:varchar(255)" json:"city"`
	Province   *string    `gorm:"type:varchar(255)" json:"province"`
	Detail     *string    `gorm:"type:text" json:"detail"`
	Longitude  *float64   `gorm:"type:double" json:"longitude"`
	Latitude   *float64   `gorm:"type:double" json:"latitude"`
	CreatedAt  *time.Time `json:"-"`
	UpdatedAt  *time.Time `json:"-"`
}
