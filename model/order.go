package model

import "time"

type Order struct {
	ID            *uint64 `gorm:"primary_key:auto_increment" json:"id"`
	UserID        *uint64
	User          *User `gorm:"foreignKey:UserID"`
	UserAddressID *uint64
	UserAddress   *Address        `gorm:"foreignKey:UserAddressID"`
	TotalPrice    *int64          `gorm:"type:integer" json:"total_price"`
	IsPaid        *bool           `gorm:"type:boolean;default:false" json:"is_paid"`
	OrderDate     *time.Time      `json:"-"`
	PaidDate      *time.Time      `json:"-"`
	EstimateDate  *time.Time      `json:"-"`
	OrderProduct  *[]OrderProduct `json:"order_products,omitempty"`
	CreatedAt     *time.Time      `json:"-"`
	UpdatedAt     *time.Time      `json:"-"`
}

type OrderProduct struct {
	ID         *uint64 `gorm:"primary_key:auto_increment" json:"id"`
	OrderID    *uint64
	Order      *Order `gorm:"foreignKey:OrderID"`
	ProductID  *uint64
	Product    *Product   `gorm:"foreignKey:ProductID"`
	Quantity   *int64     `gorm:"type:integer" json:"qty"`
	Price      *int64     `gorm:"type:integer" json:"price"`
	TotalPrice *int64     `gorm:"type:integer" json:"total_price"`
	CreatedAt  *time.Time `json:"-"`
	UpdatedAt  *time.Time `json:"-"`
}
