package router

import (
	"economic_api/user_web/api"
	"economic_api/user_web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.GET("list", middlewares.JWTAuth(), api.GetUserList)
		userRouter.POST("login_password", api.LoginByPassword)
		userRouter.POST("register", api.RegisterUser)
	}
}
