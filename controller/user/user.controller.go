package user_controller

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

	userRegister := new(entity.RegisterUserRequest)

	if errReq := ctx.ShouldBind(&userRegister); errReq != nil {

		ctx.AbortWithStatusJSON(400, gin.H{
			"message": errReq.Error(),
		})

		return
	}

	password, errBcrypt := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.MinCost)
	if errBcrypt  != nil {

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
		Name    : &userRegister.Name,
		Email   : &userRegister.Email,
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
		"data"   : user,
		"token"  : token,
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
		"token"  : token,
	})

	return

}

func Address(ctx *gin.Context)  {
	
	
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
		UserID: &UserMiddleware.ID,
		PostalCode: &addressRequest.PostalCode,
		City: &addressRequest.City,
		Province: &addressRequest.Province,
		Detail: &addressRequest.Detail,
		Longitude: &addressRequest.Longitude,
		Latitude: &addressRequest.Latitude,
	}

	if errDB := database.DB.Table("addresses").Create(&address).Error; errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errDB.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data"  : address,
	})

	return
}



