package entity

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
		AddressID string  `json:"address_id" form:"address_id" binding:"required"`
	}
)