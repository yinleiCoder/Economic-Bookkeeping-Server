package middlewares

import (
	"economic_api/user_web/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomJwtClaims)
		if currentUser.RoleId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "you don't have permission",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
