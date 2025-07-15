package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	REDIS_HOST     string `mapstructure:"REDIS_HOST"`
	REDIS_PORT     string `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	APP_PORT       string `mapstructure:"APP_PORT"`
	APP_ENV        string `mapstructure:"APP_ENV"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	env := "local"
	envConfigFileName := fmt.Sprintf(".env.%s", env)

	viper.AutomaticEnv()

	viper.AddConfigPath("./.secrets")
	viper.SetConfigName(envConfigFileName)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Using environment variables.")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal config: %w", err)
	}

	fmt.Println("config:", config)
	return config, nil
}
