package route

import (
	supplier_controller "rest-api/controller/supplier"
	"rest-api/middleware"

	"github.com/gin-gonic/gin"
)

func SupplierRoute(app *gin.RouterGroup) {

	route := app.Group("/supplier")

	route.POST("register", supplier_controller.Register)
	route.POST("login", supplier_controller.Login)

	supplierRoute := route.Group("/", middleware.AuthMiddleware)
	supplierRoute.POST("/store", supplier_controller.RegisterStore)
	supplierRoute.POST("/store-selling-are", supplier_controller.SellingArea)
	supplierRoute.GET("/product", supplier_controller.GetProduct)
	supplierRoute.POST("/product", supplier_controller.CreateProduct)
	supplierRoute.GET("/product-list", supplier_controller.ListProduct)
}