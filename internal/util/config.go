package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds env variables for the app
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig loads the env variables from app.env
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
