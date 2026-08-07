package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig represents the strongly typed configuration structure
type AppConfig struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	AppPort  int    `mapstructure:"APP_PORT"`
	DBHost   string `mapstructure:"DB_HOST"`
	DBPort   int    `mapstructure:"DB_PORT"`
	DBUser   string `mapstructure:"DB_USER"`
	DBPass   string `mapstructure:"DB_PASS"`
	DBName   string `mapstructure:"DB_NAME"`
}

// LoadConfig initializes Viper and unmarshals the config into AppConfig
func LoadConfig() (*AppConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found or readable, relying on environment variables: %v", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
