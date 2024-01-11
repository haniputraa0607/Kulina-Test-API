package user_controller

import (
	"errors"
	"fmt"
	"net/http"
	"rest-api/database"
	"rest-api/entity"
	"rest-api/model"
	"rest-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {

	userRegister := new(entity.RegisterUserRequest)

	if errReq := ctx.ShouldBind(&userRegister); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	password, errBcrypt := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.MinCost)
	if errBcrypt != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errBcrypt.Error(),
		})

		return
	}

	userEmailExist := new(model.User)
	if database.DB.Table("users").Where("email = ?", userRegister.Email).Find(&userEmailExist); userEmailExist.Email != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email already exist",
		})
		return
	}

	passwordString := string(password)
	user := model.User{
		Name:     &userRegister.Name,
		Email:    &userRegister.Email,
		Password: &passwordString,
	}

	if errDB := database.DB.Table("users").Create(&user).Error; errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errDB.Error(),
		})
		return
	}

	token, errToken := utils.GenerateTokenUser(user)

	if errToken != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "Failed generetad token",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
		"token":   token,
	})

	return

}

func Login(ctx *gin.Context) {

	userLogin := new(entity.LoginUserRequest)

	if errReq := ctx.ShouldBind(&userLogin); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	user := new(model.User)
	if errUser := database.DB.Table("users").Where("email = ?", userLogin.Email).Find(&user).Error; errUser != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Credential not valid",
		})

		return
	}

	if user.ID == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if errPassword := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(userLogin.Password)); errPassword != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Invalid password",
		})

		return
	}

	token, errToken := utils.GenerateTokenUser(*user)

	if errToken != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "Failed generetad token",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"token":   token,
	})

	return

}

func Address(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	addressRequest := new(entity.AddressUserRequest)

	if errReq := ctx.ShouldBind(&addressRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	address := model.Address{
		UserID:     &UserMiddleware.ID,
		PostalCode: &addressRequest.PostalCode,
		City:       &addressRequest.City,
		Province:   &addressRequest.Province,
		Detail:     &addressRequest.Detail,
		Longitude:  &addressRequest.Longitude,
		Latitude:   &addressRequest.Latitude,
	}

	if errDB := database.DB.Table("addresses").Create(&address).Error; errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errDB.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    address,
	})

	return
}

func GetAddress(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	address := new([]model.Address)

	if errUser := database.DB.Table("addresses").Where("user_id = ?", UserMiddleware.ID).Find(&address).Error; errUser != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Error",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    address,
	})

	return
}

func GetProduct(ctx *gin.Context) {

	listProductUserRequest := new(entity.ListProductUserRequest)

	if errReq := ctx.ShouldBind(&listProductUserRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	address := new(model.Address)

	if errUser := database.DB.Table("addresses").Where("id = ?", listProductUserRequest.AddressID).Find(&address).Error; errUser != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Error",
		})

		return
	}

	products := new([]model.Product)
	haversine := fmt.Sprintf(
		"6371 * acos(cos(radians(%f)) * cos(radians(Latitude)) * cos(radians(Longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(Latitude)))",
		*address.Latitude,
		*address.Longitude,
		*address.Latitude,
	)

	if err := database.DB.Joins("JOIN store_selling_areas ON products.store_id = store_selling_areas.store_id").Select("products.*, " + haversine + " AS distance").Order("name").Order("distance").Table("products").Find(&products).Error; err != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Error",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    products,
	})

	return
}

func CreateOrder(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	orderRequest := new(entity.OrderRequest)

	if errReq := ctx.ShouldBind(&orderRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	currentTime := time.Now()
	order := model.Order{
		UserID:        &UserMiddleware.ID,
		UserAddressID: &orderRequest.AddressID,
		OrderDate:     &currentTime,
	}

	errDBTrx := database.DB.Transaction(func(tx *gorm.DB) error {

		if errDB := tx.Table("orders").Create(&order).Error; errDB != nil {

			return errDB
		}

		var totalPrice int64 = 0
		orderProduct := make([]model.OrderProduct, 0)

		for _, productOrder := range orderRequest.Products {

			product := new(model.Product)
			if tx.Table("products").Where("id = ?", &productOrder.ID).Find(&product); product.ID == nil {
				return errors.New("Product not found")
			}

			if product.IsPurchased != nil && *product.IsPurchased == true {
				existOrder := new(model.Order)
				if tx.Joins("JOIN order_products ON orders.id = order_products.order_id").Table("orders").Where("user_id = ?", UserMiddleware.ID).Where("DATE(order_date) = ?", time.Now().UTC().Format("2006-01-02")).Where("orders.is_cancelled = 0").Where("order_products.product_id = ?", product.ID).Find(&existOrder); existOrder.ID != nil {
					return errors.New("Product one-time purchase")
				}
			}

			totalPrice += productOrder.Quantity * *product.Price
			order.TotalPrice = &totalPrice

			total := productOrder.Quantity * *product.Price

			eachOrderProduct := model.OrderProduct{
				OrderID:    order.ID,
				ProductID:  product.ID,
				Quantity:   &productOrder.Quantity,
				Price:      product.Price,
				TotalPrice: &total,
			}

			orderProduct = append(orderProduct, eachOrderProduct)

			if errUpdateData := tx.Table("orders").Where("id = ?", &order.ID).Updates(&order).Error; errUpdateData != nil {
				return errUpdateData
			}
		}

		if errDBProduct := tx.Table("order_products").Create(&orderProduct).Error; errDBProduct != nil {
			return errDBProduct
		}

		return nil
	})

	if errDBTrx != nil {

		ctx.AbortWithStatusJSON(404, gin.H{
			"message": errDBTrx.Error(),
		})

		return

	} else {

		ctx.JSON(200, gin.H{
			"message": "Success",
			"data":    order,
		})

		return
	}

}

func GetOrder(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	orderResponse := new([]entity.OrderResponse)
	if err := database.DB.
		Joins("JOIN addresses ON orders.user_address_id = addresses.id").
		Select("orders.*, addresses.detail AS address").
		Where("orders.user_id = ?", UserMiddleware.ID).
		Table("orders").
		Find(&orderResponse).Error; err != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": err.Error(),
		})

		return
	}

	unpaid := make([]entity.OrderResponse, 0)
	cancelled := make([]entity.OrderResponse, 0)
	waiting := make([]entity.OrderResponse, 0)
	received := make([]entity.OrderResponse, 0)
	cancelDelivery := make([]entity.OrderResponse, 0)

	for _, order := range *orderResponse {

		products := new([]entity.ProductResponse)
		if err := database.DB.
			Joins("JOIN products ON order_products.product_id = products.id").
			Select("order_products.*, products.name AS name").
			Where("order_id = ?", order.ID).
			Table("order_products").
			Find(&products).Error; err != nil {
			ctx.AbortWithStatusJSON(404, gin.H{
				"message": err.Error(),
			})

			return
		}

		order.Products = *products

		if order.PaidDate == nil && !order.IsCancelled {
			unpaid = append(unpaid, order)
		} else if order.IsCancelled {
			cancelled = append(cancelled, order)
		}

		currentTime := time.Now()
		if order.IsPaid && currentTime.After(*order.EstimateDate) && !order.CancelDelivery {
			received = append(received, order)
		} else if order.IsPaid && currentTime.Before(*order.EstimateDate) && !order.CancelDelivery {
			waiting = append(waiting, order)
		} else if order.CancelDelivery {
			cancelDelivery = append(cancelDelivery, order)
		}
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data": gin.H{
			"unpaid":          unpaid,
			"cancelled":       cancelled,
			"waiting":         waiting,
			"received":        received,
			"cancel_delivery": cancelDelivery,
		},
	})

	return
}

func CancelOrder(ctx *gin.Context) {

	cancelOrder := new(entity.CancelOrderRequest)

	if errReq := ctx.ShouldBind(&cancelOrder); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	order := new(model.Order)
	if database.DB.Table("orders").Where("id = ?", cancelOrder.ID).Find(&order); order.ID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Order not found",
		})
		return
	}

	if order.IsPaid != nil && *order.IsPaid == true {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Order has been paid",
		})
		return
	}

	newValue := true
	order.IsCancelled = &newValue
	if errUpdateData := database.DB.Table("orders").Where("id = ?", cancelOrder.ID).Updates(&order).Error; errUpdateData != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errUpdateData.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
	})

	return

}

func PayOrder(ctx *gin.Context) {

	payOrder := new(entity.PayOrderRequest)

	if errReq := ctx.ShouldBind(&payOrder); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	order := new(model.Order)
	if database.DB.Table("orders").Where("id = ?", payOrder.ID).Find(&order); order.ID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Order not found",
		})
		return
	}

	if order.IsCancelled != nil && *order.IsCancelled == true {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Order has been cancelled",
		})
		return
	}

	if *order.TotalPrice > payOrder.Price {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Not enough",
		})
		return
	}

	newValue := true
	currentTime := time.Now()
	today5pm := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 17, 0, 0, 0, currentTime.Location())
	var deliveryDate time.Time
	if currentTime.After(today5pm) {
		// Get the date for tomorrow at 7 am
		deliveryDate = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+4, 17, 0, 0, 0, currentTime.Location())
	} else {
		deliveryDate = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+3, 17, 0, 0, 0, currentTime.Location())
	}

	order.IsPaid = &newValue
	order.PaidDate = &currentTime
	order.EstimateDate = &deliveryDate

	fmt.Print(deliveryDate)

	if errUpdateData := database.DB.Table("orders").Where("id = ?", payOrder.ID).Updates(&order).Error; errUpdateData != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errUpdateData.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
	})

	return
}

func CancelDelivery(ctx *gin.Context) {

	cancelDelivery := new(entity.CancelDeliveryRequest)

	if errReq := ctx.ShouldBind(&cancelDelivery); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	order := new(model.Order)
	if database.DB.Table("orders").Where("id = ?", cancelDelivery.ID).Find(&order); order.ID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Order not found",
		})
		return
	}

	if order.IsPaid != nil && *order.IsPaid != true {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Order unpaid",
		})
		return
	}

	if order.CancelDelivery != nil && *order.CancelDelivery == true {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Already cancelled delivery",
		})
		return
	}

	newValue := true
	order.CancelDelivery = &newValue
	order.CancelReason = &cancelDelivery.Reason

	if errUpdateData := database.DB.Table("orders").Where("id = ?", cancelDelivery.ID).Updates(&order).Error; errUpdateData != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errUpdateData.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
	})

	return

}
