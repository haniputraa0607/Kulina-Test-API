package app

import (
	"rest-api/config/log_config"
	"rest-api/database"
	"rest-api/route"

	"github.com/gin-gonic/gin"
)

func App() {

	database.ConnectDatabase()

	log_config.DefaultLogging("logs/file.log")

	app := gin.Default()
	
	route.InitRoute(app)

	app.Run(":8080")
}