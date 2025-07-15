package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string   `mapstructure:"environment"`
	Server      Server   `mapstructure:"server"`
	Database    Database `mapstructure:"database"`
	Auth        Auth     `mapstructure:"auth"`
}

type Server struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Auth struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpiration int    `mapstructure:"jwt_expiration"`
}

func Load() (*Config, error) {
	// Check for custom config file from environment
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config"
	} else {
		// Remove .yaml extension if present
		if len(configFile) > 5 && configFile[len(configFile)-5:] == ".yaml" {
			configFile = configFile[:len(configFile)-5]
		}
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "controlme")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("auth.jwt_expiration", 86400) // 24 hours

	// Read environment variables
	viper.AutomaticEnv()

	// Bind specific environment variables to config keys
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.username", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
