package config

import (
	
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

const (
	JwtSecretKey      = "your-secret-key"
	JwtExpirationTime = 24 * time.Hour // Example: JWT token expires in 24 hours
)

// LoadDBConfig loads database configuration from environment variables or .env file
func LoadDBConfig() DBConfig {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file. Using default values or values from the environment.")
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
		log.Printf("Using default value for %s: %d\n", key, defaultValue)
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


const ExpireMinute = 30 // Atur sesuai kebutuhan Anda
const DiffUTC = 7       // Atur sesuai kebutuhan Anda
// GetTimeWithZone mengembalikan waktu dengan zona waktu tertentu
func GetTimeWithZone(zone string) time.Time {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		// Handle kesalahan jika zona waktu tidak valid
		return time.Now() // Mengembalikan waktu saat ini jika terjadi kesalahan
	}
	return time.Now().In(loc)
}

// CalculateExpirationTime menghitung waktu kedaluwarsa berdasarkan menit yang diberikan
func CalculateExpirationTime(minutes int) time.Time {
	return time.Now().Add(time.Minute * time.Duration(minutes))
}
