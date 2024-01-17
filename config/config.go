package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// DBConfig struct to hold database configuration parameters
type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
	APIKey   string
}

// LoadDBConfig loads database configuration from environment variables or .env file
func LoadDBConfig() DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT", 3308),
		DBName:   os.Getenv("DB_NAME"),
		APIKey:   os.Getenv("OPENAI_API_KEY"),
	}
}

// getEnvAsInt retrieves an environment variable as an integer or returns the default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v\n", key, err)
		return defaultValue
	}

	return value
}

// GetDBConnectionString returns the formatted database connection string
func (c *DBConfig) GetDBConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.DBName)
}
