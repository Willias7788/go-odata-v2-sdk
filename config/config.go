package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	SAPHost     string `mapstructure:"SAP_HOST"`
	SAPUsername string `mapstructure:"SAP_USERNAME"`
	SAPPassword string `mapstructure:"SAP_PASSWORD"`
	SAPClient   string `mapstructure:"SAP_CLIENT"` // Optional: sap-client param
}

// LoadConfig reads configuration from environment variables or .env file 
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Try to read .env file, but don't fail if it doesn't exist (Docker/Prod runtime)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
