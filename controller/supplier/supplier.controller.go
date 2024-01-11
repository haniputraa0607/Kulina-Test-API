package supplier_controller

import (
	"fmt"
	"net/http"
	"rest-api/database"
	"rest-api/entity"
	"rest-api/model"
	"rest-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {

	supplierRegister := new(entity.RegisterSupplierRequest)

	if errReq := ctx.ShouldBind(&supplierRegister); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	password, errBcrypt := bcrypt.GenerateFromPassword([]byte(supplierRegister.Password), bcrypt.MinCost)
	if errBcrypt != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errBcrypt.Error(),
		})

		return
	}

	supplierEmailExist := new(model.Supplier)
	if database.DB.Table("suppliers").Where("email = ?", supplierRegister.Email).Find(&supplierEmailExist); supplierEmailExist.Email != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email already exist",
		})
		return
	}

	passwordString := string(password)
	supplier := model.Supplier{
		Name:     &supplierRegister.Name,
		Email:    &supplierRegister.Email,
		Password: &passwordString,
	}

	if errDB := database.DB.Table("suppliers").Create(&supplier).Error; errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errDB.Error(),
		})
		return
	}

	token, errToken := utils.GenerateTokenSupplier(supplier)

	if errToken != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "Failed generetad token",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    supplier,
		"token":   token,
	})

	return

}

func Login(ctx *gin.Context) {

	supplierLogin := new(entity.LoginSupplierRequest)

	if errReq := ctx.ShouldBind(&supplierLogin); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	supplier := new(model.Supplier)
	if errSupplier := database.DB.Table("suppliers").Where("email = ?", supplierLogin.Email).Find(&supplier).Error; errSupplier != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Credential not valid",
		})

		return
	}

	if supplier.ID == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Supplier not found",
		})
		return
	}

	if errPassword := bcrypt.CompareHashAndPassword([]byte(*supplier.Password), []byte(supplierLogin.Password)); errPassword != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"message": "Invalid password",
		})

		return
	}

	token, errToken := utils.GenerateTokenSupplier(*supplier)

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

func RegisterStore(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	addressRequest := new(entity.StoreSupplierRequest)

	if errReq := ctx.ShouldBind(&addressRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	store := model.Store{
		SupplierID: &UserMiddleware.ID,
		PostalCode: &addressRequest.PostalCode,
		City:       &addressRequest.City,
		Province:   &addressRequest.Province,
		Detail:     &addressRequest.Detail,
		Longitude:  &addressRequest.Longitude,
		Latitude:   &addressRequest.Latitude,
	}

	storeExist := new(model.Store)
	if database.DB.Table("stores").Where("supplier_id = ?", UserMiddleware.ID).Find(&storeExist); storeExist.ID != nil {
		if errUpdateData := database.DB.Table("stores").Where("id = ?", storeExist.ID).Updates(&store).Error; errUpdateData != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errUpdateData.Error(),
			})
			return
		}
		store.ID = storeExist.ID
	} else {
		if errDB := database.DB.Table("stores").Create(&store).Error; errDB != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errDB.Error(),
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    store,
	})

	return
}

func SellingArea(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	sellingAreaRequest := new(entity.StoreSellingAreaSupplierRequest)

	if errReq := ctx.ShouldBind(&sellingAreaRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	store := new(model.Store)
	if database.DB.Table("stores").Where("supplier_id = ?", UserMiddleware.ID).Find(&store); store.ID != nil {
		if errUpdateData := database.DB.Table("stores").Where("id = ?", store.ID).Updates(&store).Error; errUpdateData != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errUpdateData.Error(),
			})
			return
		}
	}

	sellingArea := model.StoreSellingArea{
		StoreID:    store.ID,
		PostalCode: &sellingAreaRequest.PostalCode,
		City:       &sellingAreaRequest.City,
		Province:   &sellingAreaRequest.Province,
		Detail:     &sellingAreaRequest.Detail,
		Longitude:  &sellingAreaRequest.Longitude,
		Latitude:   &sellingAreaRequest.Latitude,
	}

	existSellingArea := new(model.StoreSellingArea)
	if database.DB.Table("store_selling_areas").Where("store_id = ?", &sellingArea.StoreID).Find(&existSellingArea); existSellingArea.ID != nil {
		fmt.Print("12131233")
		if errUpdateData := database.DB.Table("store_selling_areas").Where("store_id = ?", &sellingArea.StoreID).Updates(&sellingArea).Error; errUpdateData != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errUpdateData.Error(),
			})
			return
		}

		sellingArea.ID = existSellingArea.ID
	} else {
		fmt.Print("123")
		if errDB := database.DB.Table("store_selling_areas").Create(&sellingArea).Error; errDB != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errDB.Error(),
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    sellingArea,
	})

	return
}

func CreateProduct(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	store := new(model.Store)
	if database.DB.Table("stores").Where("supplier_id = ?", UserMiddleware.ID).Find(&store); store.ID != nil {
		if errUpdateData := database.DB.Table("stores").Where("id = ?", store.ID).Updates(&store).Error; errUpdateData != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errUpdateData.Error(),
			})
			return
		}
	}

	productRequest := new(entity.ProductRequest)

	if errReq := ctx.ShouldBind(&productRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	if productRequest.ID != 0 {

		existProduct := new(model.Product)
		if database.DB.Table("products").Where("id = ?", productRequest.ID).Find(&existProduct); existProduct.ID == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Product not found",
			})
			return
		}

		productData := model.Product{
			StoreID:     store.ID,
			Name:        &productRequest.Name,
			Price:       &productRequest.Price,
			IsPurchased: existProduct.IsPurchased,
			ParentID:    existProduct.ID,
		}

		if errDB := database.DB.Table("products").Create(&productData).Error; errDB != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errDB.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data": productData,
		})

		return

	} else {

		productData := model.Product{
			StoreID:     store.ID,
			Name:        &productRequest.Name,
			Price:       &productRequest.Price,
			IsPurchased: &productRequest.IsPurchased,
		}

		if errDB := database.DB.Table("products").Create(&productData).Error; errDB != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errDB.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data": productData,
		})

		return
	}

}

func GetProduct(ctx *gin.Context) {

	decodeToken := ctx.MustGet("decode_token").(jwt.MapClaims)
	UserMiddleware, errUserMiddleware := utils.UserMid(decodeToken)

	if errUserMiddleware {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Error",
		})

		return
	}

	store := new(model.Store)
	if database.DB.Table("stores").Where("supplier_id = ?", UserMiddleware.ID).Find(&store); store.ID != nil {
		if errUpdateData := database.DB.Table("stores").Where("id = ?", store.ID).Updates(&store).Error; errUpdateData != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errUpdateData.Error(),
			})
			return
		}
	}

	product := new([]model.Product)
	errProduct := database.DB.Table("products").Where("store_id = ?", store.ID).Find(&product).Error

	if errProduct != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	otherProduct := new([]model.Product)
	errOtherProduct := database.DB.Preload("Store.Supplier").Table("products").Where("store_id != ?", store.ID).Find(&otherProduct).Error
	if errOtherProduct != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"product_store":       product,
		"product_other_store": otherProduct,
	})

	return
}
