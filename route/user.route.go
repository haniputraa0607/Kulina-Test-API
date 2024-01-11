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
	userRoute.GET("/address", user_controller.GetAddress)
	userRoute.POST("/product", user_controller.GetProduct)
	userRoute.POST("/order", user_controller.CreateOrder)
	userRoute.GET("/order", user_controller.GetOrder)

}
