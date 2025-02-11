package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBPort      string
	JWTSecret   string
	LogLevel    string
}

func LoadConfig() *Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBUser:      getEnv("DB_USER", "coins_user"),
		DBPassword:  getEnv("DB_PASSWORD", "coins_pass"),
		DBName:      getEnv("DB_NAME", "coins_db"),
		DBPort:      getEnv("DB_PORT", "5432"),
		JWTSecret:   getEnv("JWT_SECRET", "mysecretkey"),
		LogLevel: 	 getEnv("LOG_lEVEL", "WARN"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}