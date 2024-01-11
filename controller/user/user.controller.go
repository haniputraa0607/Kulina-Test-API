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
				if tx.Joins("JOIN order_products ON orders.id = order_products.order_id").Table("orders").Where("user_id = ?", UserMiddleware.ID).Where("DATE(order_date) = ?", time.Now().UTC().Format("2006-01-02")).Where("order_products.product_id = ?", product.ID).Find(&existOrder); existOrder.ID != nil {
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

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    orderResponse,
	})

	return
}
