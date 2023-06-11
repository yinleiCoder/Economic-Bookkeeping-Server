package initialize

import (
	"economic_api/user_web/middlewares"
	router2 "economic_api/user_web/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.Cors())
	apiGroup := router.Group("/user/v1")
	router2.InitUserRouter(apiGroup)
	router2.InitBaseRouter(apiGroup)
	return router
}
