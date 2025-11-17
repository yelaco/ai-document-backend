package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppEnv string

const (
	DevEnv  AppEnv = "dev"
	StgEnv  AppEnv = "stg"
	ProdEnv AppEnv = "prod"
)

type Config struct {
	AppEnv AppEnv `mapstructure:"app.env"`
	Host   string `mapstructure:"app.host"`
	Port   string `mapstructure:"app.port"`
	DBName string `mapstructure:"database.name"`
	DBHost string `mapstructure:"database.host"`
	DBUser string `mapstructure:"database.user"`
	DBPass string `mapstructure:"database.password"`
}

// MustLoadConfig reads configurations from file or environment variables
func MustLoadConfig(path string) *Config {
	cfg := new(Config)
	viper.AddConfigPath(path)
	viper.SetConfigName("server")
	viper.SetConfigType("yaml")

	viper.SetDefault("app.env", DevEnv)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config file: %v", err))
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v", err))
	}
	return cfg
}
