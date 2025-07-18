// load configurations to --
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB_USER           string `mapstructure:"DB_USER"`
	DB_HOST           string `mapstructure:"DB_HOST"`
	DB_PORT           string `mapstructure:"DB_PORT"`
	DB_PASS           string `mapstructure:"DB_PASS"`
	DB_NAME           string `mapstructure:"DB_NAME"`
	DB_SSLMODE        string `mapstructure:"DB_SSLMODE"`
	APP_ENV           string `mapstructure:"APP_ENV"`
	APP_PORT          string `mapstructure:"APP_PORT"`
	USER_PORT         string `mapstructure:"USER_PORT"`
	GRPC_PORT         string `mapstructure:"GRPC_PORT"`
	REDIS_HOST        string `mapstructure:"REDIS_HOST"`
	REDIS_PORT        string `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD    string `mapstructure:"REDIS_PASSWORD"`
	NOTIFICATION_PORT string `mapstructure:"NOTIFICATION_PORT"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	env := "local"
	envConfigFileName := fmt.Sprintf(".env.%s", env)

	// ^viper is a module for handling application configuration from multiple sources in a unified way(for ex yaml, json, env variable)
	// ^if we use getenv instead of viper, we need to load every credential one by one, but using viper, we can load all credential at once, by using struct of config credentials
	viper.AutomaticEnv() //Reading environment variables

	viper.AddConfigPath("./.secrets")
	viper.SetConfigName(envConfigFileName)
	viper.SetConfigType("env")

	err := viper.ReadInConfig() //this is the main function which is reading the credentials automatcially without loading them
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Using environment variables.")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err) //agar kuch error aaya to wo struct ke taur pr nil bhejega or error ki value bhej dega
		}
	}

	err = viper.Unmarshal(&config) //convert json into struct

	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal config: %w", err)
	}

	fmt.Println("config:", config)
	return config, nil

}
