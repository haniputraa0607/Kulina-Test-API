package entity

import (
	"time"
)

type (
	RegisterUserRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Name     string `json:"name" form:"name" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Address  string `json:"address" form:"address"`
	}

	LoginUserRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	AddressUserRequest struct {
		PostalCode string  `json:"postal_code" form:"postal_code" binding:"required"`
		City       string  `json:"city" form:"city" binding:"required"`
		Province   string  `json:"province" form:"province" binding:"required"`
		Detail     string  `json:"detail" form:"detail" binding:"required"`
		Latitude   float64 `json:"latitude" form:"latitude" binding:"required"`
		Longitude  float64 `json:"longitude" form:"longitude" binding:"required"`
	}

	ListProductUserRequest struct {
		AddressID uint64 `json:"address_id" form:"address_id" binding:"required"`
	}

	OrderRequest struct {
		AddressID uint64                `json:"address_id" form:"address_id" binding:"required"`
		Products  []OrderProductRequest `json:"products" form:"products" binding:"required"`
	}

	OrderProductRequest struct {
		ID       uint64 `json:"id" form:"id" binding:"required"`
		Quantity int64  `json:"qty" form:"qty" binding:"required"`
	}

	OrderResponse struct {
		ID             uint64
		Address        string
		TotalPrice     int64
		IsPaid         bool
		IsCancelled    bool
		CancelDelivery bool
		Products       []ProductResponse
		OrderDate      *time.Time
		PaidDate       *time.Time
		EstimateDate   *time.Time
	}

	ProductResponse struct {
		ID              uint64
		OrderResponseId uint64
		Name            string
		Quantity        int64
		Price           int64
		TotalPrice      int64
	}

	CancelOrderRequest struct {
		ID uint64 `json:"id" form:"id" binding:"required"`
	}
	PayOrderRequest struct {
		ID    uint64 `json:"id" form:"id" binding:"required"`
		Price int64  `json:"price" form:"price" binding:"required"`
	}
	CancelDeliveryRequest struct {
		ID     uint64 `json:"id" form:"id" binding:"required"`
		Reason string `json:"reason" form:"reason" binding:"required"`
	}
)
