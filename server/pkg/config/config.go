// // pkg/config/config.go
package config

// import (
// 	"fmt"
// 	"time"

// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	// Server
// 	ServerPort int
// 	ServerHost string

// 	// JWT
// 	JWTSecret        string
// 	JWTExpiration    time.Duration
// 	JWTRefreshExpiry time.Duration

// 	// Database
// 	DBHost     string
// 	DBPort     int
// 	DBUser     string
// 	DBPassword string
// 	DBName     string
// 	DBSSLMode  string

// 	// Redis (if needed)
// 	RedisHost     string
// 	RedisPort     int
// 	RedisPassword string
// }

// func LoadConfig() (*Config, error) {
// 	viper.SetConfigName("config")    // config file name without extension
// 	viper.SetConfigType("env")       // config file type
// 	viper.AddConfigPath(".")         // look for config in current directory
// 	viper.AddConfigPath("./config/") // look for config in ./config directory
// 	viper.AutomaticEnv()             // read from environment variables

// 	// Set default values
// 	viper.SetDefault("SERVER_PORT", 8080)
// 	viper.SetDefault("SERVER_HOST", "0.0.0.0")
// 	viper.SetDefault("JWT_EXPIRATION", "24h")
// 	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")
// 	viper.SetDefault("DB_SSL_MODE", "disable")

// 	// Read configuration
// 	if err := viper.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
// 			return nil, fmt.Errorf("error reading config file: %w", err)
// 		}
// 	}

// 	// Parse configuration
// 	config := &Config{
// 		// Server
// 		ServerPort: viper.GetInt("SERVER_PORT"),
// 		ServerHost: viper.GetString("SERVER_HOST"),

// 		// JWT
// 		JWTSecret:        viper.GetString("JWT_SECRET"),
// 		JWTExpiration:    viper.GetDuration("JWT_EXPIRATION"),
// 		JWTRefreshExpiry: viper.GetDuration("JWT_REFRESH_EXPIRY"),

// 		// Database
// 		DBHost:     viper.GetString("DB_HOST"),
// 		DBPort:     viper.GetInt("DB_PORT"),
// 		DBUser:     viper.GetString("DB_USER"),
// 		DBPassword: viper.GetString("DB_PASSWORD"),
// 		DBName:     viper.GetString("DB_NAME"),
// 		DBSSLMode:  viper.GetString("DB_SSL_MODE"),

// 		// Redis
// 		RedisHost:     viper.GetString("REDIS_HOST"),
// 		RedisPort:     viper.GetInt("REDIS_PORT"),
// 		RedisPassword: viper.GetString("REDIS_PASSWORD"),
// 	}

// 	// Validate required fields
// 	if config.JWTSecret == "" {
// 		return nil, fmt.Errorf("JWT_SECRET is required")
// 	}

// 	if config.DBHost == "" || config.DBUser == "" || config.DBName == "" {
// 		return nil, fmt.Errorf("database configuration is incomplete")
// 	}

// 	return config, nil
// }

// // GetDBConnString returns the database connection string
// func (c *Config) GetDBConnString() string {
// 	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
// 		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
// }
