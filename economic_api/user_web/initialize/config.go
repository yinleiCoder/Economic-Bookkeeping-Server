package initialize

import (
	"economic_api/user_web/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// read env information from windows env.
	debug := GetEnvInfo("Economic_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user_web/%s-production.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user_web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("application config %s has changed!: %v", e.Name, global.ServerConfig)
	})
}
