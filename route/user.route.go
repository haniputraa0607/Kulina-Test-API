package route

import (
	user_controller "rest-api/controller/user"
	"rest-api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(app *gin.RouterGroup) {

	route := app.Group("/user")

	route.POST("register", user_controller.Register)
	route.POST("login", user_controller.Login)

	userRoute := route.Group("/", middleware.AuthMiddleware)
	userRoute.POST("/address", user_controller.Address)



}