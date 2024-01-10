package route

import "github.com/gin-gonic/gin"

func InitRoute(app *gin.Engine) {


	route := app.Group("/api")

	UserRoute(route)
	SupplierRoute(route)

}