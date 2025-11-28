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
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Token    TokenConfig    `mapstructure:"token"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type AppConfig struct {
	Env  AppEnv `mapstructure:"env"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type TokenConfig struct {
	PasetoV4LocalKey string `mapstructure:"PASETO_V4_LOCAL_KEY"`
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
