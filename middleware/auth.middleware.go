package middleware

import (
	"net/http"
	"rest-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware(ctx *gin.Context) {

	bearerToken := ctx.GetHeader("Authorization")

	if !strings.Contains(bearerToken, "Bearer") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}

	token := strings.Replace(bearerToken, "Bearer ", "", -1)

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}

	decodeToken, errDecode := utils.DecodeToken(token)
	if errDecode != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}

	ctx.Set("decode_token", decodeToken)

	ctx.Next()

}