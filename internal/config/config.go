package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
	Redis    RedisConfig    `mapstructure:"redis"`
	Token    TokenConfig    `mapstructure:"token"`
	AI       AIConfig       `mapstructure:"ai"`
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

type AIConfig struct {
	GeminiAPIKey string
}

type TokenConfig struct {
	PasetoV4LocalKey string
}

// MustLoadConfig reads configurations from file or environment variables
func MustLoadConfig(path string) *Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env file not found")
	}

	cfg := new(Config)
	viper.AddConfigPath(path)
	viper.SetConfigName("server")
	viper.SetConfigType("yaml")

	viper.SetDefault("app.env", DevEnv)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config file: %v", err))
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v", err))
	}

	cfg.AI.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	cfg.Token.PasetoV4LocalKey = os.Getenv("PASETO_V4_LOCAL_KEY")

	return cfg
}
