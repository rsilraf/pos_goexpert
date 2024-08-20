package config

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type config struct {
	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTTTL    int    `mapstructure:"JWT_TTL_SEC"`
	JWT       *jwtauth.JWTAuth
}

func Load() (*config, error) {
	var cfg *config
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	cfg.JWT = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, nil
}
