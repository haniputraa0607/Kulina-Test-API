package supplier_controller

import (
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
	if errBcrypt  != nil {

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
		Name    : &supplierRegister.Name,
		Email   : &supplierRegister.Email,
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
		"data"   : supplier,
		"token"  : token,
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
		"token"  : token,
	})

	return

}

func RegisterStore(ctx *gin.Context)  {
	
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
		City      : &addressRequest.City,
		Province  : &addressRequest.Province,
		Detail    : &addressRequest.Detail,
		Longitude : &addressRequest.Longitude,
		Latitude  : &addressRequest.Latitude,
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
	} else{
		if errDB := database.DB.Table("stores").Create(&store).Error; errDB != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": errDB.Error(),
			})
			return
		}
	}


	ctx.JSON(200, gin.H{
		"message": "Success",
		"data"  : store,
	})

	return
}

func SellingArea(ctx *gin.Context)  {
	
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
		StoreID	  : store.ID,
		PostalCode: &sellingAreaRequest.PostalCode,
		City      : &sellingAreaRequest.City,
		Province  : &sellingAreaRequest.Province,
		Detail    : &sellingAreaRequest.Detail,
		Longitude : &sellingAreaRequest.Longitude,
		Latitude  : &sellingAreaRequest.Latitude,
	}

	if errDB := database.DB.Table("store_selling_areas").Create(&sellingArea).Error; errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errDB.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data"  : sellingArea,
	})

	return
}

func GetProduct(ctx *gin.Context)  {

	product := new([]model.Product)

	err := database.DB.Table("products").Find(&product).Error

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": product,
	})

	return
}

func CreateProduct(ctx *gin.Context)  {
	
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

	productRequest := new(entity.CreateProductRequest)

	if errReq := ctx.ShouldBind(&productRequest); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	products := productRequest.Product
	
	productData := make([]model.StoreProduct, 0)

	for _, product := range products {
		productExist := new(model.StoreProduct)
		if database.DB.Table("store_products").Where("store_id = ?", store.ID).Where("product_id = ?", product.ID).Find(&productExist); productExist.ID != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Product Exist",
			})
			return
		}

		newProduct := model.StoreProduct{
			StoreID:      store.ID,
			ProductID:    &product.ID,
			Price:        &product.Price,
			IsPurchased:  &product.IsPurchased,
		}

		productData = append(productData, newProduct)
	}

	if errDB := database.DB.Table("store_products").Create(&productData).Error; errDB != nil {
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

func ListProduct(ctx *gin.Context)  {
	
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

	storeProduct := new([]model.StoreProduct)

	err := database.DB.Preload("Product").Table("store_products").Where("store_id = ?", store.ID).Find(&storeProduct).Error
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": storeProduct,
	})

	return
}