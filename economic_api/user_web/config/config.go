package config

type UserServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type AliSMSConfig struct {
	AccessKeyId     string `mapstructure:"key"`
	AccessKeySecret string `mapstructure:"secret"`
	SignName        string `mapstructure:"sign_name"`
	TemplateCode    string `mapstructure:"template_code"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Expire   int    `mapstructure:"expire"`
}

type ServerConfig struct {
	Name            string            `mapstructure:"name"`
	Port            int               `mapstructure:"port"`
	UserServiceInfo UserServiceConfig `mapstructure:"user_service"`
	JWTInfo         JWTConfig         `mapstructure:"jwt"`
	AliSMSInfo      AliSMSConfig      `mapstructure:"sms"`
	RedisInfo       RedisConfig       `mapstructure:"redis"`
}
