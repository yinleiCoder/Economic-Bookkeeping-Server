package main

import (
	"economic_api/user_web/global"
	"economic_api/user_web/initialize"
	"fmt"
	"go.uber.org/zap"
)

/*
*
Web API Layer.
*/
func main() {
	initialize.InitLogger()

	initialize.InitConfig()

	router := initialize.Routers()

	zap.S().Debugf("launch web port: %d", global.ServerConfig.Port)
	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("launch web api failed: ", err.Error())
	}
}
