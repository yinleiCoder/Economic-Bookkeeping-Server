package router

import (
	"economic_api/user_web/api"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("common")
	{
		baseRouter.GET("captcha", api.GetCaptcha)
		baseRouter.POST("sms", api.SendSMS)
	}
}
